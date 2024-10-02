use alloy_provider::ReqwestProvider;
use clap::Parser;
use reth_primitives::B256;
use rsp_client_executor::{io::ClientExecutorInput, ChainVariant, CHAIN_ID_ETH_MAINNET};
use rsp_host_executor::HostExecutor;
use sp1_sdk::{SP1Proof, HashableKey, utils, ProverClient, SP1Stdin, SP1ProofWithPublicValues, SP1VerifyingKey};
use std::path::PathBuf;
mod cli;
use cli::ProviderArgs;
use url::Url;
use polccint_lib::{BridgeCommit, BlockCommit, BlockAggregationInput, BlockAggregationCommit, BridgeInput, FinalAggregationInput};
use alloy_rpc_types::BlockNumberOrTag;
use alloy_primitives::{address, Address};
use alloy::hex;
use sp1_cc_host_executor::HostExecutor  as StaticCallHostExecutor;
use sp1_cc_client_executor::{ContractInput};
use polccint_lib::constants::CALLER;
use alloy_sol_macro::sol;


#[derive(Parser, Debug)]
struct Args {
    /// The block number of the block to execute.
    #[clap(long)]
    block_number: u64,
    #[clap(flatten)]
    provider: ProviderArgs,

    /// Whether or not to generate a proof.
    #[arg(long, default_value_t = false)]
    prove: bool,
}


const ELF_BLOCK: &[u8] = include_bytes!("../../../../elf/block");
const ELF_BLOCK_AGGREGATION: &[u8] = include_bytes!("../../../../elf/block-aggregation");
const ELF_BRIDGE: &[u8] = include_bytes!("../../../../elf/bridge");
const ELF_FINAL_AGGREGATION: &[u8] = include_bytes!("../../../../elf/aggregation-final");

/// Agg Layer contract
const CONTRACT_GER_L1: Address = address!("0506B9383477F682DDB3701CD43eD30B9958099b");
const CONTRACT_GER_L2: Address = address!("0506B9383477F682DDB3701CD43eD30B9958099b");

// try what happens if the calls revert?¿
sol! (
    interface GlobalExitRootManagerL2SovereignChain {
        function injectedGERCount() public view returns (uint256 injectedGerCount);
        function checkInjectedGERsAndReturnLER(uint256 previousInjectedGERCount, bytes32[] injectedGERs) public view returns (bool success, bytes32 localExitRoot);
    }

    interface GlobalExitRootScrapper {
        function checkGERsExistance(bytes32[] calldata globalExitRoots) public view returns (bool success);
    }
);


/// An input to the aggregation program.
///
/// Consists of a proof and a verification key.
struct AggregationInput {
    pub proof: SP1ProofWithPublicValues,
    pub vk: SP1VerifyingKey,
}

