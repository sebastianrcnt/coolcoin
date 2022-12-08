package blockchain

import (
	"encoding/hex"
	"encoding/json"

	"golang.org/x/crypto/sha3"
)

type Hashable interface{}

type Block struct {
	Height        uint64 `json:"height"`
	PrevBlockHash string `json:"prevBlockHash"`
	StateRoot     string `json:"stateRoot"`
	TxRoot        string `json:"txRoot"`
	TxCount       uint64 `json:"txCount"`
	Producer      string `json:"producer"`
	Signature     string `json:"signature"`
}

type Account struct {
	PublicKey string `json:"publicKey"`
	Nonce     uint64 `json:"nonce"`
	Balance   uint64 `json:"balance"`
}

type Transaction struct {
	PrevTransactionHash string `json:"prevTransactionHash"`
	Sender              string `json:"sender"`
	Receiver            string `json:"receiver"`
	Value               string `json:"value"`
	Signature           string `json:"signature"`
}

func Hash(hashable *Hashable) (string, error) {
	encoded, err := json.Marshal(hashable)
	if err != nil {
		return "", err
	}

	digest := sha3.Sum256(encoded)
	return hex.EncodeToString(digest[:]), nil
}
