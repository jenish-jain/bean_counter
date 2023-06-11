package invoice

import (
	"bean_counter/internal/types/transaction"
	"time"
)

type invoiceType string

const (
	Purchase invoiceType = "Purchase"
	Sales    invoiceType = "Sales"
)

type Invoice struct {
	invoiceType invoiceType             `json:"invoiceType"`
	date        time.Time               `json:"date"`
	invoiceNo   string                  `json:"invoiceNumber"`
	partyName   string                  `json:"partyName"`
	gstNo       string                  `json:"gstNumber"`
	transaction transaction.Transaction `json:"transaction"`
	tag         string                  `json:"tag"`
}

func New(invoiceType invoiceType, date time.Time, invoiceNo string, partyName string, gst string, transaction transaction.Transaction, tag string) Invoice {
	return Invoice{
		invoiceType: invoiceType,
		date:        date,
		invoiceNo:   invoiceNo,
		partyName:   partyName,
		gstNo:       gst,
		transaction: transaction,
		tag:         tag,
	}
}

func (i Invoice) GetTransaction() transaction.Transaction {
	return i.transaction
}

func (i Invoice) GetDate() time.Time {
	return i.date
}

func (i Invoice) GetGstNo() string {
	return i.gstNo
}

func (i Invoice) GetInvoiceType() invoiceType {
	return i.invoiceType
}
