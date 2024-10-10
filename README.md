# succint-zk-residency

Code done on the Succint ZK Residenct


// block
./your_script.sh 13 8

cargo run --release -- --chain-id 42069 --block-number 13 --prove 

// aggregation
cargo run --release -- --chain-id 42069 --block-number 13 --block-range 8 --prove 

// Bridge
cargo run --release -- --chain-id-l1 11155111 --chain-id-l2 42069 --block-number-l1 6846763 --block-number-l2 13 --block-range 8 --contract-ger-l1 "0xBa36ee0dBDC8fe4c2f82dD75506CF836E0205974" --contract-ger-l2 "0xB62FbB67aef805c9A8812534e91971A7b5605D56" --imported-gers-hex "0x1902b9cd537c9a0d4167adbbd4173c04b02b573349196d76d96f1a2475f7c5b0" --prove 


// aggregation final
cargo run --release -- --chain-id 42069 --block-number 13 --block-range 8 --prove 


