package attachment

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Attachment struct {
	Id       primitive.ObjectID
	Filename string
}
