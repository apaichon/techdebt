// Initial problematic implementation
package loan

import "time"

type Loan struct {
	ID           string
	Amount       float64
	Status       string
	InterestRate float64
	CustomerID   string
	CreatedAt    time.Time
}

func (l *Loan) Approve() error {
	l.Status = "approved"
	return nil
}

func (l *Loan) CalculateInterest() float64 {
	// Technical debt: Magic numbers
	if l.Amount > 10000 {
		return l.Amount * 0.15
	}
	return l.Amount * 0.12
}
