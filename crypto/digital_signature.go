package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
)

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
