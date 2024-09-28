use sp1_build::{build_program_with_args, BuildArgs};

fn main() {
    build_program_with_args(
        &format!("../{}", "fep-type-1/block"),
        BuildArgs {
            ignore_rust_version: true,
            ..Default::default()
        },
    );
}
