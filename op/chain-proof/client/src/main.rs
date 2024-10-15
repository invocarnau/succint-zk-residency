//! A program that aggregates the proofs of EVM blocks

#![no_main]
sp1_zkvm::entrypoint!(main);

use polccint_lib::ChainProof;
use polccint_lib::bridge::BridgeCommit;
use polccint_lib::op::{ChainProofOPInput, OPConsensusCommit};
use alloy_primitives::B256;

use sha2::{Digest,Sha256};
use bincode;
use polccint_lib::constants::{BRIDGE_VK, OP_CONSENSUS_VK};

pub fn main() {
    // Read the input.

    // could pass just an array of blockhashes
    let input: ChainProofOPInput = sp1_zkvm::io::read::<ChainProofOPInput>();

    // Verify the consensus proof
    // Recreate the consensus commit
    let consensus_commit: OPConsensusCommit = OPConsensusCommit {
        game_factory_address: input.game_factory_address,
        l1_block_hash: input.l1_block_hash,
        prev_l2_block_hash: input.prev_l2_block_hash,
        new_l2_block_hash: input.new_l2_block_hash,
    };
    

    let serialized_public_values_consensus = bincode::serialize(&consensus_commit).unwrap();
    let public_values_digest_consensus = Sha256::digest(serialized_public_values_consensus);
    sp1_zkvm::lib::verify::verify_sp1_proof(
        &OP_CONSENSUS_VK, 
        &public_values_digest_consensus.into()
    );

    // Recreate the bridge commit
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
    sp1_zkvm::lib::verify::verify_sp1_proof(
        &BRIDGE_VK, 
        &public_values_digest_bridge.into()
    );

    // Commit the full input, could be optimized
    sp1_zkvm::io::commit(&ChainProof{
        prev_l2_block_hash: input.prev_l2_block_hash,
        new_l2_block_hash: input.new_l2_block_hash,
        l1_block_hash: input.l1_block_hash,
        new_ler: input.new_ler,
        l1_ger_addr: input.l1_ger_addr,
        l2_ger_addr: input.l2_ger_addr,
        consensus_hash: B256::new(Sha256::digest(input.game_factory_address).into()),
    });
}
