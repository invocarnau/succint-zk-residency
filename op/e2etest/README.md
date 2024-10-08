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
6. Run the script to get info about the settlement: `go run .` (from `op/e2etest`)

In order to stop the running containers, from the `optimism` repo run: `make devnet-down`