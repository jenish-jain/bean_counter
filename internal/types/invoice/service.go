package invoice

import (
	"bean_counter/internal/types/tax"
	"bean_counter/internal/types/transaction"
	"bean_counter/pkg/gsheets"
	"context"
	"fmt"
	"strconv"
	"time"
)

type Service interface {
	GetAllPurchaseInvoices(ctx context.Context) []Invoice
	GetAllSalesInvoices(ctx context.Context) []Invoice
}

type serviceImpl struct {
	repository gsheets.Repository
}

func (s serviceImpl) GetAllPurchaseInvoices(ctx context.Context) []Invoice {
	// https://docs.google.com/spreadsheets/d/<SPREADSHEETID>/edit#gid=<SHEETID>
	sheetId := 0
	spreadsheetId := "1PyN1aHbYWM8d0gUWKdXCZCpNq8SJnG1B_LeMJ1Pcafw"

	offlineSalesRecords := s.repository.GetAllRecords(spreadsheetId, sheetId)
	var purchaseInvoices []Invoice

	for index, row := range offlineSalesRecords {
		if index > 0 {
			invoice := getInvoiceFromRow(row)
			purchaseInvoices = append(purchaseInvoices, invoice)
		}
	}
	return purchaseInvoices

}

func (s serviceImpl) GetAllSalesInvoices(ctx context.Context) []Invoice {
	sheetId := 0
	spreadsheetId := "192CsSjGPrkxFkoUTg5_TrQIB8tez_UgiH5XHgA7ITKA"

	offlinePurchaseRecords := s.repository.GetAllRecords(spreadsheetId, sheetId)
	var purchaseInvoices []Invoice
	for index, row := range offlinePurchaseRecords {
		if index > 0 {
			invoice := getInvoiceFromRow(row)
			purchaseInvoices = append(purchaseInvoices, invoice)
		}
	}
	return purchaseInvoices

}

func getInvoiceFromRow(row []interface{}) Invoice {
	totalAmount := parseFloat(fmt.Sprintf("%s", row[6]))
	invoiceDate, _ := time.Parse("2006-01-02", fmt.Sprintf("%s", row[1]))
	invoiceNo := fmt.Sprintf("%s", row[2])
	partyName := fmt.Sprintf("%s", row[3])
	gstNo := fmt.Sprintf("%s", row[5])
	cgst := parseFloat(fmt.Sprintf("%s", row[7]))
	sgst := parseFloat(fmt.Sprintf("%s", row[8]))
	igst := parseFloat(fmt.Sprintf("%s", row[9]))

	invoiceTransaction := transaction.New(totalAmount, tax.New(cgst, sgst, igst))
	invoice := New(Purchase, invoiceDate, invoiceNo, partyName, gstNo, invoiceTransaction, "offline")
	return invoice
}

func parseFloat(val string) float64 {
	if val == "" {
		return 0
	}
	output, err := strconv.ParseFloat(val, 64)
	if err != nil {
		panic(err)
	}
	return output
}

func NewService(repository gsheets.Repository) Service {
	return &serviceImpl{
		repository: repository,
	}
}
