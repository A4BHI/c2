package register

import (
	"c2/server/models"
	"crypto/rand"
	"encoding/hex"
	"io"

	"github.com/google/uuid"
)

func GenerateBotCredentials() (models.BotCreds, error) {
	buf := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		return models.BotCreds{}, err
	}
	secretkey := hex.EncodeToString(buf)
	id := uuid.NewString()

	return models.BotCreds{
		ID:        id,
		SecretKey: secretkey,
	}, nil

}

func CreateChallenge() string { return "" }
