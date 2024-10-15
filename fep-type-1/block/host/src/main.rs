use alloy_provider::ReqwestProvider;
use clap::Parser;
use rsp_client_executor::{ChainVariant, CHAIN_ID_ETH_MAINNET};
use rsp_host_executor::HostExecutor;
use sp1_sdk::{utils, ProverClient, SP1Stdin, SP1ProofWithPublicValues};
use std::path::PathBuf;
mod cli;
use cli::ProviderArgs;
use url::Url;
use polccint_lib::fep_type_1::BlockCommit;

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


const ELF: &[u8] = include_bytes!("../../../../elf/block");


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
    let host_executor = HostExecutor::new(provider);
    let variant = match provider_config.chain_id {
        CHAIN_ID_ETH_MAINNET => ChainVariant::Ethereum,
        _ => ChainVariant::CliqueShanghaiChainID,
    };

    // Construct the path for saving the proof
    let proof_path = PathBuf::from(env!("CARGO_MANIFEST_DIR"))
        .join(format!("../../proof/chain{}/block_{}_proof.bin", provider_config.chain_id, args.block_number));

    // Create all necessary directories
    if let Some(parent) = proof_path.parent() {
        std::fs::create_dir_all(parent)?;
    }

    // println!("ChainID: {:?}", provider_config.chain_id);
    // println!("Executing block number: {:?}", args.block_number);

    let client_input = host_executor
        .execute(args.block_number, variant)
        .await
        .expect("failed to execute host");

    println!("{:?}", client_input.current_block.parent_hash);
    println!("{:?}", client_input.current_block.hash_slow());

    //let client_input = load_input_from_cache(CHAIN_ID_ETH_MAINNET, 20526624);

    // Generate the proof.
    let client = ProverClient::new();

    // Setup the proving key and verification key.
    let (pk, vk) = client.setup(ELF);

    // Write the block to the program's stdin.
    let mut stdin = SP1Stdin::new();
    let buffer = bincode::serialize(&client_input).unwrap();

    // write cleintInput and chainId to stdin
    stdin.write_vec(buffer);
    stdin.write(&provider_config.chain_id);


    // Only execute the program.
    let (mut public_values, execution_report) =
        client.execute(&pk.elf, stdin.clone()).run().unwrap();
    println!(
        "Finished executing the block in {} cycles",
        execution_report.total_instruction_count()
    );

    let decoded_public_values = public_values.read::<BlockCommit>();

    // Assert outputs
    // Check that the check was successful
    assert_eq!(decoded_public_values.prev_block_hash, client_input.parent_header().hash_slow());
    assert_eq!(decoded_public_values.new_block_hash, client_input.current_block.header.hash_slow());

    // If the `prove` argument was passed in, actually generate the proof.
    // It is strongly recommended you use the network prover given the size of these programs.
    if args.prove {
        println!("Starting proof generation.");
        let proof: SP1ProofWithPublicValues = client.prove(&pk, stdin.clone()).compressed().run().expect("Proving should work.");
        println!("Proof generation finished.");

        client.verify(&proof, &vk).expect("proof verification should succeed");
        // Handle the result of the save operation
        match proof.save(PathBuf::from(env!("CARGO_MANIFEST_DIR")).join(format!("../../proof/chain{}/block_{}_proof.bin", provider_config.chain_id, args.block_number))) {
            Ok(_) => println!("Proof saved successfully."),
            Err(e) => eprintln!("Failed to save proof: {}", e),
        }

        println!("Starting proof generation.");
        // let proof: SP1ProofWithPublicValues = client.prove(&pk, stdin.clone()).compressed().run().expect("Proving should work.");
        // println!("Proof generation finished.");
        // println!("{:?}", proof);

        // client.verify(&proof, &vk).expect("proof verification should succeed");

        // save_fixture(vk.clone().bytes32(), &proof, args.block_number, provider_config.chain_id);
        // let fixture = read_fixture(args.block_number, provider_config.chain_id);

    }
    Ok(())
}
