//! A program that aggregates the proofs of EVM blocks

#![no_main]
sp1_zkvm::entrypoint!(main);

use polccint_lib::FinalAggregationInput;
use sha2::{Digest,Sha256};
use bincode;
use polccint_lib::PublicValuesFinalAggregationSolidity;
use alloy_sol_types::SolType;
use alloy_primitives::FixedBytes;

pub fn main() {
    // Read the input.

    // could pass just an array of blockhashes
    let input: FinalAggregationInput = sp1_zkvm::io::read::<FinalAggregationInput>();

    // Confirm that that both inputs are consistent
    assert!(input.bridge_commit.prev_l2_block_hash == input.block_aggregation_commit.prev_l2_block_hash );
    assert!(input.bridge_commit.new_l2_block_hash == input.block_aggregation_commit.new_l2_block_hash );

    // Verify the aggregated block proof
    let serialized_public_values_aggregation = bincode::serialize(&input.block_aggregation_commit).unwrap();
    let public_values_digest_aggregation = Sha256::digest(serialized_public_values_aggregation);
    sp1_zkvm::lib::verify::verify_sp1_proof(
        &input.block_vkey_aggregation, 
        &public_values_digest_aggregation.into()
    );

    // Verify the bridge proof
    let serialized_public_values_bridge = bincode::serialize(&input.bridge_commit).unwrap();
    let public_values_digest_bridge = Sha256::digest(serialized_public_values_bridge);
    sp1_zkvm::lib::verify::verify_sp1_proof(
        &input.block_vkey_bridge, 
        &public_values_digest_bridge.into()
    );

    let public_values_solidity: PublicValuesFinalAggregationSolidity = PublicValuesFinalAggregationSolidity {
        block_vkey_aggregation: FixedBytes::<32>::from_slice(input.block_vkey_aggregation.iter().flat_map(|&x| x.to_be_bytes()).collect::<Vec<u8>>().as_slice()),
        block_vkey: FixedBytes::<32>::from_slice(input.block_aggregation_commit.block_vkey.iter().flat_map(|&x| x.to_be_bytes()).collect::<Vec<u8>>().as_slice()),
        block_vkey_bridge: FixedBytes::<32>::from_slice(input.block_vkey_bridge.iter().flat_map(|&x| x.to_be_bytes()).collect::<Vec<u8>>().as_slice()),
        prev_l2_block_hash: input.bridge_commit.prev_l2_block_hash,
        new_l2_block_hash: input.bridge_commit.new_l2_block_hash,
        l1_block_hash: input.bridge_commit.l1_block_hash,
        new_ler: input.bridge_commit.new_ler,
        l1_ger_addr: input.bridge_commit.l1_ger_addr,
        l2_ger_addr: input.bridge_commit.l2_ger_addr, 
    };

    let public_values_solidity_encoded = PublicValuesFinalAggregationSolidity::abi_encode(&public_values_solidity);

    // Commit the full input, could be optimized
    sp1_zkvm::io::commit_slice(&public_values_solidity_encoded);
}
