package app

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/pkg/errors"
)

func generateAppGlobalValuesandAppend(globalValues map[string]interface{}) error {
	token, err := randomTokenGeneration()
	if err != nil {
		return err
	}
	globalValues["NatsToken"] = token
	return nil
}

func randomTokenGeneration() (string, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", errors.WithMessage(err, "error while generating random key")
	}
	randomString := base64.RawURLEncoding.EncodeToString(randomBytes)[:32]
	return randomString, nil
}
