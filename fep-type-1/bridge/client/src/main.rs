//! A program that verifies the bridge integrity

#![no_main]
sp1_zkvm::entrypoint!(main);

use polccint_lib::{BridgeCommit, BridgeInput};
use sp1_cc_client_executor::{ClientExecutor, ContractInput};
use alloy_sol_types::SolCall;
use alloy_primitives::{address, Address};
use alloy_sol_macro::sol;

sol! (
    function getLastInjectedGER() public view returns (bytes32);
    function checkGERsAreConsecutiveAndReturnLastLER(bytes32[] GERs) public view returns (bool, bytes32);
    function checkGERsExistance(bytes32[] calldata GERs) public view returns (bool);
);

/// Address of the caller.
const CALLER: Address = address!("0000000000000000000000000000000000000000");

pub fn main() {
    // Read the input.
    let input = sp1_zkvm::io::read::<BridgeInput>();

    // Verify bridge:
    // 1. Get the the last GER of the previous block on L2
    let executor = ClientExecutor::new(
        input.get_last_injected_ger_l2_prev_block_call
    ).unwrap();
    let get_last_injected_ger_l2_prev_block_input = ContractInput {
        contract_address: input.l2_ger_addr,
        caller_address: CALLER,
        calldata: getLastInjectedGERCall {},
    };
    let get_last_injected_ger_l2_prev_block_output = executor.execute(
        get_last_injected_ger_l2_prev_block_input
    ).unwrap();
    let initial_ger = getLastInjectedGERCall::abi_decode_returns(
        &get_last_injected_ger_l2_prev_block_output.contractOutput, true
    ).unwrap()._0;
    if !input.injected_gers.is_empty() {
        assert_eq!(initial_ger, input.injected_gers[0]);
    }
    assert_eq!(
        get_last_injected_ger_l2_prev_block_output.blockHash,
        input.prev_l2_block_hash
    );

    // 2. Check that the GERs are consecutive on L2 at the new block
    let executor: ClientExecutor = ClientExecutor::new(
        input.check_gers_are_consecutive_and_return_last_ler_call_l2_new_block_call
    ).unwrap();
    let check_gers_are_consecutive_and_return_last_ler_call_l2_new_block_input = ContractInput {
        contract_address: input.l2_ger_addr,
        caller_address: CALLER,
        calldata: checkGERsAreConsecutiveAndReturnLastLERCall { GERs: input.injected_gers.clone()},
    };
    let check_gers_are_consecutive_and_return_last_ler_call_l2_new_block_output = executor.execute(
        check_gers_are_consecutive_and_return_last_ler_call_l2_new_block_input
    ).unwrap();
    let call_result = checkGERsAreConsecutiveAndReturnLastLERCall::abi_decode_returns(
        &check_gers_are_consecutive_and_return_last_ler_call_l2_new_block_output.contractOutput, true
    ).unwrap();
    assert_eq!(call_result._0, true);
    assert_eq!(call_result._1, input.new_ler);
    assert_eq!(
        check_gers_are_consecutive_and_return_last_ler_call_l2_new_block_output.blockHash,
        input.new_l2_block_hash
    );

    // 3. Check that the GERs exist on L1
    let executor: ClientExecutor = ClientExecutor::new(
        input.check_gers_existance_l1_call
    ).unwrap();
    let check_gers_existance_l1_input = ContractInput {
        contract_address: input.l1_ger_addr,
        caller_address: CALLER,
        calldata: checkGERsExistanceCall { GERs: input.injected_gers.clone() },
    };
    let check_gers_existance_l1_output = executor.execute(check_gers_existance_l1_input).unwrap();
    let gers_exist = checkGERsExistanceCall::abi_decode_returns(
        &check_gers_existance_l1_output.contractOutput, true
    ).unwrap()._0;
    assert_eq!(gers_exist, true);
    assert_eq!(check_gers_existance_l1_output.blockHash, input.l1_block_hash);

    // Commit the bridge proof.
    let bridge_commit = BridgeCommit {
        l1_block_hash: input.l1_block_hash,
        prev_l2_block_hash: input.prev_l2_block_hash,
        new_l2_block_hash: input.new_l2_block_hash,
        new_ler: input.new_ler,
        l1_ger_addr: input.l1_ger_addr,
        l2_ger_addr: input.l2_ger_addr,
    };   
    sp1_zkvm::io::commit(&bridge_commit);
}
