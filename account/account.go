package account

import (
	"coolcoin/crypto"
	"encoding/json"
	"errors"

	"github.com/syndtr/goleveldb/leveldb"
)

type Account struct {
	Nonce   uint64
	Balance uint64
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

func (accountManager *AccountManager) CreateAccount() (addressString string, privateKeyString string, err error) {
	priv, publ, _ := crypto.GenerateKey()
	privateKeyString, _ = crypto.HexifyPrivateKey(priv)
	encodedPubl, _ := crypto.HexifyPublicKey(publ)
	addressString, err = crypto.EncodeAddress(encodedPubl)

	accountEncoded, _ := json.Marshal(Account{
		Nonce:   0,
		Balance: 0,
	})
	accountManager.database.Put([]byte(addressString), accountEncoded, nil)
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