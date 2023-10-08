package invoice

import (
	"bean_counter/internal/types/tax"
	"bean_counter/internal/types/transaction"
	"bean_counter/pkg/gsheets"
	"bean_counter/pkg/utils"
	"context"
	"fmt"
	"os"
	"strconv"
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
	sheetID, _ := strconv.Atoi(os.Getenv("PURCHASE_SHEET_ID"))
	spreadsheetID := os.Getenv("PURCHASE_SPREADSHEET_ID")

	offlineSalesRecords := s.repository.GetAllRecords(spreadsheetID, sheetID)
	var purchaseInvoices []Invoice

	for index, row := range offlineSalesRecords {
		if index > 0 {
			invoice := getInvoiceFromRow(row, Purchase)
			purchaseInvoices = append(purchaseInvoices, invoice)
		}
	}
	return purchaseInvoices

}

func (s serviceImpl) GetAllSalesInvoices(ctx context.Context) []Invoice {
	sheetID, _ := strconv.Atoi(os.Getenv("SALES_SHEET_ID"))
	spreadsheetID := os.Getenv("SALES_SPREADSHEET_ID")

	offlinePurchaseRecords := s.repository.GetAllRecords(spreadsheetID, sheetID)
	var purchaseInvoices []Invoice
	for index, row := range offlinePurchaseRecords {
		if index > 0 {
			invoice := getInvoiceFromRow(row, Sales)
			purchaseInvoices = append(purchaseInvoices, invoice)
		}
	}
	return purchaseInvoices

}

func getInvoiceFromRow(row []interface{}, transactionType Type) Invoice {
	totalAmount := parseFloat(fmt.Sprintf("%s", row[6]))
	invoiceDate := utils.ParseMDYYYYToDate(fmt.Sprintf("%s", row[1]))
	invoiceNo := fmt.Sprintf("%s", row[2])
	partyName := fmt.Sprintf("%s", row[3])
	gstNo := fmt.Sprintf("%s", row[5])
	cgst := parseFloat(fmt.Sprintf("%s", row[7]))
	sgst := parseFloat(fmt.Sprintf("%s", row[8]))
	igst := parseFloat(fmt.Sprintf("%s", row[9]))

	invoiceTransaction := transaction.New(totalAmount, tax.New(cgst, sgst, igst))
	invoice := New(transactionType, invoiceDate, invoiceNo, partyName, gstNo, invoiceTransaction, "offline")
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
