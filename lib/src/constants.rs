use alloy_primitives::{address, Address};

/// Address of the caller.
pub const CALLER: Address = address!("70997970c51812dc3a010c7d01b50e0d17dc79c8");

// 006ec047f79b977379bf77e5d307fcba13d35ce1f082daddab6f653ff6f1c642
pub const BLOCK_VK: [u32; 8] = [
    0x006ec047, 0xf79b9773, 0x79bf77e5, 0xd307fcba,
    0x13d35ce1, 0xf082dadd, 0xab6f653f, 0xf6f1c642,
];

// 00b5890ad47a7173b6796a01ea19229b23b5578e591f43b8975c52bea8edb645
pub const BRIDGE_VK: [u32; 8] = [
    0x00b5890a, 0xd47a7173, 0xb6796a01, 0xea19229b,
    0x23b5578e, 0x591f43b8, 0x975c52be, 0xa8edb645,
];

// aggregation vk
// 0028cf30add79337cef8fbc67e8922c214f6bf278fb2fc78a9cf9dfcacd6ec11
pub const AGGREGATION_VK: [u32; 8] = [
    0x0028cf30, 0xadd79337, 0xcef8fbc6, 0x7e8922c2,
    0x14f6bf27, 0x8fb2fc78, 0xa9cf9dfc, 0xacd6ec11,
];


