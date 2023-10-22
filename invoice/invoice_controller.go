package invoice

import (
	"net/http"

	"github.com/ali-ghn/AlighnSolution_go/auth"
	"github.com/ali-ghn/AlighnSolution_go/store"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
)

type InvoiceController struct {
	ir   IInvoiceRepository
	sr   store.IStoreRepository
	auth auth.Auth
}

func NewInvoiceController(
	ir InvoiceRepository, sr store.StoreRepository, auth auth.Auth) InvoiceController {
	return InvoiceController{
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

	c.JSON(http.StatusCreated, "Created Not")
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
