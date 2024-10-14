#![no_main]
sp1_zkvm::entrypoint!(main);

use polccint_lib::BlockCommit;
use rsp_client_executor::{io::ClientExecutorInput, ClientExecutor, CliqueShanghaiChainIDVariant};

pub fn main() {
    // Read the input.
    let input = sp1_zkvm::io::read_vec();
    let input = bincode::deserialize::<ClientExecutorInput>(&input).unwrap();
    let chain_id = sp1_zkvm::io::read::<u64>();

   // Execute the block.
   let executor = ClientExecutor;
   let header = executor.execute_with_chain_id::<CliqueShanghaiChainIDVariant>(chain_id, input).expect("failed to execute client");

    // Commit.
    let block_commit = BlockCommit {
        prev_block_hash: header.parent_hash,
        new_block_hash: header.hash_slow(),
    };
    sp1_zkvm::io::commit(&block_commit);
}