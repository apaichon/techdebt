package main

import (
	"fmt"
	"i-loan/loan"
)

func main() {
	loan := loan.Loan{
		Amount: 10000,
	}
	loan.Approve()
	fmt.Println(loan.Status)
	fmt.Println(loan.CalculateInterest())

}
