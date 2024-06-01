package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

// TxInput es una entrada en una transacción
// tiene ID como identificación
// Out es el monto
// Sig es el destinatario
type TxInput struct {
	ID  []byte
	Out int
	Sig string
}

// CanUnlock nos permite saber si el usuario puede ver la transacción
// porque fue el destinatario de esta
func (in *TxInput) CanUnlock(data string) bool {
	return in.Sig == data
}

// TxOutput es una salida en una transacción
// tiene Value como el monto
// PubKey es el usuario que está enviando la transacción
type TxOutput struct {
	Value  int
	PubKey string
}

// CanBeUnlocked nos permite saber si el usuario puede ver la transacción
// porque fue el emisor de esta
func (out *TxOutput) CanBeUnlocked(data string) bool {
	return out.PubKey == data
}

// Transaction y sus métodos

// Transaction representa la estructura de una transacción
// tiene un ID, y arreglos para sus entradas y salidas
type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

// SetID define el id para la transacción
// qué básicamente es el hash de su información
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	encode := gob.NewEncoder(&encoded)
	err := encode.Encode(tx)
	Handle(err)

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

func (tx *Transaction) IsCoinBase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

// métodos

func CoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Divisas a %s", to)
	}

	txIn := TxInput{
		ID:  []byte{},
		Out: -1,
		Sig: data,
	}
	txOut := TxOutput{
		Value:  100,
		PubKey: to,
	}

	tx := Transaction{
		ID:      nil,
		Inputs:  []TxInput{txIn},
		Outputs: []TxOutput{txOut},
	}
	tx.SetID()

	return &tx
}

func NewTransaction(from, to string, amount int, blockchain *BlockChain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	acc, validOutputs := blockchain.FindSpendableOutputs(from, amount)

	if acc < amount {
		log.Panic("Error: no hay suficientes fondos")
	}

	for txId, outs := range validOutputs {
		txID, err := hex.DecodeString(txId)
		Handle(err)

		for _, out := range outs {
			input := TxInput{
				ID:  txID,
				Out: out,
				Sig: from,
			}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TxOutput{
		Value:  amount,
		PubKey: to,
	})

	if acc > amount {
		outputs = append(outputs, TxOutput{
			Value:  acc - amount,
			PubKey: from,
		})
	}

	tx := Transaction{
		ID:      nil,
		Inputs:  inputs,
		Outputs: outputs,
	}
	tx.SetID()

	return &tx
}
