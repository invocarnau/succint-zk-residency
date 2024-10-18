use crate::types::*;

use base64::{prelude::BASE64_STANDARD, Engine};
use core::str;
use sha2::{Digest, Sha256};

use alloy_primitives::{Address, FixedBytes};
use reth_primitives::recover_signer_unchecked;

// Verifies if the signature is indeed signed by the expected signer or not
pub fn verify_signature(signature: &str, message_hash: &[u8; 32], expected_signer: Address) {
    // Decode the tendermint signature using standard base64 decoding
    let decoded_signature = BASE64_STANDARD
        .decode(signature)
        .expect("unable to decode signature");

    // Construct the byte array from the decoded signature for recovery
    let mut sig = [0u8; 65];
    sig.copy_from_slice(decoded_signature.as_slice());

    // Recover the signer address
    let recovered_signer = recover_signer_unchecked(&sig, message_hash).unwrap_or_default();
    let recovered_signer_alloy = Address::from_slice(recovered_signer.as_slice());

    assert_eq!(
        expected_signer, recovered_signer_alloy,
        "recovered and expected signature mismatch"
    );
}

// Verifies if the transaction data matches with the given transaction hash or not. It also
// extracts the milestone message from the transaction data and returns it.
pub fn verify_tx_data(
    tx_data: &str,
    expected_hash: &FixedBytes<32>,
) -> heimdall_types::MilestoneMsg {
    // Decode the transaction data
    let mut decoded_tx_data = BASE64_STANDARD
        .decode(tx_data)
        .expect("tx_data decoding failed");

    // Calculate the hash of decoded data
    let tx_hash = sha256(decoded_tx_data.as_slice());

    assert_eq!(*expected_hash, tx_hash);

    // Deserialize the message to extract the milestone bytes
    let decoded_message =
        deserialize_msg(&mut decoded_tx_data).expect("tx_data deserialization failed");

    decoded_message.msg.unwrap()
}

// Verifies if the precommit message includes the milestone side transaction or not by deserialising
// the encoded precommit message. It also checks if the validator voted yes on transaction or not.
pub fn verify_precommit(precommit_message: &mut Vec<u8>, expected_hash: &FixedBytes<32>) {
    // Decode the precommit message
    let precommit =
        deserialize_precommit(precommit_message).expect("precommit deserialization failed");
    let side_tx = precommit.side_tx_results;

    // If the validator didn't vote on the side transaction, the object will be empty
    let side_tx = side_tx.expect("side_tx is empty");

    assert_eq!(
        expected_hash.to_vec(),
        side_tx.tx_hash,
        "tx_hash in precommit doesn't match with milestone tx_hash"
    );
    assert_eq!(side_tx.result, 1, "no yes vote on the side tx");
}

fn sha256(decoded_tx_data: &[u8]) -> FixedBytes<32> {
    // Create a new Sha256 instance
    let mut hasher = Sha256::new();

    // Write the tx data
    hasher.update(decoded_tx_data);

    // Read hash digest and consume hasher
    let result = hasher.finalize();

    FixedBytes::from_slice(result.as_slice())
}
