package pkg

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"net"
	"os"
	"time"
)

const (
	certPath = "cmd\\shortener\\srv.crt"
	keyPath  = "cmd\\shortener\\srv.key"
)

var cert = &x509.Certificate{
	// указываем уникальный номер сертификата
	SerialNumber: big.NewInt(1444),
	// заполняем базовую информацию о владельце сертификата
	Subject: pkix.Name{
		Organization: []string{"Shortener"},
		Country:      []string{"RU"},
	},
	// разрешаем использование сертификата для 127.0.0.1 и ::1
	IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
	// сертификат верен, начиная со времени создания
	NotBefore: time.Now(),
	// время жизни сертификата — 10 лет
	NotAfter:     time.Now().AddDate(10, 0, 0),
	SubjectKeyId: []byte{1, 2, 3, 4, 6},
	// устанавливаем использование ключа для цифровой подписи, а также клиентской и серверной авторизации
	ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
	KeyUsage:    x509.KeyUsageDigitalSignature,
}

//  GetCertX509Files returns pathes for x509 cert and key files for server SSL
//  If files not exsits, generates new.
func GetCertX509Files() (string, string, error) {
	filesExists, err := certFilesExists()
	if err != nil {
		return "", "", err
	}

	if !filesExists {
		if err := genCertX509Files(); err != nil {
			return "", "", err
		}
	}
	return certPath, keyPath, nil
}

//  genCertX509Files Generates x509 cert and key files for shortener server.
func genCertX509Files() error {
	// создаём новый приватный RSA-ключ длиной 4096 бит
	// обратите внимание, что для генерации ключа и сертификата используется rand.Reader в качестве источника случайных данных
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}

	// создаём x.509-сертификат
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		return err
	}

	// кодируем сертификат и ключ в формате PEM, который используется для хранения и обмена криптографическими ключами
	var certPEM bytes.Buffer
	if err := pem.Encode(&certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}); err != nil {
		return err
	}

	if err := writeFile(certPEM, certPath); err != nil {
		return err
	}

	var privateKeyPEM bytes.Buffer
	if err := pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}); err != nil {
		return err
	}

	if err := writeFile(privateKeyPEM, keyPath); err != nil {
		return err
	}

	return nil
}

func writeFile(buf bytes.Buffer, filename string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	wr := bufio.NewWriter(file)
	if _, err := wr.Write(buf.Bytes()); err != nil {
		return err
	}
	return wr.Flush()
}

func certFilesExists() (bool, error) {
	crtExist, err := filesExist(certPath)
	if err != nil {
		return false, err
	}

	keyExist, err := filesExist(keyPath)
	if err != nil {
		return false, err
	}

	return (crtExist && keyExist), nil
}

func filesExist(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
