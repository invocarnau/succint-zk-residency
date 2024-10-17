use clap::Parser;
use polccint_lib::constants::{FEP, FEP_CHAIN_VK, OP, OP_CHAIN_VK, POS, POS_CHAIN_VK};
use polccint_lib::{AggLayerProofInput, ChainProof, ChainProofSolidity};
use sp1_sdk::{
    utils, HashableKey, ProverClient, SP1Proof, SP1ProofWithPublicValues, SP1Stdin, SP1VerifyingKey,
};

use alloy::hex;
use alloy_sol_types::SolType;
use polccint_lib::AggLayerProofSolidity;
use serde::{Deserialize, Serialize};
use std::{collections::HashMap, path::PathBuf};

const FEP_CHAIN_ELF: &[u8] = include_bytes!("../../../elf/chain-proof-fep");
const OP_CHAIN_ELF: &[u8] = include_bytes!("../../../elf/chain-proof-op"); // TODO: add correct elf
const POS_CHAIN_ELF: &[u8] = include_bytes!("../../../elf/chain-proof-fep"); // TODO: add correct elf

#[derive(Parser, Debug)]
struct Args {
    /// The block number of the block to execute.
    #[clap(long)]
    network_ids: String,

    /// Whether or not to generate a proof.
    #[arg(long, default_value_t = false)]
    prove: bool,
}

const ELF_AGGLAYER_PROOF: &[u8] = include_bytes!("../../../elf/agglayer-proof");

#[derive(Debug, Clone, Serialize, Deserialize)]
struct SP1FinalAggregationProofFixture {
    pub chain_proofs: Vec<ChainProofSolidity>,
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
    let vks = HashMap::from([(FEP, FEP_CHAIN_VK), (OP, OP_CHAIN_VK), (POS, POS_CHAIN_VK)]);
    let elfs = HashMap::from([
        (FEP, FEP_CHAIN_ELF),
        (OP, OP_CHAIN_ELF),
        (POS, POS_CHAIN_ELF),
    ]);
    // Intialize the environment variables.
    dotenv::dotenv().ok();
    utils::setup_logger();
    let args = Args::parse();
    let net_ids_str: Vec<&str> = args.network_ids.split(',').collect();
    println!("{:?}", net_ids_str);
    let mut net_ids = Vec::new();
    assert!(net_ids_str.len() > 0);
    for net_id in net_ids_str {
        // Assuming net_id is of type String
        let net_id_uint: u32 = net_id.parse().expect("Failed to parse net_id as u32");
        assert!(vks.contains_key(&net_id_uint));
        net_ids.push(net_id_uint);
    }

    // Generate the proof.
    let client = ProverClient::new();

    // Setup the proving and verifying keys.
    let (agglayer_proof_pk, agglayer_proof_vk) = client.setup(ELF_AGGLAYER_PROOF);
    let mut chain_proof_inputs: Vec<AggregationInput> = Vec::new();
    let mut input_vks: Vec<[u32; 8]> = Vec::new();

    for network_id in net_ids {
        println!("adding inputs for network id: {}", network_id);
        let elf_chain = elfs[&network_id];
        let vk_chain = vks[&network_id];
        let (_, actual_chain_vk) = client.setup(elf_chain);
        assert_eq!(vk_chain, actual_chain_vk.hash_u32());
        let proof: SP1ProofWithPublicValues = SP1ProofWithPublicValues::load(
            PathBuf::from(env!("CARGO_MANIFEST_DIR"))
                .join(format!("../../chain-proofs/proof_chain_{}.bin", network_id)),
        )
        .expect("failed to load proof");

        chain_proof_inputs.push(AggregationInput {
            proof: proof,
            vk: actual_chain_vk,
        });
        input_vks.push(vk_chain);
    }

    // encode aggregation input and write to stdin
    let mut stdin = SP1Stdin::new();
    let aggregation_input = AggLayerProofInput {
        chain_proofs: chain_proof_inputs
            .iter()
            .map(|input| input.proof.public_values.clone().read::<ChainProof>())
            .collect::<Vec<_>>(),
        vks: input_vks,
    };

    stdin.write(&aggregation_input);

    // write proofs
    for input in chain_proof_inputs {
        let SP1Proof::Compressed(proof) = input.proof.proof else {
            panic!()
        };
        stdin.write_proof(proof, input.vk.vk);
    }

    // Only execute the program.
    let (public_values, execution_report) = client
        .execute(&agglayer_proof_pk.elf, stdin.clone())
        .run()
        .unwrap();
    println!(
        "Finished executing the block in {} cycles",
        execution_report.total_instruction_count()
    );

    if args.prove {
        println!("Starting proof generation.");
        let proof: SP1ProofWithPublicValues = client
            .prove(&agglayer_proof_pk, stdin.clone())
            .plonk()
            .run()
            .expect("Proving should work.");
        println!("Proof generation finished.");

        // client
        //     .verify(&proof, &agglayer_proof_vk)
        //     .expect("proof verification should succeed");
        // Handle the result of the save operation

        let fixture_path: PathBuf =
            PathBuf::from(env!("CARGO_MANIFEST_DIR")).join("../agglayer_proofs");
        std::fs::create_dir_all(&fixture_path).expect("failed to create fixture path");

        let network_ids_string = args.network_ids.replace(",", "_");
        match proof.save(PathBuf::from(env!("CARGO_MANIFEST_DIR")).join(format!(
            "../agglayer_proofs/aggregation_for_{}.bin",
            network_ids_string
        ))) {
            Ok(_) => println!("Proof saved successfully."),
            Err(e) => eprintln!("Failed to save proof: {}", e),
        }

        let public_values_solidity_encoded = proof.public_values.as_slice();
        let decoded_values: AggLayerProofSolidity =
            AggLayerProofSolidity::abi_decode(public_values_solidity_encoded, true).unwrap();

        let fixture = SP1FinalAggregationProofFixture {
            chain_proofs: decoded_values.chain_proofs,
            vkey: agglayer_proof_vk.bytes32().to_string(),
            public_values: format!("0x{}", hex::encode(public_values_solidity_encoded)),
            proof: format!("0x{}", hex::encode(proof.bytes())),
        };

        std::fs::write(
            fixture_path.join(
                format!(
                    "aggregation_for_{}.json",
                    network_ids_string
                )
                .to_lowercase(),
            ),
            serde_json::to_string_pretty(&fixture).unwrap(),
        )
        .expect("failed to write fixture");
    }
    Ok(())
}
