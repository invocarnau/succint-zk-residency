//! A script to verify the proof locally using sp1 and also send it for on-chain verification.
//!
//! You can run this script using the following command:
//! ```shell
//! RUST_LOG=info cargo run --package operator --bin verify --release
//! ```

use alloy_sol_types::{SolCall, SolType};
use pos_consensus_proof_client::milestone::{ConsensusProofVerifier, PublicValuesStruct};
use pos_consensus_proof_host::{contract::ContractClient, ConsensusProver};
use sp1_sdk::SP1ProofWithPublicValues;

#[tokio::main]
async fn main() -> eyre::Result<()> {
    dotenv::dotenv().ok();

    // Load proof from `proof.bin`
    let proof = SP1ProofWithPublicValues::load("proof.bin").unwrap();
    println!("Proof loaded from proof.bin");

    let prover = ConsensusProver::new();
    prover.verify_consensus_proof(&proof);
    println!("Proof verified locally, sending for on-chain verification!");

    send_proof_onchain(proof).await?;
    println!("Done!");

    Ok(())
}

pub async fn send_proof_onchain(proof: SP1ProofWithPublicValues) -> eyre::Result<()> {
    // Setup the default contract client to interact with on-chain verifier
    let contract_client = ContractClient::default();

    // Decode the public values from the proof
    let vals = PublicValuesStruct::abi_decode(&proof.public_values.to_vec(), true).unwrap();

    // Construct the on-chain call and relay the proof to the contract.
    let call_data = ConsensusProofVerifier::verifyConsensusProofCall {
        _proofBytes: proof.bytes().into(),
        new_bor_block_hash: vals.new_bor_block_hash,
        l1_block_hash: vals.l1_block_hash,
    }
    .abi_encode();
    let result = contract_client.send(call_data).await;

    if result.is_err() {
        println!("error sending proof: err={:?}", result.err().unwrap());
    } else {
        println!("Successfully verified proof on-chain!");
    }

    Ok(())
}
