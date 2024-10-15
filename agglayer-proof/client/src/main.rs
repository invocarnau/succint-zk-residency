//! A program that aggregates the proofs of EVM blocks

#![no_main]
sp1_zkvm::entrypoint!(main);

use polccint_lib::{ChainProof, AggLayerProof, AggLayerProofSolidity, ChainProofSolidity};
use polccint_lib::constants::{CHAIN_VK};
use alloy_sol_types::SolType;

use sha2::{Digest,Sha256};
use bincode;

pub fn main() {
    // Read the input.

    // could pass just an array of blockhashes
    let input: AggLayerProof = sp1_zkvm::io::read::<AggLayerProof>();

    // Verify the aggregated block proof

    let mut public_values_solidity: AggLayerProofSolidity = AggLayerProofSolidity {
        chain_proofs: Vec::new(),
    };

    // Recreate the block aggregation commit
    for chain_proof in input.chain_proofs {
        let serialized_public_values_chain= bincode::serialize(&chain_proof).unwrap();
        let public_values_digest_chain = Sha256::digest(serialized_public_values_chain);
        sp1_zkvm::lib::verify::verify_sp1_proof(
            &CHAIN_VK, 
            &public_values_digest_chain.into()
        );
        public_values_solidity.chain_proofs.push(ChainProofSolidity{
             prev_l2_block_hash:chain_proof.prev_l2_block_hash,
             new_l2_block_hash:chain_proof.new_l2_block_hash,
             l1_block_hash:chain_proof.l1_block_hash,
             new_ler:chain_proof.new_ler,
             l1_ger_addr:chain_proof.l1_ger_addr,
             l2_ger_addr: chain_proof.l2_ger_addr, 
             consensus_hash: chain_proof.consensus_hash, 
        });
    }
    
    // note that consensus hash is not used 
    let public_values_solidity_encoded = AggLayerProofSolidity::abi_encode(&public_values_solidity);

    // Commit the full input, could be optimized
    sp1_zkvm::io::commit_slice(&public_values_solidity_encoded);
}
