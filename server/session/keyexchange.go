package session

import (
	"crypto/ecdh"
	"crypto/rand"
	"log"
)

func CalculateKeys() []byte {
	curve := ecdh.X25519()
	serverpriv, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		log.Println("Error generating key ", err)
		return nil
	}

	return serverpriv.Bytes()

}
