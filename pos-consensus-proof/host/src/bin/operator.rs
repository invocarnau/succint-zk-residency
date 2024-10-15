use base64::{prelude::BASE64_STANDARD, Engine};
use clap::Parser;

use prost_types::Timestamp;
use std::str::FromStr;
use url::Url;

use ethers::providers::{Http, Middleware, Provider};

use alloy_primitives::FixedBytes;
use alloy_primitives::{address, Address};
use alloy_provider::ReqwestProvider;
use alloy_rpc_types::BlockNumberOrTag;
use reth_primitives::{hex, Header};

use sp1_cc_client_executor::ContractInput;
use sp1_cc_host_executor::HostExecutor;

use pos_consensus_proof_client::milestone::{ConsensusProofVerifier, MilestoneProofInputs};
use pos_consensus_proof_client::{types, types::heimdall_types};
use pos_consensus_proof_host::types::{Precommit, Validator};
use pos_consensus_proof_host::utils::PosClient;
use pos_consensus_proof_host::ConsensusProver;

/// The arguments for the command.
#[derive(Parser, Debug)]
#[clap(author, version, about, long_about = None)]
pub struct Args {
    #[clap(long)]
    milestone_id: u64,

    #[clap(long)]
    milestone_hash: String,
}

#[tokio::main]
async fn main() -> eyre::Result<()> {
    dotenv::dotenv().ok();

    let args = Args::parse();

    // Setup the logger.
    sp1_sdk::utils::setup_logger();
    let prover = ConsensusProver::new();

    println!("Assembling data for generating proof...");
    let inputs: MilestoneProofInputs = generate_inputs(args).await?;

    println!("Starting to generate proof...");
    let proof = prover.generate_consensus_proof(inputs);

    println!("Successfully generated proof: {:?}", proof.bytes());
    println!("Public values: {:?}", proof.public_values.to_vec());

    proof.save("proof.bin").expect("saving proof failed");
    println!("Proof saved to proof.bin");

    Ok(())
}

pub async fn generate_inputs(args: Args) -> eyre::Result<MilestoneProofInputs> {
    let client = PosClient::default();

    let milestone = client
        .fetch_milestone_by_id(args.milestone_id)
        .await
        .expect("unable to fetch milestone");
    let tx = client
        .fetch_tx_by_hash(args.milestone_hash)
        .await
        .expect("unable to fetch milestone tx");

    let number: u64 = tx.result.height.parse().unwrap();
    let block = client
        .fetch_block_by_number(number + 2)
        .await
        .expect("unable to fetch block");

    let block_precommits = block.result.block.last_commit.precommits;
    let mut precommits: Vec<Vec<u8>> = [].to_vec();
    let mut sigs: Vec<String> = [].to_vec();
    let mut signers: Vec<Address> = [].to_vec();

    let heimdall_chain_id = std::env::var("HEIMDALL_CHAIN_ID").expect("HEIMDALL_CHAIN_ID not set");
    for precommit in block_precommits.iter() {
        // Only add if the side tx result is non empty
        if precommit.side_tx_results.is_some() {
            let serialized_precommit = serialize_precommit(precommit, &heimdall_chain_id);
            precommits.push(serialized_precommit);
            sigs.push(precommit.signature.clone());
            signers.push(Address::from_str(&precommit.validator_address).unwrap());
        }
    }

    // Use the host executor to fetch the required bor block
    let bor_block_number = BlockNumberOrTag::Number(milestone.result.end_block);
    let bor_header = client
        .fetch_bor_header_by_number(bor_block_number)
        .await
        .unwrap();

    // Fetch the validator set
    let validator_set = client
        .fetch_validator_set_by_height(number + 2)
        .await
        .expect("unable to fetch validator set");

    let rpc_url =
        std::env::var("ETH_RPC_URL").unwrap_or_else(|_| panic!("Missing ETH_RPC_URL in env"));

    // Calculate the best l1 block to choose from the last_updated field in validator set
    let l1_block_number = find_best_l1_block(validator_set.result.validators, &rpc_url).await;

    // The L1 block number against which the transaction is executed
    let block_number = BlockNumberOrTag::Number(l1_block_number);

    // Read the verifier contract
    let verifier = std::env::var("VERIFIER").expect("VERIFIER not set");
    let verifier_contract: Address =
        Address::from_str(&verifier).expect("invalid verifier address");

    // Prepare the host executor.
    //
    // Use `ETH_RPC_URL` to get all of the necessary state for the smart contract call.
    let provider = ReqwestProvider::new_http(Url::parse(&rpc_url)?);
    let mut host_executor = HostExecutor::new(provider.clone(), block_number).await?;

    // Keep track of the block hash. Later, validate the client's execution against this.
    let l1_block_hash = host_executor.header.hash_slow();

    // Make the call to the getEncodedValidatorInfo function.
    let call = ConsensusProofVerifier::getEncodedValidatorInfoCall {};
    let _response: ConsensusProofVerifier::getEncodedValidatorInfoReturn = host_executor
        .execute(ContractInput {
            contract_address: verifier_contract,
            caller_address: address!("0000000000000000000000000000000000000000"),
            calldata: call,
        })
        .await?;

    // Make another call to fetch the last verified bor block hash
    let call = ConsensusProofVerifier::lastVerifiedBorBlockHashCall {};
    let response: ConsensusProofVerifier::lastVerifiedBorBlockHashReturn = host_executor
        .execute(ContractInput {
            contract_address: verifier_contract,
            caller_address: address!("0000000000000000000000000000000000000000"),
            calldata: call,
        })
        .await?;

    // Now that we've executed all of the calls, get the `EVMStateSketch` from the host executor.
    let input = host_executor.finalize().await?;
    let state_sketch_bytes = bincode::serialize(&input)?;

    // Fetch the bor block again the block hash read
    let prev_bor_block_hash = response.lastVerifiedBorBlockHash;

    // If the hash is zero, use a default header
    let mut prev_bor_header = Header::default();

    if !prev_bor_block_hash.is_zero() {
        let prev_bor_block_number = client
            .fetch_bor_number_by_hash(prev_bor_block_hash)
            .await
            .unwrap();

        // Fetch the bor header using the number read
        prev_bor_header = client
            .fetch_bor_header_by_number(BlockNumberOrTag::Number(prev_bor_block_number))
            .await
            .unwrap();

        // Check if the hash matches with the original one because a mismatch can happen if block
        // read is not canonical
        assert_eq!(
            prev_bor_header.hash_slow(),
            prev_bor_block_hash,
            "prev bor block hash mismatch"
        );
    }

    Ok(MilestoneProofInputs {
        tx_data: tx.result.tx,
        tx_hash: FixedBytes::from_str(&tx.result.hash).unwrap(),
        precommits,
        sigs,
        signers,
        bor_header,
        prev_bor_header,
        state_sketch_bytes,
        l1_block_hash,
    })
}

