use std::path::PathBuf;
use std::str::FromStr;

use alloy::hex;
use alloy_primitives::Address;
use alloy_provider::ReqwestProvider;
use alloy_rpc_types::BlockNumberOrTag;
use alloy_sol_macro::sol;
use sp1_cc_client_executor::ContractInput;
use sp1_cc_host_executor::HostExecutor;
use sp1_sdk::{utils, ProverClient, SP1ProofWithPublicValues, SP1Stdin};
use url::Url;
mod cli;
use clap::Parser;

use polccint_lib::bridge::BridgeInput;
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

/// The ELF we want to execute inside the zkVM.
const ELF: &[u8] = include_bytes!("../../../elf/bridge");

#[derive(Parser, Debug)]
struct Args {
    /// Whether or not to generate a proof.
    #[arg(long, default_value_t = false)]
    prove: bool,

    /// The chain ID. If not provided, requires the rpc_url argument to be provided.
    #[clap(long)]
    chain_id_l1: Option<u64>,

    /// The chain ID. If not provided, requires the rpc_url argument to be provided.
    #[clap(long)]
    chain_id_l2: Option<u64>,

    #[clap(long)]
    block_number_l1: u64,

    #[clap(long)]
    prev_block_number_l2: u64,

    #[clap(long)]
    new_block_number_l2: u64,

    /// The contract address for GER on L1.
    #[clap(long)]
    contract_ger_l1: String,

    /// The contract address for GER on L2.
    #[clap(long)]
    contract_ger_l2: String,

    /// The hex bytes for imported GERs.
    #[clap(long)]
    imported_gers_hex: String,
}

