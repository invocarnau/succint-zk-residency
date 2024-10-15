use prost::Message;
use std::{io::Cursor, ops::Sub};

// Include the `types` module, which is generated from types.proto.
pub mod heimdall_types {
    include!(concat!(env!("OUT_DIR"), "/types.rs"));
}

// Serialize the wrapped milestone message into a byte buffer.
pub fn serialize_precommit(m: &heimdall_types::Vote) -> Vec<u8> {
    let mut buf = Vec::with_capacity(m.encoded_len());

    // Unwrap is safe, since we have reserved sufficient capacity in the vector.
    m.encode_length_delimited(&mut buf).unwrap();
    buf
}

// Deserialize the wrapped milestone message fromt the given buffer. It does byte manipulation
// to handle the decoding of message generated from the go code.
pub fn deserialize_precommit(
    buf: &mut Vec<u8>,
) -> Result<heimdall_types::Vote, prost::DecodeError> {
    heimdall_types::Vote::decode_length_delimited(&mut Cursor::new(buf))
}

// Serialize the wrapped milestone message into a byte buffer.
pub fn serialize_msg(m: &heimdall_types::StdTx) -> Vec<u8> {
    let mut buf = Vec::with_capacity(m.encoded_len());

    // Unwrap is safe, since we have reserved sufficient capacity in the vector.
    m.encode_length_delimited(&mut buf).unwrap();
    buf
}

// Deserialize the wrapped milestone message fromt the given buffer. It does byte manipulation
// to handle the decoding of message generated from the go code.
pub fn deserialize_msg(buf: &mut Vec<u8>) -> Result<heimdall_types::StdTx, prost::DecodeError> {
    // This is a hack to handle decoding of message generated from the go code. Old prefix
    // represents the encoded info for the cosmos message interface. Because it's not possible
    // to represent that info in the proto file, we need to replace the prefix with simple bytes
    // which can be decoded into the milestone message generated in rust.
    let old_prefix: Vec<u8> = vec![1, 240, 98, 93, 238, 10, 158, 1, 210, 203, 62, 102];
    let mut new_prefix: Vec<u8> = vec![1, 10, 154, 1];

    let matches = buf.len() > old_prefix.len()
        && old_prefix[..].iter().enumerate().all(|(i, &b)| {
            if i == 6 {
                new_prefix[2] = buf[i + 1].sub(4);
                true
            } else {
                b == buf[i + 1]
            }
        });

    if matches {
        buf.drain(1..1 + old_prefix.len());
        buf.splice(1..1, new_prefix.iter().cloned());
        buf[0] = buf[0].sub(8);
    } else {
        return Err(prost::DecodeError::new("Invalid prefix"));
    }

    heimdall_types::StdTx::decode_length_delimited(&mut Cursor::new(buf))
}

pub fn deserialize_validator_set(
    buf: &mut Vec<u8>,
) -> Result<heimdall_types::ValidatorSet, prost::DecodeError> {
    heimdall_types::ValidatorSet::decode(&mut Cursor::new(buf))
}

pub fn serialize_validator_set(m: &heimdall_types::ValidatorSet) -> Vec<u8> {
    let mut buf = Vec::with_capacity(m.encoded_len());

    // Unwrap is safe, since we have reserved sufficient capacity in the vector.
    m.encode_length_delimited(&mut buf).unwrap();
    buf
}

#[cfg(test)]
mod tests {
    use super::*;
    // use alloy_primitives::hex;
    use base64::{prelude::BASE64_STANDARD, Engine};
    use prost_types::Timestamp;
    use reth_primitives::hex;
    use std::str::FromStr;

    #[test]
    fn test_deserialize_msg() {
        let a = "6gHwYl3uCqAB0ss+ZgoUCSB6bv7jRss+SlSsGFI+NxXTiz8Q/KaMBhj6p4wGIiBhk6zRTSThGAsyswmIseJyY9Eg8rrnHi4vXGNFGJ/r9SoFODAwMDIyUTUwMjM0MTM1LWQ5YmUtNGU0YS04NGY3LTM1OTZjZmIwN2EwZCAtIDB4ODhiMWUyNzI2M2QxMjBmMmJhZTcxZTJlMmY1YzYzNDUxODlmZWJmNRJB5jp3Zv4MQiiaOQ612UlPgyJzjt3v5YAJs9sqArSSsXVnssdRf5as1uuwettRNPGFlohE8saPapGLQxF74mHm/AE=".to_string();
        let mut decoded_tx_data = BASE64_STANDARD.decode(a).expect("tx_data decoding failed");
        let decoded_message = deserialize_msg(&mut decoded_tx_data).unwrap();

        let m = heimdall_types::MilestoneMsg {
            proposer: hex::decode("09207a6efee346cb3e4a54ac18523e3715d38b3f")
                .unwrap()
                .to_vec(),
            start_block: 12784508,
            end_block: 12784634,
            hash: hex::decode("6193acd14d24e1180b32b30988b1e27263d120f2bae71e2e2f5c6345189febf5")
                .unwrap()
                .to_vec(),
            bor_chain_id: "80002".to_string(),
            milestone_id:
                "50234135-d9be-4e4a-84f7-3596cfb07a0d - 0x88b1e27263d120f2bae71e2e2f5c6345189febf5"
                    .to_string(),
        };
        let sig = hex::decode("0xe63a7766fe0c42289a390eb5d9494f8322738eddefe58009b3db2a02b492b17567b2c7517f96acd6ebb07adb5134f185968844f2c68f6a918b43117be261e6fc01").unwrap().to_vec();
        let msg = heimdall_types::StdTx {
            msg: Some(m),
            signature: sig,
            memo: "".to_string(),
        };
        assert_eq!(decoded_message, msg);
    }

