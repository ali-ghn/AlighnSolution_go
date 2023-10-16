package settings

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SiteSettings struct {
	Id              primitive.ObjectID
	PrivacyPolicies []PrivacyPolicy
	TermsOfUses     []TermsOfUse
	Fees            []FeeSettings
	AboutUs         string
	CreatedAt       int64
}

type FeeSettings struct {
	Id                primitive.ObjectID `bson:"_id"`
	SourceSymbol      string
	DestinationSymbol string
	FeeRate           primitive.Decimal128
}

type PrivacyPolicy struct {
	Id      primitive.ObjectID `bson:"_id"`
	Title   string
	Content string
}

type TermsOfUse struct {
	Id      primitive.ObjectID `bson:"_id"`
	Title   string
	Content string
}