#[tokio::main]
async fn main() -> eyre::Result<()> {
    // Intialize the environment variables.
    dotenv::dotenv().ok();

    if std::env::var("RUST_LOG").is_err() {
        std::env::set_var("RUST_LOG", "info");
    }

    // Setup logging.
    utils::setup_logger();

    // Parse the command line arguments.
    let args = Args::parse();

    // Convert the contract addresses from strings to Address type
    let contract_ger_l1: Address =
        Address::from_str(&args.contract_ger_l1).expect("Invalid address");
    let contract_ger_l2: Address =
        Address::from_str(&args.contract_ger_l2).expect("Invalid address");

    let mut imported_gers: Vec<alloy_primitives::FixedBytes<32>> = [].to_vec();
    if args.imported_gers_hex != "" {
        let imported_gers_splitted: Vec<&str> = args.imported_gers_hex.split(',').collect();
        println!("{:?}", imported_gers_splitted);
        imported_gers = Vec::with_capacity(imported_gers_splitted.len());
        for ger_hex in imported_gers_splitted {
            imported_gers.push(
                alloy_primitives::FixedBytes::from_slice(&hex::decode(&ger_hex).unwrap()));
        }
    }
    println!("imported GERs: {:?}", imported_gers);

    // Load the input from the cache.
    // TODO return differnet providers
    let rpc_url_l1 = std::env::var(format!("RPC_{}", args.chain_id_l1.unwrap_or_default()))
        .expect("RPC URL must be defined")
        .parse::<Url>()
        .expect("Invalid URL format");

    let rpc_url_l2 = std::env::var(format!("RPC_{}", args.chain_id_l2.unwrap_or_default()))
        .expect("RPC URL must be defined")
        .parse::<Url>()
        .expect("Invalid URL format");

    let block_number_initial = BlockNumberOrTag::Number(args.prev_block_number_l2);
    let block_number_final = BlockNumberOrTag::Number(args.new_block_number_l2);

    // 1. Get the the last injecter GER of the previous block on L2

    // Setup the provider and host executor for initial GER
    let provider_l1 = ReqwestProvider::new_http(rpc_url_l1);
    let provider_l2 = ReqwestProvider::new_http(rpc_url_l2);

    let mut executor_injected_ger_count =
        HostExecutor::new(provider_l2.clone(), block_number_initial).await?;

    // Make the call to the slot0 function.
    println!("Calling injectedGERCount on L2 @ {}", contract_ger_l2);
    let injected_ger_count = executor_injected_ger_count
        .execute(ContractInput {
            contract_address: contract_ger_l2,
            caller_address: CALLER,
            calldata: GlobalExitRootManagerL2SovereignChain::injectedGERCountCall {},
        })
        .await?
        .injectedGerCount;

    // Now that we've executed all of the calls, get the `EVMStateSketch` from the host executor.
    println!("Getting injectedGERCount sketch");
    let executor_injected_ger_count_sketch = executor_injected_ger_count.finalize().await?;

    // 2. Check that the GERs are consecutive on L2 at the new block
    let mut executor_check_injected_gers_and_return_ler =
        HostExecutor::new(provider_l2.clone(), block_number_final).await?;

    println!("Checking injectedGERs on L2");
    // Make the call to the slot0 function.
    let check_injected_gers_and_return_ler_call_output_decoded =
        executor_check_injected_gers_and_return_ler
            .execute(ContractInput {
                contract_address: contract_ger_l2,
                caller_address: CALLER,
                calldata:
                    GlobalExitRootManagerL2SovereignChain::checkInjectedGERsAndReturnLERCall {
                        previousInjectedGERCount: injected_ger_count,
                        injectedGERs: imported_gers.clone(),
                    },
            })
            .await?;
    println!("Checking injectedGERs on L2 finished");

    // Check that the check was successful
    assert_eq!(
        check_injected_gers_and_return_ler_call_output_decoded.success,
        true
    );

    // Now that we've executed all of the calls, get the `EVMStateSketch` from the host executor.
    let executor_check_injected_gers_and_return_ler_sketch =
        executor_check_injected_gers_and_return_ler
            .finalize()
            .await?;

    // 3. Check that the GERs exist on L1
    let mut executor_check_gers_existance = HostExecutor::new(
        provider_l1.clone(),
        BlockNumberOrTag::from(args.block_number_l1),
    )
    .await?;

    // Make the call to the slot0 function.
    println!("Checking injectedGERs on L1");
    let check_injected_gers_existance_decoded = executor_check_gers_existance
        .execute(ContractInput {
            contract_address: contract_ger_l1,
            caller_address: CALLER,
            calldata: GlobalExitRootScrapper::checkGERsExistanceCall {
                globalExitRoots: imported_gers.clone(),
            },
        })
        .await?;
    println!("Checking injectedGERs on L1 finished");
    // Check that the check was successful
    assert_eq!(check_injected_gers_existance_decoded.success, true);

    // Now that we've executed all of the calls, get the `EVMStateSketch` from the host executor.
    let executor_check_injected_gers_existance = executor_check_gers_existance.finalize().await?;

    // Feed the sketch into the client.

    // Commit the bridge proof.
    let bridge_input: BridgeInput = BridgeInput {
        l1_ger_addr: contract_ger_l1,
        l2_ger_addr: contract_ger_l2,
        injected_gers: imported_gers,
        injected_ger_count_sketch: executor_injected_ger_count_sketch.clone(),
        check_injected_gers_and_return_ler_sketch:
            executor_check_injected_gers_and_return_ler_sketch.clone(),
        check_gers_existance_sketch: executor_check_injected_gers_existance.clone(),
    };

    //let input_bytes = bincode::serialize(&bridge_input)?;
    let mut stdin = SP1Stdin::new();
    stdin.write(&bridge_input);

    // Execute the program using the `ProverClient.execute` method, without generating a proof.
    let client = ProverClient::new();
    // Setup the proving key and verification key.
    let (pk, vk) = client.setup(ELF);

    let (_, report) = client.execute(ELF, stdin.clone()).run().unwrap();
    println!(
        "executed program with {} cycles",
        report.total_instruction_count()
    );

    if args.prove {
        println!("Starting proof generation.");
        let proof: SP1ProofWithPublicValues = client
            .prove(&pk, stdin.clone())
            .compressed()
            .run()
            .expect("Proving should work.");
        println!("Proof generation finished.");

        client
            .verify(&proof, &vk)
            .expect("proof verification should succeed");
        // Handle the result of the save operation
        let fixture_path = PathBuf::from(env!("CARGO_MANIFEST_DIR"))
            .join(format!("../proof/chain{}/", args.chain_id_l2.unwrap()));
        std::fs::create_dir_all(&fixture_path).expect("failed to create fixture path");

        match proof.save(PathBuf::from(env!("CARGO_MANIFEST_DIR")).join(format!(
            "../proof/chain{}/bridge_block_{}_to_{}_proof.bin",
            args.chain_id_l2.unwrap(),
            args.prev_block_number_l2,
            args.new_block_number_l2
        ))) {
            Ok(_) => println!("Proof saved successfully."),
            Err(e) => eprintln!("Failed to save proof: {}", e),
        }

        println!("Starting proof generation.");
    }
    Ok(())
}
