// Модуль crypto отвечает за хэширование и шифрование данные.
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
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

// EncryptData разшифровывает данные
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
