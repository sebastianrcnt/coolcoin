package digital_signature_test

import (
	"coolcoin/digital_signature"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	privateKey, publicKey, err := digital_signature.GenerateKey()

	if !privateKey.PublicKey.Equal(publicKey) {
		t.Error("Private Key and Public Key don't match", err)
	}

	if err != nil {
		t.Error(err)
	}
}

func TestEncodeAndDecodePrivateKey(t *testing.T) {
	privateKey, _, _ := digital_signature.GenerateKey()
	encoded, err := digital_signature.HexifyPrivateKey(privateKey)

	if err != nil {
		t.Error(err)
	}

	if len(encoded) == 0 {
		t.Error("encoded string is empty", err)
	}

	decoded, err := digital_signature.DecodePrivateKeyHex(encoded)

	if err != nil {
		t.Error(err)
	}

	if !decoded.Equal(privateKey) {
		t.Error("decoded does not match encoded", err)
	}
}

func TestEncodeAndDecoePublicKey(t *testing.T) {
	_, publicKey, _ := digital_signature.GenerateKey()
	encoded, err := digital_signature.HexifyPublicKey(publicKey)

	if err != nil {
		t.Error(err)
	}

	if len(encoded) == 0 {
		t.Error("encoded string is empty")
	}

	decoded, err := digital_signature.DecodePublicKeyHex(encoded)

	if err != nil {
		t.Error(err)
	}

	if !decoded.Equal(publicKey) {
		t.Error("decoded does not match encoded", err)
	}
}

func TestSign(t *testing.T) {
	privKey, publKey, _ := digital_signature.GenerateKey()
	message := []byte("Message Data")

	signature, err := digital_signature.CreateDigitalSignature(message, privKey)

	if err != nil {
		t.Error("sign failed", err)
	}

	isValid := digital_signature.VerifyDigitalSignature(message, publKey, signature)

	if !isValid {
		t.Error("signature is not valid", isValid)
	}
}
