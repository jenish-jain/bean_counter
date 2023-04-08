package reporter

import (
	"bean_counter/internal/types/invoice"
	"bean_counter/internal/types/tax"
	"bean_counter/internal/types/transaction"
	"fmt"
	"time"
)

type Reporter interface {
	MonthlyReporter(month time.Month, year int, purchases []invoice.Invoice, sales []transaction.Transaction) transaction.Transaction
}

type reporterImpl struct {
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

func (r reporterImpl) MonthlyReporter(month time.Month, year int, purchases []invoice.Invoice, sales []transaction.Transaction) transaction.Transaction {
	var totalPurchase transaction.Transaction
	for _, purchaseInvoice := range purchases {
		// make a method for this add transaction
		if purchaseInvoice.GetDate().Month() == month && purchaseInvoice.GetDate().Year() == year {
			totalPurchase = totalPurchase.Add(purchaseInvoice.GetTransaction())
		}
	}
	println("TOTAL PURCHASE AMOUNT")
	fmt.Printf("FINAL VAL %f \n", totalPurchase.TaxableAmount)
	return totalPurchase
}

func NewReporter() Reporter {
	return &reporterImpl{}
}
