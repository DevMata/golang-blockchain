package main

import (
	"flag"
	"fmt"
	"github.com/devmata/golang-blockchain/blockchain"
	"log"
	"os"
	"runtime"
	"strconv"
)

// CommandLine y sus métodos

// CommandLine es una estructura muy básica para tener un CLI
type CommandLine struct {
}

// printUsage muestra en pantalla el manual para uso
func (cli *CommandLine) printUsage() {
	fmt.Println("Uso: ")
	fmt.Println(" getbalance -address ADDRESS - obtiene el balance para dirección address")
	fmt.Println(" createblockchain -address ADDRESS - crea la cadena y envía el reward para el bloque génesis")
	fmt.Println(" printchain - imprime los bloques en la cadena")
	fmt.Println(" send -from FROM -to TO -amount AMOUNT - envía una cantidad de divisas")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) printChain() {
	chain := blockchain.ContinueBlockChain()
	defer blockchain.HandleDBClose(chain.Database)

	iterator := chain.Iterator()

	for {
		block := iterator.Next()

		fmt.Printf("\nHash de bloque anterior: %x\n", block.PrevHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n\n", strconv.FormatBool(pow.Validate()))

		if len(block.PrevHash) == 0 {
			break
		}
	}
}

func (cli *CommandLine) createBlockChain(address string) {
	chain := blockchain.InitBlockChain(address)
	blockchain.HandleDBClose(chain.Database)
	fmt.Println("¡Blockchain creada!")
}

func (cli *CommandLine) getBalance(address string) {
	chain := blockchain.ContinueBlockChain()
	defer blockchain.HandleDBClose(chain.Database)

	balance := 0
	UTXOs := chain.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s: %d\n", address, balance)
}

func (cli *CommandLine) send(from, to string, amount int) {
	chain := blockchain.ContinueBlockChain()
	defer blockchain.HandleDBClose(chain.Database)

	tx := blockchain.NewTransaction(from, to, amount, chain)
	chain.AddBlock([]*blockchain.Transaction{tx})
	fmt.Println("Éxito")
}

// run contiene la lógica para usar la CLI
// valida que sus diferentes funciones hayan sido llamadas correctamente
func (cli *CommandLine) run() {
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("getBalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "La dirección(usuario) para ver su balance")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "La dirección(usuario) para enviar el reward del bloque Génesis")
	sendFrom := sendCmd.String("from", "", "La wallet(usuario) origen")
	sendTo := sendCmd.String("to", "", "La wallet(usuario) destino")
	sendAmount := sendCmd.Int("amount", 0, "Monto a enviar")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			runtime.Goexit()
		}
		cli.getBalance(*getBalanceAddress)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			runtime.Goexit()
		}
		cli.createBlockChain(*createBlockchainAddress)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}

		cli.send(*sendFrom, *sendTo, *sendAmount)
	}
}

func main() {
	defer os.Exit(0)
	cli := CommandLine{}
	cli.run()
}
