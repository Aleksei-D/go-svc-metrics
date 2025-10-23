// Модуль crypto отвечает за хэширование и шифрование данные.
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"hash"
	"os"
)

// GetHash возврашает хэш
func GetHash(key string, src []byte) []byte {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(src)
	return h.Sum(nil)
}

// EncryptData шифрует данные
func EncryptData(key string, src []byte) ([]byte, error) {
	cipherText, err := getCipherText(key)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, cipherText.NonceSize())
	dst := cipherText.Seal(nil, nonce, src, nil)
	return dst, nil
}

// DecryptData разшифровывает данные
func DecryptData(key string, src []byte) ([]byte, error) {
	cipherText, err := getCipherText(key)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, cipherText.NonceSize())
	dst, err := cipherText.Open(nil, nonce, src, nil)
	if err != nil {
		return nil, err
	}
	return dst, nil
}

func getCipherText(key string) (cipher.AEAD, error) {
	keyS := sha256.Sum256([]byte(key))
	aesblock, err := aes.NewCipher(keyS[:])
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, err
	}
	return aesgcm, nil
}

// ValidMAC сверяет хэш
func ValidMAC(message, messageMAC []byte, key string) bool {
	expectedMAC := GetHash(key, message)
	return hmac.Equal(messageMAC, expectedMAC)
}

// GetPrivateKey получает Приватный ключ из файла
func GetPrivateKey(filePath string) (*rsa.PrivateKey, error) {
	privateKeyPEM, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading private key file: %v", err)
	}

	block, _ := pem.Decode(privateKeyPEM)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode PEM block or block is not an RSA PRIVATE KEY")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

// GetPublickKey получает Публичный ключ из файла
func GetPublickKey(filePath string) (*rsa.PublicKey, error) {
	publicKeyPEM, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading private key file: %v", err)
	}

	block, _ := pem.Decode(publicKeyPEM)
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return nil, fmt.Errorf("failed to decode PEM block or block is not an RSA PUBLIC KEY")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}

	rrsaPubKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is not an RSA public key")
	}
	return rrsaPubKey, nil
}

// DecryptRSAData разшифровывает данные c помлщью rsa
func DecryptRSAData(privateKey *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	msgLen := len(ciphertext)
	step := privateKey.PublicKey.Size()
	var decryptedBytes []byte

	hash := sha256.New()

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		decryptedBlockBytes, err := rsa.DecryptOAEP(hash, rand.Reader, privateKey, ciphertext[start:finish], nil)
		if err != nil {
			return nil, err
		}

		decryptedBytes = append(decryptedBytes, decryptedBlockBytes...)
	}

	return decryptedBytes, nil
}

// EncryptRSAData зашифровывает данные c помлщью rsa
func EncryptRSAData(hash hash.Hash, publicKey *rsa.PublicKey, data []byte) ([]byte, error) {
	msgLen := len(data)
	step := publicKey.Size() - 2*hash.Size() - 2
	var encryptedBytes []byte

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		encryptedBlockBytes, err := rsa.EncryptOAEP(hash, rand.Reader, publicKey, data[start:finish], nil)
		if err != nil {
			return nil, err
		}

		encryptedBytes = append(encryptedBytes, encryptedBlockBytes...)
	}

	return encryptedBytes, nil
}
