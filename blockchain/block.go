package blockchain

import (
	"bytes"
	"crypto/sha256"
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
	Blocks []*Block
}

func (blockchain *BlockChain) AddBlock(data string) {
	previousBlock := blockchain.Blocks[len(blockchain.Blocks)-1]
	newBlock := CreateBlock(data, previousBlock.Hash)
	blockchain.Blocks = append(blockchain.Blocks, newBlock)
}

func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

func InitBlockChain() *BlockChain {
	return &BlockChain{Blocks: []*Block{Genesis()}}
}
