use clap::Parser;
use sp1_sdk::{SP1Proof, HashableKey, utils, ProverClient, SP1Stdin, SP1ProofWithPublicValues};
use cli::ProviderArgs;
use polccint_lib::{ChainProof};
use polccint_lib::bridge::{BridgeCommit};
use polccint_lib::op::ChainProofOPInput;

use alloy::hex;
use std::path::PathBuf;
use polccint_lib::ChainProofSolidity;
use alloy_sol_types::SolType;
// import constants from lib
use polccint_lib::constants::{BRIDGE_VK, OP_CONSENSUS_VK};
use serde::{Serialize, Deserialize};

#[derive(Parser, Debug)]
struct Args {
    /// The block number of the block to execute.
    #[clap(long)]
    block_number: u64,
    #[clap(flatten)]
    provider: ProviderArgs,

    /// Whether or not to generate a proof.
    #[arg(long, default_value_t = false)]
    prove: bool,
    /// Whether or not to generate a proof.
    #[arg(long, default_value_t = false)]
    prove: bool,

    #[clap(long)]
    chain_id_l1: Option<u64>,

    #[clap(long)]
    chain_id_l2: Option<u64>,

    #[clap(long)]
    block_number_l1: u64,

    #[clap(long)]
    game_index: u64,

    #[clap(long)]
    game_factory_address: String,

    #[clap(long)]
    prev_block_number_l2: u64,

    #[clap(long)]
    new_block_number_l2: u64,

    #[clap(long)]
    claim_block_hash: String,

    #[clap(long)]
    claim_state_root: String,

    #[clap(long)]
    claim_message_passer_storage_root: String,
}

const ELF_CONSENSUS: &[u8] = include_bytes!("../../../../elf/op-consensus");
const ELF_BRIDGE: &[u8] = include_bytes!("../../../../elf/bridge");
const ELF_CHAIN: &[u8] = include_bytes!("../../../../elf/op-chain");


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
    dotenv::dotenv().ok();
    utils::setup_logger();
    let args = Args::parse();

    let provider_config = args.provider.into_provider().await?;
    let client = ProverClient::new();

    // Setup the proving and verifying keys.
    let (_,consensus_vk) = client.setup(ELF_CONSENSUS);
    let (_, bridge_vk) = client.setup(ELF_BRIDGE);
    let (chain_pk, chain_vk) = client.setup(ELF_CHAIN);

    let initial_block_number = args.block_number;

    // assert constant vk with elf vk 
    assert!(bridge_vk.hash_u32() == BRIDGE_VK);
    assert!(consensus_vk.hash_u32() == OP_CONSENSUS_VK);

    let proof_aggregation: SP1ProofWithPublicValues = SP1ProofWithPublicValues::load(
        PathBuf::from(env!("CARGO_MANIFEST_DIR")).join(format!(
            "../../proof/chain{}/game_{}.bin",
            provider_config.chain_id,
            args.game_index,
        ))
    ).expect("failed to load proof");


    let proof_bridge: SP1ProofWithPublicValues = SP1ProofWithPublicValues::load(
        PathBuf::from(env!("CARGO_MANIFEST_DIR")).join(format!(
            "../../proof/chain{}/bridge_block_{}_to_{}_proof.bin",
            provider_config.chain_id,
            initial_block_number,
            final_block_number,
        ))
    ).expect("failed to load proof");

    // println!("proof_bridge: {:?}", proof_bridge.public_values.clone().read::<BridgeCommit>());
    // println!("proof_aggregation: {:?}", proof_aggregation.public_values.clone().read::<BlockAggregationCommit>());

    // encode aggregation input and write to stdin
    let mut stdin_final_aggregation = SP1Stdin::new();
 
    // First, read the necessary values from proof_aggregation and proof_bridge
    let block_aggregation_commit = proof_aggregation.public_values.clone().read::<BlockAggregationCommit>();
    let bridge_commit = proof_bridge.public_values.clone().read::<BridgeCommit>();

    assert!(bridge_commit.prev_l2_block_hash == block_aggregation_commit.prev_l2_block_hash);
    assert!(bridge_commit.new_l2_block_hash == block_aggregation_commit.new_l2_block_hash);

    // Now, fill the ChainProof struct using the values we just read
    let final_aggregation_input: ChainProof = ChainProof {
        prev_l2_block_hash: bridge_commit.prev_l2_block_hash,
        new_l2_block_hash: bridge_commit.new_l2_block_hash,
        l1_block_hash: bridge_commit.l1_block_hash,
        new_ler: bridge_commit.new_ler,
        l1_ger_addr: bridge_commit.l1_ger_addr,
        l2_ger_addr: bridge_commit.l2_ger_addr,
        consensus_hash: bridge_commit.prev_l2_block_hash,  // TODO this is mocked!!
    };

    stdin_final_aggregation.write(&final_aggregation_input);

    // write proofs
    let SP1Proof::Compressed(proof) = proof_aggregation.proof else {
        panic!()
    };
    stdin_final_aggregation.write_proof(proof, consensus_vk.vk);
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
        client.execute(&chain_pk.elf, stdin_final_aggregation.clone()).run().unwrap();
    println!(
        "Finished executing the block in {} cycles",
        execution_report.total_instruction_count()
    );

    if args.prove {
        println!("Starting proof generation.");
        let proof: SP1ProofWithPublicValues = client.prove(&chain_pk, stdin_final_aggregation.clone()).plonk().run().expect("Proving should work.");
        println!("Proof generation finished.");

        client.verify(&proof, &chain_vk).expect("proof verification should succeed");
        // Handle the result of the save operation
        match proof.save(PathBuf::from(env!("CARGO_MANIFEST_DIR")).join(format!("../../proof/chain{}/aggregation_final_{}_to_{}_proof.bin", provider_config.chain_id, initial_block_number, final_block_number))) {
            Ok(_) => println!("Proof saved successfully."),
            Err(e) => eprintln!("Failed to save proof: {}", e),
        }

        let public_values_solidity_encoded = proof.public_values.as_slice();
        println!("public valiues in scirpt {:?}", hex::encode(public_values_solidity_encoded));
        let decoded_values = ChainProofSolidity::abi_decode(public_values_solidity_encoded, true).unwrap();

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
            vkey: chain_vk.bytes32().to_string(),
            public_values: format!("0x{}", hex::encode(public_values_solidity_encoded)),
            proof: format!("0x{}", hex::encode(proof.bytes())),
        };

        let fixture_path = PathBuf::from(env!("CARGO_MANIFEST_DIR")).join("../fixtures");
        std::fs::create_dir_all(&fixture_path).expect("failed to create fixture path");
        std::fs::write(
            fixture_path.join(format!("{:?}-fixture.json", "proof_final_aggregation").to_lowercase()),
            serde_json::to_string_pretty(&fixture).unwrap(),
        )
        .expect("failed to write fixture");

    }
    Ok(())
}



