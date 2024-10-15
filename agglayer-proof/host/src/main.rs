use clap::Parser;
use sp1_sdk::{SP1Proof, HashableKey, utils, ProverClient, SP1Stdin, SP1ProofWithPublicValues, SP1VerifyingKey};
mod cli;
use cli::ProviderArgs;
use polccint_lib::{ChainProof, AggLayerProof};
use polccint_lib::bridge::{BridgeCommit};
use polccint_lib::fep_type_1::{BlockAggregationCommit};

use alloy::hex;
use std::path::PathBuf;
use polccint_lib::AggLayerProofSolidity;
use alloy_sol_types::SolType;
// import constants from lib
use polccint_lib::constants::{CHAIN_VK};
use serde::{Serialize, Deserialize};

#[derive(Parser, Debug)]
struct Args {
    /// The block number of the block to execute.
    #[clap(long)]
    proof_number: u64,
    #[clap(flatten)]
    provider: ProviderArgs,
    #[clap(long)]
    proof_range: u64,

    /// Whether or not to generate a proof.
    #[arg(long, default_value_t = false)]
    prove: bool,
}

const ELF_CHAIN_PROOF: &[u8] = include_bytes!("../../../elf/aggregation-final");
const ELF_AGGLAYER_PROOF: &[u8] = include_bytes!("../../../elf/agglayer-proof");

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
struct AggregationInput {
    pub proof: SP1ProofWithPublicValues,
    pub vk: SP1VerifyingKey,
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
    let (_,chain_vk) = client.setup(ELF_CHAIN_PROOF);
    let (agglayer_proof_pk, agglayer_proof_vk) = client.setup(ELF_AGGLAYER_PROOF);

    let initial_proof_number = args.proof_number;
    let proof_range = args.proof_range; // hardcode for now TODO
    let final_proof_number = initial_proof_number + proof_range;

    // assert constant vk with elf vk 
    assert!(chain_vk.hash_u32() == CHAIN_VK);

    let mut inputs: Vec<AggregationInput> = Vec::new();

    for proof_number in initial_proof_number..final_proof_number + 1 {
        let proof: SP1ProofWithPublicValues = SP1ProofWithPublicValues::load(
            PathBuf::from(env!("CARGO_MANIFEST_DIR"))
                .join(format!(
                    "../../proof/proof_{}.bin",
                    proof_number
                ))
        ).expect("failed to load proof");
        
        inputs.push(
            AggregationInput {
                proof: proof,
                vk: chain_vk.clone(),
            }
        );
    }

     // encode aggregation input and write to stdin
     let mut stdin = SP1Stdin::new();
     let aggregation_input = AggLayerProof{
        chain_proofs: inputs
         .iter()
         .map(|input| input.proof.public_values.clone().read::<ChainProof>())
         .collect::<Vec<_>>(),
     };
     stdin.write(&aggregation_input);
 
     // write proofs
     for input in inputs {
         let SP1Proof::Compressed(proof) = input.proof.proof else {
             panic!()
         };
         stdin.write_proof(proof, input.vk.vk);
     }
     
     // Only execute the program.
     let (mut public_values, execution_report) =
         client.execute(&agglayer_proof_pk.elf, stdin.clone()).run().unwrap();
     println!(
         "Finished executing the block in {} cycles",
         execution_report.total_instruction_count()
     );

    if args.prove {
        println!("Starting proof generation.");
        let proof: SP1ProofWithPublicValues = client.prove(&agglayer_proof_pk, stdin.clone()).plonk().run().expect("Proving should work.");
        println!("Proof generation finished.");

        client.verify(&proof, &agglayer_proof_vk).expect("proof verification should succeed");
        // Handle the result of the save operation
        match proof.save(PathBuf::from(env!("CARGO_MANIFEST_DIR")).join(format!("../../agglayer_proofs/aggregation_from_{}_to_{}_proof.bin", initial_proof_number, final_proof_number))) {
            Ok(_) => println!("Proof saved successfully."),
            Err(e) => eprintln!("Failed to save proof: {}", e),
        }

        let public_values_solidity_encoded = proof.public_values.as_slice();
        println!("public valiues in scirpt {:?}", hex::encode(public_values_solidity_encoded));
        let decoded_values = AggLayerProofSolidity::abi_decode(public_values_solidity_encoded, true).unwrap();

        // println!("Decoded public values:");
        // println!("prev_l2_block_hash: 0x{}", decoded_values.prev_l2_block_hash);
        // println!("new_l2_block_hash: 0x{}", decoded_values.new_l2_block_hash);
        // println!("l1_block_hash: 0x{}", decoded_values.l1_block_hash);
        // println!("new_ler: 0x{}", decoded_values.new_ler);
        // println!("l1_ger_addr: {}", decoded_values.l1_ger_addr);
        // println!("l2_ger_addr: {}", decoded_values.l2_ger_addr);

        // let fixture = SP1FinalAggregationProofFixture {
        //     prev_l2_block_hash: format!("{}", decoded_values.prev_l2_block_hash),
        //     new_l2_block_hash: format!("{}", decoded_values.new_l2_block_hash),
        //     l1_block_hash: format!("{}", decoded_values.l1_block_hash),
        //     new_ler: format!("{}", decoded_values.new_ler),
        //     l1_ger_addr: decoded_values.l1_ger_addr.to_string(),
        //     l2_ger_addr: decoded_values.l2_ger_addr.to_string(),
        //     vkey: final_aggregation_vk.bytes32().to_string(),
        //     public_values: format!("0x{}", hex::encode(public_values_solidity_encoded)),
        //     proof: format!("0x{}", hex::encode(proof.bytes())),
        // };

        // let fixture_path = PathBuf::from(env!("CARGO_MANIFEST_DIR")).join("../fixtures");
        // std::fs::create_dir_all(&fixture_path).expect("failed to create fixture path");
        // std::fs::write(
        //     fixture_path.join(format!("{:?}-fixture.json", "proof_final_aggregation").to_lowercase()),
        //     serde_json::to_string_pretty(&fixture).unwrap(),
        // )
        // .expect("failed to write fixture");

    }
    Ok(())
}



