package tax

import "math"

type Tax struct {
	cgst  float64
	sgst  float64
	igst  float64
	total float64
}

func New(cgst float64, sgst float64, igst float64) Tax {
	/*
		Add validation later that you can either have only igst or both sgst + cgst
	*/
	return Tax{
		cgst:  cgst,
		sgst:  sgst,
		igst:  igst,
		total: cgst + sgst + igst,
	}
}

func (t1 Tax) GetTotalTax() float64 {
	return t1.total
}

func (t1 Tax) GetCGST() float64 {
	return t1.cgst
}

func (t1 Tax) GetSGST() float64 {
	return t1.sgst
}

func (t1 Tax) GetIGST() float64 {
	return t1.igst
}

// Add : sum tax object t1 and t2
func (t1 Tax) Add(t2 Tax) Tax {
	return Tax{
		cgst:  t1.cgst + t2.cgst,
		sgst:  t1.sgst + t2.sgst,
		igst:  t1.igst + t2.igst,
		total: t1.total + t2.total,
	}
}

// Subtract : subtracts tax object t1 from t2
func (t1 Tax) Subtract(t2 Tax) Tax {
	return Tax{
		cgst:  t1.cgst - t2.cgst,
		sgst:  t1.sgst - t2.sgst,
		igst:  t1.igst - t2.igst,
		total: t1.total - t2.total,
	}
}

func (t1 Tax) GetAbsoluteValue() Tax {
	return Tax{
		cgst:  math.Abs(t1.cgst),
		sgst:  math.Abs(t1.sgst),
		igst:  math.Abs(t1.igst),
		total: math.Abs(t1.total),
	}
}
