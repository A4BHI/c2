package authentication

import (
	"crypto/rand"
	"encoding/hex"
	"log"
)

func CreateChallenge() ([]byte, error) {
	challengeBytes := make([]byte, 32)
	_, err := rand.Read(challengeBytes)
	if err != nil {
		log.Println("Error creating challenge : ", err)
		return nil, err
	}
	var challenge []byte
	hex.Encode(challenge, challengeBytes)

	return challenge, err
}
