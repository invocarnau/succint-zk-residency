#![no_main]
sp1_zkvm::entrypoint!(main);

use polccint_lib::{BlockCommit, BlockInput};
use rsp_client_executor::{ClientExecutor, EthereumVariant};

pub fn main() {
    // Read the input.
    let input = sp1_zkvm::io::read_vec();
    let input = bincode::deserialize::<BlockInput>(&input).unwrap();

    // Execute the block.
    let executor = ClientExecutor;
    let header = executor
        .execute::<EthereumVariant>(input.header)
        .expect("failed to execute client");
    assert_eq!(
        input.prev_block_hash.as_slice(),
        header.parent_hash.as_slice()
    );
    let block_hash = header.hash_slow();

    // Commit.
    let block_commit = BlockCommit {
        prev_block_hash: input.prev_block_hash,
        new_block_hash: block_hash,
    };
    sp1_zkvm::io::commit(&block_commit);
}
