package cryptography

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestEncryption(t *testing.T) {
	eh := NewEncryptionHelper("Some key")
	data := []byte("Hello there")
	encData, err := eh.Encrypt(data)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Printf("%x\n", encData)
}

func TestDecryption(t *testing.T) {
	eh := NewEncryptionHelper("Some key")
	data := []byte("Hello there")
	encData, _ := eh.Encrypt(data)
	hexData := fmt.Sprintf("%x", encData)
	strData, err := hex.DecodeString(hexData)
	if err != nil {
		t.Errorf(err.Error())
	}
	decData, err := eh.Decrypt(strData)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Printf("%v\n", string(decData))
}
