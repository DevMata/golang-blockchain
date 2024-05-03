package main

import (
	"fmt"
	"github.com/devmata/golang-blockchain/blockchain"
)

func main() {
	blockChain := blockchain.InitBlockChain()
	blockChain.AddBlock("First block after Genesis")
	blockChain.AddBlock("Second block after Genesis")
	blockChain.AddBlock("Third block after Genesis")

	for _, block := range blockChain.Blocks {
		fmt.Printf("Previous hash: %x\n", block.PrevHash)
		fmt.Printf("Data in block: %s\n", block.Data)
		fmt.Printf("Hash: %x\n\n", block.Hash)
	}
}
