package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const (
	directoryPerm         = 0o750
	filePerm              = 0o600
	keyByteSize           = 4096
	privateRSAKeyFileName = "rsa_key"
	publicRSAKeyFileName  = "rsa_key.pub"
	rsaKeysDirName        = "rsa_keys"
)

func generateRsaKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, keyByteSize)
	if err != nil {
		return nil, nil, fmt.Errorf("error in generate rsa keys pair: %w", err)
	}
	return privateKey, &privateKey.PublicKey, nil
}

func exportRSAPrivateKeyAsPemStr(privateKey *rsa.PrivateKey) string {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privateKeyBytes,
		},
	)
	return string(privateKeyPem)
}

// Func parseRSAPrivateKeyFromPemStr(privateKeyPEM string) (*rsa.PrivateKey, error) {
//	block, _ := pem.Decode([]byte(privateKeyPEM))
//	if block == nil {
//		return nil, errors.New("failed to parse PEM block containing the key")
//	}
//
//	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
//	if err != nil {
//		return nil, fmt.Errorf("error in parse RSA Private key from bytes: %w", err)
//	}
//
//	return privateKey, nil
// }.

func exportRSAPublicKeyAsPemStr(pubkey *rsa.PublicKey) (string, error) {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(pubkey)
	if err != nil {
		return "", fmt.Errorf("error in export RSA public key in bytes: %w", err)
	}
	publicKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: publicKeyBytes,
		},
	)

	return string(publicKeyPEM), nil
}

func createRSAKeysDir() (string, error) {
	currentPath, err := os.Executable()
	if err != nil {
		return "", errors.New("error in define current path")
	}
	currentDir := filepath.Dir(currentPath)
	rsaKeysDir := filepath.Join(currentDir, rsaKeysDirName)
	if _, err = os.Stat(rsaKeysDir); errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(rsaKeysDir, directoryPerm)
		if err != nil {
			return "", fmt.Errorf("error in create rsa keys directory: %w", err)
		}
	}
	return rsaKeysDir, nil
}

func writeRSAKeys(privateKey string, publicKey string, keysDir string) (err error) {
	err = os.WriteFile(filepath.Join(keysDir, publicRSAKeyFileName), []byte(publicKey), filePerm)
	if err != nil {
		return fmt.Errorf("error in write public RSA key in file: %w", err)
	}
	err = os.WriteFile(filepath.Join(keysDir, privateRSAKeyFileName), []byte(privateKey), filePerm)
	if err != nil {
		return fmt.Errorf("error in write private RSA key in file: %w", err)
	}
	return nil
}

func main() {
	// создаём директорию для ключей
	keysDir, err := createRSAKeysDir()
	if err != nil {
		log.Fatal(err)
	}

	// проверяем наличие rsa ключей, если их нет, то создаем новые
	publicRSAKeyPath := filepath.Join(keysDir, publicRSAKeyFileName)
	privateRSAKeyPath := filepath.Join(keysDir, privateRSAKeyFileName)
	_, publicKeyErr := os.Stat(publicRSAKeyPath)
	_, privateKeyErr := os.Stat(privateRSAKeyPath)

	switch {
	case errors.Is(publicKeyErr, os.ErrNotExist) || errors.Is(privateKeyErr, os.ErrNotExist):
		fmt.Println("Create new RSA keys")
		// создаём новые приватный и публичный RSA-ключи
		privateKey, publicKey, err := generateRsaKeyPair()
		if err != nil {
			log.Fatal(err)
		}
		publicKeyString, err := exportRSAPublicKeyAsPemStr(publicKey)
		if err != nil {
			log.Fatal(err)
		}
		privateKeyString := exportRSAPrivateKeyAsPemStr(privateKey)

		err = writeRSAKeys(privateKeyString, publicKeyString, keysDir)
		if err != nil {
			log.Fatal(err)
		}
	case publicKeyErr != nil || privateKeyErr != nil:
		errorMsg := fmt.Sprintf("error in check public RSA key: %s\n"+
			"error in check private RSA key: %s\n", publicKeyErr.Error(), privateKeyErr.Error())
		log.Fatal(errorMsg)
	default:
		fmt.Println("All OK")
	}
}
