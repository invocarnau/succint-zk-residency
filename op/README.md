# OP

The goal of this PoC is to integrate an existing OP chain to the uLxLy, using:

![](./proof.drawio.png)

## Proofs

### Execution proof

- Executes OP blocks
- Aggregates them
- Asserts the validity of injected GERs
- Commits to prev block hash + new block hash + new LER

### Consensus

- Asserts that the state root published on L1 is part of the new block hash that is beeing committed on the execution proof

## Relevant information

- [0x473300df21d047806a082244b417f96b32f13a33](https://explorer.optimism.io/address/0x473300df21d047806a082244b417f96b32f13a33) is the currrent proposer for OP mainent
- [0xe5965Ab5962eDc7477C8520243A95517CD252fA9](https://etherscan.io/address/0xe5965Ab5962eDc7477C8520243A95517CD252fA9) is the address where games are created
    - this happens once per hour
    - an event is emitted when a game is created: `DisputeGameCreated`, this event contains the address of the game
- A created game includes the l2 block number and the state root

All of this is subject to change. As the OP protocol is still evolving, and in particular fraud proofs have been recently introduced and later on halted...

## Doubts:

- It's unclear to me if the game factory is shared across different chains, and therefore we should check chain ID on each game.

## Usage

In order to generate the `chain proof` for an OP Stack chain, run the following command (from the `./op` directory):

```bash
go run . proof \
--l1-chain-id A \
--l2-chain-id B \
--game-factory C \
--prev-l2-block-num D \
--ger-l1 E \
--ger-l2 F \
--network-id G
```

Make sure to replace the following:

- A: Chain ID of the L1. If you want to use our testnet env, use `11155111` (Sepolia)
- B: Chain ID of the L1. If you want to use our testnet env, use `11155420` (OP Sepolia)
- C: Address of the game factory proxy contract, deployed on L1. If you want to use our testnet env, use `0x05f9613adb30026ffd634f38e5c4dfd30a197fa1`
- D: Last L2 block number that has been verified. This is intentionally non inforced, so you can choose any value greater than the block number where the L2 contracts were deployed, for instance `18655973` (recommended to use a recent but not too recent block number such as current block num - 1000, as this will incure calls on the L2 rpc)
- E: Addrees of the Global Exit Root contract on L1. If you want to use our testnet env, use `0xBa36ee0dBDC8fe4c2f82dD75506CF836E0205974`
- F: Addrees of the Global Exit Root contract on L2. If you want to use our testnet env, use `0x0518576bC94CF5C4078f5f7bAAea8f2DF5fe61FC`
- G: Network ID of the rollup, as per the rollup manager. If you want to use our testnet env, use `TBD`

**IMPORTANT NOTE:** make sure that you have a valid `.env` file on the root of the repo with RPC URLs for the corresponding Chain IDs.