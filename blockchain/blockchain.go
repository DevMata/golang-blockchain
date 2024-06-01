package blockchain

import (
	"errors"
	"fmt"
	"github.com/dgraph-io/badger/v4"
)

const (
	dbPath = "./db/blocks"
)

// Blockchain y sus métodos

// BlockChain es la estructura para la cadena
// contiene el hash del último bloque insertado
// y una referencia a la BDD
type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

// AddBlock nos permite insertar un bloque nuevo a la cadena
// espera data - que puede ser cualquier cosa que queramos insertar
func (blockchain *BlockChain) AddBlock(data string) {
	var lastHash []byte

	// buscamos el hash del último bloque - lastHash
	err := blockchain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		return err
	})

	Handle(err)

	// instanciamos el nuevo bloque
	newBlock := CreateBlock(data, lastHash)

	// actualizamos la info en la BDD, y el valor del último hash en la cadena
	err = blockchain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)

		blockchain.LastHash = newBlock.Hash

		return err
	})
	Handle(err)
}

// Iterator nos permite obtener el iterador para la cadena
// con este podremos recorrer a lo largo de la cadena
func (blockchain *BlockChain) Iterator() *Iterator {
	iter := &Iterator{
		CurrentHash: blockchain.LastHash,
		Database:    blockchain.Database,
	}
	return iter
}

// Iterator y sus métodos

// Iterator es una estructura que nos ayudará a iterar(recorrer) a lo largo de la cadena
// contiene una referencia a la BDD
type Iterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

// Next función que permite llamar al siguiente bloque en la cadena
// mientras la estamos recorriendo
func (iterator *Iterator) Next() *Block {
	var block *Block

	// buscamos el bloque cuyo hash tenemos como referencia
	err := iterator.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iterator.CurrentHash)
		Handle(err)
		err = item.Value(func(val []byte) error {
			block = Deserialize(val)
			return nil
		})
		return err
	})
	Handle(err)

	// guardamos la referencia al siguiente bloque(bloque anterior según BlockChain)
	iterator.CurrentHash = block.PrevHash

	return block
}

// Métodos

// InitBlockChain nos permite instanciar la cadena
// Conectará la BDD
// E insertará el bloque Génesis si hace falta
func InitBlockChain() *BlockChain {
	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath
	opts.Logger = nil

	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		// validamos si ya hay cadena, sino se creará un bloque Génesis
		if _, err := txn.Get([]byte("lh")); errors.Is(err, badger.ErrKeyNotFound) {
			fmt.Println("No se ha encontrado una cadena")
			genesis := Genesis()
			fmt.Println("Bloque Génesis ha sido provisto")
			err = txn.Set(genesis.Hash, genesis.Serialize())
			Handle(err)
			err = txn.Set([]byte("lh"), genesis.Hash)

			lastHash = genesis.Hash

			return err
		} else { // si hay ya una cadena, guardaremos el hash del último bloque que contiene
			item, err := txn.Get([]byte("lh"))
			Handle(err)
			err = item.Value(func(val []byte) error {
				lastHash = val
				return nil
			})
			return err
		}
	})

	Handle(err)

	blockchain := BlockChain{
		LastHash: lastHash, Database: db}
	return &blockchain
}
