package transaction

import (
	"bean_counter/internal/types/tax"
)

//type transactionType int
//
//const (
//	Purchase transactionType = iota
//	Sales
//)

type Transaction struct {
	//Type               transactionType
	//Channels           []string
	TaxableAmount      float64
	Tax                tax.Tax
	TotalAmountWithTax float64
}

func New(taxableAmount float64, tax tax.Tax) Transaction {
	return Transaction{
		TaxableAmount:      taxableAmount,
		Tax:                tax,
		TotalAmountWithTax: taxableAmount + tax.GetTotalTax(),
	}
}

// Add : adds transaction t1 with t2
func (t1 Transaction) Add(t2 Transaction) Transaction {
	//if t1.Type == t2.Type {
	return Transaction{
		//Channels:           append(append([]string{}, t1.Channels...), t2.Channels...), // TODO replace this with set later
		TaxableAmount:      t1.TaxableAmount + t2.TaxableAmount,
		Tax:                t1.Tax.Add(t2.Tax),
		TotalAmountWithTax: t1.TotalAmountWithTax + t2.TotalAmountWithTax,
	}

}

//func (t1 Transaction) AddPurchase(purchase Transaction) Transaction {
//	if purchase.Type != Purchase {
//		panic("invalid input")
//	}
//	var output Transaction
//
//	if t1.Type == Purchase {
//		output = Transaction{
//			Type:          Purchase,
//			Channels:      append(append([]string{}, t1.Channels...), purchase.Channels...), // TODO replace this with set later
//			TaxableAmount: t1.TaxableAmount + purchase.TaxableAmount,
//			Tax:           t1.Tax.Add(purchase.Tax),
//		}
//	} else {
//		output = Transaction{
//			Type:          t1.Type,
//			Channels:      append(append([]string{}, t1.Channels...), purchase.Channels...), // TODO replace this with set later
//			TaxableAmount: purchase.TaxableAmount - t1.TaxableAmount,
//			Tax:           t1.Tax.Add(purchase.Tax),
//		}
//		if output.TaxableAmount < 0 {
//			output.Type = Sales
//			output.TaxableAmount = math.Abs(output.TaxableAmount)
//			output.Tax = output.Tax.GetAbsoluteValue()
//		}
//	}
//	return output
//
//}
