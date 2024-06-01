package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

// Block - y métodos

// Block Es la estructura base de la cadena
// Hash del bloque
// Transactions contiene la info de las transacciones
// PrevHash hash al bloque anterior
// Nonce
type Block struct {
	Hash         []byte
	Transactions []*Transaction
	PrevHash     []byte
	Nonce        int
}

// Serialize Permite convertir una instancia de Block en Byte
// para almacenarlo en la BDD
func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)
	Handle(err)

	return res.Bytes()
}

// HashTransactions nos permite obtener el hash en bytes
// de las transacciones para un bloque
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}

// Métodos para la lógica de Block

// CreateBlock permite instaciar un bloque Block
// pasamos las transacciones contenidas en el bloque
// espera también el hash del bloque previo
func CreateBlock(txs []*Transaction, prevHash []byte) *Block {
	block := &Block{[]byte{}, txs, prevHash, 0}
	pow := NewProof(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// Genesis es una llamada específica para CreateBlock
// que nos permite crear el bloque génesis
func Genesis(coinbase *Transaction) *Block {
	return CreateBlock([]*Transaction{coinbase}, []byte{})
}

// Deserialize nos permite pasar de Byte a una instancia de bloque
func Deserialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)
	Handle(err)

	return &block
}

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
