package reporter

import (
	"bean_counter/internal/types/invoice"
	"bean_counter/internal/types/tax"
	"bean_counter/internal/types/transaction"
	"bean_counter/pkg/files"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type Reporter interface {
	GetTaxReportOfMonth(month time.Month, year int) taxReport
	GetSheetValuesToPublishReport(taxReport taxReport, month time.Month, year int) [][]interface{}
}

type reporterImpl struct {
	invoiceService invoice.Service
}

type stateTax struct {
	StateName     string
	StateCode     string
	PayableTax    tax.Tax
	TaxOnPurchase tax.Tax
}

type taxReport struct {
	Month           time.Month
	Year            int
	TotalPurchases  transaction.Transaction
	TotalSales      transaction.Transaction
	TotalInputTax   float64
	TotalPayableTax float64
	NetTax          float64
	StateTaxBreakup map[string]stateTax
}

type stateDetails struct {
	StateName string `json:"stateName"`
	StateCode string `json:"stateCode"`
}

func (r reporterImpl) GetTaxReportOfMonth(month time.Month, year int) taxReport {
	ctx := context.Background()
	purchases := r.invoiceService.GetAllPurchaseInvoices(ctx)
	sales := r.invoiceService.GetAllSalesInvoices(ctx)

	tinToStateDetailsMap := getTinToStateDetailsMap()
	stateTaxMap := map[string]stateTax{}

	var purchasesOfMonth transaction.Transaction
	var salesOfMonth transaction.Transaction
	for _, purchaseInvoice := range purchases {
		if purchaseInvoice.GetDate().Month() == month && purchaseInvoice.GetDate().Year() == year {
			purchasesOfMonth = purchasesOfMonth.Add(purchaseInvoice.GetTransaction())
			stateTaxMap = checkAndPopulateStateTax(purchaseInvoice, stateTaxMap, tinToStateDetailsMap)
		}
	}
	for _, salesInvoice := range sales {
		if salesInvoice.GetDate().Month() == month && salesInvoice.GetDate().Year() == year {
			salesOfMonth = salesOfMonth.Add(salesInvoice.GetTransaction())
			stateTaxMap = checkAndPopulateStateTax(salesInvoice, stateTaxMap, tinToStateDetailsMap)
		}
	}
	report := populateTaxReport(purchasesOfMonth, salesOfMonth)
	report.StateTaxBreakup = stateTaxMap
	report.Month = month
	report.Year = year
	file, _ := json.MarshalIndent(report, "", " ")
	if err := os.WriteFile("report.json", file, 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n REPORT %+v \n ", report)
	return report

}

//go:embed assets/*
var res embed.FS

func getTinToStateDetailsMap() map[string]stateDetails {
	var tinToStateDetailsMap map[string]stateDetails
	byteValue := files.GetAll(res, "assets/tinToStateDetails.json")
	json.Unmarshal(byteValue, &tinToStateDetailsMap)
	return tinToStateDetailsMap
}

func checkAndPopulateStateTax(invoiceRecord invoice.Invoice, stateTaxMap map[string]stateTax, tinToStateDetailsMap map[string]stateDetails) map[string]stateTax {
	stateTin := invoiceRecord.GetGstNo()[0:2]
	fmt.Printf("\n TIN IS %s \n", stateTin)
	if stateTin != "" || stateTin != "**" {
		currentStateTax := stateTaxMap[stateTin]
		currentStateTax.StateCode = tinToStateDetailsMap[stateTin].StateCode
		currentStateTax.StateName = tinToStateDetailsMap[stateTin].StateName
		if invoiceRecord.GetInvoiceType() == invoice.Sales {
			currentStateTax.PayableTax = currentStateTax.PayableTax.Add(invoiceRecord.GetTransaction().Tax)
		} else {
			currentStateTax.TaxOnPurchase = currentStateTax.TaxOnPurchase.Add(invoiceRecord.GetTransaction().Tax)
		}
		stateTaxMap[stateTin] = currentStateTax
	}
	return stateTaxMap
}

func populateTaxReport(totalPurchase transaction.Transaction, totalSales transaction.Transaction) taxReport {
	report := taxReport{
		TotalPurchases:  totalPurchase,
		TotalSales:      totalSales,
		TotalInputTax:   totalPurchase.Tax.GetTotalTax(),
		TotalPayableTax: totalSales.Tax.GetTotalTax(),
	}
	report.NetTax = report.TotalInputTax - report.TotalPayableTax
	return report
}

func (r reporterImpl) GetSheetValuesToPublishReport(taxReport taxReport, month time.Month, year int) [][]interface{} {
	values := [][]interface{}{
		{"GSTIN", "24ABWPJ2263R1ZW"},
		{"Name", "Padamchand Jain c/o Jainco Textile"},
		{},
		{"Month", fmt.Sprintf("%s %d", month, year)},
		{},
		{"A )", "Purchases"},
		{},
		{nil, "Total Taxable Amount", nil, taxReport.TotalPurchases.TaxableAmount},
		{nil, "Total C.G.S.T", nil, taxReport.TotalPurchases.Tax.GetCGST()},
		{nil, "Total S.G.S.T", nil, taxReport.TotalPurchases.Tax.GetSGST()},
		{nil, "Total I.G.S.T", nil, taxReport.TotalPurchases.Tax.GetIGST()},
		{},
		{nil, "Total Purchases Amount with tax", nil, taxReport.TotalPurchases.TotalAmountWithTax},
		{},
		{"B )", "Input Tax amount of Purchases"},
		{},
		{nil, "Total C.G.S.T", nil, taxReport.TotalPurchases.Tax.GetCGST()},
		{nil, "Total S.G.S.T", nil, taxReport.TotalPurchases.Tax.GetSGST()},
		{nil, "Total I.G.S.T", nil, taxReport.TotalPurchases.Tax.GetIGST()},
		{},
		{nil, "Total Input Tax amount", nil, taxReport.TotalPurchases.Tax.GetTotalTax()},
		{},
		{"C )", "Sales"},
		{},
		{nil, "Total Sales", nil, taxReport.TotalSales.TotalAmountWithTax},
		{},
		{nil, "Total Sales with Tax amount", nil, taxReport.TotalSales.TotalAmountWithTax},
		{},
		{"D )", "Total Sales as per Following"},
		{},
		{nil, "Total Taxable Amount", nil, taxReport.TotalSales.TaxableAmount},
		{nil, "Total C.G.S.T", nil, taxReport.TotalSales.Tax.GetCGST()},
		{nil, "Total S.G.S.T", nil, taxReport.TotalSales.Tax.GetSGST()},
		{nil, "Total I.G.S.T", nil, taxReport.TotalSales.Tax.GetIGST()},
		{},
		{nil, "Total Sales with tax amount", nil, taxReport.TotalSales.TotalAmountWithTax},
		{},
		{"E )", "Total payable tax amount of sales"},
		{},
		{nil, "Total payable C.G.S.T", nil, taxReport.TotalSales.Tax.GetCGST()},
		{nil, "Total payable S.G.S.T", nil, taxReport.TotalSales.Tax.GetSGST()},
		{nil, "Total payable I.G.S.T", nil, taxReport.TotalSales.Tax.GetIGST()},
		{},
		{nil, "Total payable tax amount", nil, taxReport.TotalSales.Tax.GetTotalTax()},
		{},
		{"F )", "Total payable tax liability"},
		{},
		{nil, "Total input tax amount", nil, taxReport.TotalPurchases.Tax.GetTotalTax()},
		{nil, "Total payable tax amount", nil, taxReport.TotalSales.Tax.GetTotalTax()},
		{},
		{nil, "Total Tax Credit/ Payable", nil, taxReport.NetTax},
		{},
	}

	tinToStateDetailsMap := getTinToStateDetailsMap()
	stateTaxBlock := [][]interface{}{
		{"G )", "Total payable tax as per following states"},
		{},
		{"sr no.", "Name of state", "I.G.S.T", "C.G.S.T", "S.G.S.T", "Total Tax"},
		{},
	}
	for stateCode, taxBreakup := range taxReport.StateTaxBreakup {
		var stateEntry []interface{}
		payableTax := taxBreakup.PayableTax
		if stateCode != "**" && payableTax.GetTotalTax() != 0 {
			stateEntry = append(stateEntry, nil, tinToStateDetailsMap[stateCode].StateName, payableTax.GetIGST(), payableTax.GetCGST(), payableTax.GetSGST(), payableTax.GetTotalTax())
			stateTaxBlock = append(stateTaxBlock, stateEntry)
		}
	}
	totalSalesTax := taxReport.TotalSales.Tax
	stateTaxBlock = append(stateTaxBlock,
		[]interface{}{},
		[]interface{}{nil, "TOTAL", totalSalesTax.GetIGST(), totalSalesTax.GetCGST(), totalSalesTax.GetSGST(), totalSalesTax.GetTotalTax()})

	values = append(values, stateTaxBlock...)

	return values
}

func NewReporter(invoiceService invoice.Service) Reporter {
	return &reporterImpl{
		invoiceService: invoiceService,
	}
}
