// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test} from "forge-std/Test.sol";
import {CompoundDepositor} from "../src/CompoundDepositor.sol";

contract CompoundDepositorTest is Test {
    CompoundDepositor depositor;
    MockUSDC usdc;
    MockMessageTransmitter transmitter;
    MockComet comet;
    address operator;
    address shop = address(0x5408);

    function setUp() public {
        operator = address(this);
        usdc = new MockUSDC();
        transmitter = new MockMessageTransmitter(address(usdc));
        comet = new MockComet(address(usdc));
        depositor = new CompoundDepositor(
            address(usdc),
            operator,
            address(transmitter),
            address(comet)
        );
        transmitter.setMintTarget(address(depositor));
    }

    function test_relayAndSupply() public {
        transmitter.setMintAmount(5e6);

        depositor.relayAndSupply("message", "attestation", shop);

        assertEq(usdc.balanceOf(address(depositor)), 0);
        assertEq(comet.lastDst(), shop);
        assertEq(comet.lastAsset(), address(usdc));
        assertEq(comet.lastAmount(), 5e6);
    }

    function test_event() public {
        transmitter.setMintAmount(3e6);

        vm.expectEmit(true, false, false, true);
        emit CompoundDepositor.Supplied(shop, 3e6);
        depositor.relayAndSupply("att", "sig", shop);
    }

    function test_revert_notOperator() public {
        transmitter.setMintAmount(1e6);

        vm.prank(address(0xBAD));
        vm.expectRevert("only operator");
        depositor.relayAndSupply("att", "sig", shop);
    }

    function test_revert_receiveMessageFails() public {
        transmitter.setShouldFail(true);

        vm.expectRevert("receiveMessage failed");
        depositor.relayAndSupply("msg", "att", shop);
    }

    function test_revert_noUSDC() public {
        transmitter.setMintAmount(0);

        vm.expectRevert("no USDC received");
        depositor.relayAndSupply("att", "sig", shop);
    }

    function test_constructorApproval() public view {
        uint256 allowance = usdc.allowance(address(depositor), address(comet));
        assertEq(allowance, type(uint256).max);
    }
}

// ---------------------------------------------------------------------------
// Mocks
// ---------------------------------------------------------------------------

contract MockUSDC {
    mapping(address => uint256) public balanceOf;
    mapping(address => mapping(address => uint256)) public allowance;

    function mint(address to, uint256 amount) external {
        balanceOf[to] += amount;
    }

    function approve(address spender, uint256 amount) external returns (bool) {
        allowance[msg.sender][spender] = amount;
        return true;
    }

    function transferFrom(address from, address to, uint256 amount) external returns (bool) {
        require(balanceOf[from] >= amount, "insufficient");
        require(allowance[from][msg.sender] >= amount, "no allowance");
        allowance[from][msg.sender] -= amount;
        balanceOf[from] -= amount;
        balanceOf[to] += amount;
        return true;
    }
}

contract MockMessageTransmitter {
    MockUSDC immutable usdc;
    address mintTarget;
    uint256 mintAmount;
    bool shouldFail;

    constructor(address _usdc) { usdc = MockUSDC(_usdc); }
    function setMintTarget(address t) external { mintTarget = t; }
    function setMintAmount(uint256 a) external { mintAmount = a; }
    function setShouldFail(bool f) external { shouldFail = f; }

    function receiveMessage(bytes calldata, bytes calldata) external returns (bool) {
        if (shouldFail) return false;
        if (mintAmount > 0) usdc.mint(mintTarget, mintAmount);
        return true;
    }
}

contract MockComet {
    MockUSDC immutable usdc;
    address public lastDst;
    address public lastAsset;
    uint256 public lastAmount;

    constructor(address _usdc) { usdc = MockUSDC(_usdc); }

    function supplyTo(address dst, address asset, uint256 amount) external {
        MockUSDC(asset).transferFrom(msg.sender, address(this), amount);
        lastDst = dst;
        lastAsset = asset;
        lastAmount = amount;
    }
}
