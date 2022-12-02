package main

import (
	"coolcoin/account"
	"fmt"
	"log"

	"github.com/syndtr/goleveldb/leveldb"
)

func main() {
	db, _ := leveldb.OpenFile("./db/account", nil)
	accountManager := account.NewAccountManager(db)
	address, privateKey, err := accountManager.CreateAccount()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Your Address is ", address, "\nYour Private Key is ", privateKey)
	account, _ := accountManager.GetAccount(address)
	fmt.Println(account)

	account2, err := accountManager.GetAccount("fault address")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(account2)
}
