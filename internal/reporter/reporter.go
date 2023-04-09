package reporter

import (
	"bean_counter/internal/types/invoice"
	"bean_counter/internal/types/tax"
	"bean_counter/internal/types/transaction"
	"context"
	"time"
)

type Reporter interface {
	GetTaxReportOfMonth(month time.Month, year int) taxReport
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

func NewReporter(invoiceService invoice.Service) Reporter {
	return &reporterImpl{
		invoiceService: invoiceService,
	}
}
