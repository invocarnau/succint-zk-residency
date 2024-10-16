use alloy_primitives::B256;
use serde::{Deserialize, Serialize};

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
