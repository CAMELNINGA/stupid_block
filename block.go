package main

import (
	"crypto"
	"crypto/rand"
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

type Block struct {
	priv         *rsa.PrivateKey
	pub          *rsa.PublicKey
	data         map[string]interface{}
	hash         [32]byte
	salt         *rsa.PSSOptions
	sign         []byte
	previousHash [32]byte
	timestamp    time.Time
	pow          int
}

type Blockchain struct {
	genesisBlock Block
	chain        []Block
	difficulty   int
}

func (b Block) signHash(privateKey *rsa.PrivateKey, hashed [32]byte) []byte {
	newhash := crypto.SHA256
	fmt.Println(hashed)
	sign, err := rsa.SignPKCS1v15(rand.Reader, privateKey, newhash, hashed[:])

	if err != nil {
		fmt.Printf("sign :%s", err)
		os.Exit(1)
	}
	return sign

}
func (b Block) calculateHash() [32]byte {
	data, _ := json.Marshal(b.data)
	Bdata := fmt.Sprintf("%x", data)
	blockData := fmt.Sprintf("%x", b.previousHash) + Bdata + b.timestamp.String() + strconv.Itoa(b.pow)
	hashed := sha256.Sum256([]byte(blockData))

	return hashed
}

func (b *Block) mine(privateKey *rsa.PrivateKey, difficulty int) {
	var hashed [32]byte
	for !strings.HasPrefix(fmt.Sprintf("%x", b.hash), strings.Repeat("0", difficulty)) {
		b.pow++
		hashed = b.calculateHash()
		b.hash = b.calculateHash()
	}
	b.sign = b.signHash(privateKey, hashed)
}

func (b *Block) decrypt() error {
	return rsa.VerifyPKCS1v15(b.pub, crypto.SHA256, b.hash[:], b.sign)
}

func (b *Blockchain) addBlock(user User, amount float64) {
	blockData := map[string]interface{}{
		"amount": amount,
	}
	lastBlock := b.chain[len(b.chain)-1]
	newBlock := Block{
		priv:         user.Private,
		data:         blockData,
		previousHash: lastBlock.hash,
		timestamp:    time.Now(),
		pow:          0,
	}
	newBlock.mine(user.Private, b.difficulty)
	b.chain = append(b.chain, newBlock)
}

func (b Blockchain) isValid(user User) bool {

	for i := range b.chain[1:] {
		previousBlock := b.chain[i]
		currentBlock := b.chain[i+1]
		if err := currentBlock.decrypt(); err != nil {
			fmt.Println(err)
			fmt.Println("Who are U? Verify Signature failed")
			os.Exit(1)
		}

		if currentBlock.hash != currentBlock.calculateHash() || currentBlock.previousHash != previousBlock.hash {
			return false
		}
	}
	return true
}
