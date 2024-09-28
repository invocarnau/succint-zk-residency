# FEP type I

The goal behind this PoC is to integrate an EVM type I prover to the uLxLy, using the full execution proof (FEP) path. The challenge is to be able to proof the bridge (consumed GERs + produced LERs), while using 100% vanilla client (gETH, rETH, ...)

![](./proof.drawio.png)

In order to achieve this the strategy will be to:

- proof vanilla blocks, and aggregate them using [RSP](https://github.com/succinctlabs/rsp/) as is
- have an aggreagation proof that:
    - proof the blocks proofs
    - get the last injected GER index from the previos block header (init GER index) using a storage proof
    - get the last injected GER index from the new block header (last GER index) using a storage proof
    - get all the GERs between init/last GER index, and assert them with public input (smart contreact will make sure that those exist)
    - get the last LER from the new block header using a storage proof, include it as output