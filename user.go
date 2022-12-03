package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
)

type User struct {
	Public  *rsa.PublicKey
	Private *rsa.PrivateKey
}

func NewUser() (User, error) {

	singkey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Print("unable to read private key")
		return User{}, err
	}

	public := singkey.PublicKey

	return User{
		Private: singkey,
		Public:  &public,
	}, nil

}
