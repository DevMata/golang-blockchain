package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
}

func (b *Block) DeriveHash() {
	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
	hash := sha256.Sum256(info)
	b.Hash = hash[:]
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash}
	block.DeriveHash()
	return block
}

type BlockChain struct {
	blocks []*Block
}

func (blockchain *BlockChain) AddBlock(data string) {
	previousBlock := blockchain.blocks[len(blockchain.blocks)-1]
	newBlock := CreateBlock(data, previousBlock.Hash)
	blockchain.blocks = append(blockchain.blocks, newBlock)
}

func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

func InitBlockChain() *BlockChain {
	return &BlockChain{blocks: []*Block{Genesis()}}
}

func main() {
	blockChain := InitBlockChain()

	blockChain.AddBlock("First block after Genesis")
	blockChain.AddBlock("Second block after Genesis")
	blockChain.AddBlock("Third block after Genesis")

	for _, block := range blockChain.blocks {
		fmt.Printf("Previous hash: %x\n", block.PrevHash)
		fmt.Printf("Data in block: %s\n", block.Data)
		fmt.Printf("Hash: %x\n\n", block.Hash)
	}
}
