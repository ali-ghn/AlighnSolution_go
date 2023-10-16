package settings

import (
	"context"
	"time"

	"github.com/ali-ghn/Coinopay_Go/shared"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	userCollectionName         = "Users"
	siteSettingsCollectionName = "Settings"
)

type ISettingsRepository interface {
	CreateSettings(siteSettings *SiteSettings) (*SiteSettings, error)
	GetSettings(id string) (*SiteSettings, error)
	GetLatestSettings() (*SiteSettings, error)
	UpdateSettings(siteSettings *SiteSettings) (*SiteSettings, error)
}

type SettingsRepository struct {
	Client *mongo.Client
}

func NewSettingsRepository(client *mongo.Client) SettingsRepository {
	return SettingsRepository{
		Client: client,
	}
}

func (str SettingsRepository) CreateSettings(siteSettings *SiteSettings) (*SiteSettings, error) {
	siteSettings.Id = primitive.NewObjectID()
	siteSettings.CreatedAt = time.Now().UTC().Unix()
	res, err := str.Client.Database(shared.DATABASE_NAME).Collection(siteSettingsCollectionName).InsertOne(context.TODO(), siteSettings)
	if err != nil {
		return nil, err
	}
	siteSettings.Id = res.InsertedID.(primitive.ObjectID)
	return siteSettings, nil
}

func (str SettingsRepository) GetSettings(id string) (*SiteSettings, error) {
	bId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{"_id", bId}}
	var settings SiteSettings
	err = str.Client.Database(shared.DATABASE_NAME).Collection(siteSettingsCollectionName).FindOne(context.TODO(), filter).Decode(&settings)
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

func (str SettingsRepository) GetLatestSettings() (*SiteSettings, error) {
	options := options.Find().SetSort(bson.D{{"createdat", -1}}).SetLimit(1)
	var siteSettings []SiteSettings
	cur, err := str.Client.Database(shared.DATABASE_NAME).Collection(siteSettingsCollectionName).Find(context.TODO(), bson.D{{}}, options)
	if err != nil {
		return nil, err
	}
	err = cur.All(context.TODO(), &siteSettings)
	if err != nil {
		return nil, err
	}
	return &siteSettings[0], nil
}

func (str SettingsRepository) UpdateSettings(siteSettings *SiteSettings) (*SiteSettings, error) {
	filter := bson.D{{Key: "_id", Value: siteSettings.Id}}
	err := str.Client.Database(shared.DATABASE_NAME).Collection(siteSettingsCollectionName).FindOneAndReplace(context.TODO(), filter, siteSettings).Decode(&siteSettings)
	if err != nil {
		return nil, err
	}
	return siteSettings, nil
}
