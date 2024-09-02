const { ethers } = require("hardhat");
const { expect } = require("chai");
const {
    loadFixture,
} = require("@nomicfoundation/hardhat-toolbox/network-helpers");

describe("AirdropImplByMerkleTree Contract", function () {


    async function deployAirdropImplByMerkleTreeFixture() {

        const AirdropImplByMerkleTree = await ethers.getContractFactory("AirdropImplByMerkleTree");
        const merkleRoot = "0xccbdde18eb6ea041e7c408ad5b77def68fe2b3c534195436ecb043f197271cff";

        const airdropImplByMerkleTreeContract = await AirdropImplByMerkleTree.deploy(merkleRoot);

        return {
            airdropImplByMerkleTreeContract
        };
    }


    describe("Claim", function () {

        it("Deployment", async function () {
            const { airdropImplByMerkleTreeContract } = await loadFixture(deployAirdropImplByMerkleTreeFixture);
            const merkleRoot = await airdropImplByMerkleTreeContract.merkleRoot();
            expect(merkleRoot).to.equal("0xccbdde18eb6ea041e7c408ad5b77def68fe2b3c534195436ecb043f197271cff");

        });

        it("Should verify before claim", async function () {
            const { airdropImplByMerkleTreeContract } = await loadFixture(deployAirdropImplByMerkleTreeFixture);

            const merkleProof = [
                "0x7980ab3658af943c225e5c8841ad65c3c4e5936f5ccbad9e4ebe4ad358e81601",
                "0xe9707d0e6171f728f7473c24cc0432a9b07eaaf1efed6a137a4a8c12c79552d9"
            ]

            const verifyAccount = "0xEe014d7DfeB2e46Fef57CA4aDa42e79397edA76e";


           await airdropImplByMerkleTreeContract.claim(merkleProof, verifyAccount, BigInt("1000000000000000000"));

        });

    });

});