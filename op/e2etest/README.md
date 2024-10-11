# E2E test

## Pre (do it once)

1. Clone the op repo in whatever path u want: `git clone https://github.com/ethereum-optimism/optimism.git`
2. [Optional] Edit `packages/contracts-bedrock/deploy-config/devnetL1-template.json` from the optimism repo:
    - l2BlockTime: from `2` to `20`
    - l1BlockTime: from `1` to `21`
3. Clone the zkevm-contract repo in whatever path u want: `git clone https://github.com/0xPolygonHermez/zkevm-contracts.git`
4. `cd` into the zkevm-contracts repo and checkout the branch: `git checkout test/succint`
5. from the zkevm-contracts repo, run `npm i`
6. run:
    ```
    cp docker/scripts/v2/create_rollup_parameters_docker.json deployment/v2/create_rollup_parameters.json && \
    cp docker/scripts/v2/deploy_parameters_docker.json deployment/v2/deploy_parameters.json
    ```
    
## Run (do it every time)

1. Go to `optimism` repo and run `make devnet-up` (first time is gonna take a long time)
2. Go back to `zkevm-contracts` repo and run `npm run deploy:testnet:v2:localhost`
3. Edit `op/e2etest/main_test.go` with values found on `zkevm-contracts` repo:
    - `gerAddrL1` -> `polygonZkEVMGlobalExitRootAddress` @ `deployment/v2/deploy_output.json`
    - `bridgeAddrL1` -> `polygonZkEVMBridgeAddress` @ `deployment/v2/deploy_output.json`
    - `rollupManagerAddrL1` -> `polygonRollupManagerAddress` @ `deployment/v2/deploy_output.json`
    - `rollupAddrL1` -> `rollupAddress` @ `deployment/v2/create_rollup_output.json`
4. Run the test (from `op/e2etest`): `go test -v -count=1 ./main_test.go`
5. Edit `op/e2etest/main.go` with the info from `.devnet/addresses.json` of the `optimism` repo:
    - `disputeGameFactoryProxyAddr` -> `DisputeGameFactoryProxy`
6. Run the script to get info about the settlement: `go run .` (from `op/e2etest`). If the script doesn't output anything, just wait for a while until the creation of the first game
7. Let's generate the consensus proof, `cd` to `op/consensus/host` and run: `cargo run --release -- --block-number-l1 A --game-index B --game-factory-address C --prev-block-number-l2 D --new-block-number-l2 E --claim-block-hash F --claim-state-root G --claim-message-passer-storage-root H --chain-id-l1 I --chain-id-l2 J`. Replace the following:
    - A: any block number on L1 that contains the game. Just get the last blcok number: `curl --location --request POST 'http://localhost:8545' --header 'Content-Type: application/json' --data-raw '{"jsonrpc": "2.0", "method": "eth_blockNumber", "params": [], "id": 1}'`. Remember to change the hex to decimal
    - B: Game index from the output of `main.go` run on previosu step
    - C: Address of the game factory, the one also added in `main.go` on prev step
    - D: The last block number that was verified (rignt now can be any old block num, in the future the one that matches with the last block hash verified on the contract)
    - E: The block number on the `Final / Block number` mathing the game index picked for `B`
    - F: The block hash on the `Final / Block hash` mathing the game index picked for `B`
    - G: The block root on the `Final / Block root` mathing the game index picked for `B`
    - H: The MessagePasserStorageRoot on the `Final / MessagePasserStorageRoot` mathing the game index picked for `B`
    - I: Chain ID of L1, if running local should be 900
    - j: Chain ID of L2, if running local should be 901
    - example: `cargo run --release -- --block-number-l1 2379 --game-index 444 --game-factory-address 0x5e8176772863842cd9d692c9E793dc4958626E69 --prev-block-number-l2 1 --new-block-number-l2 38071 --claim-block-hash 0xd95abbfd5ead213824bd8bfcff1f56d8b00c32e3b50adc3f59e1d464219a6dac --claim-state-root 0xdbb20475223c40978a4ae36957eb91bb9d02466a2ae0f18d4566c2b186c5d7dc --claim-message-passer-storage-root 0x8ed4baae3a927be3dea54996b4d5899f8c01e7594bf50b17dc1e741388ce3d12 --chain-id-l1 900 --chain-id-l2 901`

In order to stop the running containers, from the `optimism` repo run: `make devnet-down`