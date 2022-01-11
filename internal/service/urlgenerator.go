package service

import (
	"crypto/sha256"
	"encoding/base64"
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
