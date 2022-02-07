package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

func sha256Of(input string) []byte {
	result := sha256.New()
	result.Write([]byte(input))

	return result.Sum(nil)
}

func GenerateShortLink(srcLink string, salt string) string {
	urlHashBytes := sha256Of(srcLink + salt)
	base64Encoded := base64.StdEncoding.EncodeToString(urlHashBytes)
	slashReplaced := strings.ReplaceAll(base64Encoded[:8], "/", "+")

	return slashReplaced
}

func GenUserID() (string, error) {
	rand, err := GenerateRandom(16)
	if err != nil {
		return "", fmt.Errorf("ошибка генерации ID пользователя:%w", err)
	}

	hexToken := make([]byte, hex.EncodedLen(len(rand)))
	hex.Encode(hexToken, rand)

	return string(ToHEX(rand)[:16]), nil
}

func ToHEX(src []byte) []byte {
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)

	return dst
}

func GenerateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
