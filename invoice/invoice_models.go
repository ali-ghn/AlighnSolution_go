package invoice

import (
	"github.com/ali-ghn/Coinopay_Go/transaction"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Invoice struct {
	Id                 primitive.ObjectID `bson:"_id"`
	StoreId            primitive.ObjectID
	PaymentProcessorId string
	Description        string
	Amount             primitive.Decimal128
	Currency           string
	Status             string
	CallbackUrl        string
	Transactions       []transaction.Transaction
}

type CreateInvoiceRequest struct {
	Amount      decimal.Decimal
	Currency    string
	Token       string
	Description string
}

type CreateInvoiceResponse struct {
	Id           string
	StoreId      string
	Description  string
	Amount       decimal.Decimal
	Currency     string
	Status       string
	Transactions []transaction.Transaction
}

type GetInvoiceRequest struct {
	Token     string
	InvoiceId string
}

type GetInvoiceResponse struct {
	Id           string
	StoreId      string
	Amount       decimal.Decimal
	Currency     string
	Status       string
	Description  string
	Transactions []transaction.Transaction
}

type GetInvoicesRequest struct {
	Token string
	Skip  int64
	Limit int64
}
type GetInvoicesResponse struct {
	Invoices []GetInvoiceResponse
}
