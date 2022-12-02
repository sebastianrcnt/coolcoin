package crypto_test

import (
	crypto "coolcoin/crypto"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	privateKey, publicKey, err := crypto.GenerateKey()

	if !privateKey.PublicKey.Equal(publicKey) {
		t.Error("Private Key and Public Key don't match", err)
	}

	if err != nil {
		t.Error(err)
	}
}

func TestEncodeAndDecodePrivateKey(t *testing.T) {
	privateKey, _, _ := crypto.GenerateKey()
	encoded, err := crypto.HexifyPrivateKey(privateKey)

	if err != nil {
		t.Error(err)
	}

	if len(encoded) == 0 {
		t.Error("encoded string is empty", err)
	}

	decoded, err := crypto.DecodePrivateKeyHex(encoded)

	if err != nil {
		t.Error(err)
	}

	if !decoded.Equal(privateKey) {
		t.Error("decoded does not match encoded", err)
	}
}

func TestEncodeAndDecoePublicKey(t *testing.T) {
	_, publicKey, _ := crypto.GenerateKey()
	encoded, err := crypto.HexifyPublicKey(publicKey)

	if err != nil {
		t.Error(err)
	}

	if len(encoded) == 0 {
		t.Error("encoded string is empty")
	}

	decoded, err := crypto.DecodePublicKeyHex(encoded)

	if err != nil {
		t.Error(err)
	}

	if !decoded.Equal(publicKey) {
		t.Error("decoded does not match encoded", err)
	}
}

func TestSign(t *testing.T) {
	privKey, publKey, _ := crypto.GenerateKey()
	message := []byte("Message Data")

	signature, err := crypto.CreateDigitalSignature(message, privKey)

	if err != nil {
		t.Error("sign failed", err)
	}

	isValid := crypto.VerifyDigitalSignature(message, publKey, signature)

	if !isValid {
		t.Error("signature is not valid", isValid)
	}
}
