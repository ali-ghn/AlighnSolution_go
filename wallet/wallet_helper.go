package wallet

import (
	"fmt"
	"math"

	"github.com/go-resty/resty/v2"
)

type IWalletHelper interface {
	GetWalletOverview(cryptoCode string) (*GetWalletOverviewResponse, error)
	GetWalletFeeRate(cryptoCode string) (*GetWalletFeeRateResponse, error)
	GetWalletAddress(cryptoCode string) (*GetWalletAddressResponse, error)
	GetWalletTransactions(cryptoCode string, skip int64, limit int64) (*[]WalletTransaction, error)
	CreateWalletTransaction(cryptoCode string, cwtr *CreateWalletTransactionRequest) (*CreateWalletTransactionResponse, error)
}

type WalletHelper struct {
	Client  resty.Client
	StoreId string
	Token   string
}

func NewWalletHelper(token string, hostname string, storeId string, client resty.Client) WalletHelper {
	return WalletHelper{
		Client:  *client.SetBaseURL(hostname),
		StoreId: storeId,
		Token:   token,
	}
}

func (wh WalletHelper) GetWalletOverview(cryptoCode string) (*GetWalletOverviewResponse, error) {
	url := fmt.Sprintf("/stores/%v/payment-methods/onchain/%v/wallet", wh.StoreId, cryptoCode)
	var walletOverview GetWalletOverviewResponse
	_, err := wh.Client.R().SetResult(&walletOverview).SetHeader("Authorization", wh.Token).Get(url)
	if err != nil {
		return nil, err
	}
	return &walletOverview, nil
}

func (wh WalletHelper) GetWalletFeeRate(cryptoCode string) (*GetWalletFeeRateResponse, error) {
	url := fmt.Sprintf("/stores/%v/payment-methods/onchain/%v/wallet/feerate", wh.StoreId, cryptoCode)
	var feeRate GetWalletFeeRateResponse
	_, err := wh.Client.R().SetResult(&feeRate).SetHeader("Authorization", wh.Token).Get(url)
	if err != nil {
		return nil, err
	}
	return &feeRate, nil
}

func (wh WalletHelper) GetWalletAddress(cryptoCode string) (*GetWalletAddressResponse, error) {
	url := fmt.Sprintf("/stores/%v/payment-methods/onchain/%v/wallet/address", wh.StoreId, cryptoCode)
	var walletAddress GetWalletAddressResponse
	_, err := wh.Client.R().SetResult(&walletAddress).SetHeader("Authorization", wh.Token).Get(url)
	if err != nil {
		return nil, err
	}
	return &walletAddress, nil
}

func (wh WalletHelper) GetWalletTransactions(cryptoCode string, skip int64, limit int64) (*[]WalletTransaction, error) {
	url := fmt.Sprintf("/stores/%v/payment-methods/onchain/%v/wallet/transactions", wh.StoreId, cryptoCode)
	if limit == 0 {
		limit = math.MaxInt64
	}
	res, err := wh.Client.R().SetResult(&[]WalletTransaction{}).SetHeader("Authorization", wh.Token).SetQueryParam("skip", fmt.Sprint(skip)).SetQueryParam("limit", fmt.Sprint(limit)).Get(url)
	if err != nil {
		return nil, err
	}
	return res.Result().(*[]WalletTransaction), nil
}

func (wh WalletHelper) CreateWalletTransaction(cryptoCode string, cwtr *CreateWalletTransactionRequest) (*CreateWalletTransactionResponse, error) {
	url := fmt.Sprintf("/stores/%v/payment-methods/onchain/%v/wallet/transactions", wh.StoreId, cryptoCode)
	var transactionResult CreateWalletTransactionResponse
	_, err := wh.Client.R().SetBody(cwtr).SetResult(&transactionResult).SetHeader("Authorization", wh.Token).Post(url)
	if err != nil {
		return nil, err
	}
	return &transactionResult, nil
}
