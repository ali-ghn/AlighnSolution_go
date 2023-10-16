package invoice

import (
	"context"
	"fmt"
	"testing"

	"github.com/ali-ghn/Coinopay_Go/shared"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ir     InvoiceRepository
	client *mongo.Client
)

func init() {
	client, _ = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	ir = NewInvoiceRepository(client)
}

func TestCreate(t *testing.T) {
	bId, err := primitive.ObjectIDFromHex("632f5caa63bb2d85621a4241")
	if err != nil {
		t.Errorf(err.Error())
	}
	amount, err := primitive.ParseDecimal128("5.0")
	if err != nil {
		t.Errorf(err.Error())
	}
	repoInvoice := Invoice{
		StoreId:     bId,
		Amount:      amount,
		Currency:    "USD",
		Status:      shared.InvoiceStatusNew,
		Description: "Invoice description",
	}
	res, err := ir.Create(&repoInvoice)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(res.Id)
}

func TestGetInvoice(t *testing.T) {
	invoiceId := "633076799ef4f2714f3d988c"
	invoice, err := ir.GetInvoice(invoiceId)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(invoice.Id)
	fmt.Println(invoice.Amount)
	fmt.Println(invoice.Currency)
}

func TestGetInvoices(t *testing.T) {
	filter := bson.D{{Key: "currency", Value: "USD"}}
	invoices, err := ir.GetInvoices(filter, 0, 0)
	if err != nil {
		t.Errorf(err.Error())
	}
	for _, v := range *invoices {
		fmt.Println(v.Id.Hex())
	}
}
