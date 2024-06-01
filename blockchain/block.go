package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
)

// Block - y métodos

// Block Es la estructura base de la cadena
// Hash del bloque
// Data contiene la info del bloque
// PrevHash hash al bloque anterior
// Nonce
type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
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

// Métodos para la lógica de Block

// CreateBlock permite instaciar un bloque Block
// data es cualquier cosa que se quiera almacenar
// espera también el hash del bloque previo
func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash, 0}
	pow := NewProof(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// Genesis es una llamada específica para CreateBlock
// que nos permite crear el bloque génesis
func Genesis() *Block {
	return CreateBlock("Bloque Génesis", []byte{})
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
