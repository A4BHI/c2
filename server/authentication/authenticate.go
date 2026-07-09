package authentication

import (
	"crypto/rand"
	"crypto/sha256"
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

	h := sha256.New()
	h.Write([]byte(registerkey))
	h.Write([]byte(challenge))

	hash := h.Sum(nil)

	hashHex := hex.EncodeToString(hash)

	if hashHex == botresponse {

	}

}
