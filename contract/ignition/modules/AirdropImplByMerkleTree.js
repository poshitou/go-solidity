const { buildModule } = require("@nomicfoundation/hardhat-ignition/modules");

const AirdropImplByMerkleTreeModule = buildModule("AirdropImplByMerkleTreeModule", (m) => {

    const _merkleRoot = "0x48eac99e5c170876825a939679d1d2b1a567a6b9a90e218d0d9721bded4c6ca9";

    const token = m.contract("AirdropImplByMerkleTree",[_merkleRoot]);

    return { token };
});

module.exports = AirdropImplByMerkleTreeModule;