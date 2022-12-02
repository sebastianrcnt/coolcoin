package block

import (
	"coolcoin/account"
	"coolcoin/crypto"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type BlockData struct {
	Height            uint64
	Timestamp         uint64
	ProducerAddress   string
	PreviousBlockHash string
	Transactions      []Transaction
}

type Block struct {
	Data      BlockData
	Signature string
}

func CreateGenesisBlock(
	producerPrivateKey *ecdsa.PrivateKey,
) (*Block, error) {
	producerPublicKeyHex, err := crypto.HexifyPublicKey(&producerPrivateKey.PublicKey)
	if err != nil {
		return nil, err
	}

	producerAddress, err := crypto.EncodeAddress(producerPublicKeyHex)
	if err != nil {
		return nil, err
	}

	blockData := BlockData{
		Height:            0,
		Timestamp:         uint64(time.Now().UnixMilli()),
		ProducerAddress:   producerAddress,
		PreviousBlockHash: "genesis",
		Transactions:      make([]Transaction, 0),
	}

	encodedBlockData, err := json.Marshal(blockData)

	if err != nil {
		return nil, err
	}

	signature, err := crypto.CreateDigitalSignature(encodedBlockData, producerPrivateKey)

	if err != nil {
		return nil, err
	}

	return &Block{
		Signature: string(signature),
		Data:      blockData,
	}, nil
}

var ErrInsufficientBalance = errors.New("insufficient balance")
var ErrAccountNotFound = errors.New("account not found")
var ErrInvalidPublicKey = errors.New("invalid public key")
var ErrInvalidSignature = errors.New("invalid signature")
var ErrInvalidNonce = errors.New("invalid nonce")

func (transaction *Transaction) ValidateSignature() {
}

func CheckTransactionValidity(transaction *Transaction, accountManager *account.AccountManager) (bool, error) {
	senderAccount, err := accountManager.GetAccount(transaction.Body.SenderAddress)

	if senderAccount == nil {
		return false, ErrAccountNotFound
	}

	if err != nil {
		return false, err
	}

	senderPublicKey, senderPublicKeyErr := crypto.DecodePublicKeyHex(transaction.Body.SenderAddress)
	_, receiverPublicKeyErr := crypto.DecodePublicKeyHex(transaction.Body.ReceiverAddress)

	// Check PublicKey Validity
	if senderPublicKeyErr != nil || receiverPublicKeyErr != nil {
		fmt.Println(senderPublicKeyErr, receiverPublicKeyErr)
		return false, ErrInvalidPublicKey
	}

	// Check Balance
	if senderAccount.Balance < transaction.Body.Value {
		return false, ErrInsufficientBalance
	}

	// Check Nonce
	if senderAccount.Nonce != transaction.Body.Nonce+1 {
		return false, ErrInvalidNonce
	}

	// Check Signature
	originalMessage, err := json.Marshal(transaction.Body)

	if err != nil {
		return false, err
	}

	isValidSignature := crypto.VerifyDigitalSignature(
		originalMessage, senderPublicKey, []byte(transaction.SenderSignature),
	)

	if !isValidSignature {
		return false, ErrInvalidSignature
	}

	return true, nil
}

func (transaction *Transaction) ToAccountStateMutation(accountManager *account.AccountManager) (*account.AccountStateMutation, error) {

	accountStateMutation := account.AccountStateMutation{}
	return &accountStateMutation, nil
}

func (block *Block) Hash() string {
	encoded, err := json.Marshal(block)

	if err != nil {
		panic(err)
	}

	hash := sha256.Sum256(encoded)
	return hex.EncodeToString(hash[:])
}

func (previousBlock *Block) ProduceNextBlock(
	transactions []Transaction,
	producerPrivateKey ecdsa.PrivateKey,
	accountManager *account.AccountManager,
) (*Block, []*account.AccountStateMutation, error) {
	var accountStateMutations []*account.AccountStateMutation = make([]*account.AccountStateMutation, len(transactions))

	for idx, transaction := range transactions {
		isTransactionValid, err := CheckTransactionValidity(&transaction, accountManager)

		if !isTransactionValid {
			return nil, nil, err
		}

		accountStateMutation, err := transaction.ToAccountStateMutation(
			accountManager,
		)

		if err != nil {
			return nil, nil, err
		}

		accountStateMutations[idx] = accountStateMutation
	}

	producerPublicKey, err := crypto.HexifyPublicKey(&producerPrivateKey.PublicKey)

	if err != nil {
		panic(err)
	}

	nextBlockData := BlockData{
		Height:            previousBlock.Data.Height + 1, // increment block height by 1
		Timestamp:         uint64(time.Now().Nanosecond()),
		ProducerAddress:   string(producerPublicKey),
		PreviousBlockHash: previousBlock.Hash(),
		Transactions:      transactions,
	}

	encodedNextBlockData, _ := json.Marshal(nextBlockData)
	signature, err := crypto.CreateDigitalSignature(encodedNextBlockData, &producerPrivateKey)

	if err != nil {
		return nil, accountStateMutations, err
	}

	nextBlock := Block{
		Data:      nextBlockData,
		Signature: string(signature),
	}

	return &nextBlock, accountStateMutations, nil
}
