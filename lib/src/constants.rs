use alloy_primitives::{address, Address};

// final aggregation vk: 0x00c7dca51c03c7b4db25b4c342d4178b8e7e1107dbcbf246c372a91c2950a068

/// Address of the caller.
pub const CALLER: Address = address!("70997970c51812dc3a010c7d01b50e0d17dc79c8");
pub const CALLER_L1: Address = address!("0000000000000000000000000000000000000000");

pub const BLOCK_VK: [u32; 8] = [929047547, 1726340318, 938409146, 813681569, 513468175, 1108044662, 1457441407, 1995556418];

pub const BRIDGE_VK: [u32; 8] = [849736881, 658232624, 1674334477, 1541155920, 438872218, 596269535, 1057449073, 523602103];

// aggregation vk
pub const AGGREGATION_VK: [u32; 8] = [1480055870, 451662487, 1899574304, 729955651, 300839141, 1354052588, 1891154660, 47361319];

// OP
pub const OP_CONSENSUS_VK: [u32; 8] = [0, 0, 0, 0, 0, 0, 0, 0]; // TODO: add correct vkey


