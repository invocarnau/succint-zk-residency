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
            block_hash: b256!("c89e39471b02783985e6c5a9304ce219fe24fe888f4b175414d8e156eb2c66ed"),
            state_root: b256!("a46cfc8f0a9eb0d222463cef0d85015c07ece30e747c23852023ac17551c3682"),
            message_passer_storage_root: b256!("8ed4baae3a927be3dea54996b4d5899f8c01e7594bf50b17dc1e741388ce3d12"),
        };
        let hash = pre_image.hash();
        assert_eq!(hash, b256!("f58f94a10efbd83c675189b3738c8862556bc868fb1130de4fca4188a636ab19"));
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
