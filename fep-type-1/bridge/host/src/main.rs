use std::path::PathBuf;

use alloy::hex;
use alloy_primitives::{address, Address};
use alloy_provider::ReqwestProvider;
use alloy_rpc_types::BlockNumberOrTag;
use alloy_sol_macro::sol;
use sp1_cc_client_executor::ContractInput;
use sp1_cc_host_executor::HostExecutor;
use sp1_sdk::{utils, ProverClient, SP1ProofWithPublicValues, SP1Stdin};
use url::Url;
mod cli;
use cli::ProviderArgs;
use clap::Parser;

use polccint_lib::{BridgeInput, SP1CCProofFixture};
use polccint_lib::constants::CALLER;

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

/// Agg Layer contract
const CONTRACT_GER_L1: Address = address!("0506B9383477F682DDB3701CD43eD30B9958099b");
const CONTRACT_GER_L2: Address = address!("0506B9383477F682DDB3701CD43eD30B9958099b");

/// The ELF we want to execute inside the zkVM.
const ELF: &[u8] = include_bytes!("../../../../elf/bridge");



/// Generate a `SP1CCProofFixture`, and save it as a json file.
///
/// This is useful for verifying the proof of contract call execution on chain.
fn save_fixture(vkey: String, proof: &SP1ProofWithPublicValues) {
    let fixture = SP1CCProofFixture {
        vkey,
        public_values: format!("0x{}", hex::encode(proof.public_values.as_slice())),
        proof: format!("0x{}", hex::encode(proof.bytes())),
    };

    let fixture_path = PathBuf::from(env!("CARGO_MANIFEST_DIR")).join("../contracts/src/fixtures");
    std::fs::create_dir_all(&fixture_path).expect("failed to create fixture path");
    std::fs::write(
        fixture_path.join("plonk-fixture.json".to_lowercase()),
        serde_json::to_string_pretty(&fixture).unwrap(),
    )
    .expect("failed to write fixture");
}

#[derive(Parser, Debug)]
struct Args {
    /// The block number of the block to execute.
    #[clap(flatten)]
    provider: ProviderArgs,

    /// Whether or not to generate a proof.
    #[arg(long, default_value_t = false)]
    prove: bool,
}

#[tokio::main]
async fn main() -> eyre::Result<()> {
    // Intialize the environment variables.
    dotenv::dotenv().ok();

    // Setup logging.
    utils::setup_logger();

    // Parse the command line arguments.
    let args = Args::parse();

    // load imported gers
    let imported_gers: Vec<alloy_primitives::FixedBytes<32>> = vec![
        alloy_primitives::FixedBytes::from_slice(
            &hex::decode("b00b00b000b00b00b00b00b00b00b000b00b00b00b00b000b00b00b00b00b00b").unwrap()
        )
    ];
    
    // Load the input from the cache.
    // TODO return differnet providers
    let provider_config = args.provider.into_provider().await?;
    let rpc_url: Url = provider_config.rpc_url.expect("URL must be defined");

    let block_number_initial = BlockNumberOrTag::Number(6797411);
    let block_number_final = BlockNumberOrTag::Number(6797428);

     // 1. Get the the last injecter GER of the previous block on L2

    // Setup the provider and host executor for initial GER
    let provider = ReqwestProvider::new_http(rpc_url);
    let mut executor_injected_ger_count = HostExecutor::new(provider.clone(), block_number_initial).await?;

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
        HostExecutor::new(provider.clone(), block_number_final).await?;

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
     HostExecutor::new(provider.clone(), block_number_final).await?;

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
    let mut stdin = SP1Stdin::new();
    stdin.write(&bridge_input);

    // Execute the program using the `ProverClient.execute` method, without generating a proof.
    let client = ProverClient::new();
    let (_, report) = client.execute(ELF, stdin.clone()).run().unwrap();
    println!("executed program with {} cycles", report.total_instruction_count());

    // // Generate the proof for the given program and input.
    // let (pk, vk) = client.setup(ELF);
    // let proof = client.prove(&pk, stdin).plonk().run().unwrap();
    // println!("generated proof");

    // // Save the proof, public values, and vkey to a json file.
    // save_fixture(vk.bytes32(), &proof);
    // println!("saved proof to plonk-fixture.json");

    // // Verify proof and public values.
    // client.verify(&proof, &vk).expect("verification failed");
    // println!("successfully generated and verified proof for the program!");
    Ok(())
}
