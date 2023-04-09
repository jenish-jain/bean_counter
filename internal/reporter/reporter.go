package reporter

import (
	"bean_counter/internal/types/invoice"
	"bean_counter/internal/types/tax"
	"bean_counter/internal/types/transaction"
	"context"
	"fmt"
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
	stateName string
	stateCode string
	tax       tax.Tax
}

type taxReport struct {
	Month           time.Month
	Year            int
	TotalPurchases  transaction.Transaction
	TotalSales      transaction.Transaction
	TotalInputTax   float64
	TotalPayableTax float64
	NetTax          float64
	StateTaxBreakup []stateTax
}

func (r reporterImpl) GetTaxReportOfMonth(month time.Month, year int) taxReport {
	ctx := context.Background()
	purchases := r.invoiceService.GetAllPurchaseInvoices(ctx)
	sales := r.invoiceService.GetAllSalesInvoices(ctx)

	var purchasesOfMonth transaction.Transaction
	var salesOfMonth transaction.Transaction
	for _, purchaseInvoice := range purchases {
		if purchaseInvoice.GetDate().Month() == month && purchaseInvoice.GetDate().Year() == year {
			purchasesOfMonth = purchasesOfMonth.Add(purchaseInvoice.GetTransaction())
		}
	}
	for _, salesInvoice := range sales {
		if salesInvoice.GetDate().Month() == month && salesInvoice.GetDate().Year() == year {
			salesOfMonth = salesOfMonth.Add(salesInvoice.GetTransaction())
		}
	}
	report := populateTaxReport(purchasesOfMonth, salesOfMonth)
	report.Month = month
	report.Year = year
	return report

}

func populateTaxReport(totalPurchase transaction.Transaction, totalSales transaction.Transaction) taxReport {
	report := taxReport{
		TotalPurchases:  totalPurchase,
		TotalSales:      totalSales,
		TotalInputTax:   totalPurchase.Tax.GetTotalTax(),
		TotalPayableTax: totalSales.Tax.GetTotalTax(),
	}
	report.NetTax = report.TotalInputTax - report.TotalPayableTax
	// check on state tax report
	return report
}

func (r reporterImpl) GetSheetValuesToPublishReport(taxReport taxReport, month time.Month, year int) [][]interface{} {
	values := [][]interface{}{
		{"GSTIN", "GSTNO"},
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
		{nil, "Amazon Sales", nil, 0},
		{nil, "Flipkart Sales", nil, 0},
		{nil, "Offline Sales", nil, taxReport.TotalSales.TotalAmountWithTax},
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
	}
	return values
}

func NewReporter(invoiceService invoice.Service) Reporter {
	return &reporterImpl{
		invoiceService: invoiceService,
	}
}
