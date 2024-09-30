use alloy_primitives::{B256, Address};
use rsp_client_executor::io::ClientExecutorInput;
use sp1_cc_client_executor::io::EVMStateSketch;
use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct BlockInput {
    pub prev_block_hash: B256,
    pub header: ClientExecutorInput,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct BlockCommit {
    pub prev_block_hash: B256,
    pub new_block_hash: B256,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct BlockAggregationInput {
    pub block_commits: Vec<BlockCommit>,
    pub block_vkey: [u32; 8],
    pub prev_l2_block_hash: B256,
    // proving the bridge
    pub l1_block_hash: B256,
    pub l1_multi_ger_assertor: Address,
    pub l2_ger: Address,
    pub get_l2_ger_index_prev_block: EVMStateSketch,
    pub get_l2_gers: EVMStateSketch,
    pub check_l1_gers_existance: EVMStateSketch,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct BlockAggregationCommit {
    pub prev_l2_block_hash: B256,
    pub new_l2_block_hash: B256,
    pub l1_block_hash: B256,
}
