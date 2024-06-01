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

func doesDBExist() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

func InitBlockChain(address string) *BlockChain {
	if doesDBExist() {
		fmt.Println("La cadena ya existe")
		runtime.Goexit()
	}

	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath
	//	TODO: add settings for the db logger

	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		cbTx := CoinbaseTx(address, genesisData)
		genesis := Genesis(cbTx)
		fmt.Println("Genesis creada")
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

func ContinueBlockChain(address string) *BlockChain {
	if doesDBExist() == false {
		fmt.Println("¡La cadena no existe, crea una!")
		runtime.Goexit()
	}

	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath
	//	TODO: add settings for the db logger

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
// espera data - que puede ser cualquier cosa que queramos insertar
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

func (blockchain *BlockChain) FindUnspentTransactions(address string) []Transaction {
	var unspentTxs []Transaction

	spentTXOs := make(map[string][]int)

	iter := blockchain.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			txId := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Outputs {
				if spentTXOs[txId] != nil {
					for _, spentOut := range spentTXOs[txId] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				if out.CanBeUnlocked(address) {
					unspentTxs = append(unspentTxs, *tx)
				}
			}

			if !tx.IsCoinBase() {
				for _, in := range tx.Inputs {
					if in.CanUnlock(address) {
						inTxID := hex.EncodeToString(in.ID)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Out)
					}
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return unspentTxs
}

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

func (blockchain *BlockChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unspentTxs := blockchain.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Outputs {
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

// métodos

// HandleDBClose permite cerrar la BDD de manera controlada
func HandleDBClose(Database *badger.DB) {
	err := Database.Close()
	Handle(err)
}
