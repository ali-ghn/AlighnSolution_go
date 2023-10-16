package paymentProcessor

import "github.com/shopspring/decimal"

type PaymentProcessorInvoice struct {
	MetaData
	Checkout
	Receipt
	Id                                string
	StoreId                           string
	Amount                            decimal.Decimal
	Currency                          string
	Type                              string
	CheckOutLink                      string
	CreatedTime                       int64
	ExpirationTime                    int64
	MonitoringTime                    int64
	Status                            string
	AdditionalStatus                  string
	AvailableStatusesForManualMarking []string
	Archived                          bool
}

type MetaData struct {
	OrderId  string
	OrderUrl string
}

type Checkout struct {
	SpeedPolicy           string
	PaymentMethods        []string
	DefaultPaymentMethod  string
	ExpirationMinutes     int
	MonitoringMinutes     int
	PaymentTolerance      int
	RedirectUrl           string
	RedirectAutomatically bool
	RequiresRefundEmail   bool
	DefaultLanguage       string
}

type Receipt struct {
	Enabled      bool
	ShowQR       bool
	ShowPayments bool
}

type PaymentMethod struct {
	PaymentMethod     string
	CryptoCode        string
	Destination       string
	PaymentLink       string
	Rate              string
	PaymentMethodPaid string
	TotalPaid         string
	Due               string
	Amount            string
	NetworkFee        string
	Payments          []Payment
	Activated         bool
	AdditionalData
}

type Payment struct {
	Id           string
	ReceivedDate int64
	Value        string
	Fee          string
	Status       string
	Destination  string
}

type AdditionalData struct {
	ProvidedComment          string
	ConsumedLightningAddress string
}

type PaymentProcessorCreateRequest struct {
	Amount   decimal.Decimal
	Currency string
}
