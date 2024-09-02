package main

import (
	"encoding/hex"
	"hash"
	"log"
	"net/http"

	"github.com/cbergoon/merkletree"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/sha3"
)

/**
 * Main function
	1.定时任务每天跑一次来生成空投列表
	2.空投列表生成merkle tree，
	3.生成merkle root
	4.生成merkle proof
	5.将空投列表和merkle root保存到数据库
	6.将merkle proof保存到redis
	7.如果合约中增加了新的空投，需要重新生成merkle tree，merkle root，merkle proof，然后更新数据库和redis
	8.用户在前端领取空投时，验证用户地址是否在merkle tree中，如果在则给前端返回merkle proof
	9.前端将merkle proof ，用户地址，空投金额发送到合约，合约验证merkle proof，用户地址，空投金额是否正确，如果正确则给用户转账
*/

type AirdropItem struct {
	Address      string
	Amount       string
	MerkleProof  []string
	MerkleRoot   string
	MerkleTreeId string
}

var airdropList = []AirdropItem{
	{Address: "0xEe014d7DfeB2e46Fef57CA4aDa42e79397edA76e", Amount: "1000000000000000000"},
	{Address: "0xe16C1623c1AA7D919cd2241d8b36d9E79C1Be2A2", Amount: "2000000000000000000"},
	{Address: "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266", Amount: "3000000000000000000"},
}

type MerkleTreeNode struct {
	Hash  string          `json:"hash"`
	Left  *MerkleTreeNode `json:"left,omitempty"`
	Right *MerkleTreeNode `json:"right,omitempty"`
}

type MerkleTree struct {
	Root *MerkleTreeNode `json:"root"`
}

var merkleTreeList []MerkleTree

func convertToCustomMerkleTree(node *merkletree.Node) *MerkleTreeNode {
	if node == nil {
		return nil
	}
	return &MerkleTreeNode{
		Hash:  hex.EncodeToString(node.Hash),
		Left:  convertToCustomMerkleTree(node.Left),
		Right: convertToCustomMerkleTree(node.Right),
	}
}

func (a AirdropItem) CalculateHash() ([]byte, error) {
	data := a.Address
	keccak256 := sha3.NewLegacyKeccak256()
	keccak256.Write([]byte(data))
	return keccak256.Sum(nil), nil
}

func (a AirdropItem) Equals(other merkletree.Content) (bool, error) {
	return a.Address == other.(AirdropItem).Address && a.Amount == other.(AirdropItem).Amount, nil
}

func generateProofs(tree *merkletree.MerkleTree, items []merkletree.Content) {
	proofs := make(map[string][][]byte)
	for _, item := range items {
		proof, _, err := tree.GetMerklePath(item)
		if err != nil {
			log.Fatal(err)
		}
		proofs[item.(AirdropItem).Address] = proof
	}

	for address, proof := range proofs {
		for _, p := range proof {
			hexProof := hex.EncodeToString(p)
			//更新airdropList，填充MerkleProof字段
			for i, item := range airdropList {
				if item.Address == address {
					airdropList[i].MerkleProof = append(airdropList[i].MerkleProof, hexProof)
				}
			}
		}
	}
}

func keccak256HashStrategy() hash.Hash {
	return sha3.NewLegacyKeccak256()
}

func generateMerkleTree(items []merkletree.Content) *merkletree.MerkleTree {

	// 构建merkleTreeList
	merkleTree, err := merkletree.NewTreeWithHashStrategy(items, keccak256HashStrategy)
	if err != nil {
		log.Fatal(err)
	}

	// 更新 merkleTreeList
	merkleTreeList = append(merkleTreeList, MerkleTree{
		Root: convertToCustomMerkleTree(merkleTree.Root),
	})
	return merkleTree
}

func generateMerkleRoot(tree *merkletree.MerkleTree) []byte {
	//为AirdropItem中的merkle root赋值
	root := tree.MerkleRoot()
	for i := range airdropList {
		airdropList[i].MerkleRoot = hex.EncodeToString(root)
		airdropList[i].MerkleTreeId = "1"
	}
	return tree.MerkleRoot()
}

func getAirdropList(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, airdropList)
}

func getMerkleTreeList(c *gin.Context) {
	var customMerkleTreeList []MerkleTree
	for _, tree := range merkleTreeList {
		customMerkleTreeList = append(customMerkleTreeList, MerkleTree{
			Root: tree.Root,
		})
	}

	c.IndentedJSON(http.StatusOK, merkleTreeList)
}

func getMerkleProofByAddress(c *gin.Context) {

	address := c.Param("address")

	var proofs []string
	for i := range airdropList {
		if airdropList[i].Address == address {
			proofs = append(proofs, airdropList[i].MerkleProof...)
		}
	}
	c.IndentedJSON(http.StatusOK, proofs)
}

func getMerkleRootByTreeId(c *gin.Context) {

	treeId := c.Param("treeId")

	var root string
	for _, item := range airdropList {
		if item.MerkleTreeId == treeId {
			root = item.MerkleRoot
		}
	}
	c.IndentedJSON(http.StatusOK, root)
}

func main() {
	var list []merkletree.Content
	for _, item := range airdropList {
		list = append(list, item)
	}

	tree := generateMerkleTree(list)

	generateMerkleRoot(tree)

	generateProofs(tree, list)

	// 设置 Gin 服务器
	router := gin.Default()
	router.GET("/airdrops/list", getAirdropList)
	router.GET("/airdrops/merkletree/list", getMerkleTreeList)
	router.GET("/airdrop/merkleproof/:address", getMerkleProofByAddress)
	router.GET("/airdrop/merkleroot/:treeId", getMerkleRootByTreeId)

	// 启动 Gin 服务器
	err := router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
