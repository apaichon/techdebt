package loan

import (
	"context"
)

// Technical Debt - Architectural Debt:
// - Missing proper layered architecture
// - No clear separation of concerns
// - Missing repository interface
// - No dependency injection
// - Missing proper error handling
// - No context usage for timeouts and cancellation

// LoanRepository interface for data persistence
type LoanRepository interface {
	Save(ctx context.Context, loan *Loan) error
	FindByID(ctx context.Context, id string) (*Loan, error)
	Update(ctx context.Context, loan *Loan) error
}

// LoanService handles loan business logic
type LoanService struct {
	repo LoanRepository
}

// NewLoanService creates a new loan service
func NewLoanService(repo LoanRepository) *LoanService {
	return &LoanService{
		repo: repo,
	}
}

// ProcessLoanApplication handles the loan application process
func (s *LoanService) ProcessLoanApplication(ctx context.Context, loan *Loan) error {
	if err := loan.Validate(); err != nil {
		return err
	}

	// Technical Debt - Missing Features:
	// - Credit score check
	// - Risk assessment
	// - Fraud detection
	// - Compliance checks
	// - Automated approval rules

	return s.repo.Save(ctx, loan)
} 