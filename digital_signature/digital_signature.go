package digital_signature

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
	println("dcd", publKeyHex, decodedData)

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

func CreateDigitalSignature(message []byte, privKey *ecdsa.PrivateKey) ([]byte, error) {
	hash := sha256.Sum256(message)
	signature, err := ecdsa.SignASN1(rand.Reader, privKey, hash[:])

	if err != nil {
		return []byte{}, err
	}

	return signature, nil
}

func VerifyDigitalSignature(message []byte, publKey *ecdsa.PublicKey, signature []byte) bool {
	hash := sha256.Sum256((message))
	isValid := ecdsa.VerifyASN1(publKey, hash[:], signature)

	return isValid
}
