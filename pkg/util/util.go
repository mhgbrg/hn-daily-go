package util

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/pkg/errors"
)

func RandomHexString(length int) (string, error) {
	b := make([]byte, length/2)
	_, err := rand.Read(b)
	if err != nil {
		return "", errors.Wrapf(err, "failed to read %d bytes from rand", len(b))
	}
	return hex.EncodeToString(b), nil
}
