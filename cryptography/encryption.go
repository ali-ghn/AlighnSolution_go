package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

type IEncryptionHelper interface {
	Encrypt(data []byte) ([]byte, error)
	EncryptString(data string) ([]byte, error)
	Decrypt(encryptedData []byte) ([]byte, error)
	DecryptString(encryptedData string) ([]byte, error)
}

type EncryptionHelper struct {
	Key []byte
}

func NewEncryptionHelper(key string) EncryptionHelper {
	hashedKey := Hash([]byte(key))
	return EncryptionHelper{
		Key: hashedKey,
	}
}

func (eh EncryptionHelper) EncryptString(data string) ([]byte, error) {
	bData := []byte(data)
	c, err := aes.NewCipher(eh.Key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)

	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, bData, nil), nil
}

func (eh EncryptionHelper) Encrypt(data []byte) ([]byte, error) {
	c, err := aes.NewCipher(eh.Key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)

	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, data, nil), nil
}

func (eh EncryptionHelper) Decrypt(encryptedData []byte) ([]byte, error) {
	c, err := aes.NewCipher(eh.Key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)

	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(encryptedData) < nonceSize {
		return nil, fmt.Errorf("nonce size doesn't match")
	}

	nonce, encryptedData := encryptedData[:nonceSize], encryptedData[nonceSize:]

	data, err := gcm.Open(nil, nonce, encryptedData, nil)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (eh EncryptionHelper) DecryptString(encryptedData string) ([]byte, error) {
	bEncryptedData := []byte(encryptedData)
	c, err := aes.NewCipher(eh.Key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)

	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(bEncryptedData) < nonceSize {
		return nil, fmt.Errorf("nonce size doesn't match")
	}

	nonce, bEncryptedData := bEncryptedData[:nonceSize], bEncryptedData[nonceSize:]

	data, err := gcm.Open(nil, nonce, bEncryptedData, nil)

	if err != nil {
		return nil, err
	}

	return data, nil
}
