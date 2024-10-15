// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {ISP1Verifier} from "@sp1-contracts/ISP1Verifier.sol";

/// @dev Interface for the PoS Stake Manager contract with require methods to be used.
interface StakeManager {
    // Borrowed from the StakeManager contracts
    enum Status {Inactive, Active, Locked, Unstaked}
    function NFTCounter() external view returns (uint256);
    function validators(uint256) external view returns (uint256, uint256, uint256, uint256, uint256, address, address, Status, uint256, uint256, uint256, uint256, uint256);
    function validatorState() external view returns (uint256, uint256);
}

/// @title PoS Consensus Proof Verifier.
/// @author Manav Darji (manav2401)
/// @notice This contract verifies a consensus proof representing that a bor block has been 
///         voted upon by >2/3 of validators.
contract ConsensusProofVerifier {
    /// @notice The address of the SP1 verifier contract.
    /// @dev This can either be a specific SP1Verifier for a specific version, or the
    ///      SP1VerifierGateway which can be used to verify proofs for any version of SP1.
    ///      For the list of supported verifiers on each chain, see:
    ///      https://github.com/succinctlabs/sp1-contracts/tree/main/contracts/deployments
    address public verifier;

    /// @notice The verification key for the consensus proof program.
    bytes32 public consensusProofVKey;

    /// @notice The address of the PoS Stake Manager contract.
    address public posStakeManager;

    /// @notice The last verified bor block hash.
    bytes32 public lastVerifiedBorBlockHash;

    constructor(address _verifier, bytes32 _consensusProofVKey, address _posStakeManager) {
        verifier = _verifier;
        consensusProofVKey = _consensusProofVKey;
        posStakeManager = _posStakeManager;
    }

    /// @notice Fetches the active validator info like signer address, respective stake, and 
    ///         total stake and returns the encoded data.
    /// @return encodedValidatorInfo abi encoded data of all active validators
    function getEncodedValidatorInfo() public view returns(address[] memory, uint256[] memory, uint256) {
        // Get the total number of validators stored by fetching the NFT count. The count is
        // assigned to the next validator and hence we subtract 1 from it.
        uint256 length = StakeManager(posStakeManager).NFTCounter() - 1;

        address[] memory signers = new address[](length);
        uint256[] memory stakes = new uint256[](length);
        bool[] memory isActive = new bool[](length);
        uint256 totalActive = 0;

        // Validator index starts from 1.
        for (uint256 i = 1; i <= length; i++) {
            uint256 selfStake;
            uint256 delegatedStake;
            address signer;
            StakeManager.Status status;
            (selfStake, , , , , signer, , status, , , ,delegatedStake,) = StakeManager(posStakeManager).validators(i);
            signers[i-1] = signer;
            stakes[i-1] = selfStake + delegatedStake;
            isActive[i-1] = status == StakeManager.Status.Active;
            if (isActive[i-1]) {
                totalActive += 1;
            }
        }

        address[] memory activeSigners = new address[](totalActive);
        uint256[] memory activeStakes = new uint256[](totalActive);

        uint256 j = 0;
        for (uint256 i = 0; i < length; i++) {
            if (isActive[i]) {
                activeSigners[j] = signers[i];
                activeStakes[j] = stakes[i] / 1e18;
                j++;
            }
        }

        uint256 totalStake;
        (totalStake, ) = StakeManager(posStakeManager).validatorState();
        return (activeSigners, activeStakes, totalStake / 1e18);
    }

    /// @notice The entrypoint for the verifier.
    /// @param _proofBytes The encoded proof.
    /// @param bor_block_hash The bor block hash to be verified.
    /// @param l1_block_hash The l1 block hash for anchor.
    function verifyConsensusProof(
        bytes calldata _proofBytes, 
        bytes32 bor_block_hash, 
        bytes32 l1_block_hash
    ) public {
        bytes memory publicValues = abi.encodePacked(bor_block_hash, l1_block_hash);
        ISP1Verifier(verifier).verifyProof(consensusProofVKey, publicValues, _proofBytes);
        lastVerifiedBorBlockHash = bor_block_hash;
    }
}