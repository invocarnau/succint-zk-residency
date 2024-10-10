//! A program that verifies the bridge integrity

#![no_main]
sp1_zkvm::entrypoint!(main);

use polccint_lib::{BridgeCommit, BridgeInput};
use polccint_lib::constants::{CALLER};

use sp1_cc_client_executor::{ClientExecutor, ContractInput};
use alloy_sol_types::SolCall;
use alloy_sol_macro::sol;


// try what happens if the calls revert?¿
sol! (
    interface GlobalExitRootManagerL2SovereignChain {
        function injectedGERCount() public view returns (uint256 injectedGerCount);
        function checkInjectedGERsAndReturnLER(uint256 previousInjectedGERCount, bytes32[] injectedGERs) public view returns (bool success, bytes32 localExitRoot);
    }

    interface GlobalExitRootScrapper {
        function checkGERsExistance(bytes32[] calldata globalExitRoots) public view returns (bool success);
    }
);


pub fn main() {
    // Read the input.
    //let sbridge_input_bytes = sp1_zkvm::io::read::<Vec<u8>>();
    // let input = bincode::deserialize::<BridgeInput>(&sbridge_input_bytes).unwrap();
    let input: BridgeInput = sp1_zkvm::io::read::<BridgeInput>();

    // Verify bridge:
    // 1. Get the the last injecter GER of the previous block on L2

    // Load exeutor with the previous L2 block sketch
    let executor_injected_ger_count: ClientExecutor = ClientExecutor::new(
        input.injected_ger_count_sketch
    ).unwrap();

    let get_last_injected_ger_input = ContractInput {
        contract_address: input.l2_ger_addr,
        caller_address: CALLER,
        calldata: GlobalExitRootManagerL2SovereignChain::injectedGERCountCall {},
    };

    // Execute the static call
    let injected_ger_count_call_output = executor_injected_ger_count.execute(
        get_last_injected_ger_input
    ).unwrap();

    // Decode ger count from the result
    let initial_ger_count = GlobalExitRootManagerL2SovereignChain::injectedGERCountCall::abi_decode_returns(
        &injected_ger_count_call_output.contractOutput, true
    ).unwrap().injectedGerCount;

    // 2. Check that the GERs are consecutive on L2 at the new block
    let executor_injected_gers_and_return_ler: ClientExecutor = ClientExecutor::new(
        input.check_injected_gers_and_return_ler_sketch
    ).unwrap();

    let check_injected_gers_and_return_ler_sketch_input = ContractInput {
        contract_address: input.l2_ger_addr,
        caller_address: CALLER,
        calldata: GlobalExitRootManagerL2SovereignChain::checkInjectedGERsAndReturnLERCall { 
            previousInjectedGERCount: initial_ger_count, 
            injectedGERs: input.injected_gers.clone()
        },
    };

    // Execute the static call
    let check_injected_gers_and_return_ler_call_output = 
    executor_injected_gers_and_return_ler.execute(
        check_injected_gers_and_return_ler_sketch_input
    ).unwrap();

    // Decode the call result
    let check_injected_gers_and_return_ler_call_output_decoded = 
    GlobalExitRootManagerL2SovereignChain::checkInjectedGERsAndReturnLERCall::abi_decode_returns(
        &check_injected_gers_and_return_ler_call_output.contractOutput, true
    ).unwrap();

    // Check that the check was successful
    assert_eq!(check_injected_gers_and_return_ler_call_output_decoded.success, true);

    // 3. Check that the GERs exist on L1
    let executor_check_gers_existance: ClientExecutor = ClientExecutor::new(
        input.check_gers_existance_sketch
    ).unwrap();

    let check_gers_existance_l1_input = ContractInput {
        contract_address: input.l1_ger_addr,
        caller_address: CALLER,
        calldata: GlobalExitRootScrapper::checkGERsExistanceCall { globalExitRoots: input.injected_gers.clone() },
    };
    let check_gers_existance_l1_call_output = executor_check_gers_existance
        .execute(check_gers_existance_l1_input)
        .unwrap();

    // Decode the call result
    let check_gers_existance_l1_output_decoded = GlobalExitRootScrapper::checkGERsExistanceCall::
    abi_decode_returns(
        &check_gers_existance_l1_call_output.contractOutput, true
    ).unwrap();

    // Check that the check was successful
    assert_eq!(check_gers_existance_l1_output_decoded.success, true);

    // Commit the bridge proof.
    let bridge_commit: BridgeCommit = BridgeCommit {
        l1_block_hash: check_gers_existance_l1_call_output.blockHash,
        prev_l2_block_hash: injected_ger_count_call_output.blockHash,
        new_l2_block_hash: check_injected_gers_and_return_ler_call_output.blockHash,
        new_ler: check_injected_gers_and_return_ler_call_output_decoded.localExitRoot,
        l1_ger_addr: input.l1_ger_addr,
        l2_ger_addr: input.l2_ger_addr,
    };   

    sp1_zkvm::io::commit(&bridge_commit);
}
