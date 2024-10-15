use serde::Deserialize;

#[derive(Debug, Deserialize)]
pub struct MilestoneResponse {
    pub result: Milestone,
}

#[derive(Debug, Deserialize)]
pub struct Milestone {
    pub proposer: String,
    pub start_block: u64,
    pub end_block: u64,
    pub hash: String,
    pub bor_chain_id: String,
    pub milestone_id: String,
    pub timestamp: u64,
}

#[derive(Debug, Deserialize)]
pub struct TxResponse {
    pub result: TxResponseResult,
}

#[derive(Debug, Deserialize)]
pub struct TxResponseResult {
    pub hash: String,
    pub height: String,
    pub tx: String,
}

#[derive(Debug, Deserialize)]
pub struct BlockResponse {
    pub result: BlockResponseResult,
}

#[derive(Debug, Deserialize)]
pub struct BlockResponseResult {
    pub block: Block,
}

#[derive(Debug, Deserialize)]
pub struct Block {
    pub last_commit: LastCommit,
}

#[derive(Debug, Deserialize)]
pub struct LastCommit {
    pub precommits: Vec<Precommit>,
}

#[derive(Debug, Deserialize)]
pub struct Precommit {
    #[serde(rename = "type")]
    pub type_field: u32,
    pub height: String,
    pub round: String,
    pub block_id: BlockId,
    pub timestamp: String,
    pub validator_address: String,
    pub validator_index: String,
    pub signature: String,
    pub side_tx_results: Option<Vec<SideTxResult>>,
}

#[derive(Debug, Deserialize)]
pub struct BlockId {
    pub hash: String,
    pub parts: Parts,
}

#[derive(Debug, Deserialize)]
pub struct Parts {
    pub total: u32,
    pub hash: String,
}

#[derive(Debug, Deserialize)]
pub struct SideTxResult {
    #[serde(rename = "tx_hash")]
    pub tx_hash: String,
    pub result: i32,
    pub sig: Option<String>,
}

#[derive(Debug, Deserialize)]

pub struct ValidatorSetResponse {
    pub result: Validators,
}

#[derive(Debug, Deserialize)]

pub struct Validators {
    pub validators: Vec<Validator>,
}

#[derive(Debug, Deserialize)]

pub struct Validator {
    pub last_updated: String,
}
