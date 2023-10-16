package paymentProcessor

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type IPaymentProcessor interface {
	GetInvoices() (*[]PaymentProcessorInvoice, error)
	CreateInvoice(ppcr PaymentProcessorCreateRequest) (*PaymentProcessorInvoice, error)
	GetInvoice(invoiceId string) (*PaymentProcessorInvoice, error)
	GetInvoicePaymentMethods(invoiceId string) (*[]PaymentMethod, error)
}

type PaymentProcessor struct {
	Token    string
	Hostname string
	StoreId  string
	Client   resty.Client
}

func NewPaymentProcessor(token string, hostname string, storeId string, client resty.Client) PaymentProcessor {
	return PaymentProcessor{
		Token:    token,
		Hostname: hostname,
		StoreId:  storeId,
		Client:   *client.SetBaseURL(hostname),
	}
}

func (pp PaymentProcessor) GetInvoices() (*[]PaymentProcessorInvoice, error) {
	url := fmt.Sprintf("/stores/%v/invoices", pp.StoreId)
	var invoices []PaymentProcessorInvoice
	_, err := pp.Client.R().SetResult(&invoices).SetHeader("Authorization", pp.Token).Get(url)
	if err != nil {
		return nil, err
	}
	return &invoices, nil
}

func (pp PaymentProcessor) CreateInvoice(ppcr PaymentProcessorCreateRequest) (*PaymentProcessorInvoice, error) {
	url := fmt.Sprintf("/stores/%v/invoices", pp.StoreId)
	var invoice PaymentProcessorInvoice
	_, err := pp.Client.R().SetBody(ppcr).SetResult(&invoice).SetHeader("Authorization", pp.Token).Post(url)
	if err != nil {
		return nil, err
	}
	return &invoice, nil
}

func (pp PaymentProcessor) GetInvoice(invoiceId string) (*PaymentProcessorInvoice, error) {
	url := fmt.Sprintf("/stores/%v/invoices/%v", pp.StoreId, invoiceId)
	var invoice PaymentProcessorInvoice
	_, err := pp.Client.R().SetResult(&invoice).SetHeader("Authorization", pp.Token).Get(url)
	if err != nil {
		return nil, err
	}
	return &invoice, nil
}

func (pp PaymentProcessor) GetInvoicePaymentMethods(invoiceId string) (*[]PaymentMethod, error) {
	url := fmt.Sprintf("/stores/%v/invoices/%v/payment-methods", pp.StoreId, invoiceId)
	var paymentMethods []PaymentMethod
	_, err := pp.Client.R().SetResult(&paymentMethods).SetHeader("Authorization", pp.Token).Get(url)
	if err != nil {
		return nil, err
	}
	return &paymentMethods, nil
}
