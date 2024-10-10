# E2E test

1. Run `go test -v -count=1 ./main_test.go` from the `e2etest` folder
2. On a separated terminal, run `docker compose logs -f test-fep-type1-cdk 2>&1 | grep "inject GER tx"` from the `e2etest` folder. This will output the block number at which the GER was injected on L2.
3. run `docker compose logs -f test-fep-type1-cdk 2>&1 | grep "claim tx with id"` from the `e2etest` folder. This will output the block number at which the bridge claim tx was mined on L2
4. ==TODO==: explain how to generate a proof that includes both bridge injection and claim from here
