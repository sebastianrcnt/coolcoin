package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
)

func GenerateKey() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	privateKey, err := ecdsa.GenerateKey(
		elliptic.P256(),
		rand.Reader,
	)

	if err != nil {
		return nil, nil, err
	}

	return privateKey, &privateKey.PublicKey, nil
}

func HexifyPrivateKey(privKey *ecdsa.PrivateKey) (string, error) {
	encoded, err := x509.MarshalECPrivateKey(privKey)
	if err != nil {
		return "", err
	}

	var hexEncoded []byte = make([]byte, hex.EncodedLen(len(encoded)))
	hex.Encode(hexEncoded, encoded)

	return string(hexEncoded), nil
}

func DecodePrivateKeyHex(privKeyHex string) (*ecdsa.PrivateKey, error) {
	decodedData, err := hex.DecodeString(privKeyHex)

	if err != nil {
		return nil, err
	}

	privKey, err := x509.ParseECPrivateKey(decodedData)

	if err != nil {
		return nil, err
	}

	return privKey, nil
}

func HexifyPublicKey(publKey *ecdsa.PublicKey) (string, error) {
	encoded, err := x509.MarshalPKIXPublicKey(publKey)
	if err != nil {
		return "", nil
	}

	var hexEncoded []byte = make([]byte, hex.EncodedLen(len(encoded)))
	hex.Encode(hexEncoded, encoded)

	return string(hexEncoded), nil
}

func DecodePublicKeyHex(publKeyHex string) (*ecdsa.PublicKey, error) {
	decodedData, err := hex.DecodeString(publKeyHex)

	if err != nil {
		return nil, err
	}

	var publKeyRaw interface{}

	publKeyRaw, err = x509.ParsePKIXPublicKey(decodedData)

	if err != nil {
		return nil, err
	}

	publKey := publKeyRaw.(*ecdsa.PublicKey)

	return publKey, nil
}

func EncodeAddress(hexEncodedPublicKey string) (string, error) {
	hash := sha256.Sum256([]byte(hexEncodedPublicKey))
	return hex.EncodeToString(hash[:]), nil // [:]는 array -> slice 로 바꿔준다.
}
