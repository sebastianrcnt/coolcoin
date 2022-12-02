package block

type TransactionData struct {
	SenderAddress   string
	ReceiverAddress string
	Value           uint64
	Nonce           uint64
}

type Transaction struct {
	SenderSignature string
	Body            TransactionData
}