    #[test]
    fn test_precommit_msg() {
        let hex_msg = "9701080211327b30010000000022480a20fd648de965c020911f2bcfa3825fe2bd6698aa93009f0e63348ad74506221fae12240a20218d85717b5904942ce7c7b89b201aa1c2711dddb6e380cd0357c4647f35ac9b10012a0c08f29ee6b50610abfdeacc03320c6865696d64616c6c2d31333742240a204c6bb9c1426cef3b0252efadfbd09b88350f508cc2a4ec0c837612958ad37c851001";
        let mut bytes_msg = hex::decode(hex_msg).unwrap();
        let decoded = deserialize_precommit(&mut bytes_msg).unwrap();

        let timestamp = Timestamp::from_str("2024-08-12T04:28:34.966442667Z").unwrap();
        let parts_header = heimdall_types::CanonicalPartSetHeader {
            total: 1,
            hash: hex::decode("218D85717B5904942CE7C7B89B201AA1C2711DDDB6E380CD0357C4647F35AC9B")
                .unwrap(),
        };
        let block_id = Some(heimdall_types::CanonicalBlockId {
            hash: hex::decode("FD648DE965C020911F2BCFA3825FE2BD6698AA93009F0E63348AD74506221FAE")
                .unwrap(),
            parts_header: Some(parts_header),
        });
        let side_tx = heimdall_types::SideTxResult {
            tx_hash: hex::decode(
                "4c6bb9c1426cef3b0252efadfbd09b88350f508cc2a4ec0c837612958ad37c85",
            )
            .unwrap(),
            result: 1,
            sig: [].to_vec(),
        };
        let vote = heimdall_types::Vote {
            r#type: 2,
            height: 19954482,
            round: 0,
            block_id,
            timestamp: Some(timestamp),
            chain_id: "heimdall-137".to_string(),
            data: [].to_vec(),
            side_tx_results: Some(side_tx),
        };

        assert_eq!(decoded, vote);
    }

    #[test]
    fn test_validator_set() {
        let hex_msg = "0a600801200128904e324104dc19fdf9a82fd5c4327f31b96b6bbe0b9d44564ad89c2139db47c5cb2def87ac584fc05117663de2f17ae5ee50eced7283a596e10aaf33fb34c4cf5f98e4fda73a146ab3d36c46ecfb9b9c0bd51cb1c3da5a2c81cea612600801200128904e324104dc19fdf9a82fd5c4327f31b96b6bbe0b9d44564ad89c2139db47c5cb2def87ac584fc05117663de2f17ae5ee50eced7283a596e10aaf33fb34c4cf5f98e4fda73a146ab3d36c46ecfb9b9c0bd51cb1c3da5a2c81cea6";
        let mut bytes_msg = hex::decode(hex_msg).unwrap();

        let decoded = deserialize_validator_set(&mut bytes_msg);

        let validator_set = heimdall_types::ValidatorSet {
                    validators: vec![
                        heimdall_types::Validator {
                            id: 1,
                            start_epoch: 0,
                            end_epoch: 0,
                            nonce: 1,
                            voting_power: 10000,
                            pub_key: hex::decode("0x04dc19fdf9a82fd5c4327f31b96b6bbe0b9d44564ad89c2139db47c5cb2def87ac584fc05117663de2f17ae5ee50eced7283a596e10aaf33fb34c4cf5f98e4fda7").unwrap().to_vec(),
                            signer: hex::decode("0x6ab3d36c46ecfb9b9c0bd51cb1c3da5a2c81cea6").unwrap().to_vec(),
                            last_updated: "".to_string(),
                            jailed: false,
                            proposer_priority: 0,
                        },
                    ],
                    proposer: Some(heimdall_types::Validator {
                        id: 1,
                        start_epoch: 0,
                        end_epoch: 0,
                        nonce: 1,
                        voting_power: 10000,
                        pub_key: hex::decode("0x04dc19fdf9a82fd5c4327f31b96b6bbe0b9d44564ad89c2139db47c5cb2def87ac584fc05117663de2f17ae5ee50eced7283a596e10aaf33fb34c4cf5f98e4fda7").unwrap().to_vec(),
                        signer: hex::decode("0x6ab3d36c46ecfb9b9c0bd51cb1c3da5a2c81cea6").unwrap().to_vec(),
                        last_updated: "".to_string(),
                        jailed: false,
                        proposer_priority: 0,
            }),
        };

        assert_eq!(decoded.unwrap(), validator_set);
    }
}
