use alloy_provider::ReqwestProvider;
use clap::Parser;
use rsp_client_executor::{ChainVariant, CHAIN_ID_ETH_MAINNET};
use rsp_host_executor::HostExecutor;
use sp1_sdk::{SP1Proof, HashableKey, utils, ProverClient, SP1Stdin, SP1ProofWithPublicValues, SP1VerifyingKey};
mod cli;
use cli::ProviderArgs;
use url::Url;
use polccint_lib::{BridgeCommit, BlockCommit, BlockAggregationInput, BlockAggregationCommit, FinalAggregationInput, u32_array_to_hex};
use alloy_rpc_types::BlockNumberOrTag;
use alloy_primitives::{address, Address};
use alloy::hex;
use sp1_cc_host_executor::HostExecutor  as StaticCallHostExecutor;
use sp1_cc_client_executor::{ContractInput};
use std::path::PathBuf;
use polccint_lib::PublicValuesFinalAggregationSolidity;
use alloy_sol_types::SolType;
// import constants from lib
use polccint_lib::constants::{BRIDGE_VK, AGGREGATION_VK};
use serde::{Serialize, Deserialize};

#[derive(Parser, Debug)]
struct Args {
    /// The block number of the block to execute.
    #[clap(long)]
    block_number: u64,
    #[clap(flatten)]
    provider: ProviderArgs,
    #[clap(long)]
    block_range: u64,

    /// Whether or not to generate a proof.
    #[arg(long, default_value_t = false)]
    prove: bool,
}

const ELF_BLOCK_AGGREGATION: &[u8] = include_bytes!("../../../../elf/block-aggregation");
const ELF_BRIDGE: &[u8] = include_bytes!("../../../../elf/bridge");
const ELF_FINAL_AGGREGATION: &[u8] = include_bytes!("../../../../elf/aggregation-final");


