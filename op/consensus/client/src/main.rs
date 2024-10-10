//! A program that verifies the consensus of an OP chain, in particular, it checks:
//! - A given L2 block hash exists in a L1 game contract => using the "rootClaim"
//! - The factory of this contract is a given address
//! - The block number is greater than the one of the previous block
//! Note that this means that the correct execution of this blocks are not checked, and therefor
//! those blocks could be reorged via fraud proofs.

#![no_main]
sp1_zkvm::entrypoint!(main);

use polccint_lib::op::{OPConsensusCommit, OpConsensusInput};
use polccint_lib::constants::CALLER;
use sp1_cc_client_executor::{ClientExecutor, ContractInput};
use alloy_sol_types::SolCall;
use alloy_sol_macro::sol;

sol! (
    interface GameFactory {
        function gameAtIndex(uint256 _index) view returns(uint32 gameType_, uint64 timestamp_, address game);
    }

    interface Game {
        function rootClaim() pure returns(bytes32 rootClaim);
    }
);

pub fn main() {
    let input = sp1_zkvm::io::read::<OpConsensusInput>();

    // 1. Check that the new block header matches with the block hash of the root claim
    assert_eq!(input.new_l2_block_header.hash_slow(), input.root_claim_pre_image.block_hash);

    // 2. Get the game address from the game contract factory
    // Load exeutor with the L1 block sketch
    let executor_get_game_from_factory: ClientExecutor = ClientExecutor::new(
        input.get_game_from_factory_sketch
    ).unwrap();
    let get_game_from_factory_input = ContractInput {
        contract_address: input.game_factory_address,
        caller_address: CALLER,
        calldata: GameFactory::gameAtIndexCall {_index: input.game_index},
    };
    // Execute the static call
    let get_game_from_factory_output = executor_get_game_from_factory.execute(
        get_game_from_factory_input
    ).unwrap();
    // Decode root claim from the result
    let game_address = GameFactory::gameAtIndexCall::abi_decode_returns(
        &get_game_from_factory_output.contractOutput, true
    ).unwrap().game;

    // 3. Get the root claim
    // Load exeutor with the L1 block sketch
    let executor_get_root_claim: ClientExecutor = ClientExecutor::new(
        input.get_root_claim_sketch
    ).unwrap();
    let get_root_claim_input = ContractInput {
        contract_address: game_address,
        caller_address: CALLER,
        calldata: Game::rootClaimCall {},
    };
    // Execute the static call
    let get_root_claim_output = executor_get_root_claim.execute(
        get_root_claim_input
    ).unwrap();
    // Decode root claim from the result
    let root_claim = Game::rootClaimCall::abi_decode_returns(
        &get_root_claim_output.contractOutput, true
    ).unwrap().rootClaim;
    let expected_root_claim = input.root_claim_pre_image.hash();
    assert_eq!(root_claim, expected_root_claim);

    // 4. Check that the current block number is older than the previous one
    assert!(input.new_l2_block_header.number > input.prev_l2_block_header.number);

    // 5. Assert that same block has been used for both calls
    assert_eq!(get_game_from_factory_output.blockHash, get_root_claim_output.blockHash);

    // Commit
    let consensus_commit = OPConsensusCommit {
        game_factory_address: input.game_factory_address,
        l1_block_hash: get_root_claim_output.blockHash,
        prev_l2_block_hash: input.prev_l2_block_header.hash_slow(),
        new_l2_block_hash: input.root_claim_pre_image.block_hash,
    };   

    sp1_zkvm::io::commit(&consensus_commit);
}
