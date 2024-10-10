#!/bin/sh
set -e

# MULTI GER ASSERTOR
docker run --platform linux/amd64 --rm -v $(pwd):/output -v $(pwd)/../contracts:/contracts ethereum/solc:0.8.20-alpine - /contracts/multi-ger-getter.sol -o /output --abi --bin --overwrite --optimize
mv -f MultiGERAssertor.abi abi/multigerassertor.abi
mv -f MultiGERAssertor.bin bin/multigerassertor.bin
rm -f IPolygonZkEVMGlobalExitRootV2.abi
rm -f IPolygonZkEVMGlobalExitRootV2.bin

# CONTRACTS FROM REPO
PATH_TO_CONTRACTS_REPO=/Users/arnaub/Documents/polygon/zkevm-contracts

## GER L2
PATH_TO_GER_L2=$PATH_TO_CONTRACTS_REPO/artifacts/contracts/v2/sovereignChains/GlobalExitRootManagerL2SovereignChain.sol/GlobalExitRootManagerL2SovereignChain.json
cat $PATH_TO_GER_L2 | jq .abi > abi/gerl2.abi
cat $PATH_TO_GER_L2 | jq .bytecode | sed 's/^"0x//; s/"$//g' > bin/gerl2.bin

## BRIDGE L2
PATH_TO_BRIDGE_L2=$PATH_TO_CONTRACTS_REPO/artifacts/contracts/v2/sovereignChains/BridgeL2SovereignChain.sol/BridgeL2SovereignChain.json
cat $PATH_TO_BRIDGE_L2 | jq .abi > abi/bridgel2.abi
cat $PATH_TO_BRIDGE_L2 | jq .bytecode | sed 's/^"0x//; s/"$//g' > bin/bridgel2.bin

gen() {
    local package=$1

    abigen --bin bin/${package}.bin --abi abi/${package}.abi --pkg=${package} --out=${package}/${package}.go
}

gen multigerassertor
gen gerl2
gen bridgel2
gen opfaultdisputegame
