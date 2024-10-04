use alloy_provider::ReqwestProvider;
use clap::Parser;
use reth_primitives::B256;
use rsp_client_executor::{io::ClientExecutorInput, ChainVariant, CHAIN_ID_ETH_MAINNET};
use rsp_host_executor::HostExecutor;
use sp1_sdk::{HashableKey, utils, ProverClient, SP1Stdin, SP1ProofWithPublicValues};
use std::path::PathBuf;
mod cli;
use cli::ProviderArgs;
use url::Url;
use polccint_lib::{BlockCommit, SP1CCProofFixture};
use alloy::hex;
use std::io::Read;

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
        _ => {
            eyre::bail!("unknown chain ID: {}", provider_config.chain_id);
        }
    };

    // println!("ChainID: {:?}", provider_config.chain_id);
    // println!("Executing block number: {:?}", args.block_number);

    let client_input = host_executor
        .execute(args.block_number, variant)
        .await
        .expect("failed to execute host");

    //let client_input = load_input_from_cache(CHAIN_ID_ETH_MAINNET, 20526624);

    // Generate the proof.
    let client = ProverClient::new();

    // Setup the proving key and verification key.
    let (pk, vk) = client.setup(ELF);

    // Write the block to the program's stdin.
    let mut stdin = SP1Stdin::new();
    let buffer = bincode::serialize(&client_input).unwrap();
    stdin.write_vec(buffer);

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
    assert_eq!(decoded_public_values.new_block_hash, client_input.current_block.hash_slow());

    // If the `prove` argument was passed in, actually generate the proof.
    // It is strongly recommended you use the network prover given the size of these programs.
    if args.prove {
        println!("Starting proof generation.");
        let proof: SP1ProofWithPublicValues = client
            .prove(&pk, stdin)
            .compressed()
            .run()
            .expect("Proving should work.");
        println!("Proof generation finished.");

        save_fixture(vk.clone().bytes32(), &proof, args.block_number, provider_config.chain_id);

        client
            .verify(&proof, &vk)
            .expect("proof verification should succeed");
    }
    Ok(())
}

/// Generate a `SP1CCProofFixture`, and save it as a json file.
///
/// This is useful for verifying the proof of contract call execution on chain.
fn save_fixture(vkey: String, proof: &SP1ProofWithPublicValues, block_number: u64, chain_id: u64) {
    let fixture = SP1CCProofFixture {
        vkey,
        public_values: format!("0x{}", hex::encode(proof.public_values.as_slice())),
        proof: format!("0x{}", hex::encode(proof.bytes())),
    };

    let fixture_path = PathBuf::from(env!("CARGO_MANIFEST_DIR")).join(format!("../proof/chain{}/block_{}_proof_fixture.jso", chain_id, block_number));
    std::fs::create_dir_all(&fixture_path).expect("failed to create fixture path");
    std::fs::write(
        fixture_path.join("plonk-fixture.json".to_lowercase()),
        serde_json::to_string_pretty(&fixture).unwrap(),
    )
    .expect("failed to write fixture");
}

fn read_fixture(block_number: u64, chain_id: u64) -> SP1CCProofFixture {
    let fixture_path = PathBuf::from(env!("CARGO_MANIFEST_DIR")).join(format!("../proof/chain{}/block_{}_proof_fixture.json", chain_id, block_number));
    let mut file = std::fs::File::open(fixture_path).expect("failed to open fixture file");
    let mut contents = String::new();
    file.read_to_string(&mut contents).expect("failed to read fixture file");
    serde_json::from_str(&contents).expect("failed to deserialize fixture")
}
