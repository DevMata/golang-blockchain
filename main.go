package main

import (
	"flag"
	"fmt"
	"github.com/dgraph-io/badger/v4"
	"os"
	"runtime"
	"strconv"

	"github.com/devmata/golang-blockchain/blockchain"
)

// CommandLine y sus métodos

// CommandLine es una estructura muy básica para tener un CLI
// contiene una referencia a la cadena
type CommandLine struct {
	blockchain *blockchain.BlockChain
}

// printUsage muestra en pantalla el manual para uso
func (cli *CommandLine) printUsage() {
	fmt.Println("Uso: ")
	fmt.Println(" add -block BLOCK_DATA - agrega un bloque a la cadena")
	fmt.Println(" print - muestra en pantalla el contenido de la cadena")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) addBlock(data string) {
	cli.blockchain.AddBlock(data)
	fmt.Println("Bloque agregado")
}

// printChain recorre la cadena e imprime en pantalla
// su contenido
func (cli *CommandLine) printChain() {
	iterator := cli.blockchain.Iterator()

	for {
		block := iterator.Next()

		fmt.Printf("\nHash de bloque anterior: %x\n", block.PrevHash)
		fmt.Printf("Info en el bloque: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Nonce: %d\n", block.Nonce)

		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevHash) == 0 {
			break
		}
	}
}

// run contiene la lógica para usar la CLI
// valida que sus diferentes funciones hayan sido llamadas correctamente
func (cli *CommandLine) run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	addBlockData := addBlockCmd.String("block", "", "Block data")

	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			runtime.Goexit()
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}

func main() {
	defer os.Exit(0)
	chain := blockchain.InitBlockChain()
	defer func(Database *badger.DB) {
		err := Database.Close()
		blockchain.Handle(err)
	}(chain.Database)

	cli := CommandLine{blockchain: chain}
	cli.run()
}
