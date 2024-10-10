use alloy_primitives::{U256,B256, Address};
use sp1_cc_client_executor::io::EVMStateSketch;
use serde::{Deserialize, Serialize};
use tiny_keccak::{Keccak,Hasher};
use reth_primitives::Header;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct RootClaimPreImage {
    pub block_hash: B256,
    pub state_root: B256,
    pub message_passer_storage_root: B256,
}

impl RootClaimPreImage {
    fn marshal(&self) -> [u8; 128] {
        let mut buf = [0u8; 128];
        let version = B256::ZERO;
        
        buf[..32].copy_from_slice(version.as_ref());
        buf[32..64].copy_from_slice(self.state_root.as_ref());
        buf[64..96].copy_from_slice(self.message_passer_storage_root.as_ref());
        buf[96..128].copy_from_slice(self.block_hash.as_ref());
        
        buf
    }

    pub fn hash(&self) -> B256 {
        let marshaled = self.marshal();
        let mut output = [0u8; 32];
        let mut hasher = Keccak::v256();
        hasher.update(&marshaled);
        hasher.finalize(&mut output);
        output.into()
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use alloy_primitives::b256;

    #[test]
    fn test_root_claim_pre_image() {
        // Create a test instance
        let pre_image = RootClaimPreImage {
            block_hash: b256!("4f8f97b1f7fbc30e8315d320ae93942f8230c8a8a9c0543bfd6afbc60aa863c2"),
            state_root: b256!("0f0e5aed12699c19b5c82fd0247a07821f18c0aa39a739faf549d56e77210fae"),
            message_passer_storage_root: b256!("8ed4baae3a927be3dea54996b4d5899f8c01e7594bf50b17dc1e741388ce3d12"),
        };
        let hash = pre_image.hash();
        assert_eq!(hash, b256!("b598ce68aab8c040740fb183cce963ba651fb6b643972eb9b3831382e381935d"));
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct OpConsensusInput {
    pub get_root_claim_sketch: EVMStateSketch,
    pub get_game_from_factory_sketch: EVMStateSketch,
    pub game_index: U256,
    pub root_claim_pre_image: RootClaimPreImage,
    pub prev_l2_block_header: Header,
    pub new_l2_block_header: Header,
    pub game_factory_address: Address,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct OPConsensusCommit {
    pub game_factory_address: Address,
    pub l1_block_hash: B256,
    pub prev_l2_block_hash: B256,
    pub new_l2_block_hash: B256,
}
