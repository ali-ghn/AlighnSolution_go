package wallet

import (
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Wallet struct {
	Id      primitive.ObjectID `bson:"_id"`
	Name    string
	Symbol  string
	Balance primitive.Decimal128
}

type GetWalletOverviewResponse struct {
	balance            decimal.Decimal
	UnconfirmedBalance decimal.Decimal
	ConfirmedBalance   decimal.Decimal
}

type GetWalletFeeRateResponse struct {
	FeeRate decimal.Decimal
}

type GetWalletAddressResponse struct {
	Address     string
	KeyPath     string
	PaymentLink string
}

type WalletTransaction struct {
	TransactionHash string          `json:"transactionHash"`
	Comment         string          `json:"comment"`
	Amount          decimal.Decimal `json:"amount"`
	BlockHash       string          `json:"blockHash"`
	BlockHeight     string          `json:"blockHeight"`
	Confirmations   string          `json:"confirmations"`
	Timestamp       int             `json:"timestamp"`
	Status          string          `json:"status"`
}

type TransactionDestination struct {
	Destination        string          `json:"destination"`
	Amount             decimal.Decimal `json:"amount"`
	SubtractFromAmount bool            `json:"subtractFromAmount"`
}

type CreateWalletTransactionRequest struct {
	Destinations []TransactionDestination `json:"destinations"`
}

type CreateWalletTransactionResponse struct {
	TransactionHash string          `json:"transactionHash"`
	Comment         string          `json:"comment"`
	Amount          decimal.Decimal `json:"amount"`
	BlockHash       string          `json:"blockHash"`
	BlockHeight     string          `json:"blockHeight"`
	Confirmations   string          `json:"confirmations"`
	Timestamp       int64           `json:"timestamp"`
	Status          string          `json:"status"`
}

type GetWalletRequest struct {
	CryptoCode string
}

type GetWalletsResponse struct {
	Id         string
	CryptoCode string
	Symbol     string
	Name       string
	Balance    decimal.Decimal
}

type CreateTransactionResponse struct {
	TransactionHash string          `json:"transactionHash"`
	Comment         string          `json:"comment"`
	Amount          decimal.Decimal `json:"amount"`
	BlockHash       string          `json:"blockHash"`
	BlockHeight     string          `json:"blockHeight"`
	Confirmations   string          `json:"confirmations"`
	Timestamp       int64           `json:"timestamp"`
	Status          string          `json:"status"`
}

type CreateTransactionRequest struct {
	CryptoCode  string
	Destination string
	Amount      decimal.Decimal
	SubtractFee bool
}
