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
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct BlockAggregationCommit {
    pub prev_l2_block_hash: B256,
    pub new_l2_block_hash: B256,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct BridgeInput {
    pub l1_block_hash: B256,
    pub prev_l2_block_hash: B256,
    pub new_l2_block_hash: B256,
    pub new_ler: B256,
    pub l1_ger_addr: Address,
    pub l2_ger_addr: Address,

    pub injected_gers: Vec<B256>,
    pub get_last_injected_ger_l2_prev_block_call: EVMStateSketch,
    pub check_gers_are_consecutive_and_return_last_ler_call_l2_new_block_call: EVMStateSketch,
    pub check_gers_existance_l1_call: EVMStateSketch,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct BridgeCommit {
    pub l1_block_hash: B256,
    pub prev_l2_block_hash: B256,
    pub new_l2_block_hash: B256,
    pub new_ler: B256,
    pub l1_ger_addr: Address,
    pub l2_ger_addr: Address,
}
