package invoice

import (
	"bean_counter/internal/types/tax"
	"bean_counter/internal/types/transaction"
	"context"
	"fmt"
	"google.golang.org/api/sheets/v4"
	"strconv"
	"time"
)

type Repository interface {
	getAllPurchases()
}

type repositoryImpl struct {
	service *sheets.Service
}

func (r repositoryImpl) GetAllPurchaseInvoices(ctx context.Context) []Invoice {
	// https://docs.google.com/spreadsheets/d/<SPREADSHEETID>/edit#gid=<SHEETID>
	sheetId := 0
	spreadsheetId := "1PyN1aHbYWM8d0gUWKdXCZCpNq8SJnG1B_LeMJ1Pcafw"

	// Convert sheet ID to sheet name.
	response1, err := r.service.Spreadsheets.
		Get(spreadsheetId).
		Fields("sheets(properties(sheetId,title))").
		Do()
	if err != nil || response1.HTTPStatusCode != 200 {
		fmt.Println(err) //panic
	}

	sheetName := ""
	for _, v := range response1.Sheets {
		prop := v.Properties
		if prop.SheetId == int64(sheetId) {
			sheetName = prop.Title
			break
		}
	}

	response, err := r.service.Spreadsheets.Values.
		Get(spreadsheetId, sheetName).
		Context(ctx).Do()
	if err != nil || response.HTTPStatusCode != 200 {
		fmt.Println(err) //panic
		//return
	}
	var purchaseInvoices []Invoice

	for index, row := range response.Values {
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

func NewRepository(service *sheets.Service) repositoryImpl {
	return repositoryImpl{
		service: service,
	}
}
