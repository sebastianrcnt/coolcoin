package account

import (
	"coolcoin/crypto"
	"encoding/json"
	"errors"

	"github.com/syndtr/goleveldb/leveldb"
)

type Account struct {
	Nonce   uint64 `json:"nonce"`
	Balance uint64 `json:"balance"`
}

type AccountStateMutation struct {
	PublicKey string
	Nonce     uint64
	Balance   uint64
}

type AccountManager struct {
	database *leveldb.DB
}

var ErrAccountNotFound = errors.New("Account Not Found")

func NewAccountManager(database *leveldb.DB) *AccountManager {
	return &AccountManager{
		database,
	}
}

func (accountManager *AccountManager) CreateAccount() (publicKey string, privateKeyString string, err error) {
	priv, publ, _ := crypto.GenerateKey()
	privateKeyString, _ = crypto.HexifyPrivateKey(priv)
	encodedPubl, _ := crypto.HexifyPublicKey(publ)
	publicKey, err = crypto.EncodeAddress(encodedPubl)

	accountEncoded, _ := json.Marshal(Account{
		Nonce:   0,
		Balance: 0,
	})
	accountManager.database.Put([]byte(publicKey), accountEncoded, nil)
	return
}

func (accountManager *AccountManager) HasAccount(address string) bool {
	exists, _ := accountManager.database.Has([]byte(address), nil)
	return exists
}

func (accountManager *AccountManager) GetAccount(address string) (*Account, error) {
	accountRaw, err := accountManager.database.Get([]byte(address), nil)
	if err == leveldb.ErrNotFound {
		return nil, ErrAccountNotFound
	}

	var account Account
	json.Unmarshal(accountRaw, &account)
	return &account, nil
}

func (accountManger *AccountManager) ApplyStateMutation(accountStateMutation *AccountStateMutation) error {
	accountEncoded, _ := json.Marshal(Account{
		Nonce:   accountStateMutation.Nonce,
		Balance: accountStateMutation.Balance,
	})

	err := accountManger.database.Put(
		[]byte(accountStateMutation.PublicKey),
		accountEncoded,
		nil,
	)

	if err != nil {
		return err
	}

	return nil
}
