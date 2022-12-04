package server

import (
	"coolcoin/account"
	"coolcoin/block"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/syndtr/goleveldb/leveldb"
)

type TransactionRequestType struct {
	SenderAddress   string `json:"sender"`
	ReceiverAddress string `json:"receiver"`
	Value           uint64 `json:"value"`
	Nonce           uint64 `json:"nonce"`
	Signature       string `json:"signature"`
}

var database, databaseErr = leveldb.OpenFile("./db/accounts", nil)

var accountManager = account.NewAccountManager(database)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func CreateServer() *gin.Engine {
	if databaseErr != nil {
		panic(databaseErr)
	}
	r := gin.Default()

	r.POST("/api/v1/create-account", func(c *gin.Context) {
		address, privateKey, err := accountManager.CreateAccount()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}
		c.JSON(
			http.StatusCreated,
			gin.H{
				"address":    address,
				"privateKey": privateKey,
			},
		)
	})

	r.GET("/api/v1/accounts/:address", func(ctx *gin.Context) {
		address := ctx.Param("address")
		found, err := accountManager.GetAccount(address)

		if err == account.ErrAccountNotFound {
			ctx.JSON(http.StatusNotFound, "")
			return
		}

		ctx.JSON(
			http.StatusOK,
			gin.H{
				"balance": found.Balance,
				"nonce":   found.Nonce,
			},
		)
	})

	r.POST("/api/v1/transactions/", func(ctx *gin.Context) {
		var req TransactionRequestType
		err := ctx.BindJSON(req)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		}

		// Create Transaction
		tx := block.Transaction{
			SenderSignature: req.Signature,
			Body: block.TransactionBody{
				SenderAddress: req.SenderAddress,
				Value:         req.Value,
				Nonce:         req.Nonce,
			},
		}

		isValid, err := block.CheckTransactionValidity(&tx, accountManager)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		if !isValid {
			ctx.JSON(400, "Invalid Transaction")
			return
		}

		sender, err := accountManager.GetAccount(req.SenderAddress)
		if err != nil {
			panic(err)
		}
		receiver, err := accountManager.GetAccount(req.ReceiverAddress)
		if err != nil {
			panic(err)
		}

		senderBalance := sender.Balance
		receiverBalance := receiver.Balance

		senderStateMutation := account.AccountStateMutation{
			Address: req.SenderAddress,
			Nonce:   req.Nonce,
			Balance: senderBalance - req.Value,
		}

		receiverStateMutation := account.AccountStateMutation{
			Address: req.ReceiverAddress,
			Nonce:   req.Nonce,
			Balance: receiverBalance + req.Value,
		}

		err = accountManager.ApplyStateMutation(&senderStateMutation)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
		err = accountManager.ApplyStateMutation(&receiverStateMutation)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		ctx.JSON(http.StatusCreated, "")
	})

	type FaucetRequestBodyType struct {
		Address string `json:"address"`
	}

	r.POST("/api/v1/faucet", func(ctx *gin.Context) {
		var body *FaucetRequestBodyType
		err := ctx.BindJSON(&body)

		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusBadRequest, err)
			return
		}

		found, err := accountManager.GetAccount(body.Address)

		if found == nil {
			ctx.JSON(http.StatusNotFound, err)
			return
		}

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		err = accountManager.ApplyStateMutation(
			&account.AccountStateMutation{
				Address: body.Address,
				Nonce:   found.Nonce,
				Balance: found.Balance + 100,
			},
		)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		found, err = accountManager.GetAccount(body.Address)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		ctx.JSON(http.StatusCreated, found)
	})

	// r.POST("/api/v1/mempool", func(ctx *gin.Context) {
	// 	// validate transaction
	// 	rawData, _ := ioutil.ReadAll(ctx.Request.Body)
	// 	var request MempoolRequest
	// 	json.Unmarshal(rawData, &request)
	// 	mempool.Requests = append(mempool.Requests, request)
	// 	ctx.JSON(http.StatusOK, "")
	// })

	// r.GET("/api/v1/mempool", func(ctx *gin.Context) {
	// 	ctx.JSON(
	// 		http.StatusOK,
	// 		mempool.Requests,
	// 	)
	// })

	// // Seal a Block
	// r.POST("/api/v1/block", func(ctx *gin.Context) {

	// })

	return r
}
