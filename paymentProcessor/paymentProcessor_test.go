package paymentProcessor

import (
	"fmt"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/shopspring/decimal"
)

var pp PaymentProcessor

func init() {
	hostname := "https://testnet.demo.btcpayserver.org/api/v1"
	token := "token 4b9bf39f9b18dff14285fce8e8de9937a30fdd3a"
	storeId := "ChsKJKCKheQLeYdoG6Y14RjLr8AWGt9axQmuB8XGNo8u"
	pp = NewPaymentProcessor(token, hostname, storeId, *resty.New())
}

func TestGetInvoices(t *testing.T) {
	invoices, err := pp.GetInvoices()
	if err != nil {
		t.Errorf(err.Error())
	}
	for _, v := range *invoices {
		fmt.Println(v.Id)
	}
}

func TestCreateInvoice(t *testing.T) {
	amount, err := decimal.NewFromString("5.3")
	if err != nil {
		t.Errorf(err.Error())
	}
	invoice, err := pp.CreateInvoice(PaymentProcessorCreateRequest{
		Amount:   amount,
		Currency: "USD",
	})

	if err != nil {
		t.Errorf(err.Error())
	}

	fmt.Println(invoice.Id)
}

func TestGetInvoice(t *testing.T) {
	invoiceId := "Pf3vEDHqT8uUBxysVQPJxn"
	invoice, err := pp.GetInvoice(invoiceId)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(invoice.Id)
	fmt.Println(invoice.PaymentMethods)
}

func TestGetPaymentMethods(t *testing.T) {
	invoiceId := "JhUnEeVGRNz3d1RNaz1CgY"
	paymentMethods, err := pp.GetInvoicePaymentMethods(invoiceId)
	if err != nil {
		t.Errorf(err.Error())
	}
	for _, v := range *paymentMethods {
		fmt.Println(v.Destination)
		fmt.Println(v.Amount)
	}
}
