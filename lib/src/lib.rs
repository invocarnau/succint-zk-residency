use alloy_primitives::{B256};
use serde::{Deserialize, Serialize};
use alloy_sol_types::sol;


pub mod constants;
pub mod bridge;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct BlockCommit {
    pub prev_block_hash: B256,
    pub new_block_hash: B256,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct BlockAggregationInput {
    pub block_commits: Vec<BlockCommit>,
    pub block_vkey: [u32; 8],
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct BlockAggregationCommit {
    pub prev_l2_block_hash: B256,
    pub new_l2_block_hash: B256,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct FinalAggregationInput {
    pub block_aggregation_commit: BlockAggregationCommit,
    pub bridge_commit: bridge::BridgeCommit,
}

/// A fixture that can be used to test the verification of SP1 zkVM proofs inside Solidity.
#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct SP1CCProofFixture {
    pub vkey: String,
    pub public_values: String,
    pub proof: String,
}


sol! {
    #[derive(Debug, Serialize, Deserialize)]
    struct PublicValuesFinalAggregationSolidity {
        bytes32 prev_l2_block_hash;
        bytes32 new_l2_block_hash;
        bytes32 l1_block_hash;
        bytes32 new_ler;
        address l1_ger_addr;
        address l2_ger_addr;   
    }
}
