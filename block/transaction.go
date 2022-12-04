package block

type TransactionBody struct {
	SenderAddress   string `json:"sender"`
	ReceiverAddress string `json:"receiver"`
	Value           uint64 `json:"value"`
	Nonce           uint64 `json:"nonce"`
}

type Transaction struct {
	SenderSignature string
	Body            TransactionBody
}
