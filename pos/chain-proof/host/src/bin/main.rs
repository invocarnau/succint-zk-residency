use std::path::PathBuf;

use clap::Parser;
use polccint_lib::{
    bridge::BridgeCommit,
    constants::{BRIDGE_VK, POS_CONSENSUS_VK},
    pos::{ChainProofPoSInput, PoSConsensusCommit},
};
use sp1_sdk::{HashableKey, ProverClient, SP1Proof, SP1ProofWithPublicValues, SP1Stdin};

pub const POS_CONSENSUS_PROOF_ELF: &[u8] = include_bytes!("../../../../../elf/pos-consensus-proof");
pub const POS_BRIDGE_PROOF_ELF: &[u8] = include_bytes!("../../../../../elf/bridge");
pub const POS_CHAIN_PROOF_ELF: &[u8] = include_bytes!("../../../../../elf/pos-chain");

/// The arguments for the command.
#[derive(Parser, Debug)]
#[clap(author, version, about, long_about = None)]
pub struct Args {
    #[clap(long)]
    prev_l2_block_number: u64,

    #[clap(long)]
    new_l2_block_number: u64,

    #[arg(long, default_value_t = false)]
    prove: bool,
}

#[tokio::main]
async fn main() -> eyre::Result<()> {
    dotenv::dotenv().ok();

    // Setup the logger.
    sp1_sdk::utils::setup_logger();

    // Parse the command line arguments.
    let args = Args::parse();

    // Init the prover client and load vkeys
    println!("Initializing SP1 ProverClient...");
    let client = ProverClient::new();
    let (_, consensus_vk) = client.setup(POS_CONSENSUS_PROOF_ELF);
    let (_, bridge_vk) = client.setup(POS_CONSENSUS_PROOF_ELF);
    let (chain_proof_pk, chain_proof_vk) = client.setup(POS_CONSENSUS_PROOF_ELF);

    println!("consensus vk {:?}", consensus_vk.hash_u32());
    println!("bridge vk {:?}", bridge_vk.hash_u32());
    assert!(consensus_vk.hash_u32() == POS_CONSENSUS_VK);
    assert!(bridge_vk.hash_u32() == BRIDGE_VK);

    // Load the proofs
    let l2_chain_id = std::env::var("L2_CHAIN_ID").expect("L2_CHAIN_ID not set");
    println!(
        "Loading consensus proof for chain: {} from block: {} to block: {}",
        l2_chain_id, args.prev_l2_block_number, args.new_l2_block_number,
    );
    let proof_consensus: SP1ProofWithPublicValues =
        SP1ProofWithPublicValues::load(PathBuf::from(env!("CARGO_MANIFEST_DIR")).join(format!(
            "../../proof/chain{}/consensus_block_{}_to_{}.bin",
            l2_chain_id, args.prev_l2_block_number, args.new_l2_block_number,
        )))
        .expect("failed to load consensus proof");

    println!(
        "loading bridge proof of the chain: {} from block: {} to block: {}",
        l2_chain_id, args.prev_l2_block_number, args.new_l2_block_number,
    );
    let proof_bridge: SP1ProofWithPublicValues =
        SP1ProofWithPublicValues::load(PathBuf::from(env!("CARGO_MANIFEST_DIR")).join(format!(
            "../../../bridge/proof/chain{}/bridge_block_{}_to_{}_proof.bin",
            l2_chain_id, args.prev_l2_block_number, args.new_l2_block_number,
        )))
        .expect("failed to load proof");

    // Prepare inputs for chain proof
    let mut stdin_chain_proof = SP1Stdin::new();
    let consensus_commit = proof_consensus
        .public_values
        .clone()
        .read::<PoSConsensusCommit>();
    let bridge_commit = proof_bridge.public_values.clone().read::<BridgeCommit>();

    assert!(bridge_commit.prev_l2_block_hash == consensus_commit.prev_bor_hash);
    assert!(bridge_commit.new_l2_block_hash == consensus_commit.new_bor_hash);
    assert!(bridge_commit.l1_block_hash == consensus_commit.l1_block_hash);

    let chain_proof_input = ChainProofPoSInput {
        prev_l2_block_hash: bridge_commit.prev_l2_block_hash,
        new_l2_block_hash: bridge_commit.new_l2_block_hash,
        l1_block_hash: bridge_commit.l1_block_hash,
        new_ler: bridge_commit.new_ler,
        l1_ger_addr: bridge_commit.l1_ger_addr,
        l2_ger_addr: bridge_commit.l2_ger_addr,
        stake_manager_address: consensus_commit.stake_manager_address,
    };
    stdin_chain_proof.write(&chain_proof_input);

    // write proofs
    let SP1Proof::Compressed(proof) = proof_consensus.proof else {
        panic!()
    };
    stdin_chain_proof.write_proof(proof, consensus_vk.vk);
    println!("Finished writing consensus proof",);

    let SP1Proof::Compressed(proof) = proof_bridge.proof else {
        panic!()
    };
    stdin_chain_proof.write_proof(proof, bridge_vk.vk);
    println!("Finished writing bridge proof",);

    // Only execute the program.
    let (_, execution_report) = client
        .execute(&chain_proof_pk.elf, stdin_chain_proof.clone())
        .run()
        .unwrap();
    println!(
        "Finished executing the block in {} cycles",
        execution_report.total_instruction_count()
    );

    if args.prove {
        println!("Starting proof generation.");
        let proof: SP1ProofWithPublicValues = client
            .prove(&chain_proof_pk, stdin_chain_proof.clone())
            .compressed()
            .run()
            .expect("Proving should work.");
        println!("Proof generation finished.");

        client
            .verify(&proof, &chain_proof_vk)
            .expect("proof verification should succeed");
        // Handle the result of the save operation
        let fixture_path = PathBuf::from(env!("CARGO_MANIFEST_DIR")).join("../../../chain-proofs");
        std::fs::create_dir_all(&fixture_path).expect("failed to create fixture path");

        match proof.save(PathBuf::from(env!("CARGO_MANIFEST_DIR")).join(format!(
            "../../../chain-proofs/proof_chain_{}.bin",
            l2_chain_id
        ))) {
            Ok(_) => println!("Proof saved successfully."),
            Err(e) => eprintln!("Failed to save proof: {}", e),
        }
        println!("Proof saved.");
    }

    Ok(())
}
