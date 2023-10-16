package wallet

import (
	"fmt"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/shopspring/decimal"
)

var wh WalletHelper

func init() {
	hostname := "https://testnet.demo.btcpayserver.org/api/v1"
	token := "token 4b9bf39f9b18dff14285fce8e8de9937a30fdd3a"
	storeId := "ChsKJKCKheQLeYdoG6Y14RjLr8AWGt9axQmuB8XGNo8u"
	wh = NewWalletHelper(token, hostname, storeId, *resty.New())
}

func TestGetWalletOverview(t *testing.T) {
	walletOverview, err := wh.GetWalletOverview("BTC")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(walletOverview.balance)
	fmt.Println(walletOverview.ConfirmedBalance)
	fmt.Println(walletOverview.UnconfirmedBalance)
}

func TestGetWalletFeeRate(t *testing.T) {
	walletFeeRate, err := wh.GetWalletFeeRate("BTC")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(walletFeeRate.FeeRate)
}

func TestGetWalletAddress(t *testing.T) {
	walletAddress, err := wh.GetWalletAddress("BTC")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(walletAddress.Address)
	fmt.Println(walletAddress.KeyPath)
	fmt.Println(walletAddress.PaymentLink)
}

func TestGetWalletTransactions(t *testing.T) {
	walletTransactions, err := wh.GetWalletTransactions("BTC", 0, 0)
	if err != nil {
		t.Error(err)
	}
	for _, v := range *walletTransactions {
		fmt.Println(v.Amount)
		fmt.Println(v.TransactionHash)
		fmt.Println(v.Status)
	}
}
func TestCreateWalletTransaction(t *testing.T) {
	amount, err := decimal.NewFromString("0.0005")
	createWalletTransactionResponse, err := wh.CreateWalletTransaction("BTC", &CreateWalletTransactionRequest{
		Destinations: []TransactionDestination{
			{
				Destination:        "tb1q0gyfsfjfqwu8jaae6yhx0qj9feqgad80m9t5ru",
				Amount:             amount,
				SubtractFromAmount: true,
			},
		},
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(createWalletTransactionResponse.Amount)
	fmt.Println(createWalletTransactionResponse.TransactionHash)
	fmt.Println(createWalletTransactionResponse.Status)
}
