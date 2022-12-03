package main

import (
	"fmt"
	"time"
)

func main() {
	// create a new blockchain instance with a mining difficulty of 2
	blockchain := CreateBlockchain(2)
	// User
	user, err := NewUser()
	if err != nil {
		return
	}

	// record transactions on the blockchain for Alice, Bob, and John
	blockchain.addBlock(user, 5)
	blockchain.addBlock(user, 2)
	blockchain.addBlock(user, 10)

	// check if the blockchain is valid; expecting true
	fmt.Println(blockchain.isValid(user))
}
func CreateBlockchain(difficulty int) Blockchain {
	genesisBlock := Block{
		timestamp: time.Now(),
	}
	return Blockchain{
		genesisBlock,
		[]Block{genesisBlock},
		difficulty,
	}
}
