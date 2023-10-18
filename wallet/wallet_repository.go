package wallet

import (
	"context"

	"github.com/ali-ghn/AlighnSolution_go/shared"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	collectionName = "Wallets"
)

type IWalletRepository interface {
	CreateWallet(wallet *Wallet) (*Wallet, error)
	GetWallet(id string) (*Wallet, error)
	GetWallets(filter interface{}, options *options.FindOptions) (*[]Wallet, error)
	UpdateWallet(wallet *Wallet) (*Wallet, error)
}

type WalletRepository struct {
	Client *mongo.Client
}

func NewWalletRepository(client *mongo.Client) WalletRepository {
	return WalletRepository{
		Client: client,
	}
}

func (wr WalletRepository) CreateWallet(wallet *Wallet) (*Wallet, error) {
	wallet.Id = primitive.NewObjectID()
	res, err := wr.Client.Database(shared.DATABASE_NAME).Collection(collectionName).InsertOne(context.TODO(), wallet)
	if err != nil {
		return nil, err
	}
	wallet.Id = res.InsertedID.(primitive.ObjectID)
	return wallet, nil
}

func (wr WalletRepository) GetWallet(id string) (*Wallet, error) {
	bId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{"_id", bId}}
	var wallet Wallet
	err = wr.Client.Database(shared.DATABASE_NAME).Collection(collectionName).FindOne(context.TODO(), filter).Decode(&wallet)
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (wr WalletRepository) GetWallets(filter interface{}, options *options.FindOptions) (*[]Wallet, error) {
	var wallets []Wallet
	cur, err := wr.Client.Database(shared.DATABASE_NAME).Collection(collectionName).Find(context.TODO(), filter, options)
	if err != nil {
		return nil, err
	}
	err = cur.All(context.TODO(), &wallets)
	if err != nil {
		return nil, err
	}
	return &wallets, nil
}

func (wr WalletRepository) UpdateWallet(wallet *Wallet) (*Wallet, error) {
	filter := bson.D{{"_id", wallet.Id}}
	err := wr.Client.Database(shared.DATABASE_NAME).Collection(collectionName).FindOneAndUpdate(context.TODO(), filter, wallet).Decode(&wallet)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}
