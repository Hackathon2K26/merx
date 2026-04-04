// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {IERC20} from "./interfaces/IERC20.sol";

/// @title CompoundDepositor
/// @notice Deployed on Ethereum Sepolia. Self-relays a CCTP message (receiveMessage + mint),
///         then supplies the minted USDC into Compound V3 (Comet) on behalf of the shop.
/// @dev Only the operator (shop backend) can call relayAndSupply().
contract CompoundDepositor {
    IERC20 public immutable usdc;
    address public immutable operator;
    IMessageTransmitter public immutable messageTransmitter;
    IComet public immutable comet;

    event Supplied(address indexed beneficiary, uint256 amount);

    constructor(
        address _usdc,
        address _operator,
        address _messageTransmitter,
        address _comet
    ) {
        usdc = IERC20(_usdc);
        operator = _operator;
        messageTransmitter = IMessageTransmitter(_messageTransmitter);
        comet = IComet(_comet);
        // Pre-approve Comet to pull USDC.
        IERC20(_usdc).approve(_comet, type(uint256).max);
    }

    /// @notice Relay a CCTP message to mint USDC here, then supply into Compound.
    /// @param message     Raw CCTP message bytes (from attestation API)
    /// @param attestation Circle attestation signature (from attestation API)
    /// @param beneficiary Address to credit in Compound (receives cUSDCv3 balance)
    function relayAndSupply(
        bytes calldata message,
        bytes calldata attestation,
        address beneficiary
    ) external {
        require(msg.sender == operator, "only operator");

        // 1. Relay CCTP message — mints USDC to this contract.
        bool success = messageTransmitter.receiveMessage(message, attestation);
        require(success, "receiveMessage failed");

        // 2. Supply all received USDC into Compound on behalf of beneficiary.
        uint256 balance = usdc.balanceOf(address(this));
        require(balance > 0, "no USDC received");

        comet.supplyTo(beneficiary, address(usdc), balance);

        emit Supplied(beneficiary, balance);
    }
}

interface IMessageTransmitter {
    function receiveMessage(bytes calldata message, bytes calldata attestation) external returns (bool);
}

/// @dev Compound V3 Comet interface (subset).
interface IComet {
    function supplyTo(address dst, address asset, uint256 amount) external;
}
