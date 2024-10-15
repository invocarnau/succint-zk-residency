//! A program which generates consensus proofs for Polygon PoS chain using 
//! the milestone message asserting that majority of validators by stake
//! voted on a specific block.

#![no_main]
sp1_zkvm::entrypoint!(main); 

use alloy_primitives::{Address, FixedBytes};
use alloy_sol_types::SolType;
use reth_primitives::Header;

use milestone::{MilestoneProofInputs, MilestoneProver, PublicValuesStruct};

pub mod helper;
pub mod milestone;
pub mod types;

fn main() {
    // Read inputs from the zkVM's stdin.
    let tx_data = sp1_zkvm::io::read::<String>();
    let tx_hash = sp1_zkvm::io::read::<FixedBytes<32>>();
    let precommits = sp1_zkvm::io::read::<Vec<Vec<u8>>>();
    let sigs = sp1_zkvm::io::read::<Vec<String>>();
    let signers = sp1_zkvm::io::read::<Vec<Address>>();
    let bor_header = sp1_zkvm::io::read::<Header>();
    let state_sketch_bytes = sp1_zkvm::io::read::<Vec<u8>>();
    let l1_block_hash = sp1_zkvm::io::read::<FixedBytes<32>>();

    let inputs = MilestoneProofInputs {
        tx_data,
        tx_hash,
        precommits,
        sigs,
        signers,
        bor_header: bor_header.clone(),
        prev_bor_header: bor_header,
        state_sketch_bytes,
        l1_block_hash,
    };
    let prover = MilestoneProver::init(inputs);
    let outputs = prover.prove();

    // Encode the public values
    let bytes = PublicValuesStruct::abi_encode_packed(&PublicValuesStruct {
        bor_block_hash: outputs.bor_block_hash,
        l1_block_hash: outputs.l1_block_hash,
    });

    // Commit the values as bytes to be exposed to the verifier
    sp1_zkvm::io::commit_slice(&bytes);
}