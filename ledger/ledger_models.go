package ledger

import "go.mongodb.org/mongo-driver/bson/primitive"

type Ledger struct {
	Id        primitive.ObjectID `bson:"_id"`
	Name      string
	Symbol    string
	NetworkId string
}
