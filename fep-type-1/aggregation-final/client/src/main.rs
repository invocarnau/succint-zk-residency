//! A program that aggregates the proofs of EVM blocks

#![no_main]
sp1_zkvm::entrypoint!(main);

use polccint_lib::{ChainProof};
use polccint_lib::bridge::{BridgeCommit};
use polccint_lib::fep_type_1::{BlockAggregationCommit};

use sha2::{Digest,Sha256};
use bincode;
use polccint_lib::ChainProofSolidity;
use alloy_sol_types::SolType;
use polccint_lib::constants::{BRIDGE_VK, AGGREGATION_VK};

pub fn main() {
    // Read the input.

    // could pass just an array of blockhashes
    let input: ChainProof = sp1_zkvm::io::read::<ChainProof>();

    // Verify the aggregated block proof

    // Recreate the block aggregation commit
    let block_aggregation_commit: BlockAggregationCommit = BlockAggregationCommit {
        prev_l2_block_hash: input.prev_l2_block_hash,
        new_l2_block_hash: input.new_l2_block_hash,
    };
    

    let serialized_public_values_aggregation = bincode::serialize(&block_aggregation_commit).unwrap();
    let public_values_digest_aggregation = Sha256::digest(serialized_public_values_aggregation);
    sp1_zkvm::lib::verify::verify_sp1_proof(
        &AGGREGATION_VK, 
        &public_values_digest_aggregation.into()
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

    let public_values_solidity: ChainProofSolidity = ChainProofSolidity {
        prev_l2_block_hash: input.prev_l2_block_hash,
        new_l2_block_hash: input.new_l2_block_hash,
        l1_block_hash: input.l1_block_hash,
        new_ler: input.new_ler,
        l1_ger_addr: input.l1_ger_addr,
        l2_ger_addr: input.l2_ger_addr, 
        consensusHash: input.consensus_hash
    };

    // note that consensus hash is not used 
    let public_values_solidity_encoded = ChainProofSolidity::abi_encode(&public_values_solidity);

    // Commit the full input, could be optimized
    sp1_zkvm::io::commit_slice(&public_values_solidity_encoded);
}
