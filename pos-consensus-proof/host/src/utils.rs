use crate::types::{BlockResponse, MilestoneResponse, TxResponse, ValidatorSetResponse};

use alloy_primitives::FixedBytes;
use alloy_provider::ReqwestProvider;
use alloy_rpc_types::BlockNumberOrTag;
use eyre::Result;
use ethers::providers::{Http, Middleware, Provider};
use ethers::types::{BlockId, H256};
use reqwest::header::{HeaderMap, HeaderValue, USER_AGENT};
use reqwest::Client;
use reth_primitives::Header;
use std::env;
use url::Url;

use sp1_cc_host_executor::HostExecutor;
// PosClient holds a http client instance along with endpoints for heimdall rest-server,
// tendermint rpc server and bor's rpc server to interact with.
pub struct PosClient {
    heimdall_url: String,
    tendermint_url: String,
    bor_rpc_url: String,
    http_client: Client,
    headers: HeaderMap,
}

impl Default for PosClient {
    fn default() -> Self {
        let heimdall_url =
            env::var("HEIMDALL_REST_ENDPOINT").expect("HEIMDALL_REST_ENDPOINT not set");
        let tendermint_url = env::var("TENDERMINT_ENDPOINT").expect("TENDERMINT_ENDPOINT not set");
        let http_client = Client::new();
        let bor_rpc_url = env::var("BOR_RPC_URL").expect("BOR_RPC_URL not set");

        let mut headers = HeaderMap::new();
        headers.insert(
            USER_AGENT,
            HeaderValue::from_static("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36"),
        );

        Self {
            heimdall_url,
            tendermint_url,
            bor_rpc_url,
            http_client,
            headers,
        }
    }
}

impl PosClient {
    pub fn new(heimdall_url: String, tendermint_url: String, bor_rpc_url: String) -> Self {
        let mut headers = HeaderMap::new();
        headers.insert(
            USER_AGENT,
            HeaderValue::from_static("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36"),
        );
        Self {
            heimdall_url,
            tendermint_url,
            bor_rpc_url,
            http_client: Client::new(),
            headers,
        }
    }

    /// Fetches a heimdall milestone by id
    pub async fn fetch_milestone_by_id(&self, id: u64) -> Result<MilestoneResponse> {
        let url = format!("{}/milestone/{}", self.heimdall_url, id);
        println!("Fetching milestone from: {}", url);
        let response = self
            .http_client
            .get(url)
            .headers(self.headers.clone())
            .send()
            .await?
            .json::<MilestoneResponse>()
            .await?;
        Ok(response)
    }

    /// Fetches a tendermint transaction by hash
    pub async fn fetch_tx_by_hash(&self, hash: String) -> Result<TxResponse> {
        let url = format!("{}/tx?hash={}", self.tendermint_url, hash);
        println!("Fetching milestone tx by hash: {}", url);
        let response: TxResponse = self
            .http_client
            .get(url)
            .headers(self.headers.clone())
            .send()
            .await?
            .json::<TxResponse>()
            .await?;
        Ok(response)
    }

    /// Fetches a tendermint block by number
    pub async fn fetch_block_by_number(&self, number: u64) -> Result<BlockResponse> {
        let url = format!("{}/block?height={}", self.tendermint_url, number);
        println!("Fetching block by number: {}", url);
        let response = self
            .http_client
            .get(url)
            .headers(self.headers.clone())
            .send()
            .await?
            .json::<BlockResponse>()
            .await?;
        Ok(response)
    }

    /// Fetches the validator set from heimdall
    pub async fn fetch_validator_set(&self) -> Result<ValidatorSetResponse> {
        let url: String = format!("{}/staking/validator-set", self.heimdall_url);
        println!("Fetching validator set from: {}", url);
        let response: ValidatorSetResponse = self
            .http_client
            .get(url)
            .headers(self.headers.clone())
            .send()
            .await?
            .json::<ValidatorSetResponse>()
            .await?;
        Ok(response)
    }

    /// Fetches the validator set from heimdall at a specific block height
    pub async fn fetch_validator_set_by_height(&self, height: u64) -> Result<ValidatorSetResponse> {
        let url: String = format!(
            "{}/staking/validator-set?height={}",
            self.heimdall_url, height
        );
        println!("Fetching validator set at height from: {}", url);
        let response = self
            .http_client
            .get(url)
            .headers(self.headers.clone())
            .send()
            .await?;
        let json_response = response.json::<ValidatorSetResponse>().await;
        if json_response.is_ok() {
            Ok(json_response?)
        } else {
            println!("Failed to fetch validator set at height, please use an archive node. Using latest block instead");
            let response = self.fetch_validator_set().await?;
            Ok(response)
        }
    }

    /// Fetches a bor block header (of type `reth-primitives::Header`) by number
    pub async fn fetch_bor_header_by_number(
        &self,
        block_number: BlockNumberOrTag,
    ) -> Result<Header> {
        // Use the host executor to fetch the required bor block
        let bor_provider = ReqwestProvider::new_http(Url::parse(&self.bor_rpc_url)?);
        let bor_host_executor = HostExecutor::new(bor_provider.clone(), block_number)
            .await
            .expect("unable to fetch bor block by number");
        Ok(bor_host_executor.header)
    }

    /// Fetches a bor block number by hash
    pub async fn fetch_bor_number_by_hash(&self, block_hash: FixedBytes<32>) -> Result<u64> {
        let provider = Provider::<Http>::try_from(self.bor_rpc_url.clone())?;
        let hash: H256 = H256::from_slice(block_hash.as_ref());
        let block = provider
            .get_block(BlockId::Hash(hash))
            .await?
            .expect("unable to fetch bor block by hash");
        Ok(block.number.unwrap().as_u64())
    }
}
