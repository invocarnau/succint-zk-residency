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


### PoS 

In order to plug into agglayer, 3 kinds of proofs are being generated in PoS.

1. Consensus Proof: A proof asserting that a particular execution layer block has been voted upon by majority (>2/3+1) of the validator set.
2. Bridge Proof: Check if the global exit roots have been updated correctly on L1 based on LxLy design.
3. Chain Proof: An aggregated proof for both consensus and bridge.

The aggregated chain proof will be sent to AggLayer which will further aggregate multiple chain proofs. The aggregated proofs are also verifiable on-chain. 

#### Steps to generate proofs

1. Make sure you have the `.env` file updated.
2. Run the command below in the `pos` directory:
```bash
go run main.go proof \
    --l1 11155111 \
    --l2 80002 \
    --id 332682 \
    --hash 0x88a07f460e45ca55db87cdf00f7d0392a92bfaa276237944fef62714aad2a841 \
    --from 13296845 \
    --gl1 0xe8085E052669cA2CDeCe52123A3E77461AA31494 \
    --gl2 0x0707c0726726D2334E6E304763CBDE922170d8cf
```

The first 2 flags are l1 and l2 chain ids, id and hash are the latest milestone id and transaction
hash in heimdall, from is the last l2 block number proved (can be used the same for now to demo), 
gl1 and gl2 are the global exit roots contracts on l1 and l2 respectively. Note that by l2, it means
the execution layer of pos i.e. bor.

This script will generate the consensus proof and bridge proof in parallel and on successfull generation
of both of them, it will generate an aggregate proof and save everything in their respective locations.

Cycle count:
- Consensus: 15369480 (compressed proof; for consensus on amoy, ~20 signature verifications)
- Bridge: 1964872
- Chain: 339619

### OP

TODO
