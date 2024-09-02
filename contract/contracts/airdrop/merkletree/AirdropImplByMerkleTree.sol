// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/utils/cryptography/MerkleProof.sol";

contract AirdropImplByMerkleTree {
    //  0x5FbDB2315678afecb367f032d93F642f64180aa3
    bytes32 public merkleRoot;

    mapping(address => uint256) public claimedMap;

    constructor(bytes32 _merkleRoot) {
        merkleRoot = _merkleRoot;
    }

    /**
        在实际项目中，Merkle Tree 通常是在后端生成的，而 Merkle Proof 则可以在前端生成。具体流程如下：
            1.后端生成 Merkle Tree：
                -后端根据空投名单生成 Merkle Tree，并计算出 Merkle Root。
                -将 Merkle Root 部署到智能合约中，或者通过其他方式传递给前端。
            2.前端调用后端生成的 Merkle Proof：
                -前端将从后端获取的merkle proof和领取空投的地址和金额发送到智能合约进行验证和领取空投。
    */

    function claim(
        bytes32[] memory proof,
        address account,
        uint256 amount
    ) public returns (bool result) {
        require(claimedMap[account] == 0, "Airdrop already claimed");

        bytes32 leaf = keccak256(abi.encodePacked(account));

        result = MerkleProof.verify(proof, merkleRoot, leaf);

        require(result, "Invalid proof");
        claimedMap[account] = amount;
    }
}
