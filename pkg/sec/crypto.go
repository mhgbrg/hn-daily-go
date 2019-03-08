package sec

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"

	"golang.org/x/crypto/sha3"

	"github.com/pkg/errors"
)

type Crypter struct {
	HashKey  []byte
	CryptKey []byte
}

func (crypter *Crypter) EncryptInt(i int) (string, error) {
	encoded, err := encodeInt(i)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(encoded))
	iv := ciphertext[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	block, err := aes.NewCipher(crypter.CryptKey)
	if err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], encoded)

	str := fmt.Sprintf("%x", ciphertext)
	return str, nil
}

func (crypter *Crypter) DecryptInt(str string) (int, error) {
	ciphertext, err := hex.DecodeString(str)
	if err != nil {
		return 0, err
	}

	block, err := aes.NewCipher(crypter.CryptKey)
	if err != nil {
		return 0, err
	}

	if len(ciphertext) < aes.BlockSize {
		return 0, errors.Errorf("ciphertext is shorter than block size. length=%d, blockSize=%d", len(ciphertext), aes.BlockSize)
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	if len(ciphertext)%aes.BlockSize != 0 {
		return 0, errors.Errorf("ciphertext length is not a multiple of block size. length=%d, blockSize=%d", len(ciphertext), aes.BlockSize)
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	return decodeInt(ciphertext)
}

func (crypter *Crypter) HMACInt(i int) (string, error) {
	mac := hmac.New(sha3.New256, crypter.HashKey)
	encoded, err := encodeInt(i)
	if err != nil {
		return "", err
	}
	mac.Write(encoded)
	hash := mac.Sum(nil)
	str := hex.EncodeToString(hash)
	return str, nil
}

func encodeInt(i int) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, int32(i))
	if err != nil {
		return []byte{}, err
	}
	encoded := buf.Bytes()

	// Pad with zeroes so that plaintext is a full block in length.
	padded := encoded[:aes.BlockSize]
	return padded, nil
}

func decodeInt(encoded []byte) (int, error) {
	var i int32
	buf := bytes.NewReader(encoded)
	err := binary.Read(buf, binary.LittleEndian, &i)
	if err != nil {
		return 0, err
	}
	return int(i), nil
}
