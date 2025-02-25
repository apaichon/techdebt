package loan

import (
	"errors"
	"time"
)

// Technical Debt - Documentation Debt:
// - Missing package documentation
// - Missing type and function documentation
// - No clear explanation of valid loan statuses
// - No documentation about business rules for interest rates

// Technical Debt - Code Debt:
// - No validation for Amount, InterestRate
// - Status is using magic strings
// - No proper error handling
// - Missing important loan properties like duration, payment schedule
// - No validation for CustomerID

// Loan status constants to replace magic strings
const (
	StatusPending  = "pending"
	StatusApproved = "approved"
	StatusRejected = "rejected"
	StatusDefault  = "default"
)

// Loan represents a financial loan agreement
type Loan struct {
	ID           string
	Amount       float64
	Status       string
	InterestRate float64
	CustomerID   string
	CreatedAt    time.Time
	// Technical Debt - Missing Fields:
	// Duration     int      // Loan duration in months
	// PaymentSchedule []Payment
	// LastModified time.Time
	// ApprovedBy   string
	// Purpose      string
}

// Validate checks if the loan data is valid
func (l *Loan) Validate() error {
	if l.Amount <= 0 {
		return errors.New("loan amount must be positive")
	}
	if l.CustomerID == "" {
		return errors.New("customer ID is required")
	}
	if l.InterestRate < 0 {
		return errors.New("interest rate cannot be negative")
	}
	return nil
}

// Approve changes the loan status to approved
func (l *Loan) Approve() error {
	// Technical Debt - Code Debt:
	// - No validation before approval
	// - No audit trail
	// - No check for valid state transitions
	if err := l.Validate(); err != nil {
		return err
	}
	l.Status = StatusApproved
	return nil
}

// CalculateInterest calculates the interest amount for the loan
func (l *Loan) CalculateInterest() float64 {
	// Technical Debt - Code Debt:
	// - Hard-coded interest rates
	// - No consideration of loan duration
	// - Oversimplified calculation
	// - No risk assessment
	if l.Amount > 10000 {
		return l.Amount * 0.15
	}
	return l.Amount * 0.12
}
