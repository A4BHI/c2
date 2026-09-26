package session

import (
	"crypto/ecdh"
	"crypto/rand"
	"log"
)

func CalculatePublicKey() (privatekey *ecdh.PrivateKey, publickey []byte) {
	curve := ecdh.X25519()
	serverpriv, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		log.Println("Error generating key ", err)
		return nil, nil
	}

	return serverpriv, serverpriv.PublicKey().Bytes()

}
func calculateSessionKey(pubkey []byte) []byte {
	curve := ecdh.X25519()
	sessionkey, err := curve.NewPrivateKey(pubkey)
	if err != nil {
		log.Println("Error generating private key: ", err)
		return nil
	}
	return sessionkey.Bytes()
}
