#!/bin/sh

set -e

PATH_TO_CONTRACTS_REPO=/Users/arnaub/Documents/polygon/zkevm-contracts

docker run --platform linux/amd64 --rm -v $(pwd):/output -v $(pwd)/../contracts:/contracts ethereum/solc:0.8.20-alpine - /contracts/multi-ger-getter.sol -o /output --abi --bin --overwrite --optimize
mv -f MultiGERAssertor.abi abi/multigerassertor.abi
mv -f MultiGERAssertor.bin bin/multigerassertor.bin
rm -f IPolygonZkEVMGlobalExitRootV2.abi
rm -f IPolygonZkEVMGlobalExitRootV2.bin

gen() {
    local package=$1

    abigen --bin bin/${package}.bin --abi abi/${package}.abi --pkg=${package} --out=${package}/${package}.go
}

gen multigerassertor