package invoice

import (
	"context"

	"github.com/ali-ghn/Coinopay_Go/shared"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	collectionName = "Invoices"
)

type IInvoiceRepository interface {
	Create(invoice *Invoice) (*Invoice, error)
	GetInvoice(id string) (*Invoice, error)
	GetInvoices(filter interface{}, skip int64, limit int64) (*[]Invoice, error)
	UpdateInvoice(invoice *Invoice) (*Invoice, error)
}

type InvoiceRepository struct {
	Client *mongo.Client
}

func NewInvoiceRepository(client *mongo.Client) InvoiceRepository {
	return InvoiceRepository{
		Client: client,
	}
}

func (ir InvoiceRepository) Create(invoice *Invoice) (*Invoice, error) {
	invoice.Id = primitive.NewObjectID()
	res, err := ir.Client.Database(shared.DATABASE_NAME).Collection(collectionName).InsertOne(context.TODO(), invoice)
	if err != nil {
		return nil, err
	}
	invoice.Id = res.InsertedID.(primitive.ObjectID)
	return invoice, nil
}

func (ir InvoiceRepository) GetInvoice(id string) (*Invoice, error) {
	bId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: bId}}
	var invoice Invoice
	err = ir.Client.Database(shared.DATABASE_NAME).Collection(collectionName).FindOne(context.TODO(), filter).Decode(&invoice)
	if err != nil {
		return nil, err
	}
	return &invoice, nil
}

func (ir InvoiceRepository) GetInvoices(filter interface{}, skip int64, limit int64) (*[]Invoice, error) {
	var invoices []Invoice
	if limit == 0 {
		limit = 9223372036854775807
	}
	options := options.Find().SetSkip(skip).SetLimit(limit)
	res, err := ir.Client.Database(shared.DATABASE_NAME).Collection(collectionName).Find(context.TODO(), filter, options)
	if err != nil {
		return nil, err
	}
	res.All(context.TODO(), &invoices)
	return &invoices, nil
}

func (ir InvoiceRepository) UpdateInvoice(invoice *Invoice) (*Invoice, error) {
	filter := bson.D{{Key: "_id", Value: invoice.Id}}
	err := ir.Client.Database(shared.DATABASE_NAME).Collection(collectionName).FindOneAndReplace(context.TODO(), filter, invoice).Decode(&invoice)
	if err != nil {
		return nil, err
	}
	return invoice, err
}
