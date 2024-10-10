#!/bin/bash

# Check if the correct number of arguments is provided
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <start_block> <range>"
    exit 1
fi

# Read start block and range from arguments
start_block=$1
range=$2

# Define the number of parallel jobs
num_jobs=6  # Adjust this number based on your CPU cores

# Function to run cargo with a specific block number
run_cargo() {
    block_number=$1
    echo "Running block number $block_number"
    cargo run --release -- --chain-id 42069 --block-number "$block_number" --prove
}

export -f run_cargo

# Run the jobs in parallel
seq "$start_block" "$((start_block + range))" | xargs -n 1 -P "$num_jobs" -I {} bash -c 'run_cargo "$@"' _ {}

echo "All jobs are done!"