#[tokio::main]
async fn main() -> eyre::Result<()> {
    // Intialize the environment variables.
    dotenv::dotenv().ok();

    // // Fallback to 'info' level if RUST_LOG is not set
    // if std::env::var("RUST_LOG").is_err() {
    //     std::env::set_var("RUST_LOG", "info");
    // }

    // Initialize the logger.
    utils::setup_logger();

    // Parse the command line arguments.
    let args = Args::parse();

    // Load the input from the cache.
    let provider_config = args.provider.into_provider().await?;

    // Setup the provider.
    // let default_url = url::Url::parse("https://example.com").unwrap();
    // let url: url::Url = some_option.unwrap_or(default_url);
    let rpc_url: Url = provider_config.rpc_url.expect("URL must be defined");

    let provider = ReqwestProvider::new_http(rpc_url);

    // Setup the host executor.
    let host_executor = HostExecutor::new(provider.clone());
    let variant = match provider_config.chain_id {
        CHAIN_ID_ETH_MAINNET => ChainVariant::Ethereum,
        _ => {
            eyre::bail!("unknown chain ID: {}", provider_config.chain_id);
        }
    };

    // Generate the proof.
    let client = ProverClient::new();

    // Setup the proving and verifying keys.
    let (aggregation_pk,aggregation_vk) = client.setup(ELF_BLOCK_AGGREGATION);
    let (block_pk, block_vk) = client.setup(ELF_BLOCK);
    let (bridge_pk, bridge_vk) = client.setup(ELF_BRIDGE);
    let (final_aggregation_pk, final_aggregation_vk) = client.setup(ELF_FINAL_AGGREGATION);

    // cargo run --release -- --chain-id 1 --block-number 18884864
    let initial_block_number = 6797425;//args.block_number;
    let block_range = 2; // hardcode for now TODO
    let final_block_number = initial_block_number + block_range;

    let mut inputs: Vec<AggregationInput> = Vec::new();

    for block_number in initial_block_number..final_block_number {
        let client_input = host_executor
            .execute(block_number, variant)
            .await
            .expect("failed to execute host");

        let mut stdin_block = SP1Stdin::new();
        let buffer = bincode::serialize(&client_input).unwrap();
        stdin_block.write_vec(buffer);
        stdin_block.write(&client_input);
        let proof = client
            .prove(&block_pk.clone(), stdin_block)
            .compressed()
            .run()
            .expect("proving failed");
        
        inputs.push(
            AggregationInput {
                proof: proof,
                vk: block_vk.clone(),
            }
        );
    }

    // encode aggregation input and write to stdin
    let mut stdin_aggregation = SP1Stdin::new();
    let aggregation_input = BlockAggregationInput{
        block_commits: inputs
        .iter()
        .map(|input| input.proof.public_values.clone().read::<BlockCommit>())
        .collect::<Vec<_>>(),
        block_vkey: block_vk.clone().hash_u32(), // probabyl worth to change interface TODO
    };
    stdin_aggregation.write(&aggregation_input);

    // write proofs
    for input in inputs {
        let SP1Proof::Compressed(proof) = input.proof.proof else {
            panic!()
        };
        stdin_aggregation.write_proof(proof, input.vk.vk);
    }
    
    let proof_aggregation = client
    .prove(&aggregation_pk.clone(), stdin_aggregation.clone())
    .compressed()
    .run()
    .expect("proving failed");

    // now compute the bridge proof! ^^

     // load imported gers
     let imported_gers: Vec<alloy_primitives::FixedBytes<32>> = vec![
        alloy_primitives::FixedBytes::from_slice(
            &hex::decode("b00b00b000b00b00b00b00b00b00b000b00b00b00b00b000b00b00b00b00b00b").unwrap()
        )
    ];
    
    // 1. Get the the last injecter GER of the previous block on L2

    // Setup the provider and host executor for initial GER
    let mut executor_injected_ger_count = 
        StaticCallHostExecutor::new(
            provider.clone(), 
            BlockNumberOrTag::Number(initial_block_number)
        ).await?;

    // Make the call to the slot0 function.
    let injected_ger_count = executor_injected_ger_count
        .execute(ContractInput {
            contract_address: CONTRACT_GER_L2,
            caller_address: CALLER,
            calldata: GlobalExitRootManagerL2SovereignChain::injectedGERCountCall {},
        })
        .await?
        .injectedGerCount;

    // Now that we've executed all of the calls, get the `EVMStateSketch` from the host executor.
    let executor_injected_ger_count_sketch = executor_injected_ger_count.finalize().await?;

    // 2. Check that the GERs are consecutive on L2 at the new block
    let mut executor_check_injected_gers_and_return_ler = 
    StaticCallHostExecutor::new(provider.clone(),  BlockNumberOrTag::Number(final_block_number)).await?;

    // Make the call to the slot0 function.
    let check_injected_gers_and_return_ler_call_output_decoded = executor_check_injected_gers_and_return_ler
        .execute(ContractInput {
            contract_address: CONTRACT_GER_L2,
            caller_address: CALLER,
            calldata: GlobalExitRootManagerL2SovereignChain::checkInjectedGERsAndReturnLERCall { 
                previousInjectedGERCount: injected_ger_count, 
                injectedGERs: imported_gers.clone()
            },
        })
        .await?;

    // Check that the check was successful
    assert_eq!(check_injected_gers_and_return_ler_call_output_decoded.success, true);

    // Now that we've executed all of the calls, get the `EVMStateSketch` from the host executor.
    let executor_check_injected_gers_and_return_ler_sketch = executor_check_injected_gers_and_return_ler.finalize().await?;

    // 3. Check that the GERs exist on L1
    let mut executor_check_gers_existance =
    StaticCallHostExecutor::new(provider.clone(),  BlockNumberOrTag::Number(final_block_number)).await?;

    // Make the call to the slot0 function.
    let check_injected_gers_existance_decoded = executor_check_gers_existance
        .execute(ContractInput {
            contract_address: CONTRACT_GER_L1,
            caller_address: CALLER,
            calldata: GlobalExitRootScrapper::checkGERsExistanceCall { globalExitRoots: imported_gers.clone() },
        })
        .await?;

    // Check that the check was successful
    assert_eq!(check_injected_gers_existance_decoded.success, true);

    // Now that we've executed all of the calls, get the `EVMStateSketch` from the host executor.
    let executor_check_injected_gers_existance = executor_check_injected_gers_and_return_ler.finalize().await?;

    // Feed the sketch into the client.

     // Commit the bridge proof.
     let bridge_input: BridgeInput = BridgeInput {
        l1_ger_addr: CONTRACT_GER_L1,
        l2_ger_addr: CONTRACT_GER_L2,
        injected_gers: imported_gers,
        injected_ger_count_sketch: executor_injected_ger_count_sketch.clone(),
        check_injected_gers_and_return_ler_sketch: executor_check_injected_gers_and_return_ler_sketch.clone(),
        check_gers_existance_sketch: executor_check_injected_gers_existance.clone(),
    };   


    //let input_bytes = bincode::serialize(&bridge_input)?;
    let mut stdin_bridge = SP1Stdin::new();
    stdin_bridge.write(&bridge_input);

    let proof_bridge = client
    .prove(&bridge_pk.clone(), stdin_bridge.clone())
    .compressed()
    .run()
    .expect("proving failed");

  


      // encode aggregation input and write to stdin
    let mut stdin_final_aggregation = SP1Stdin::new();
    let final_aggregation_input: FinalAggregationInput = FinalAggregationInput {
        block_vkey_aggregation: aggregation_vk.clone().hash_u32(),
        block_aggregation_commit: proof_aggregation.public_values.clone().read::<BlockAggregationCommit>(),
        block_vkey_bridge: bridge_vk.clone().hash_u32(),
        bridge_commit:proof_bridge.public_values.clone().read::<BridgeCommit>()
    };
    stdin_final_aggregation.write(&final_aggregation_input);

    // write proofs
    let SP1Proof::Compressed(proof) = proof_aggregation.proof else {
        panic!()
    };
    stdin_final_aggregation.write_proof(proof, aggregation_vk.vk);

    let SP1Proof::Compressed(proof) = proof_bridge.proof else {
        panic!()
    };
    stdin_final_aggregation.write_proof(proof, bridge_vk.vk);

    
    let proof_final_aggregation = client
    .prove(&final_aggregation_pk.clone(), stdin_final_aggregation.clone())
    .compressed()
    .run()
    .expect("proving failed");
    Ok(())
}


/// Generate a `SP1CCProofFixture`, and save it as a json file.
///
/// This is useful for verifying the proof of contract call execution on chain.
fn save_fixture(vkey: String, proof: &SP1ProofWithPublicValues) {
    // let fixture = SP1CCProofFixture {
    //     vkey,
    //     public_values: format!("0x{}", hex::encode(proof.public_values.as_slice())),
    //     proof: format!("0x{}", hex::encode(proof.bytes())),
    // };

    // let fixture_path = PathBuf::from(env!("CARGO_MANIFEST_DIR")).join("../contracts/src/fixtures");
    // std::fs::create_dir_all(&fixture_path).expect("failed to create fixture path");
    // std::fs::write(
    //     fixture_path.join("plonk-fixture.json".to_lowercase()),
    //     serde_json::to_string_pretty(&fixture).unwrap(),
    // )
    // .expect("failed to write fixture");
}
