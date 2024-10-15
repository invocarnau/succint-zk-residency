use clap::Parser;
use sp1_sdk::{SP1Proof, HashableKey, utils, ProverClient, SP1Stdin, SP1ProofWithPublicValues, SP1VerifyingKey};
mod cli;
use cli::ProviderArgs;
use polccint_lib::fep_type_1::{BlockCommit, BlockAggregationInput, BlockAggregationCommit};
use std::path::PathBuf;
use polccint_lib::constants::BLOCK_VK;

#[derive(Parser, Debug)]
struct Args {
    /// The block number of the block to execute.
    #[clap(long)]
    block_number: u64,
    #[arg(long)]
    block_range: u64,
    #[clap(flatten)]
    provider: ProviderArgs,

    /// Whether or not to generate a proof.
    #[arg(long, default_value_t = false)]
    prove: bool,
}


const ELF_BLOCK: &[u8] = include_bytes!("../../../../elf/block");
const ELF_BLOCK_AGGREGATION: &[u8] = include_bytes!("../../../../elf/block-aggregation");


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
    //let rpc_url: Url = provider_config.rpc_url.expect("URL must be defined");

    //let provider = ReqwestProvider::new_http(rpc_url);

    // // Setup the host executor.
    // let host_executor = HostExecutor::new(provider);
    // let variant = match provider_config.chain_id {
    //     CHAIN_ID_ETH_MAINNET => ChainVariant::Ethereum,
    //     _ => {
    //         eyre::bail!("unknown chain ID: {}", provider_config.chain_id);
    //     }
    // };

    // Generate the proof.
    let client = ProverClient::new();

    // Setup the proving and verifying keys.
    let (aggregation_pk, aggregation_vk) = client.setup(ELF_BLOCK_AGGREGATION);
    let (_, block_vk) = client.setup(ELF_BLOCK);

   
    // assert constant vk with elf vk 
    assert!(block_vk.hash_u32() == BLOCK_VK);

    let initial_block_number = args.block_number;
    let block_range = args.block_range; // hardcode for now TODO
    let final_block_number = initial_block_number + block_range;

    let mut inputs: Vec<AggregationInput> = Vec::new();

    // for block_number in initial_block_number..initial_block_number + block_range {
    //     let client_input = host_executor
    //         .execute(block_number, variant)
    //         .await
    //         .expect("failed to execute host");

    //     let mut stdin_block = SP1Stdin::new();
    //     let buffer = bincode::serialize(&client_input).unwrap();
    //     stdin_block.write_vec(buffer);
    //     stdin_block.write(&client_input);
    //     let proof = client
    //         .prove(&block_pk.clone(), stdin_block)
    //         .compressed()
    //         .run()
    //         .expect("proving failed");
        
    //     inputs.push(
    //         AggregationInput {
    //             proof: proof,
    //             vk: block_vk.clone(),
    //         }
    //     );
    // }

    for block_number in initial_block_number..final_block_number + 1 {
        let proof: SP1ProofWithPublicValues = SP1ProofWithPublicValues::load(
            PathBuf::from(env!("CARGO_MANIFEST_DIR"))
                .join(format!(
                    "../../proof/chain{}/block_{}_proof.bin",
                    provider_config.chain_id,
                    block_number
                ))
        ).expect("failed to load proof");
        
        inputs.push(
            AggregationInput {
                proof: proof,
                vk: block_vk.clone(),
            }
        );
    }

    //print the public inputs values of all the proofs
    for input in &inputs {
        println!("{:?}", input.proof.public_values.clone().read::<BlockCommit>());
    }

    // encode aggregation input and write to stdin
    let mut stdin = SP1Stdin::new();
    let aggregation_input = BlockAggregationInput{
        block_commits: inputs
        .iter()
        .map(|input| input.proof.public_values.clone().read::<BlockCommit>())
        .collect::<Vec<_>>(),
        block_vkey: block_vk.clone().hash_u32(), // probabyl worth to change interface TODO
    };
    stdin.write(&aggregation_input);

    // write proofs
    for input in inputs {
        let SP1Proof::Compressed(proof) = input.proof.proof else {
            panic!()
        };
        stdin.write_proof(proof, input.vk.vk);
    }
    
    // Only execute the program.
    let (mut public_values, execution_report) =
        client.execute(&aggregation_pk.elf, stdin.clone()).run().unwrap();
    println!(
        "Finished executing the block in {} cycles",
        execution_report.total_instruction_count()
    );

    let decoded_public_values = public_values.read::<BlockAggregationCommit>();

    // Assert outputs
    // Check that the check was successful
    assert_eq!(decoded_public_values.prev_l2_block_hash, aggregation_input.block_commits[0].prev_block_hash);
    assert_eq!(decoded_public_values.new_l2_block_hash, aggregation_input.block_commits.last().unwrap().new_block_hash);

    // // If the `prove` argument was passed in, actually generate the proof.
    // // It is strongly recommended you use the network prover given the size of these programs.
    if args.prove {
        println!("Starting proof generation.");
        let proof: SP1ProofWithPublicValues = client.prove(&aggregation_pk, stdin.clone()).compressed().run().expect("Proving should work.");
        println!("Proof generation finished.");

        client.verify(&proof, &aggregation_vk).expect("proof verification should succeed");
        // Handle the result of the save operation
        match proof.save(PathBuf::from(env!("CARGO_MANIFEST_DIR")).join(format!("../../proof/chain{}/aggregation_{}_to_{}_proof.bin", provider_config.chain_id, initial_block_number, final_block_number))) {
            Ok(_) => println!("Proof saved successfully."),
            Err(e) => eprintln!("Failed to save proof: {}", e),
        }

        println!("Starting proof generation.");
    }
    Ok(())
}

