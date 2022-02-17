package api

import (
	"crypto/aes"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

const (
	key           = "5885c4a300814d5c"
	uuidBlockSize = 16
)

type AuthCrypt struct {
}

// New token crypto
func NewAuthCrypt() *AuthCrypt {
	return &AuthCrypt{}
}

// Encode uuid to token
func (c *AuthCrypt) EncodeUUID(id uuid.UUID) (string, error) {
	if id == uuid.Nil {
		return "", errors.New("ошибка шифрации: входной UUID nil")
	}

	aesblock, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", fmt.Errorf("ошибка шифрации: %w", err)
	}

	dst := make([]byte, uuidBlockSize)
	aesblock.Encrypt(dst, id[:])

	dstHEX := make([]byte, hex.EncodedLen(len(dst)))
	hex.Encode(dstHEX, dst)

	return string(dstHEX), nil
}

// Decode uuid to token
func (c *AuthCrypt) DecodeToken(token string) (uuid.UUID, error) {
	byteToken, err := hex.DecodeString(token)
	if err != nil {
		return uuid.Nil, fmt.Errorf("ошибка дешифрации: %w", err)
	}

	aesblock, err := aes.NewCipher([]byte(key))
	if err != nil {
		return uuid.Nil, fmt.Errorf("ошибка дешифрации: %w", err)
	}

	dst := make([]byte, aes.BlockSize)
	aesblock.Decrypt(dst, byteToken)
	id, err := uuid.FromBytes(dst)
	if err != nil {
		return uuid.Nil, fmt.Errorf("ошибка дешифрации: %w", err)
	}

	return id, nil
}
