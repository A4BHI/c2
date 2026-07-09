package authentication

import (
	"crypto/rand"
	"encoding/hex"
	"log"
)

func CreateChallenge() (string, error) {
	challengeBytes := make([]byte, 32)
	_, err := rand.Read(challengeBytes)
	if err != nil {
		log.Println("Error creating challenge : ", err)
		return nil, err
	}

	return hex.EncodeToString(challengeBytes), err
}

func CheckChallenge(registerkey string, challenge string, botresponse string) {

}
