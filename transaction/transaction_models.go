package transaction

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	Id                  primitive.ObjectID `bson:"_id"`
	LedgerTransactionId string
	LedgerId            primitive.ObjectID
	Amount              primitive.Decimal128
	Currency            string
	Date                time.Time
}
