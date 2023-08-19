package invoice

import (
	"bean_counter/internal/types/transaction"
	"time"
)

type Type string

const (
	Purchase Type = "Purchase"
	Sales    Type = "Sales"
)

type Invoice struct {
	InvoiceType Type                    `json:"invoiceType"`
	Date        time.Time               `json:"date"`
	InvoiceNo   string                  `json:"invoiceNumber"`
	PartyName   string                  `json:"partyName"`
	GstNo       string                  `json:"gstNumber"`
	Transaction transaction.Transaction `json:"transaction"`
	Tag         string                  `json:"tag"`
}

func New(invoiceType Type, date time.Time, invoiceNo string, partyName string, gst string, transaction transaction.Transaction, tag string) Invoice {
	return Invoice{
		InvoiceType: invoiceType,
		Date:        date,
		InvoiceNo:   invoiceNo,
		PartyName:   partyName,
		GstNo:       gst,
		Transaction: transaction,
		Tag:         tag,
	}
}

func (i Invoice) GetTransaction() transaction.Transaction {
	return i.Transaction
}

func (i Invoice) GetDate() time.Time {
	return i.Date
}

func (i Invoice) GetGstNo() string {
	return i.GstNo
}

func (i Invoice) GetInvoiceType() Type {
	return i.InvoiceType
}
