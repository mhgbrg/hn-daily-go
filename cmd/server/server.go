package main

import (
	"encoding/hex"
	"log"
	"os"
	"strconv"

	"github.com/mhgbrg/hndaily/pkg/web"
	"github.com/pkg/errors"
)

func main() {
	hostname := os.Getenv("HOSTNAME")
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal(errors.Wrapf(err, "failed to convert environment variable PORT=%s to int", os.Getenv("PORT")))
	}

	dbURL := os.Getenv("DATABASE_URL")
	if len(dbURL) == 0 {
		log.Fatal(errors.New("environment variable DATABASE_URL not set"))
	}

	hashKey, err := hex.DecodeString(os.Getenv("HASH_KEY"))
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to decode environment variable HASH_KEY"))
	}
	if len(hashKey) != 32 {
		log.Fatal(errors.Errorf("expected len(HASH_KEY) == 32, got %d", len(hashKey)))
	}

	encryptionKey, err := hex.DecodeString(os.Getenv("ENCRYPTION_KEY"))
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to decode environment variable ENCRYPTION_KEY"))
	}
	if len(encryptionKey) != 32 {
		log.Fatal(errors.Errorf("expected len(ENCRYPTION_KEY) == 32, got %d", len(encryptionKey)))
	}

	config := web.Config{
		Hostname:    hostname,
		Port:        port,
		DatabaseURL: dbURL,
		CryptoKeys: web.CryptoKeys{
			HashKey:       hashKey,
			EncryptionKey: encryptionKey,
		},
	}

	log.Fatal(web.StartServer(config))
}
