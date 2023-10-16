package email

import (
	"testing"

	"github.com/ali-ghn/Coinopay_Go/shared"
)

func TestSendEmail(t *testing.T) {
	smtpHost := "localhost:1025"
	password := "PTA_WIFuARWZP-64Ntxd9A"
	from := shared.INFO_EMAIL
	to := "ali.daniel.1381@gmail.com"
	es := NewEmailSender(smtpHost, password)
	err := es.Send("test", "Hello from test", from, to)
	if err != nil {
		t.Errorf(err.Error())
	}
}
