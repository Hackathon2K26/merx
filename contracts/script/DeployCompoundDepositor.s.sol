// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script, console} from "forge-std/Script.sol";
import {CompoundDepositor} from "../src/CompoundDepositor.sol";

/// @notice Deploy CompoundDepositor on Ethereum Sepolia.
///
///   PRIVATE_KEY=0x... forge script script/DeployCompoundDepositor.s.sol \
///     --rpc-url https://ethereum-sepolia-rpc.publicnode.com --broadcast --skip-simulation
contract DeployCompoundDepositor is Script {
    // Ethereum Sepolia addresses.
    address constant USDC = 0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238;
    address constant MESSAGE_TRANSMITTER = 0xE737e5cEBEEBa77EFE34D4aa090756590b1CE275;
    address constant COMPOUND_COMET = 0xAec1F48e02Cfb822Be958B68C7957156EB3F0b6e;

    function run() external {
        uint256 deployerKey = vm.envUint("PRIVATE_KEY");
        address operator = vm.addr(deployerKey);

        vm.startBroadcast(deployerKey);
        CompoundDepositor depositor = new CompoundDepositor(USDC, operator, MESSAGE_TRANSMITTER, COMPOUND_COMET);
        vm.stopBroadcast();

        console.log("CompoundDepositor deployed at:", address(depositor));
        console.log("  USDC:", USDC);
        console.log("  MessageTransmitter:", MESSAGE_TRANSMITTER);
        console.log("  Compound Comet:", COMPOUND_COMET);
        console.log("  Operator:", operator);
    }
}
