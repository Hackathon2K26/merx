// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script, console} from "forge-std/Script.sol";
import {ShopPaymaster} from "../src/ShopPaymaster.sol";

/// @notice Deploy ShopPaymaster on a source chain.
///
///   PRIVATE_KEY=0x... forge script script/DeployShopPaymaster.s.sol \
///     --rpc-url https://sepolia.unichain.org --broadcast --skip-simulation
contract DeployShopPaymaster is Script {
    address constant TOKEN_MESSENGER_V2 = 0x8FE6B999Dc680CcFDD5Bf7EB0974218be2542DAA;
    // Shop wallet on Arc — CCTP mints USDC directly here.
    address constant SHOP_WALLET = 0x2A94238046B648EFF3Ec899fbe6C2B7990C52ca3;
    uint32 constant ARC_DOMAIN = 26;

    function _usdc() internal view returns (address) {
        address env = vm.envOr("USDC", address(0));
        if (env != address(0)) return env;
        uint256 chainId = block.chainid;
        if (chainId == 1301) return 0x31d0220469e10c4E71834a79b1f276d740d3768F;
        if (chainId == 84532) return 0x036CbD53842c5426634e7929541eC2318f3dCF7e;
        if (chainId == 11155111) return 0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238;
        revert("unsupported chain - set USDC env var");
    }

    function run() external {
        uint256 deployerKey = vm.envUint("PRIVATE_KEY");
        address usdc = _usdc();
        address messenger = vm.envOr("TOKEN_MESSENGER", TOKEN_MESSENGER_V2);

        vm.startBroadcast(deployerKey);
        ShopPaymaster paymaster = new ShopPaymaster(usdc, messenger, ARC_DOMAIN, SHOP_WALLET);
        vm.stopBroadcast();

        console.log("ShopPaymaster deployed at:", address(paymaster));
        console.log("  USDC:", usdc);
        console.log("  TokenMessenger:", messenger);
        console.log("  Arc domain:", ARC_DOMAIN);
        console.log("  mintRecipient (shop wallet):", SHOP_WALLET);
    }
}