/// An input to the aggregation program.
///
/// Consists of a proof and a verification key.
struct AggregationInput {
    pub proof: SP1ProofWithPublicValues,
    pub vk: SP1VerifyingKey,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
struct SP1FinalAggregationProofFixture {
    pub prev_l2_block_hash: String,
    pub new_l2_block_hash: String,
    pub l1_block_hash: String,
    pub new_ler: String,
    pub l1_ger_addr: String,
    pub l2_ger_addr: String,
    pub vkey: String,
    pub public_values: String,
    pub proof: String,
}   

#[tokio::main]
async fn main() -> eyre::Result<()> {
    // Intialize the environment variables.
    dotenv::dotenv().ok();

    // // Fallback to 'info' level if RUST_LOG is not set
    // if std::env::var("RUST_LOG").is_err() {
    //     std::env::set_var("RUST_LOG", "info");
    // }

    // Initialize the logger.
    utils::setup_logger();

    // Parse the command line arguments.
    let args = Args::parse();

    // Load the input from the cache.
    let provider_config = args.provider.into_provider().await?;

    // Generate the proof.
    let client = ProverClient::new();

    // Setup the proving and verifying keys.
    let (aggregation_pk,aggregation_vk) = client.setup(ELF_BLOCK_AGGREGATION);
    let (bridge_pk, bridge_vk) = client.setup(ELF_BRIDGE);
    let (final_aggregation_pk, final_aggregation_vk) = client.setup(ELF_FINAL_AGGREGATION);

    let initial_block_number = args.block_number;
    let block_range = args.block_range; // hardcode for now TODO
    let final_block_number = initial_block_number + block_range;

    // assert constant vk with elf vk 
    assert!(bridge_vk.hash_u32() == BRIDGE_VK);
    assert!(aggregation_vk.hash_u32() == AGGREGATION_VK);

    let proof_aggregation: SP1ProofWithPublicValues = SP1ProofWithPublicValues::load(
        PathBuf::from(env!("CARGO_MANIFEST_DIR"))
            .join(format!("../../proof/chain{}/aggregation_{}_to_{}_proof.bin", provider_config.chain_id, initial_block_number, final_block_number))
    ).expect("failed to load proof");


    let proof_bridge: SP1ProofWithPublicValues = SP1ProofWithPublicValues::load(
        PathBuf::from(env!("CARGO_MANIFEST_DIR"))
            .join(format!("../../proof/chain{}/bridge_block_{}_to_{}_proof.bin", provider_config.chain_id, initial_block_number, final_block_number))
    ).expect("failed to load proof");

    // println!("proof_bridge: {:?}", proof_bridge.public_values.clone().read::<BridgeCommit>());
    // println!("proof_aggregation: {:?}", proof_aggregation.public_values.clone().read::<BlockAggregationCommit>());


    // encode aggregation input and write to stdin
    let mut stdin_final_aggregation = SP1Stdin::new();
    let final_aggregation_input: FinalAggregationInput = FinalAggregationInput {
        block_aggregation_commit: proof_aggregation.public_values.clone().read::<BlockAggregationCommit>(),
        bridge_commit:proof_bridge.public_values.clone().read::<BridgeCommit>()
    };
    stdin_final_aggregation.write(&final_aggregation_input);

    // write proofs
    let SP1Proof::Compressed(proof) = proof_aggregation.proof else {
        panic!()
    };
    stdin_final_aggregation.write_proof(proof, aggregation_vk.vk);
    println!(
      "Finished writing proof",
    );

    let SP1Proof::Compressed(proof) = proof_bridge.proof else {
        panic!()
    };
    stdin_final_aggregation.write_proof(proof, bridge_vk.vk);
    println!(
        "Finished writing proof",
      );
    
    // Only execute the program.
    let (_, execution_report) =
        client.execute(&final_aggregation_pk.elf, stdin_final_aggregation.clone()).run().unwrap();
    println!(
        "Finished executing the block in {} cycles",
        execution_report.total_instruction_count()
    );

    if args.prove {
        println!("Starting proof generation.");
        let proof: SP1ProofWithPublicValues = client.prove(&final_aggregation_pk, stdin_final_aggregation.clone()).plonk().run().expect("Proving should work.");
        println!("Proof generation finished.");

        client.verify(&proof, &final_aggregation_vk).expect("proof verification should succeed");
        // Handle the result of the save operation
        match proof.save(PathBuf::from(env!("CARGO_MANIFEST_DIR")).join(format!("../../proof/chain{}/aggregation_final_{}_to_{}_proof.bin", provider_config.chain_id, initial_block_number, final_block_number))) {
            Ok(_) => println!("Proof saved successfully."),
            Err(e) => eprintln!("Failed to save proof: {}", e),
        }

        let public_values_solidity_encoded = proof.public_values.as_slice();
        println!("public valiues in scirpt {:?}", hex::encode(public_values_solidity_encoded));
        let decoded_values = PublicValuesFinalAggregationSolidity::abi_decode(public_values_solidity_encoded, true).unwrap();

        println!("Decoded public values:");
        println!("prev_l2_block_hash: 0x{}", decoded_values.prev_l2_block_hash);
        println!("new_l2_block_hash: 0x{}", decoded_values.new_l2_block_hash);
        println!("l1_block_hash: 0x{}", decoded_values.l1_block_hash);
        println!("new_ler: 0x{}", decoded_values.new_ler);
        println!("l1_ger_addr: {}", decoded_values.l1_ger_addr);
        println!("l2_ger_addr: {}", decoded_values.l2_ger_addr);

        let fixture = SP1FinalAggregationProofFixture {
            prev_l2_block_hash: format!("{}", decoded_values.prev_l2_block_hash),
            new_l2_block_hash: format!("{}", decoded_values.new_l2_block_hash),
            l1_block_hash: format!("{}", decoded_values.l1_block_hash),
            new_ler: format!("{}", decoded_values.new_ler),
            l1_ger_addr: decoded_values.l1_ger_addr.to_string(),
            l2_ger_addr: decoded_values.l2_ger_addr.to_string(),
            vkey: final_aggregation_vk.bytes32().to_string(),
            public_values: format!("0x{}", hex::encode(public_values_solidity_encoded)),
            proof: format!("0x{}", hex::encode(proof.bytes())),
        };

        let fixture_path = PathBuf::from(env!("CARGO_MANIFEST_DIR")).join("../contracts/src/fixtures");
        std::fs::create_dir_all(&fixture_path).expect("failed to create fixture path");
        std::fs::write(
            fixture_path.join(format!("{:?}-fixture.json", "proof_final_aggregation").to_lowercase()),
            serde_json::to_string_pretty(&fixture).unwrap(),
        )
        .expect("failed to write fixture");

    }
    Ok(())
}


// Generate a `SP1CCProofFixture`, and save it as a json file.
//
// This is useful for verifying the proof of contract call execution on chain.
//fn save_fixture(vkey: String, proof: &SP1ProofWithPublicValues) {
    // let fixture = SP1CCProofFixture {
    //     vkey,
    //     public_values: format!("0x{}", hex::encode(proof.public_values.as_slice())),
    //     proof: format!("0x{}", hex::encode(proof.bytes())),
    // };

    // let fixture_path = PathBuf::from(env!("CARGO_MANIFEST_DIR")).join("../contracts/src/fixtures");
    // std::fs::create_dir_all(&fixture_path).expect("failed to create fixture path");
    // std::fs::write(
    //     fixture_path.join("plonk-fixture.json".to_lowercase()),
    //     serde_json::to_string_pretty(&fixture).unwrap(),
    // )
    // .expect("failed to write fixture");
//}
