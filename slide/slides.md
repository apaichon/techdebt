---
theme: default
background: "#000"
class: "text-center"
highlighter: shiki
lineNumbers: false
info: |
  ## Technical Debt Presentation
  Understanding and Managing Technical Debt in Software Development
drawings:
  persist: false
transition: slide-left
title: Understanding Technical Debt
css: unocss
---

<style>
h1 {
  background-image: linear-gradient(45deg, #4F9CF9, #27C6FF);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  text-shadow: 0 0 20px rgba(79, 156, 249, 0.5);
}

h2 {
  color: #4F9CF9;
}

.slide-content {
  background: radial-gradient(circle at center, #000B1A, #000000);
}

.content {
  color: #FFFFFF;
}

.neon-box {
  border: 2px solid #4F9CF9;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 0 10px #4F9CF9;
  margin: 10px 0;
}

.grid-cols-2 > div {
  @apply neon-box;
}
</style>

# Understanding Technical Debt
<div class="text-white">Managing the Hidden Cost of Software Development</div>

---
layout: two-cols
---
# Table of Contents

1. Understanding Technical Debt
   - Definition and Origins
   - Types of Technical Debt
   - Code Examples

2. Impact Analysis
   - Development Impact
   - Cost and Metrics
   - Real-world Examples

::right::
<br/>
<br/>

3. Managing Technical Debt
   - Prevention Strategies
   - Measurement Tools
   - Best Practices

4. Implementation Strategies
   - Code Reviews
   - Refactoring
   - Testing
   - Quality Gates
---
layout: default
---

# Part 1: What is Technical Debt?

<div class="neon-box">

- A metaphor coined by Ward Cunningham in 1992
- Represents the implied cost of additional rework
- Like financial debt: you borrow time now and pay interest later
- Affects both code quality and development velocity
</div>

<div class="text-center"  style="height: 40%;">
<img src="./images/part1/intro.svg" alt="intro"  style="height: 100%; width: 100%;" />
</div>

ref: [Ward Cunningham](https://techdebtpolicy.com/introduction-to-metaphors)

---
layout: default
---

# Types of Technical Debt

<img src="./images/part1/types.svg" alt="types" style="height: 90%; width: 100%;" />


---
layout: default
---

# Technical Debt in Practice
## Code Debt

<div style="overflow-x: auto; height: 400px;">

```go
// Initial problematic implementation
package loan

type Loan struct {
    ID string
    Amount float64
    Status string // Technical debt: Using string instead of proper type
    InterestRate float64
    CustomerID string
    CreatedAt time.Time
    // Missing proper validation
}

func (l *Loan) Approve() error {
    // Technical debt: No proper state management
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
``` 
</div>

---
layout: default
---

# Architectural Debt
## Violated patterns

<div style="overflow-x: auto; height: 400px;">
```go
// Anti-pattern: Monolithic Controller
type LoanController struct {
    db *sql.DB  // Direct DB dependency
}

func (lc *LoanController) ProcessLoan(w http.ResponseWriter, r *http.Request) {
    // Anti-pattern: Everything in one place
    // Validation, business logic, data access, response handling all mixed
    loan := &Loan{}
    if err := json.NewDecoder(r.Body).Decode(loan); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Direct SQL in controller
    row := lc.db.QueryRow("SELECT credit_score FROM customers WHERE id = ?", loan.CustomerID)
    var creditScore int
    if err := row.Scan(&creditScore); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Business logic mixed with everything else
    if creditScore < 600 {
        http.Error(w, "Credit score too low", http.StatusBadRequest)
        return
    }

    // More mixed concerns...
    if err := lc.db.QueryRow("INSERT INTO loans ..."); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(loan)
}

// Anti-pattern: Global State
var (
    globalDB *sql.DB
    globalConfig *Config
)

func init() {
    globalDB, _ = sql.Open("postgres", "connection-string")
    globalConfig = &Config{/* ... */}
}

// Anti-pattern: Missing Abstractions
type LoanService struct {
    db *sql.DB  // Direct dependency on concrete type
}

// IMPROVED VERSION
// =====================================

// Proper interfaces and separation of concerns
type LoanRepository interface {
    Save(ctx context.Context, loan *Loan) error
    GetByID(ctx context.Context, id string) (*Loan, error)
}

type CustomerService interface {
    GetCreditScore(ctx context.Context, customerID string) (int, error)
}

type LoanValidator interface {
    Validate(loan *Loan) error
}

// Clean controller with dependency injection
type LoanController struct {
    loanService    LoanService
    validator      LoanValidator
    errorHandler   ErrorHandler
}

func NewLoanController(ls LoanService, v LoanValidator, eh ErrorHandler) *LoanController {
    return &LoanController{
        loanService:    ls,
        validator:      v,
        errorHandler:   eh,
    }
}

func (lc *LoanController) ProcessLoan(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    loan, err := lc.parseLoanRequest(r)
    if err != nil {
        lc.errorHandler.Handle(w, err)
        return
    }

    if err := lc.validator.Validate(loan); err != nil {
        lc.errorHandler.Handle(w, err)
        return
    }

    result, err := lc.loanService.Process(ctx, loan)
    if err != nil {
        lc.errorHandler.Handle(w, err)
        return
    }

    lc.renderResponse(w, result)
}
```

</div>

---
layout: default
---

# Architectural Debt
## Poor system design

<div style="overflow-x: auto; height: 400px;">

```go

// 2. POOR SYSTEM DESIGN
// =====================================

// Anti-pattern: Business Logic in SQL
func (ls *LoanService) CalculateInterest(loan *Loan) error {
    // Business logic in SQL
    _, err := ls.db.Exec(`
        UPDATE loans 
        SET interest_rate = 
            CASE 
                WHEN amount > 50000 THEN base_rate + 2
                WHEN amount > 25000 THEN base_rate + 1
                ELSE base_rate
            END
        WHERE id = ?`, loan.ID)
    return err
}

// Improved: Business Logic in Go
type InterestCalculator struct {
    config *Config
}

func (ic *InterestCalculator) Calculate(loan *Loan) (Rate, error) {
    baseRate := ic.config.GetBaseRate()
    
    var adjustment Rate
    switch {
    case loan.Amount > Money(50000):
        adjustment = 200 // 2.00%
    case loan.Amount > Money(25000):
        adjustment = 100 // 1.00%
    }
    
    return baseRate + adjustment, nil
}
```

</div>

---
layout: default
---

# Architectural Debt
## Tight coupling

<div style="overflow-x: auto; height: 400px;">

```go
// 3. TIGHT COUPLING
// =====================================

// Anti-pattern: Direct Database Calls
type LoanProcessor struct {
    db *sql.DB
}

func (lp *LoanProcessor) ProcessLoan(loan *Loan) error {
    tx, err := lp.db.Begin()
    if err != nil {
        return err
    }
    // Direct SQL calls everywhere
    _, err = tx.Exec("INSERT INTO loans ...")
    if err != nil {
        tx.Rollback()
        return err
    }
    return tx.Commit()
}

// Improved: Repository Pattern
type PostgresLoanRepository struct {
    db *sql.DB
}

func (r *PostgresLoanRepository) Save(ctx context.Context, loan *Loan) error {
    return r.db.RunInTx(ctx, func(tx *sql.Tx) error {
        return r.saveLoan(ctx, tx, loan)
    })
}

// Anti-pattern: Hardcoded Dependencies
type CreditCheck struct {
    externalAPI *ThirdPartyAPI
}

func (cc *CreditCheck) Check(customerID string) error {
    // Direct dependency on external service
    return cc.externalAPI.CheckCredit(customerID)
}

// Improved: Dependency Injection + Interface
type CreditChecker interface {
    Check(ctx context.Context, customerID string) error
}

type creditCheckService struct {
    checker CreditChecker
    cache   Cache
}

func NewCreditCheckService(checker CreditChecker, cache Cache) *creditCheckService {
    return &creditCheckService{
        checker: checker,
        cache:   cache,
    }
}

```

</div>

---
layout: default
---

# Architectural Debt
## Complexity

<div style="overflow-x: auto; height: 400px;">

```mermaid
graph TB
    subgraph "Client Layer"
        Web["Web UI<br>(React/TypeScript)"]
        Mobile["Mobile App<br>(iOS/Android)"]
        API_Client["API Client<br>(SDK)"]
    end

    subgraph "API Gateway Layer"
        Gateway["API Gateway<br>(Golang)"]
        Auth["Authentication<br>Service"]
        Rate["Rate Limiter"]
    end

    subgraph "Application Services"
        LoanService["Loan Service<br>(Core Domain)"]
        UserService["User Service"]
        NotifService["Notification<br>Service"]
        
        subgraph "Loan Service Components"
            LoanProcessor["Loan<br>Processor"]
            RiskAnalyzer["Risk<br>Analyzer"]
            InterestCalc["Interest<br>Calculator"]
        end
    end

    subgraph "Integration Services"
        CreditBureau["Credit Bureau<br>Integration"]
        BankAPI["Bank API<br>Integration"]
        IDVerification["ID Verification<br>Service"]
    end

    subgraph "Data Layer"
        MainDB[("Primary DB<br>(PostgreSQL)")]
        Cache[("Redis Cache")]
        DocumentDB[("Document Store<br>(MongoDB)")]
    end

    subgraph "Event Bus"
        Kafka["Message Queue<br>(Kafka)"]
    end

    subgraph "Infrastructure"
        Logging["Logging<br>(ELK Stack)"]
        Metrics["Metrics<br>(Prometheus)"]
        Tracing["Distributed<br>Tracing"]
    end

    %% Client Layer Connections
    Web --> Gateway
    Mobile --> Gateway
    API_Client --> Gateway

    %% Gateway Layer Connections
    Gateway --> Auth
    Gateway --> Rate
    Gateway --> LoanService
    Gateway --> UserService

    %% Service Connections
    LoanService --> LoanProcessor
    LoanService --> RiskAnalyzer
    LoanService --> InterestCalc
    LoanProcessor --> Kafka
    RiskAnalyzer --> CreditBureau
    LoanService --> Cache

    %% Data Connections
    LoanService --> MainDB
    UserService --> MainDB
    NotifService --> DocumentDB

    %% Integration Connections
    LoanProcessor --> BankAPI
    UserService --> IDVerification

    %% Event Flow
    Kafka --> NotifService
    
    %% Infrastructure Connections
    LoanService -.-> Logging
    LoanService -.-> Metrics
    Gateway -.-> Tracing
    UserService -.-> Logging

    %% Styling
    classDef primary fill:#f9f,stroke:#333,stroke-width:2px
    classDef secondary fill:#bbf,stroke:#333
    classDef database fill:#dfd,stroke:#333
    classDef infrastructure fill:#ffd,stroke:#333

    class LoanService,LoanProcessor,RiskAnalyzer,InterestCalc primary
    class UserService,NotifService,Gateway,Auth secondary
    class MainDB,Cache,DocumentDB database
    class Logging,Metrics,Tracing infrastructure
```

</div>
---
layout: default
---

# Intentional vs Unintentional

<div class="grid grid-cols-2 gap-4">
<div>

## Intentional
- Strategic decisions
- Time-to-market pressure
- Conscious trade-offs
- Documented decisions
</div>

<div>

## Unintentional
- Lack of knowledge
- Poor practices
- Outdated skills
- Gradual degradation
</div>
</div>

---
layout: default
---

# Intentional Technical Debt

## 1. Strategic Decisions
**Examples:**
- Using a simpler architecture for MVP launch
- Implementing quick solutions for market opportunities
- Choosing a familiar technology stack over optimal one

**Code Example:**

<div style="overflow-x: auto; height: 200px;">

```go
// Strategic Decision: Quick implementation for MVP
// Technical Debt: Direct database calls in handler
func (h *Handler) GetUserLoans(w http.ResponseWriter, r *http.Request) {
    db := database.GetConnection() // Global DB connection
    rows, err := db.Query("SELECT * FROM loans WHERE user_id = ?", userID)
    // ... handle results
}

// Proper Implementation (Planned for Post-MVP):
type LoanService interface {
    GetUserLoans(ctx context.Context, userID string) ([]Loan, error)
}
```

</div>

---
layout: default
---

### 2. Time-to-Market Pressure
**Examples:**
- Simplified validation logic
- Hardcoded configurations
- Minimal error handling

**Code Example:**

<div style="overflow-x: auto; height: 200px;">

```go
// Time Pressure Decision: Simplified error handling
// Technical Debt: Generic error responses
func ProcessLoan(loan *Loan) error {
    if err := validateLoan(loan); err != nil {
        return err // Generic error, lacks context
    }
    return saveLoan(loan)
}

// Proper Implementation (Planned):
type LoanError struct {
    Code    string
    Message string
    Details map[string]interface{}
}
```

</div>

---
layout: default
---

### 3. Conscious Trade-offs
**Examples:**
- Performance vs. maintainability decisions
- Storage vs. computation trade-offs
- Monolithic vs. microservices approach

**Code Example:**

<div style="overflow-x: auto; height: 200px;">

```go
// Trade-off: Performance over maintainability
// Technical Debt: Complex in-memory caching
var loanCache = make(map[string]*Loan) // Global cache
var cacheMutex sync.RWMutex

func GetLoan(id string) *Loan {
    cacheMutex.RLock()
    defer cacheMutex.RUnlock()
    return loanCache[id]
}

// Proper Implementation (Future):
type Cache interface {
    Get(ctx context.Context, key string) (interface{}, error)
    Set(ctx context.Context, key string, value interface{}) error
}
```

</div>

---
layout: default
---

# Unintentional Technical Debt

## 1. Lack of Knowledge
**Examples:**
- Improper use of language features
- Misunderstanding of design patterns
- Poor architectural decisions

**Code Example:**

<div style="overflow-x: auto; height: 200px;">

```go
// Knowledge Gap: Improper concurrency handling
// Technical Debt: Race conditions
type LoanProcessor struct {
    loans map[string]*Loan // Shared map without proper synchronization
}

func (lp *LoanProcessor) ProcessLoan(loan *Loan) {
    lp.loans[loan.ID] = loan // Race condition
}

// Proper Implementation:
type LoanProcessor struct {
    loans sync.Map
}
```

</div>

---
layout: default
---

## 2. Poor Practices
**Examples:**
- Missing tests
- Poor error handling
- Inconsistent coding standards

**Code Example:**

<div style="overflow-x: auto; height: 200px;">

```go
// Poor Practice: No error handling
// Technical Debt: Silent failures
func UpdateLoanStatus(id string, status string) {
    db.Exec("UPDATE loans SET status = ? WHERE id = ?", status, id)
    // Errors ignored
}

// Proper Implementation:
func UpdateLoanStatus(ctx context.Context, id, status string) error {
    if err := validateStatus(status); err != nil {
        return fmt.Errorf("invalid status: %w", err)
    }
    return db.ExecContext(ctx, "UPDATE loans SET status = ? WHERE id = ?", status, id)
}
```

</div>

--- 
layout: default
---

## 3. Outdated Skills
**Examples:**
- Old design patterns
- Deprecated libraries
- Legacy coding practices

**Code Example:**

<div style="overflow-x: auto; height: 200px;">

```go
// Outdated Skills: Old error handling pattern
// Technical Debt: Non-standard error handling
type Result struct {
    Success bool
    Error   string
}

func ProcessLoan() Result {
    return Result{Success: false, Error: "failed"}
}

// Modern Implementation:
type LoanError struct {
    Code    string
    Message string
}

func ProcessLoan() error {
    return &LoanError{
        Code:    "PROCESS_FAILED",
        Message: "Loan processing failed",
    }
}
```

</div>

---
layout: default
---
# Summary Part 1

<img src="./images/part1/sum1.svg" style="height: 400px; width: 100%;" />

---
layout: cover
---

# Part 2: Impact on Development

<img src="./images/part1/impact-technical-debt.svg" alt="types" style="height: 400px; width: 100%;" />

<!-- Note:
- Decreased velocity over time
- Higher maintenance costs
- Increased bug frequency
- Developer frustration
- Harder to implement new features
- Risk of system instability
-->

---
layout: default
---

# Measuring Technical Debt
<!-- Note:
- Velocity
- Maintainability
- Bugs
- Developers
- Stability
-->


<img src="./images/part1/tech-debt-keymetrics.svg" alt="types" style="height: 400px; width: 100%;" />

---
layout: default
---

# Technical Debt Impact Metrics

## Development Velocity Metrics

| Metric | Previous | Current | Target | Impact Score | Trend |
|--------|----------|---------|--------|--------------|-------|
| Time to Implement (hours) | 24 | 36 | 20 | 1.5 | ↑ |
| Code Review Cycles | 2 | 4 | 2 | 2.0 | ↑ |
| Deployments per Week | 10 | 6 | 12 | 1.7 | ↓ |
| Release Cycle (days) | 2 | 4 | 1 | 2.0 | ↑ |
| Sprint Velocity | 25 | 18 | 30 | 1.4 | ↓ |

---
layout: default
---

# Technical Debt Impact Metrics

## Maintenance Cost Metrics

| Metric | Previous | Current | Target | Cost Impact ($) | Trend |
|--------|----------|---------|--------|----------------|-------|
| Bug Fix Time (hours) | 4 | 8 | 3 | 2000 | ↑ |
| Support Tickets/Week | 10 | 15 | 8 | 1500 | ↑ |
| Documentation Updates | 2 | 5 | 2 | 900 | ↑ |
| Infrastructure Cost | 5000 | 7500 | 4500 | 2500 | ↑ |
| Technical Debt Hours | 20 | 35 | 15 | 3000 | ↑ |

---
layout: default
---

# Technical Debt Impact Metrics

## Bug Impact Metrics

| Severity | Count | Avg Fix Time | Customer Impact | Cost | Trend |
|----------|-------|--------------|-----------------|------|-------|
| P0 (Critical) | 2 | 24h | 1000+ | 10000 | ↑ |
| P1 (High) | 5 | 12h | 500+ | 5000 | ↑ |
| P2 (Medium) | 8 | 6h | 100+ | 2400 | → |
| P3 (Low) | 12 | 3h | 10+ | 900 | ↓ |
| Regressions | 3 | 8h | 200+ | 1500 | ↑ |


---
layout: default
---

# Technical Debt Impact Metrics

## System Stability Metrics

| Metric | Previous | Current | Target | Impact | Trend |
|--------|----------|---------|---------|---------|-------|
| Error Rate (%) | 0.5 | 2.0 | 0.1 | High | ↑ |
| Response Time (ms) | 200 | 500 | 150 | High | ↑ |
| CPU Usage (%) | 60 | 80 | 50 | Medium | ↑ |
| Memory Usage (%) | 65 | 85 | 60 | Medium | ↑ |
| Recovery Time (min) | 5 | 15 | 3 | High | ↑ |


---
layout: default
---

# Technical Debt Impact Metrics

## Developer Experience Metrics

| Metric | Score (1-5) | Previous | Target | Impact | Trend |
|--------|------------|----------|---------|---------|-------|
| Code Maintainability | 2 | 3 | 4 | High | ↓ |
| Documentation Quality | 2 | 3 | 4 | Medium | ↓ |
| Tooling Satisfaction | 3 | 4 | 4 | Low | ↓ |
| Knowledge Sharing | 2 | 3 | 4 | High | ↓ |
| Team Morale | 2 | 4 | 4 | High | ↓ |

---
layout: default
---

# Summary Part 2

<img src="./images/part2/summary.svg" style="height: 400px; width: 100%;" />

---
layout: default
---

# Part 3: Managing Technical Debt

<!-- Note:
1. Regular code reviews
2. Refactoring sprints
3. Boy Scout Rule: "Leave the code better than you found it"
4. Technical debt tracking
5. Set quality gates
6. Automated testing
-->

- Prevention Strategies
- Measurement Tools
- Best Practices

---
layout: default
---
# Prevention Strategies
<img src="./images/part3/prevention-strategies.svg" alt="types" style="height: 100%; width: 100%;" />

---
layout: default
---
# Measurement Tools

<img src="./images/part3/measurement-tools.svg" alt="types" style="height: 100%; width: 100%;" />

---
layout: default
---
# Best Practices

<img src="./images/part3/best-practices.svg" alt="types" style="height: 100%; width: 100%;" />

---
layout: default
---

# Part 4: Implementation Strategies
<img src="./images/part1/steps-to-manage-tech-debt.svg" alt="types" style="height: 400px; width: 100%;" />




---
layout: default
---

# 1. Regular Code Reviews

Code reviews serve as a preventive measure against technical debt accumulation.

<img src="./images/part1/ci-cd-code-review.webp" alt="types" style="height: 400px; width: 100%;" />

---
layout: default
---

# 2. Refactoring Sprints

<img src="./images/part1/refactoring-sprints.svg" alt="types" style="height: 400px; width: 100%;" />

---
layout: default
---

# 3. Boy Scout Rule

<img src="./images/part1/boy-scout-rule.svg" alt="types" style="height: 400px; width: 100%;" />

---
layout: default
---

<div style="overflow-x: auto; height: 600px;">
```go
// BEFORE: Original code with technical debt
package loan

type LoanApp struct {
    id string
    amt float64
    dur int
    inc float64
    cr  int
}

func (l *LoanApp) proc() (bool, error) {
    // Calculate debt-to-income ratio
    dti := l.amt / (l.inc * float64(l.dur))
    
    // Check eligibility
    if dti > 0.43 {
        return false, nil
    }
    if l.cr < 640 {
        return false, nil
    }
    
    // Calculate interest rate
    var r float64
    if l.cr >= 760 {
        r = 0.0325
    } else if l.cr >= 700 {
        r = 0.0425
    } else {
        r = 0.0525
    }
    
    // Calculate monthly payment
    mp := (l.amt * (1 + r * float64(l.dur)/12)) / float64(l.dur)
    
    if mp > l.inc * 0.28 {
        return false, nil
    }
    
    return true, nil
}

// AFTER: Improved code following Boy Scout Rule
package loan

import (
    "errors"
    "fmt"
)

// LoanApplication represents a customer's loan application with all relevant financial information
type LoanApplication struct {
    ID              string
    Amount          float64  // Loan amount in dollars
    DurationMonths  int      // Loan duration in months
    MonthlyIncome   float64  // Applicant's monthly income
    CreditScore     int      // Applicant's credit score
}

// Common errors that can occur during loan processing
var (
    ErrInvalidAmount = errors.New("loan amount must be positive")
    ErrInvalidDuration = errors.New("loan duration must be between 1 and 360 months")
    ErrInvalidIncome = errors.New("monthly income must be positive")
    ErrInvalidCreditScore = errors.New("credit score must be between 300 and 850")
)

// Constants for loan qualification criteria
const (
    MaxDebtToIncomeRatio  = 0.43
    MinCreditScore        = 640
    MaxPaymentToIncome    = 0.28
    
    ExcellentCreditScore  = 760
    GoodCreditScore       = 700
    
    ExcellentRate        = 0.0325
    GoodRate            = 0.0425
    StandardRate        = 0.0525
)

// ProcessApplication evaluates a loan application and determines if it should be approved
// Returns approval status and any validation errors encountered
func (l *LoanApplication) ProcessApplication() (bool, error) {
    // Validate input parameters
    if err := l.validate(); err != nil {
        return false, fmt.Errorf("validation failed: %w", err)
    }
    
    // Check debt-to-income ratio
    if !l.isDebtToIncomeRatioAcceptable() {
        return false, nil
    }
    
    // Verify credit score requirements
    if !l.isCreditScoreAcceptable() {
        return false, nil
    }
    
    // Calculate monthly payment and check affordability
    monthlyPayment := l.calculateMonthlyPayment()
    if !l.isPaymentAffordable(monthlyPayment) {
        return false, nil
    }
    
    return true, nil
}

// validate checks if all loan application fields have valid values
func (l *LoanApplication) validate() error {
    if l.Amount <= 0 {
        return ErrInvalidAmount
    }
    if l.DurationMonths < 1 || l.DurationMonths > 360 {
        return ErrInvalidDuration
    }
    if l.MonthlyIncome <= 0 {
        return ErrInvalidIncome
    }
    if l.CreditScore < 300 || l.CreditScore > 850 {
        return ErrInvalidCreditScore
    }
    return nil
}

// isDebtToIncomeRatioAcceptable checks if the applicant's debt-to-income ratio is within acceptable limits
func (l *LoanApplication) isDebtToIncomeRatioAcceptable() bool {
    debtToIncomeRatio := l.Amount / (l.MonthlyIncome * float64(l.DurationMonths))
    return debtToIncomeRatio <= MaxDebtToIncomeRatio
}

// isCreditScoreAcceptable verifies if the applicant's credit score meets minimum requirements
func (l *LoanApplication) isCreditScoreAcceptable() bool {
    return l.CreditScore >= MinCreditScore
}

// determineInterestRate returns the appropriate interest rate based on credit score
func (l *LoanApplication) determineInterestRate() float64 {
    switch {
    case l.CreditScore >= ExcellentCreditScore:
        return ExcellentRate
    case l.CreditScore >= GoodCreditScore:
        return GoodRate
    default:
        return StandardRate
    }
}

// calculateMonthlyPayment computes the monthly loan payment including interest
func (l *LoanApplication) calculateMonthlyPayment() float64 {
    rate := l.determineInterestRate()
    totalAmount := l.Amount * (1 + rate * float64(l.DurationMonths)/12)
    return totalAmount / float64(l.DurationMonths)
}

// isPaymentAffordable checks if the monthly payment is within the acceptable percentage of income
func (l *LoanApplication) isPaymentAffordable(monthlyPayment float64) bool {
    return monthlyPayment <= l.MonthlyIncome * MaxPaymentToIncome
}

// Additional helper methods for better maintainability

// GetAnnualPercentageRate returns the APR based on credit score
func (l *LoanApplication) GetAnnualPercentageRate() float64 {
    return l.determineInterestRate() * 100
}

// GetMonthlyPaymentEstimate provides an estimate of monthly payments
func (l *LoanApplication) GetMonthlyPaymentEstimate() float64 {
    return l.calculateMonthlyPayment()
}

// GetTotalInterest calculates the total interest paid over the loan term
func (l *LoanApplication) GetTotalInterest() float64 {
    monthlyPayment := l.calculateMonthlyPayment()
    totalPayment := monthlyPayment * float64(l.DurationMonths)
    return totalPayment - l.Amount
}
```
</div>


---
layout: default
---

# 4. Technical Debt Tracking

## Loan System Technical Debt Items

<style>
table {
    width: 100%; /* Make the table take the full width */
    border-collapse: collapse; /* Merge borders */
    font-size: 0.6em; /* Reduce font size */
}

th, td {
    padding: 8px; /* Add padding for better readability */
    text-align: left; /* Align text to the left */
    border: 1px solid #ddd; /* Add border to cells */
}

th {
    background-color: #f2f2f2; /* Light gray background for header */
}

tr:nth-child(even) {
    background-color: #f9f9f9; /* Zebra striping for rows */
}

tr:hover {
    background-color: #f1f1f1; /* Highlight row on hover */
}
</style>

<div style="overflow-x: auto; height: 400px;">

| ID    | Title                        | Description                                                                 | Severity | Area             | Status      | Effort (SP) | Risk  | Priority | Target Date | Assignee    |
|-------|------------------------------|-----------------------------------------------------------------------------|----------|------------------|-------------|-------------|-------|----------|-------------|-------------|
| TD-001| Legacy Loan Calculator        | Outdated formulas and lacks unit tests. Risk of calculation errors.        | HIGH     | Loan Processing   | OPEN        | 8           | HIGH  | P1       | 2025-Q1     | Alice Chen  |
| TD-002| Payment Batch Processing      | Inefficient batch processing causing slowdowns. Needs async implementation. | MEDIUM   | Payment           | IN_PROGRESS | 5           | MEDIUM| P2       | 2025-Q2     | Bob Smith   |
| TD-003| Password Hash Migration       | Using deprecated MD5. Need to migrate to bcrypt/Argon2.                   | CRITICAL | Security          | OPEN        | 3           | HIGH  | P0       | 2025-Q1     | Carol Wu    |
| TD-004| Code Documentation            | Missing API documentation slowing down onboarding.                         | LOW      | Documentation     | OPEN        | 2           | LOW   | P3       | 2025-Q3     | David Lee   |
| TD-005| Database Indexing            | Missing indexes causing slow searches.                                      | MEDIUM   | Performance       | PLANNED     | 3           | MEDIUM| P2       | 2025-Q2     | Emma Davis  |
| TD-006| Error Handling                | Inconsistent error handling in integration. Need standardization.          | HIGH     | Integration       | OPEN        | 5           | HIGH  | P1       | 2025-Q1     | Frank Zhang  |
| TD-007| Test Coverage                | Unit test coverage below 60% in risk assessment module.                   | HIGH     | Testing           | IN_PROGRESS | 8           | MEDIUM| P2       | 2025-Q2     | Grace Kim   |
| TD-008| Legacy Dependencies           | Running on outdated versions with known vulnerabilities.                   | CRITICAL | Security          | PLANNED     | 13          | HIGH  | P0       | 2025-Q1     | Henry Ford  |
| TD-009| Code Duplication             | Duplicate validation logic across loan types. Need to refactor.           | MEDIUM   | Code Quality      | OPEN        | 5           | LOW   | P3       | 2025-Q3     | Ivan Patel  |
| TD-010| Logging System               | Insufficient logging in fraud detection.                                   | HIGH     | Monitoring        | PLANNED     | 5           | MEDIUM| P2       | 2025-Q2     | Julia Santos |

</div>
---

## Priority Levels
- **P0**: Critical - Must be addressed immediately
- **P1**: High - Should be addressed in current quarter
- **P2**: Medium - Plan for next quarter
- **P3**: Low - Address when resources available

---

## Severity Levels
- **CRITICAL**: Major security/business risk
- **HIGH**: Significant impact on quality/performance
- **MEDIUM**: Moderate impact on development
- **LOW**: Minor inconvenience

---

## Status Types
- **OPEN**: Identified but not started
- **PLANNED**: Scheduled for specific sprint/quarter
- **IN_PROGRESS**: Currently being addressed
- **RESOLVED**: Fixed and verified
- **DEFERRED**: Postponed for valid reasons

---

## Area Categories
- Loan Processing
- Payment
- Security
- Documentation
- Performance
- Integration
- Testing
- Code Quality
- Monitoring

---

## Risk Levels
- **HIGH**: Potential for significant business impact
- **MEDIUM**: Moderate business risk
- **LOW**: Minimal business impact

---

## Effort Estimation
**Story Points (SP):**
- **1**: Very small task (~1 day)
- **2**: Small task (~2-3 days)
- **3**: Medium task (~1 week)
- **5**: Large task (~2 weeks)
- **8**: Very large task (~1 month)
- **13**: Extra large task (>1 month)

---
layout: default
---

# 5. Set Quality Gates

## Example SonarQube quality gates:

```yaml
quality_gates:
  - metric: code_coverage
    threshold: 80%
  - metric: duplicate_lines
    threshold: 3%
  - metric: maintainability_rating
    threshold: A
  - metric: reliability_rating
    threshold: A
  - metric: security_rating
    threshold: A
```

---
layout: default
---

# 6.Automated Testing

<!-- Note:
- Unit Tests
- Integration Tests
- Mock Testing
- Performance Testing
- Property-Based Testing
-->

<img src="./images/part1/automated-testing.svg" alt="types" style="height: 400px; width: 100%;" />

---
layout: default
---

# Example Unit Test

<div style="overflow-x: auto; height: 500px;">
```go
package loan

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// LoanApplication represents a loan application
type LoanApplication struct {
	ID            string
	Amount        float64
	Term          int
	APR           float64
	CreditScore   int
	MonthlyIncome float64
	Status        string
}

// LoanService handles loan processing logic
type LoanService struct {
	repo          LoanRepository
	creditChecker CreditChecker
	riskAnalyzer  RiskAnalyzer
}

// Mock interfaces for dependencies
type MockLoanRepository struct {
	mock.Mock
}

type MockCreditChecker struct {
	mock.Mock
}

type MockRiskAnalyzer struct {
	mock.Mock
}

// Unit Tests

func TestLoanApproval(t *testing.T) {
	// Test cases struct
	tests := []struct {
		name          string
		application   LoanApplication
		expectedValid bool
		expectedErr   error
	}{
		{
			name: "Valid Loan Application",
			application: LoanApplication{
				Amount:        100000,
				Term:         360,
				CreditScore:  750,
				MonthlyIncome: 5000,
			},
			expectedValid: true,
			expectedErr:   nil,
		},
		{
			name: "Invalid Amount",
			application: LoanApplication{
				Amount:        0,
				Term:         360,
				CreditScore:  750,
				MonthlyIncome: 5000,
			},
			expectedValid: false,
			expectedErr:   ErrInvalidAmount,
		},
		// Add more test cases...
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewLoanService()
			valid, err := service.ValidateApplication(&tt.application)
			
			assert.Equal(t, tt.expectedValid, valid)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

// Integration test with database
func TestLoanApplicationWorkflow(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup test database
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer db.Close()

	// Create dependencies
	repo := NewLoanRepository(db)
	creditChecker := NewCreditChecker()
	riskAnalyzer := NewRiskAnalyzer()

	// Create service
	service := NewLoanService(repo, creditChecker, riskAnalyzer)

	// Test complete loan application workflow
	ctx := context.Background()
	application := &LoanApplication{
		Amount:        200000,
		Term:         360,
		CreditScore:  780,
		MonthlyIncome: 7000,
	}

	// Step 1: Submit application
	err = service.SubmitApplication(ctx, application)
	assert.NoError(t, err)

	// Step 2: Process application
	err = service.ProcessApplication(ctx, application.ID)
	assert.NoError(t, err)

	// Step 3: Verify final status
	savedApp, err := repo.GetByID(ctx, application.ID)
	assert.NoError(t, err)
	assert.Equal(t, "APPROVED", savedApp.Status)
}

// Benchmark tests
func BenchmarkLoanCalculation(b *testing.B) {
	service := NewLoanService()
	application := &LoanApplication{
		Amount:        300000,
		Term:         360,
		APR:          0.035,
		MonthlyIncome: 8000,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.CalculateMonthlyPayment(application)
	}
}

// Table Driven Tests for Payment Calculations
func TestPaymentCalculations(t *testing.T) {
	tests := []struct {
		name            string
		amount          float64
		term            int
		apr             float64
		expectedPayment float64
	}{
		{
			name:            "30 Year Fixed Rate",
			amount:          300000,
			term:           360,
			apr:            0.035,
			expectedPayment: 1347.13,
		},
		{
			name:            "15 Year Fixed Rate",
			amount:          300000,
			term:           180,
			apr:            0.03,
			expectedPayment: 2071.74,
		},
		// Add more test cases...
	}

	calculator := NewLoanCalculator()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payment := calculator.CalculateMonthlyPayment(tt.amount, tt.term, tt.apr)
			assert.InDelta(t, tt.expectedPayment, payment, 0.01)
		})
	}
}

// Property-Based Testing with quick
func TestLoanCalculatorProperties(t *testing.T) {
	calculator := NewLoanCalculator()

	// Property: Monthly payment should always be positive for valid inputs
	assertion := func(amount float64, term int, apr float64) bool {
		if amount <= 0 || term <= 0 || apr < 0 {
			return true // Skip invalid inputs
		}
		payment := calculator.CalculateMonthlyPayment(amount, term, apr)
		return payment > 0
	}

	if err := quick.Check(assertion, nil); err != nil {
		t.Error("Property test failed:", err)
	}
}

// Test Fixtures and Helpers

// setupTestDB creates a test database
func setupTestDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", "postgres://test:test@localhost:5432/testdb?sslmode=disable")
	if err != nil {
		return nil, err
	}

	// Run migrations
	if err := runMigrations(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

// TestLoanApplication creates a test loan application
func createTestLoanApplication() *LoanApplication {
	return &LoanApplication{
		ID:            "TEST-" + time.Now().Format("20060102150405"),
		Amount:        250000,
		Term:         360,
		APR:          0.0375,
		CreditScore:  720,
		MonthlyIncome: 6000,
		Status:       "PENDING",
	}
}

// Mock Implementation Example
func (m *MockLoanRepository) SaveApplication(ctx context.Context, app *LoanApplication) error {
	args := m.Called(ctx, app)
	return args.Error(0)
}

func (m *MockCreditChecker) CheckCredit(ctx context.Context, ssn string) (int, error) {
	args := m.Called(ctx, ssn)
	return args.Int(0), args.Error(1)
}

// Example test using mocks
func TestLoanProcessingWithMocks(t *testing.T) {
	// Create mocks
	mockRepo := new(MockLoanRepository)
	mockCredit := new(MockCreditChecker)
	mockRisk := new(MockRiskAnalyzer)

	// Setup expectations
	mockRepo.On("SaveApplication", mock.Anything, mock.Anything).Return(nil)
	mockCredit.On("CheckCredit", mock.Anything, "123-45-6789").Return(750, nil)
	mockRisk.On("AnalyzeRisk", mock.Anything, mock.Anything).Return(0.2, nil)

	// Create service with mocks
	service := NewLoanService(mockRepo, mockCredit, mockRisk)

	// Create test application
	app := createTestLoanApplication()

	// Execute test
	err := service.ProcessApplication(context.Background(), app)
	assert.NoError(t, err)

	// Verify expectations
	mockRepo.AssertExpectations(t)
	mockCredit.AssertExpectations(t)
	mockRisk.AssertExpectations(t)
}

// Performance test helper
type PerformanceTest struct {
	Name     string
	Setup    func() interface{}
	Teardown func(interface{})
	Test     func(interface{})
}

func RunPerformanceTest(t *testing.T, test PerformanceTest) {
	data := test.Setup()
	defer test.Teardown(data)

	start := time.Now()
	test.Test(data)
	duration := time.Since(start)

	t.Logf("Performance test '%s' took %v", test.Name, duration)
}

// Example performance test
func TestLoanProcessingPerformance(t *testing.T) {
	perfTest := PerformanceTest{
		Name: "Batch Loan Processing",
		Setup: func() interface{} {
			// Create 1000 test applications
			apps := make([]*LoanApplication, 1000)
			for i := range apps {
				apps[i] = createTestLoanApplication()
			}
			return apps
		},
		Teardown: func(data interface{}) {
			// Cleanup
		},
		Test: func(data interface{}) {
			apps := data.([]*LoanApplication)
			service := NewLoanService()
			
			for _, app := range apps {
				err := service.ProcessApplication(context.Background(), app)
				assert.NoError(t, err)
			}
		},
	}

	RunPerformanceTest(t, perfTest)
}
```
</div>

---
layout: default
---

# Example Integration Test

<div style="overflow-x: auto; height: 500px;">
```go
// Integration test with database
func TestLoanApplicationWorkflow(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup test database
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer db.Close()

	// Create dependencies
	repo := NewLoanRepository(db)
	creditChecker := NewCreditChecker()
	riskAnalyzer := NewRiskAnalyzer()

	// Create service
	service := NewLoanService(repo, creditChecker, riskAnalyzer)

	// Test complete loan application workflow
	ctx := context.Background()
	application := &LoanApplication{    
		Amount:        200000,
		Term:         360,
		CreditScore:  780,
		MonthlyIncome: 7000,
	}
    
	// Step 1: Submit application
	err = service.SubmitApplication(ctx, application)
	assert.NoError(t, err)

	// Step 2: Process application
	err = service.ProcessApplication(ctx, application.ID)
	assert.NoError(t, err)  
```
</div>

---
layout: default
---

# Summary

<img src="./images/part1/summary.png" style="height: 400px; width: 100%;" />

---
layout: end
---

# Thank You

Learn more about managing technical debt:
- [Martin Fowler's Blog](https://martinfowler.com/bliki/TechnicalDebt.html)
- [Ward Cunningham's Wiki](http://wiki.c2.com/?TechnicalDebt)
