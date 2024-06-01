package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

// Difficulty es el parámetro que nos permite configurar
// la dificultad del reto para la prueba de trabajo (PoW)
const Difficulty = 16

// ProofOfWork y sus métodos

// ProofOfWork es una estructura básica para la prueba de trabajo
// básicamente contiene una referencia a un bloque
// y una referencia a un entero que será el Nonce
type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

// NewProof instancia un ProofOfWork y hace desplazamiento a la izquierda
// el target no es más que un uint igualado a 1 que ha sido
// desplazado a la izquierda de manera que los primeros bits estén en 0
func NewProof(b *Block) *ProofOfWork {
	// desplazamiento a la izquierda considerando la dificultad
	// buscamos que "n" bits al principio del número estén en 0
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))

	pow := &ProofOfWork{
		Block:  b,
		Target: target,
	}

	return pow
}

// InitData nos permite obtener en Byte la combinación de info que contiene un bloque
// es decir, el hash al bloque anterior, su propia data, el Nonce probado en ese momento
// e incluso la dificultad que se usaba en ese momento
func (pow *ProofOfWork) InitData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.HashTransactions(),
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{})
	return data
}


// Run tiene la lógica para encontrar el nonce
// Es tan difícil de calcular como Difficulty hayamos configurado
func (pow *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	nonce := 0

	// iteramos desde 0 hasta encontrar el nonce
	for nonce < math.MaxInt64 {
		// comprobamos que "n" bits del principio del hash
		// estén en 0
		data := pow.InitData(nonce)
		hash = sha256.Sum256(data)

		// como un AND lógico
		intHash.SetBytes(hash[:])

		// encontramos el nonce o seguimos
		// básicamente comparamos si el hash de nuestra data
		// es menor o igual al target que tenemos
		if intHash.Cmp(pow.Target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Println()

	return nonce, hash[:]
}

// Validate nos permitirá ver en pantalla que el nonce es correcto
func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int

	// obtenemos la info de nuestro bloque en Byte
	data := pow.InitData(pow.Block.Nonce)

	// lo hasheamos
	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	// comparamos si el hash de nuestra data
	// es menor o igual al target que tenemos
	return intHash.Cmp(pow.Target) == -1
}

func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}