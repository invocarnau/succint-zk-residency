//! A program that aggregates the consensus and bridge proof for pos

#![no_main]
sp1_zkvm::entrypoint!(main);

use alloy_primitives::B256;
use polccint_lib::bridge::BridgeCommit;
use polccint_lib::constants::{BRIDGE_VK, POS_CONSENSUS_VK};
use polccint_lib::pos_consensus::{ChainProofPoSInput, PoSConsensusCommit};
use polccint_lib::ChainProof;
use sha2::{Digest, Sha256};

fn main() {
    let input = sp1_zkvm::io::read::<ChainProofPoSInput>();

    // Construct the consensus commit
    let consensus_commit = PoSConsensusCommit {
        prev_bor_hash: input.prev_l2_block_hash,
        new_bor_hash: input.new_l2_block_hash,
        l1_block_hash: input.l1_block_hash,
        stake_manager_address: input.stake_manager_address,
    };

    // Verify the consensus proof
    let serialized_public_values_consensus = bincode::serialize(&consensus_commit).unwrap();
    let public_values_digest_consensus = Sha256::digest(serialized_public_values_consensus);
    sp1_zkvm::lib::verify::verify_sp1_proof(
        &POS_CONSENSUS_VK,
        &public_values_digest_consensus.into(),
    );

    // Construct the bridge commit
    let bridge_commit: BridgeCommit = BridgeCommit {
        prev_l2_block_hash: input.prev_l2_block_hash,
        new_l2_block_hash: input.new_l2_block_hash,
        l1_block_hash: input.l1_block_hash,
        new_ler: input.new_ler,
        l1_ger_addr: input.l1_ger_addr,
        l2_ger_addr: input.l2_ger_addr,
    };

    let serialized_public_values_bridge = bincode::serialize(&bridge_commit).unwrap();
    let public_values_digest_bridge = Sha256::digest(serialized_public_values_bridge);
    sp1_zkvm::lib::verify::verify_sp1_proof(&BRIDGE_VK, &public_values_digest_bridge.into());

    // Commit to the final aggregated proof
    sp1_zkvm::io::commit(&ChainProof {
        prev_l2_block_hash: input.prev_l2_block_hash,
        new_l2_block_hash: input.new_l2_block_hash,
        l1_block_hash: input.l1_block_hash,
        new_ler: input.new_ler,
        l1_ger_addr: input.l1_ger_addr,
        l2_ger_addr: input.l2_ger_addr,
        consensus_hash: B256::new(Sha256::digest(input.stake_manager_address).into()),
    });
}
