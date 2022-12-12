package main

import (
	"context"
	"fmt"
	"os"
	"time"
)

func main() {

	confDB := Db{
		Host:            "127.0.0.1",
		Port:            "5432",
		Name:            "block",
		User:            "postgres",
		Password:        "pass",
		MaxOpenConns:    25,
		MaxIdleConns:    10,
		ConnMaxLifeTime: int64(5),
	}
	db, err := New(&confDB)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// create a new blockchain instance with a mining difficulty of 2
	blockchain := CreateBlockchain(2, db)

	// User
	user, err := NewUser()
	if err != nil {
		return
	}
	have := blockchain.isValid(user)
	if !have {
		os.Exit(1)
	}
	// record transactions on the blockchain for Alice, Bob, and John
	blockchain.addBlock(user, 5)
	blockchain.addBlock(user, 2)
	blockchain.addBlock(user, 10)
	db.Close()
	// check if the blockchain is valid; expecting true
	fmt.Println(blockchain.isValid(user))
}
func CreateBlockchain(difficulty int, db *Postgres) Blockchain {
	bloks, err := db.GetBlock(context.Background())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	client := NewClient("http://89.108.115.118")
	if len(bloks) == 0 {
		genesisBlock := Block{
			timestamp: time.Now(),
		}
		fmt.Println("nul")
		return Blockchain{
			genesisBlock,
			[]Block{genesisBlock},
			difficulty,
			db,
			client,
		}
	} else {
		genesisBlock := bloks[0]
		return Blockchain{
			genesisBlock,
			bloks,
			difficulty,
			db,
			client,
		}
	}

}
