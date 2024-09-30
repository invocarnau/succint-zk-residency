//! A program that aggregates the proofs of EVM blocks

#![no_main]
sp1_zkvm::entrypoint!(main);

use alloy_sol_types::{abi,SolCall};
use polccint_lib::{BlockAggregationCommit, BlockAggregationInput};
use sha2::Digest;
use sha2::Sha256;
use alloy_primitives::{address, Address, Bytes, B256};
use alloy_sol_macro::sol;
use alloy_sol_types::SolValue;
use bincode;
use sp1_cc_client_executor::{io::EVMStateSketch, ClientExecutor, ContractInput};

sol! (
    function checkGERsExistance(bytes32[] calldata GERs) public view returns (bool);
    function getInjectedGERIndex() public view returns (uint256); // TODO: change for get last injected GER
    function getGERsFromIndex(uint256 index) public view returns (bytes32[]); // TODO: pass a list of GERs and first index, check that GERs are consecutive and go until last index
);

/// Address of the caller.
const CALLER: Address = address!("0000000000000000000000000000000000000000");

pub fn main() {
    // Read the input.
    let input = sp1_zkvm::io::read::<BlockAggregationInput>();

    // Confirm that the blocks are sequential.
    assert!(!input.block_commits.is_empty());
    assert_eq!(
        input.prev_l2_block_hash,
        input.block_commits[0].prev_block_hash
    );
    input.block_commits.windows(2).for_each(|pair| {
        let (prev_block, block) = (&pair[0], &pair[1]);
        assert_eq!(prev_block.new_block_hash, block.prev_block_hash);
    });

    // Verify the block proofs.
    for i in 0..input.block_commits.len() {
        let public_values = &input.block_commits[i];
        let serialized_public_values = bincode::serialize(public_values).unwrap();
        let public_values_digest = Sha256::digest(serialized_public_values);
        sp1_zkvm::lib::verify::verify_sp1_proof(&input.block_vkey, &public_values_digest.into());
    }
    let new_block_hash = input.block_commits.last().unwrap().new_block_hash;

    // Verify bridge:
    // 1. Get the index of the injected GER in the previous block
    let executor = ClientExecutor::new(inputs.get_l2_ger_index_prev_block).unwrap();
    let get_injected_ger_index_prev_block_call = ContractInput {
        contract_address: input.l2_ger,
        caller_address: CALLER,
        calldata: L2GER::getInjectedGERIndex {},
    };
    let get_injected_ger_index_prev_block_call_output = executor.execute(get_injected_ger_index_prev_block_call).unwrap();
    assert_eq!(get_injected_ger_index_prev_block_call_output.blockHash, input.prev_l2_block_hash);
    // TODO: convert bytes to uint256
    let initial_ger_index = get_injected_ger_index_prev_block_call_output.contractOutput;

    // 2. Get the injected GERs
    let get_gers_from_index_call = ContractInput {
        contract_address: input.l2_ger,
        caller_address: CALLER,
        calldata: L2GER::getGERsFromIndex { index: initial_ger_index },
    };
    let get_gers_from_index_call_output = executor.execute(get_gers_from_index_call).unwrap();
    assert_eq!(get_gers_from_index_call_output.blockHash, new_block_hash);
    let new_injected_gers = getGERsFromIndexCall::abi_decode_returns(&get_gers_from_index_call_output.contractOutput, false).unwrap()._0;

    // 3. Check that the GERs exist on L1
    let check_gers_existance_call = ContractInput {
        contract_address: input.l1_multi_ger_assertor,
        caller_address: CALLER,
        calldata: checkGERsExistance { gers: new_injected_gers },
    };
    let check_gers_existance_call_output = executor.execute(check_gers_existance_call).unwrap();
    assert_eq!(check_gers_existance_call_output.blockHash, input.l1_block_hash);
    // TODO: convert bytes to bool
    assert_eq!(check_gers_existance_call_output.contractOutput, true);

    // Commit the block aggregation proof.
    let block_aggregation_commit = BlockAggregationCommit {
        prev_l2_block_hash: input.prev_l2_block_hash,
        new_l2_block_hash: new_block_hash,
        l1_block_hash: input.l1_block_hash,
        // TODO: commit LER
    };
    sp1_zkvm::io::commit(&block_aggregation_commit);
}
