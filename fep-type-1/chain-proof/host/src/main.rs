use clap::Parser;
use sp1_sdk::{SP1Proof, HashableKey, utils, ProverClient, SP1Stdin, SP1ProofWithPublicValues};
mod cli;
use cli::ProviderArgs;
use polccint_lib::{ChainProof};
use polccint_lib::bridge::{BridgeCommit};
use polccint_lib::fep_type_1::{BlockAggregationCommit};

use alloy::hex;
use std::path::PathBuf;
use polccint_lib::ChainProofSolidity;
use alloy_sol_types::SolType;
// import constants from lib
use polccint_lib::constants::{BRIDGE_VK, AGGREGATION_VK};
use serde::{Serialize, Deserialize};

#[derive(Parser, Debug)]
struct Args {
    /// The block number of the block to execute.
    #[clap(long)]
    network_id: u64,
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
const ELF_CHAIN_PROOF: &[u8] = include_bytes!("../../../../elf/chain-proof-fep");


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
    let (_,aggregation_vk) = client.setup(ELF_BLOCK_AGGREGATION);
    let (_, bridge_vk) = client.setup(ELF_BRIDGE);
    let (chain_proof_pk, chain_proof_vk) = client.setup(ELF_CHAIN_PROOF);

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
    let mut stdin_chain_proof = SP1Stdin::new();
 
    // First, read the necessary values from proof_aggregation and proof_bridge
    let block_aggregation_commit = proof_aggregation.public_values.clone().read::<BlockAggregationCommit>();
    let bridge_commit = proof_bridge.public_values.clone().read::<BridgeCommit>();

    assert!(bridge_commit.prev_l2_block_hash == block_aggregation_commit.prev_l2_block_hash);
    assert!(bridge_commit.new_l2_block_hash == block_aggregation_commit.new_l2_block_hash);

    // Now, fill the ChainProof struct using the values we just read
    let chain_proof_input: ChainProof = ChainProof {
        prev_l2_block_hash: bridge_commit.prev_l2_block_hash,
        new_l2_block_hash: bridge_commit.new_l2_block_hash,
        l1_block_hash: bridge_commit.l1_block_hash,
        new_ler: bridge_commit.new_ler,
        l1_ger_addr: bridge_commit.l1_ger_addr,
        l2_ger_addr: bridge_commit.l2_ger_addr,
        consensus_hash: bridge_commit.prev_l2_block_hash,  // TODO this is mocked!!
    };

    stdin_chain_proof.write(&chain_proof_input);

    // write proofs
    let SP1Proof::Compressed(proof) = proof_aggregation.proof else {
        panic!()
    };
    stdin_chain_proof.write_proof(proof, aggregation_vk.vk);
    println!(
      "Finished writing proof",
    );

    let SP1Proof::Compressed(proof) = proof_bridge.proof else {
        panic!()
    };
    stdin_chain_proof.write_proof(proof, bridge_vk.vk);
    println!(
        "Finished writing proof",
      );
    
    // Only execute the program.
    let (_, execution_report) =
        client.execute(&chain_proof_pk.elf, stdin_chain_proof.clone()).run().unwrap();
    println!(
        "Finished executing the block in {} cycles",
        execution_report.total_instruction_count()
    );

    if args.prove {
        println!("Starting proof generation.");
        let proof: SP1ProofWithPublicValues = client.prove(&chain_proof_pk, stdin_chain_proof.clone()).compressed().run().expect("Proving should work.");
        println!("Proof generation finished.");

        client.verify(&proof, &chain_proof_vk).expect("proof verification should succeed");
        // Handle the result of the save operation
        let fixture_path = PathBuf::from(env!("CARGO_MANIFEST_DIR")).join(format!("../../../chain-proofs"));
        std::fs::create_dir_all(&fixture_path).expect("failed to create fixture path");

        match proof.save(PathBuf::from(env!("CARGO_MANIFEST_DIR")).join(format!("../../../chain-proofs/proof_chain_{}.bin", args.network_id))) {
            Ok(_) => println!("Proof saved successfully."),
            Err(e) => eprintln!("Failed to save proof: {}", e),
        }
        println!("Proof generation saved.");

    }
    Ok(())
}



