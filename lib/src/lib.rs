use alloy_primitives::{B256, Address};
use sp1_cc_client_executor::io::EVMStateSketch;
use serde::{Deserialize, Serialize};
use alloy_sol_types::sol;


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
    pub block_vkey: [u32; 8], 
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct BridgeInput {
    pub l1_ger_addr: Address, // this could be constant
    pub l2_ger_addr: Address, // this could be retrieve fro the Bridge L1 which is constant

    pub injected_gers: Vec<B256>,
    pub injected_ger_count_sketch: EVMStateSketch,
    pub check_injected_gers_and_return_ler_sketch: EVMStateSketch,
    pub check_gers_existance_sketch: EVMStateSketch,
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

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct FinalAggregationInput {
    pub block_vkey_aggregation: [u32; 8], 
    pub block_aggregation_commit: BlockAggregationCommit,
    pub block_vkey_bridge: [u32; 8], 
    pub bridge_commit: BridgeCommit,
}


/// A fixture that can be used to test the verification of SP1 zkVM proofs inside Solidity.
#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct SP1CCProofFixture {
    pub vkey: String,
    pub public_values: String,
    pub proof: String,
}

pub mod constants;


sol! {
    /// The public values encoded as a struct that can be easily deserialized inside Solidity.
    struct PublicValuesFinalAggregationSolidity {
        bytes32 block_vkey_aggregation;
        bytes32 block_vkey;
        bytes32 block_vkey_bridge; 
        bytes32 prev_l2_block_hash;
        bytes32 new_l2_block_hash;
        bytes32 l1_block_hash;
        bytes32 new_ler;
        address l1_ger_addr;
        address l2_ger_addr;   
    }
}