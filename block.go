package main

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	_ "crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// data,hash,previous_hash,timestamp,pow
type Block struct {
	data         map[string]interface{} `db:"data"`
	hash         string                 `db:"hash"`
	previousHash string                 `db:"previous_hash"`
	timestamp    time.Time              `db:"timestamp"`
	pow          int                    `db:"pow"`
	sign         string
}

type Blockchain struct {
	genesisBlock Block
	chain        []Block
	difficulty   int
	db           *Postgres
	client       Client
}

func (b Block) calculateHash() string {
	data, _ := json.Marshal(b.data)
	blockData := b.previousHash + string(data) + b.timestamp.String() + strconv.Itoa(b.pow)
	hashed := sha256.Sum256([]byte(blockData))
	return fmt.Sprintf("%x", hashed)
}

func (b *Block) mine(privateKey *rsa.PrivateKey, difficulty int) {

	for !strings.HasPrefix((b.hash), strings.Repeat("0", difficulty)) {
		b.pow++
		b.hash = b.calculateHash()
	}

}

func (b *Block) Verify(pub rsa.PublicKey) error {
	return rsa.VerifyPKCS1v15(&pub, crypto.SHA256, []byte(b.hash), []byte(b.sign))
}
func (b *Blockchain) addBlock(user User, amount float64) {

	blockData := map[string]interface{}{
		"amount": amount,
	}
	lastBlock := b.chain[len(b.chain)-1]
	newBlock := Block{
		data:         blockData,
		previousHash: lastBlock.hash,
		timestamp:    time.Now().UTC(),
		pow:          0,
	}
	newBlock.mine(user.Private, b.difficulty)
	data, _ := json.Marshal(blockData)
	sign, signTs := b.client.GetSign(newBlock.hash)
	if err := b.db.SetBlock(context.Background(), string(data), newBlock.hash, newBlock.previousHash, newBlock.timestamp, newBlock.pow, sign, signTs); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	b.chain = append(b.chain, newBlock)
}

func (b Blockchain) isValid(user User) bool {

	for i := range b.chain[1:] {
		previousBlock := b.chain[i]
		currentBlock := b.chain[i+1]
		if currentBlock.hash != currentBlock.calculateHash() {
			fmt.Println(currentBlock.calculateHash())
			fmt.Println(currentBlock.data)
			return false
		}
		if currentBlock.previousHash != previousBlock.hash {
			fmt.Println("test2")
			return false
		}

	}
	return true
}