pub fn serialize_precommit(precommit: &Precommit, heimdall_chain_id: &String) -> Vec<u8> {
    let timestamp = Timestamp::from_str(&precommit.timestamp).unwrap();
    let parts_header = heimdall_types::CanonicalPartSetHeader {
        total: precommit.block_id.parts.total,
        hash: hex::decode(&precommit.block_id.parts.hash).unwrap(),
    };
    let block_id = Some(heimdall_types::CanonicalBlockId {
        hash: hex::decode(&precommit.block_id.hash).unwrap(),
        parts_header: Some(parts_header),
    });
    let mut sig_bytes: Vec<u8> = [].to_vec();
    let side_tx_result = &precommit.side_tx_results.as_ref().unwrap()[0];
    let sig = side_tx_result.sig.clone().unwrap_or_default();
    if !sig.is_empty() {
        sig_bytes = BASE64_STANDARD.decode(&sig).unwrap();
    }
    let side_tx = heimdall_types::SideTxResult {
        tx_hash: BASE64_STANDARD.decode(&side_tx_result.tx_hash).unwrap(),
        result: side_tx_result.result,
        sig: sig_bytes,
    };
    let vote = heimdall_types::Vote {
        r#type: precommit.type_field,
        height: u64::from_str(&precommit.height).unwrap(),
        round: u64::from_str(&precommit.round).unwrap(),
        block_id,
        timestamp: Some(timestamp),
        chain_id: heimdall_chain_id.to_string(),
        data: [].to_vec(),
        side_tx_results: Some(side_tx),
    };
    types::serialize_precommit(&vote)
}

async fn find_best_l1_block(validator_set: Vec<Validator>, rpc_url: &str) -> u64 {
    let mut max_block = 0;
    for validator in validator_set.iter() {
        let last_updated = u64::from_str(&validator.last_updated).unwrap();
        // Block number is multipled with 100k to get the last updated value in heimdall
        let block_number = last_updated / 100000;
        if block_number > max_block {
            max_block = block_number;
        }
    }

    // Fetch the latest l1 block
    let provider = Provider::<Http>::try_from(rpc_url).unwrap();
    let latest_block = provider.get_block_number().await.unwrap().as_u64();

    // Because we can only access last 256 blocks in solidity, if the max_block is beyond that, use
    // the latest one.
    if max_block < latest_block - 256 {
        max_block = latest_block;
    }

    println!("Choosing L1 block to generate proof against: {}, latest: {}", max_block, latest_block);

    // TODO: Make sure no staking event happened after this block
    max_block
}
