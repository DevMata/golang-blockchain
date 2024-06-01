package blockchain

import (
	"encoding/hex"
	"fmt"
	"github.com/dgraph-io/badger/v4"
	"os"
	"runtime"
)

const (
	dbPath      = "./db/blocks"
	dbFile      = "./db/blocks/MANIFEST"
	genesisData = "Primera Transacción - Génesis"
)

// Blockchain y sus métodos

// BlockChain es la estructura para la cadena
// contiene el hash del último bloque insertado
// y una referencia a la BDD
type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

// InitBlockChain nos permite instaciar una nueva cadena
// address es la dirección del usuario, en nuestro caso pasaremos un nombre
func InitBlockChain(address string) *BlockChain {
	// validamos si la cadena ya existe
	if doesDBExist() {
		fmt.Println("La cadena ya existe")
		runtime.Goexit()
	}

	var lastHash []byte

	// conectamos la BDD
	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath
	opts.Logger = nil

	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		cbTx := CoinbaseTx(address, genesisData)
		genesis := Genesis(cbTx)
		fmt.Println("Bloque Génesis creado")
		err = txn.Set(genesis.Hash, genesis.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), genesis.Hash)

		lastHash = genesis.Hash

		return err
	})

	Handle(err)

	blockchain := BlockChain{
		LastHash: lastHash, Database: db}
	return &blockchain
}

// ContinueBlockChain nos permite cargar la cadena ya existente
func ContinueBlockChain() *BlockChain {
	// validamos si la cadena no existe
	if doesDBExist() == false {
		fmt.Println("Error: La cadena no existe, creála.")
		runtime.Goexit()
	}

	var lastHash []byte

	// conectamos la BDD
	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath
	opts.Logger = nil

	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		return err
	})
	Handle(err)

	blockchain := BlockChain{
		LastHash: lastHash,
		Database: db,
	}
	return &blockchain
}

// AddBlock nos permite insertar un bloque nuevo a la cadena
// espera un arreglo de transacciones
func (blockchain *BlockChain) AddBlock(transactions []*Transaction) {
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
	newBlock := CreateBlock(transactions, lastHash)

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

// FindUnspentTransactions nos permite encontrar aquellas transacciones cuyos fondos no se han gastado
// las filtraremos por address - es decir el nombre de un usuario
func (blockchain *BlockChain) FindUnspentTransactions(address string) []Transaction {
	var unspentTxs []Transaction

	spentTXOs := make(map[string][]int)

	iter := blockchain.Iterator()

	// iteramos los bloques en la cadena
	for {
		block := iter.Next()

		// iteramos las transacciones en el bloque
		for _, tx := range block.Transactions {
			txId := hex.EncodeToString(tx.ID)

		Outputs: // iteramos las salidas en la transacción
			for outIdx, out := range tx.Outputs {
				if spentTXOs[txId] != nil {
					for _, spentOut := range spentTXOs[txId] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}

				// el usuario puede ver la transacción porque fue el emisor
				if out.CanBeUnlocked(address) {
					unspentTxs = append(unspentTxs, *tx)
				}
			}

			// validamos si es una transacción monetaria
			if !tx.IsCoinBase() {
				// iteramos las entradas en la transacción
				for _, in := range tx.Inputs {
					// el usuario puede ver la transacción porque fue el destinatario
					if in.CanUnlock(address) {
						inTxID := hex.EncodeToString(in.ID)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Out)
					}
				}
			}
		}

		// si es el bloque Génesis
		if len(block.PrevHash) == 0 {
			break
		}
	}

	return unspentTxs
}

// FindUTXO nos devuelve un arreglo con transacciones de salida para un usuario
func (blockchain *BlockChain) FindUTXO(address string) []TxOutput {
	var UTXOs []TxOutput
	unspentTransactions := blockchain.FindUnspentTransactions(address)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Outputs {
			if out.CanBeUnlocked(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

// FindSpendableOutputs nos devuelve el monto total de los recursos disponibles para el usuario
// también un diccionario de transacciones y sus respectivos montos
func (blockchain *BlockChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unspentTxs := blockchain.FindUnspentTransactions(address)
	accumulated := 0

	// iteramos las transacciones con fondos no gastados
Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)

		// iteramos las salidas en la transacción
		for outIdx, out := range tx.Outputs {
			// validamos que el monto acumulado aún no nos permite gastar lo deseado
			if out.CanBeUnlocked(address) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[txID] = append(unspentOuts[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOuts
}

// doesDBExist nos deja saber si la BDD ya existe
func doesDBExist() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

// HandleDBClose permite cerrar la BDD de manera controlada
func HandleDBClose(Database *badger.DB) {
	err := Database.Close()
	Handle(err)
}
