// SPDX-License-Identifier: AGPL-3.0
pragma solidity ^0.8.20;

interface IPolygonZkEVMGlobalExitRootV2 {
    function globalExitRootMap(bytes32) external view returns (uint256);
}

contract MultiGERAssertor {
    // Global Exit Root address
    IPolygonZkEVMGlobalExitRootV2 public immutable globalExitRootManager;

    constructor(IPolygonZkEVMGlobalExitRootV2 _globalExitRootManager) {
        globalExitRootManager = _globalExitRootManager;
    }

    function CheckGERsExistance(bytes32[] calldata GERs) public view returns (bool) {
        for (uint256 i = 0; i < GERs.length; i++) {
            if (globalExitRootManager.globalExitRootMap(GERs[i]) == 0) {
                return false;
            }
        }
        return true;
    }
}