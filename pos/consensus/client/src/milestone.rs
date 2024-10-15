use crate::helper::*;
use std::{collections::HashMap, str::FromStr};

use bincode;

use alloy_primitives::{address, keccak256, Address, FixedBytes, Uint};
use alloy_sol_types::{sol, SolCall};
use reth_primitives::Header;
use sp1_cc_client_executor::{io::EVMStateSketch, ClientExecutor, ContractInput};

sol! {
    /// The public values encoded as a struct that can be easily deserialized inside Solidity.
    struct PublicValuesStruct {
        bytes32 prev_bor_block_hash;
        bytes32 new_bor_block_hash;
        bytes32 l1_block_hash;
    }
}

sol! {
    contract ConsensusProofVerifier {
        bytes32 public lastVerifiedBorBlockHash;
        function verifyConsensusProof(bytes calldata _proofBytes, bytes32 new_bor_block_hash, bytes32 l1_block_hash) public view;
        function getEncodedValidatorInfo() public view returns(address[] memory, uint256[] memory, uint256);
    }
}

#[derive(Clone, Debug)]
pub struct MilestoneProofInputs {
    // heimdall related data
    pub tx_data: String,
    pub tx_hash: FixedBytes<32>,
    pub precommits: Vec<Vec<u8>>,
    pub sigs: Vec<String>,
    pub signers: Vec<Address>,

    // bor related data
    pub bor_header: Header,
    pub prev_bor_header: Header,

    // l1 related data
    pub state_sketch_bytes: Vec<u8>,
    pub l1_block_hash: FixedBytes<32>,
}

#[derive(Clone, Debug)]
pub struct MilestoneProofOutputs {
    pub prev_bor_hash: FixedBytes<32>,
    pub new_bor_hash: FixedBytes<32>,
    pub l1_block_hash: FixedBytes<32>,
}

pub struct MilestoneProver {
    inputs: MilestoneProofInputs,
}

impl MilestoneProver {
    pub fn init(inputs: MilestoneProofInputs) -> Self {
        MilestoneProver { inputs }
    }

    pub fn prove(&self) -> MilestoneProofOutputs {
        let a: &str = "0x01Eb85F73dA540C66CE1d4262BF7F80d5BA6CF89";
        let verifier_contract: Address = Address::from_str(a).unwrap();
        let caller_address: Address = address!("0000000000000000000000000000000000000000");

        // Verify if the transaction data provided is actually correct or not
        let milestone = verify_tx_data(&self.inputs.tx_data, &self.inputs.tx_hash);

        // Calculate the bor block hash from the given header
        let bor_block_hash = self.inputs.bor_header.hash_slow();

        // Verify if the bor block header matches with the milestone or not
        assert_eq!(
            milestone.end_block, self.inputs.bor_header.number,
            "block number mismatch between milestone and bor block header"
        );
        assert_eq!(
            milestone.hash,
            bor_block_hash.to_vec(),
            "block hash mismatch between milestone and bor block header"
        );

        // Make sure that we have equal number of precommits, signatures and signers.
        assert_eq!(self.inputs.precommits.len(), self.inputs.sigs.len());
        assert_eq!(self.inputs.sigs.len(), self.inputs.signers.len());

        let state_sketch =
            bincode::deserialize::<EVMStateSketch>(&self.inputs.state_sketch_bytes).unwrap();

        // Initialize the client executor with the state sketch.
        // This step also validates all of the storage against the provided state root.
        let executor = ClientExecutor::new(state_sketch).unwrap();

        // Execute the `getEncodedValidatorInfo` call using the client executor to fetch the
        // active validator's info from L1.
        let call = ConsensusProofVerifier::getEncodedValidatorInfoCall {};
        let input = ContractInput {
            contract_address: verifier_contract,
            caller_address,
            calldata: call.clone(),
        };
        let output = executor.execute(input).unwrap();
        let response = ConsensusProofVerifier::getEncodedValidatorInfoCall::abi_decode_returns(
            &output.contractOutput,
            true,
        )
        .unwrap();

        // Extract the signers, powers, and total_power from the response.
        let signers = response._0;
        let powers = response._1;
        let total_power = response._2;

        let mut majority_power: Uint<256, 4> = Uint::from(0);
        let mut validator_stakes = HashMap::new();
        for (i, signer) in signers.iter().enumerate() {
            validator_stakes.insert(signer, powers[i]);
        }

        // Execute the `lastVerifiedBorBlockHash` call using the client executor to fetch the
        // last verified bor block hash.
        let call = ConsensusProofVerifier::lastVerifiedBorBlockHashCall {};
        let input = ContractInput {
            contract_address: verifier_contract,
            caller_address,
            calldata: call.clone(),
        };
        let output = executor.execute(input).unwrap();
        let last_verified_bor_block_hash_return =
            ConsensusProofVerifier::lastVerifiedBorBlockHashCall::abi_decode_returns(
                &output.contractOutput,
                true,
            )
            .unwrap();
        let prev_bor_hash = last_verified_bor_block_hash_return.lastVerifiedBorBlockHash;

        // If we're running prover for the first time, we won't have a previous bor block hash. Skip
        // all validations if that's the case else verify against that.
        if !prev_bor_hash.is_zero() {
            // Verify if the `prev_bor_header`s hash matches with the one we fetched from the contract.
            let prev_derived_bor_hash = self.inputs.prev_bor_header.hash_slow();
            assert_eq!(
                prev_derived_bor_hash, prev_bor_hash,
                "prev bor hash mismatch"
            );

            // Ensure that we're maintaining sequence of bor blocks and are not proving anything random
            assert!(
                self.inputs.bor_header.number > self.inputs.prev_bor_header.number,
                "bor block is not sequential"
            );
        }

        // Verify that the signatures generated by signing the precommit message are indeed signed
        // by the given validators.
        for i in 0..self.inputs.precommits.len() {
            // Validate if the signer of this precommit message is a part of the active validator
            // set or not.
            assert!(validator_stakes.contains_key(&self.inputs.signers[i]));

            // Verify if the precommit message is for the same milestone transaction or not.
            let precommit = &self.inputs.precommits[i];
            verify_precommit(&mut precommit.clone(), &self.inputs.tx_hash);

            // Verify if the message is indeed signed by the validator or not.
            verify_signature(
                self.inputs.sigs[i].as_str(),
                &keccak256(precommit),
                self.inputs.signers[i],
            );

            // Add the power of the validator to the majority power
            majority_power =
                majority_power.add_mod(validator_stakes[&self.inputs.signers[i]], Uint::MAX);
        }

        // Check if the majority power is greater than 2/3rd of the total power
        let expected_majority = total_power
            .mul_mod(Uint::from(2), Uint::MAX)
            .div_ceil(Uint::from(3));
        if majority_power <= expected_majority {
            panic!("Majority voting power is less than 2/3rd of the total power, total_power: {}, majority_power: {}", total_power, majority_power);
        }

        MilestoneProofOutputs {
            prev_bor_hash,
            new_bor_hash: bor_block_hash,
            l1_block_hash: self.inputs.l1_block_hash,
        }
    }
}
