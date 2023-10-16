package settings

import (
	"context"
	"fmt"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	sr SettingsRepository
)

func init() {
	client, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	sr = NewSettingsRepository(client)
}

func TestCreateSettings(t *testing.T) {
	fee, err := primitive.ParseDecimal128("2.0")
	if err != nil {
		t.Error(err)
	}
	settings := SiteSettings{
		PrivacyPolicies: []PrivacyPolicy{
			{
				Id:      primitive.NewObjectID(),
				Title:   "Use of service",
				Content: "We value our customers privacy.",
			},
		},
		TermsOfUses: []TermsOfUse{
			{
				Id:      primitive.NewObjectID(),
				Title:   "Terms of service",
				Content: "You should follow our terms of service",
			},
		},
		AboutUs: "About us content 4",
		Fees: []FeeSettings{
			{
				Id:                primitive.NewObjectID(),
				SourceSymbol:      "IRT",
				DestinationSymbol: "USD",
				FeeRate:           fee,
			},
		},
	}
	res, err := sr.CreateSettings(&settings)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(res.Id)
}

func TestGetLatestSettings(t *testing.T) {
	settings, err := sr.GetLatestSettings()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(settings.AboutUs)
	fmt.Println(settings.Id)
}
