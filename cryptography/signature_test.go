package cryptography

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"
)

func TestSignAndVerify(t *testing.T) {
	ec := NewEllipticCurve(elliptic.P256())
	_, _, err := ec.GenerateKeys()
	if err != nil {
		t.Errorf(err.Error())
	}
	data := []byte("Hello there")
	signedData, err := ecdsa.SignASN1(rand.Reader, ec.privateKey, data)
	if err != nil {
		t.Errorf(err.Error())
	}
	verify, err := ec.Verify(signedData, data)
	if err != nil {
		t.Errorf(err.Error())
	}
	if !verify {
		t.Errorf(err.Error())
	}
}
