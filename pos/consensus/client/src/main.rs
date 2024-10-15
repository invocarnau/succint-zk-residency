//! A program which generates consensus proofs for Polygon PoS chain using
//! the milestone message asserting that majority of validators by stake
//! voted on a specific block.

#![no_main]
sp1_zkvm::entrypoint!(main);

use alloy_sol_types::SolType;
use polccint_lib::pos_consensus::{PoSConsensusInput, PublicValuesStruct};
use pos_consensus_proof_client::milestone::prove;

fn main() {
    // Read inputs from the zkVM's stdin.
    let input = sp1_zkvm::io::read::<PoSConsensusInput>();

    // Calculate the milestone proof
    let commit = prove(input);

    // Encode the public values from the commit
    let bytes = PublicValuesStruct::abi_encode_packed(&PublicValuesStruct {
        prev_bor_block_hash: commit.prev_bor_hash,
        new_bor_block_hash: commit.new_bor_hash,
        l1_block_hash: commit.l1_block_hash,
    });

    // Commit the values as bytes to be exposed to the verifier
    sp1_zkvm::io::commit_slice(&bytes);
}
