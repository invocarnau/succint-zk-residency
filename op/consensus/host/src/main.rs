use alloy_sol_macro::sol;
use clap::Parser;
use alloy_primitives::{B256, U256, Address};

use std::path::PathBuf;
use std::str::FromStr;

use alloy_provider::ReqwestProvider;
use alloy_rpc_types::BlockNumberOrTag;
use sp1_cc_client_executor::ContractInput;
use sp1_cc_host_executor::HostExecutor;
use sp1_sdk::{utils, ProverClient, SP1ProofWithPublicValues, SP1Stdin};
use url::Url;
use polccint_lib::{op::OpConsensusInput, op::RootClaimPreImage};
use polccint_lib::constants::CALLER;

sol! (
    interface GameFactory {
        function gameAtIndex(uint256 _index) view returns(uint32 gameType_, uint64 timestamp_, address game);
    }

    interface Game {
        function rootClaim() pure returns(bytes32 rootClaim);
    }
);

#[derive(Parser, Debug)]
struct Args {
    /// Whether or not to generate a proof.
    #[arg(long, default_value_t = false)]
    prove: bool,

    #[clap(long)]
    chain_id_l1: Option<u64>,

    #[clap(long)]
    chain_id_l2: Option<u64>,

    #[clap(long)]
    block_number_l1: u64,

    #[clap(long)]
    game_index: u64,

    #[clap(long)]
    game_factory_address: String,

    #[clap(long)]
    prev_block_number_l2: u64,

    #[clap(long)]
    new_block_number_l2: u64,

    #[clap(long)]
    claim_block_hash: String,

    #[clap(long)]
    claim_state_root: String,

    #[clap(long)]
    claim_message_passer_storage_root: String,
}

/// The ELF we want to execute inside the zkVM.
const ELF: &[u8] = include_bytes!("../../../../elf/op-consensus");

#[tokio::main]
async fn main() -> eyre::Result<()> {
    // Intialize the environment variables.
    dotenv::dotenv().ok();
    if std::env::var("RUST_LOG").is_err() {
        std::env::set_var("RUST_LOG", "info");
    }
    // Setup logging.
    utils::setup_logger();
    // Parse the command line arguments, convert types.
    let args = Args::parse();
    let game_factory_address: Address = Address::from_str(&args.game_factory_address).expect("Invalid address");
    let claim_block_hash: B256 = B256::from_str(&args.claim_block_hash).expect("Invalid block_hash");
    let claim_state_root: B256 = B256::from_str(&args.claim_state_root).expect("Invalid claim_state_root");
    let claim_message_passer_storage_root: B256 = B256::from_str(&args.claim_message_passer_storage_root).expect("Invalid claim_message_passer_storage_root");
    let root_claim_pre_image = RootClaimPreImage{
        block_hash: claim_block_hash,
        state_root: claim_state_root,
        message_passer_storage_root: claim_message_passer_storage_root,
    };
    let rpc_url_l1 = std::env::var(format!("RPC_{}", args.chain_id_l1.unwrap_or_default()))
    .expect("RPC URL must be defined")
    .parse::<Url>()
    .expect("Invalid URL format");
    let rpc_url_l2 = std::env::var(format!("RPC_{}", args.chain_id_l2.unwrap_or_default()))
    .expect("RPC URL must be defined")
    .parse::<Url>()
    .expect("Invalid URL format");

    // Setup the provider and host executor
    let provider_l1 = ReqwestProvider::new_http(rpc_url_l1);
    let provider_l2 = ReqwestProvider::new_http(rpc_url_l2);


    // Get trace for get_game_from_factory_sketch call
    let mut executor_get_game_from_factory = HostExecutor::new(
        provider_l1.clone(), 
        BlockNumberOrTag::Number(args.block_number_l1)
    ).await?;
    println!("Calling gameAtIndex on L1");
    let game_index = U256::from(args.game_index);
    let game_addr = executor_get_game_from_factory
        .execute(ContractInput {
            contract_address: game_factory_address,
            caller_address:  CALLER,
            calldata: GameFactory::gameAtIndexCall { _index: game_index },
        })
        .await?
        .game;
    let get_game_from_factory_sketch = executor_get_game_from_factory.finalize().await?;

    // Get trace for get_root_claim_sketch call
    let mut executor_get_root_claim = HostExecutor::new(
        provider_l1.clone(), 
        BlockNumberOrTag::Number(args.block_number_l1)
    ).await?;
    println!("Calling rootClaim on L1");
    executor_get_root_claim
        .execute(ContractInput {
            contract_address: game_addr,
            caller_address:  CALLER,
            calldata: Game::rootClaimCall { },
        })
        .await?;
    let get_root_claim_sketch = executor_get_root_claim.finalize().await?;

    // Get prev l2 block header
    let executor_prev_l2_block_header = HostExecutor::new(
        provider_l2.clone(), 
        BlockNumberOrTag::Number(args.prev_block_number_l2)
    ).await?;
    println!("Calling rootClaim on L1");
    let prev_l2_block_header = executor_prev_l2_block_header.header;

    // Get new l2 block header
    let executor_new_l2_block_header = HostExecutor::new(
        provider_l2.clone(), 
        BlockNumberOrTag::Number(args.new_block_number_l2)
    ).await?;
    println!("Calling rootClaim on L1");
    let new_l2_block_header = executor_new_l2_block_header.header;


    // prove
    let op_consensus_input = OpConsensusInput{
        get_root_claim_sketch,
        get_game_from_factory_sketch,
        game_index,
        root_claim_pre_image,
        prev_l2_block_header,
        new_l2_block_header,
        game_factory_address,
    };

    //let input_bytes = bincode::serialize(&bridge_input)?;
    let mut stdin = SP1Stdin::new();
    stdin.write(&op_consensus_input);

    // Execute the program using the `ProverClient.execute` method, without generating a proof.
    let client = ProverClient::new();
    // Setup the proving key and verification key.
    let (pk, vk) = client.setup(ELF);

    let (_, report) = client.execute(ELF, stdin.clone()).run().unwrap();
    println!("executed program with {} cycles", report.total_instruction_count());

    if args.prove {
        println!("Starting proof generation.");
        let proof: SP1ProofWithPublicValues = client.prove(&pk, stdin.clone()).compressed().run().expect("Proving should work.");
        println!("Proof generation finished.");

        client.verify(&proof, &vk).expect("proof verification should succeed");
        // Handle the result of the save operation
        match proof.save(PathBuf::from(env!("CARGO_MANIFEST_DIR")).join(format!("../../proof/chain{}/op_consensus_game_{}.bin", args.chain_id_l2.unwrap(), args.game_index))) {
            Ok(_) => println!("Proof saved successfully."),
            Err(e) => eprintln!("Failed to save proof: {}", e),
        }

        println!("Starting proof generation.");
    }

    Ok(())
}
