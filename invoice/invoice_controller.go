package invoice

import (
	"net/http"

	"github.com/ali-ghn/AlighnSolution_go/auth"
	"github.com/ali-ghn/AlighnSolution_go/paymentProcessor"
	"github.com/ali-ghn/AlighnSolution_go/shared"
	"github.com/ali-ghn/AlighnSolution_go/store"
	"github.com/ali-ghn/AlighnSolution_go/transaction"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InvoiceController struct {
	pp   paymentProcessor.IPaymentProcessor
	ir   IInvoiceRepository
	sr   store.IStoreRepository
	auth auth.Auth
}

func NewInvoiceController(pp paymentProcessor.PaymentProcessor,
	ir InvoiceRepository, sr store.StoreRepository, auth auth.Auth) InvoiceController {
	return InvoiceController{
		pp:   pp,
		ir:   ir,
		sr:   sr,
		auth: auth,
	}
}

func (ic InvoiceController) CreateInvoice(c *gin.Context) {
	var invoiceReq CreateInvoiceRequest
	err := c.Bind(&invoiceReq)
	if err != nil {
		c.String(http.StatusBadRequest, "Couldn't parse the object.")
		return
	}
	if invoiceReq.Amount == decimal.Zero || invoiceReq.Currency == "" {
		c.String(http.StatusBadRequest, "Field 'amount' and 'currency' is necessary")
		return
	}
	amount := invoiceReq.Amount.String()
	inAmount, err := primitive.ParseDecimal128(amount)
	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again.")
		return
	}
	store, err := ic.sr.GetStoreByToken(invoiceReq.Token)
	if err != nil {
		c.String(http.StatusBadRequest, "Store doesn't exist")
		return
	}
	// TODO: Add support for IRT
	ppInvoice, err := ic.pp.CreateInvoice(paymentProcessor.PaymentProcessorCreateRequest{
		Amount:   invoiceReq.Amount,
		Currency: invoiceReq.Currency,
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again.")
		return
	}
	invoice := Invoice{
		Amount:             inAmount,
		Currency:           invoiceReq.Currency,
		StoreId:            store.Id,
		Description:        invoiceReq.Description,
		Status:             shared.InvoiceStatusNew,
		PaymentProcessorId: ppInvoice.Id,
		Transactions:       []transaction.Transaction{},
	}
	res, err := ic.ir.Create(&invoice)
	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again.")
		return
	}
	resInvoice := CreateInvoiceResponse{
		Id:           res.Id.Hex(),
		StoreId:      store.Id.Hex(),
		Description:  res.Description,
		Amount:       invoiceReq.Amount,
		Currency:     res.Currency,
		Status:       res.Status,
		Transactions: res.Transactions,
	}
	c.JSON(http.StatusCreated, resInvoice)
}

func (ic InvoiceController) GetInvoice(c *gin.Context) {
	var invoiceReq GetInvoiceRequest
	err := c.BindJSON(&invoiceReq)
	if err != nil {
		c.String(http.StatusBadRequest, "Couldn't parse the object.")
		return
	}
	store, err := ic.sr.GetStoreByToken(invoiceReq.Token)

	if err != nil {
		c.String(http.StatusForbidden, "Store doesn't exist.")
		return
	}

	invoice, err := ic.ir.GetInvoice(invoiceReq.InvoiceId)

	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again.")
		return
	}

	if invoice.StoreId != store.Id {
		c.String(http.StatusForbidden, "You don't have access to this resource.")
		return
	}

	amount := invoice.Amount.String()
	dAmount, err := decimal.NewFromString(amount)
	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again.")
		return
	}
	resInvoice := GetInvoiceResponse{
		Id:           invoice.Id.Hex(),
		StoreId:      store.Id.Hex(),
		Amount:       dAmount,
		Currency:     invoice.Currency,
		Status:       invoice.Status,
		Description:  invoice.Description,
		Transactions: invoice.Transactions,
	}
	c.JSON(http.StatusOK, resInvoice)
}

func (ic InvoiceController) GetInvoices(c *gin.Context) {
	var invoiceReq GetInvoicesRequest
	err := c.BindJSON(&invoiceReq)
	if err != nil {
		c.String(http.StatusBadRequest, "Couldn't parse the object.")
		return
	}
	store, err := ic.sr.GetStoreByToken(invoiceReq.Token)

	if err != nil || store == nil {
		c.String(http.StatusForbidden, "Store doesn't exist.")
		return
	}

	filter := bson.D{{Key: "storeid", Value: store.Id}}

	var invoicesResponse []GetInvoiceResponse
	invoices, err := ic.ir.GetInvoices(filter, invoiceReq.Skip, invoiceReq.Limit)

	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong, please try again.")
		return
	}

	for _, v := range *invoices {
		amount := v.Amount.String()
		dAmount, err := decimal.NewFromString(amount)

		if err != nil {
			c.String(http.StatusInternalServerError, "Something went wrong, please try again.")
			return
		}

		invoicesResponse = append(invoicesResponse, GetInvoiceResponse{
			Id:           v.Id.Hex(),
			StoreId:      v.StoreId.Hex(),
			Amount:       dAmount,
			Currency:     v.Currency,
			Status:       v.Status,
			Description:  v.Description,
			Transactions: v.Transactions,
		})
	}

	c.JSON(http.StatusOK, invoicesResponse)

}
