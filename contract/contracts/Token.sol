// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/utils/Context.sol";

contract Token is Context, IERC20 {

    string public name;

    string public symbol;

    uint8 public decimals;

    uint256 public override totalSupply;

    constructor(string memory _name, string memory _symbol, uint8 _decimals, uint256 _totalSupply) {
        name = _name;
        symbol = _symbol;
        decimals = _decimals;
        totalSupply = _totalSupply;
        balanceOf[_msgSender()] = _totalSupply;
    }

    mapping(address => uint256) public override balanceOf;

    mapping(address => mapping(address => uint256)) public override allowance;

    function transfer(address to, uint256 value) public override returns (bool) {
        require(balanceOf[_msgSender()] >= value, "Insufficient balance");
        balanceOf[_msgSender()] -= value;
        balanceOf[to] += value;
        emit Transfer(_msgSender(), to, value);
        return true;
    }

    function approve(address spender, uint256 value) public override returns (bool) {
        allowance[_msgSender()][spender] = value;
        emit Approval(_msgSender(), spender, value);
        return true;
    }

    function transferFrom(address from, address to, uint256 value) public override returns (bool) {
        require(balanceOf[from] >= value, "Insufficient balance");
        require(allowance[from][_msgSender()] >= value, "Insufficient allowance");
        balanceOf[from] -= value;
        balanceOf[to] += value;
        allowance[from][_msgSender()] -= value;
        emit Transfer(from, to, value);
        return true;
    }

    function mint(address to, uint256 value) public {
        totalSupply += value;
        balanceOf[to] += value;
        emit Transfer(address(0), to, value);
    }

    function burn(uint256 value) public {
        require(balanceOf[_msgSender()] >= value, "Insufficient balance");
        totalSupply -= value;
        balanceOf[_msgSender()] -= value;
        emit Transfer(_msgSender(), address(0), value);
    }

}
