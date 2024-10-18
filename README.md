# succint-zk-residency

Generate aggregated proofs using SP1. This repo generates proofs of the following things.

The goal fo the projects is to itnegrate multiples chains/proofs into de aggLayer.

This project have the following proofs

- Bridge proofs, which given a state root of an EVM outputs our bridge information reading the state with static calls
- op
  - consensus --> Proof optimism root given a open game on L1
  - chain-proof --> aggregate the Bridge proof + the consensus proof, linking blockhashes
- PoS
  - consensus --> Proof PoS tenderly fork consensus, gets the stakers from a static call from L1 and then prove the tenderly signatures
  - chain-proof --> aggregate the Bridge proof + the consensus proof, linking blockhashes
- Full execution proofs on vanilca clients using local Clique Geth
    - fep-type-1
        - Block --> generate block proofs
        - block-aggregation --> aggregate block proofs
        - chain-proof --> aggregate the Bridge proof + block aggregation proof, linking blockHashes
- AggLayer proof
    - Aggregates chain-proofs and create a plonk proof to be sent on-chin

Example of aggregating all this proofs on our current version of the agglayer on sepolia:
https://sepolia.etherscan.io/tx/0xb5cba33c2225e4890072b14e0a94a059d0e5480eab0b74e5b7b2089f2e1ba492

### Type 1 EVM

See [this](./fep-type-1/README.md) for more info.

### OP

See [this](./op/README.md) for more info.

### Polygon PoS

See [this](./pos/README.md) for more info.