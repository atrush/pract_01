package pkg

import (
	"crypto/aes"
	"fmt"
)

const (
	key       = "secretSECRETsecr"
	blockSize = aes.BlockSize
)

type Crypt struct {
}

func NewCrypt() *Crypt {
	return &Crypt{}
}

func (c *Crypt) EncodeToken(value string) ([]byte, error) {
	aesblock, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, fmt.Errorf("ошибка шифрации: %w", err)
	}
	dst := make([]byte, blockSize)
	aesblock.Encrypt(dst, []byte(value))

	return dst, nil
}

func (c *Crypt) DecodeToken(token []byte) ([]byte, error) {
	aesblock, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, fmt.Errorf("ошибка дешифрации: %w", err)
	}
	dst := make([]byte, aes.BlockSize)
	aesblock.Decrypt(dst, token)

	return dst, nil
}
