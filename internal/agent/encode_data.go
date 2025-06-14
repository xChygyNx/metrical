package agent

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

func parseRSAPublicKeyFromPemStr(pubPEM []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pubPEM)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	publicKey, ok := pub.(*rsa.PublicKey)
	if ok {
		return publicKey, nil
	}
	return nil, errors.New("key type is not RSA")
}

func encodeDataRSA(data []byte, rsaPublicKeyFile string) ([]byte, error) {
	pubPEM, err := os.ReadFile(rsaPublicKeyFile)
	if err != nil {
		return nil, fmt.Errorf("error in read file %s: %w\n", rsaPublicKeyFile, err)
	}

	publicKey, err := parseRSAPublicKeyFromPemStr(pubPEM)
	if err != nil {
		return nil, fmt.Errorf("error in parse RSA public key")
	}

	label := []byte("OAEP Encrypted")
	rng := rand.Reader
	fmt.Printf("Data length: %v\n", len(data))
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, publicKey, data, label)
	if err != nil {
		return nil, fmt.Errorf("error in cipher data: %w\n", err)
	}
	result := []byte(base64.StdEncoding.EncodeToString(ciphertext))
	return result, nil
}
