package cryptography

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
)

type EllipticCurve struct {
	pubKeyCurve elliptic.Curve
	privateKey  *ecdsa.PrivateKey
	publicKey   *ecdsa.PublicKey
}

func NewEllipticCurve(curve elliptic.Curve) *EllipticCurve {
	return &EllipticCurve{
		pubKeyCurve: curve,
		privateKey:  new(ecdsa.PrivateKey),
	}
}

func (ec *EllipticCurve) GenerateKeys() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {

	var err error
	privKey, err := ecdsa.GenerateKey(ec.pubKeyCurve, rand.Reader)

	if err == nil {
		ec.privateKey = privKey
		ec.publicKey = &privKey.PublicKey
	}

	return ec.privateKey, ec.publicKey, err
}

func (ec *EllipticCurve) EncodePrivate(privKey *ecdsa.PrivateKey) (string, error) {

	encoded, err := x509.MarshalECPrivateKey(privKey)

	if err != nil {
		return "", err
	}
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: encoded})

	return string(pemEncoded), nil
}

func (ec *EllipticCurve) EncodePublic(pubKey *ecdsa.PublicKey) (string, error) {

	encoded, err := x509.MarshalPKIXPublicKey(pubKey)

	if err != nil {
		return "", err
	}
	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: encoded})

	return string(pemEncodedPub), nil
}

func (ec *EllipticCurve) DecodePrivate(pemEncodedPriv string) (*ecdsa.PrivateKey, error) {
	blockPriv, _ := pem.Decode([]byte(pemEncodedPriv))

	x509EncodedPriv := blockPriv.Bytes

	privateKey, err := x509.ParseECPrivateKey(x509EncodedPriv)

	return privateKey, err
}

func (ec *EllipticCurve) DecodePublic(pemEncodedPub string) (*ecdsa.PublicKey, error) {
	blockPub, _ := pem.Decode([]byte(pemEncodedPub))

	x509EncodedPub := blockPub.Bytes

	genericPublicKey, err := x509.ParsePKIXPublicKey(x509EncodedPub)
	publicKey := genericPublicKey.(*ecdsa.PublicKey)

	return publicKey, err
}

func (ec *EllipticCurve) Sign(hash []byte) ([]byte, error) {
	sd, err := ecdsa.SignASN1(rand.Reader, ec.privateKey, hash)
	if err != nil {
		return nil, err
	}
	return sd, nil
}

func (ec *EllipticCurve) Verify(signature []byte, hash []byte) (bool, error) {
	result := ecdsa.VerifyASN1(ec.publicKey, hash, signature)
	return result, nil
}
