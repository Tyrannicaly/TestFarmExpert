package main

import (
	"fmt"
	"math/big"
	"strconv"
)

// Invoice struct represents the invoice data.
type Invoice struct {
	LineExtensionAmount string   `json:"LineExtensionAmount"`
	TaxTotal            TaxTotal `json:"TaxTotal"`
}

// TaxTotal struct represents tax information.
type TaxTotal struct {
	TaxAmount   string        `json:"TaxAmount"`
	TaxSubtotal []TaxSubtotal `json:"TaxSubtotal"`
}

// TaxSubtotal struct represents tax subcategories.
type TaxSubtotal struct {
	TaxableAmount string      `json:"TaxableAmount"`
	TaxAmount     string      `json:"TaxAmount"`
	TaxCategory   TaxCategory `json:"TaxCategory"`
}

// TaxCategory struct represents the tax category.
type TaxCategory struct {
	ID      string `json:"ID"`
	Percent string `json:"Percent"`
}

// parseToInt64 converts a string with two decimal places to int64 (e.g., "1500.00" -> 150000).
func parseToInt64(value string) (int64, error) {
	// Use big.Float for precise parsing.
	f, ok := new(big.Float).SetString(value)
	if !ok {
		return 0, fmt.Errorf("invalid value: %s", value)
	}

	// Multiply by 100 to convert to smallest units (cents).
	f.Mul(f, big.NewFloat(100))
	result, _ := f.Int64() // Get int64 value.
	return result, nil
}

// CalculateInvoiceTotal calculates the total invoice amount and validates data.
func CalculateInvoiceTotal(invoice Invoice) (int64, error) {
	// Parse LineExtensionAmount.
	lineExtAmount, err := parseToInt64(invoice.LineExtensionAmount)
	if err != nil {
		return 0, fmt.Errorf("invalid LineExtensionAmount: %v", err)
	}

	// Parse total tax amount.
	totalTaxAmount, err := parseToInt64(invoice.TaxTotal.TaxAmount)
	if err != nil {
		return 0, fmt.Errorf("invalid TaxAmount in TaxTotal: %v", err)
	}

	// Validate TaxSubtotal.
	var sumTaxSubtotal int64
	for i, subtotal := range invoice.TaxTotal.TaxSubtotal {
		// Parse values.
		taxableAmount, err := parseToInt64(subtotal.TaxableAmount)
		if err != nil {
			return 0, fmt.Errorf("invalid TaxableAmount in TaxSubtotal[%d]: %v", i, err)
		}
		taxAmount, err := parseToInt64(subtotal.TaxAmount)
		if err != nil {
			return 0, fmt.Errorf("invalid TaxAmount in TaxSubtotal[%d]: %v", i, err)
		}
		percent, err := strconv.Atoi(subtotal.TaxCategory.Percent) // Parse Percent as integer (e.g., "20" -> 20).
		if err != nil {
			return 0, fmt.Errorf("invalid Percent in TaxSubtotal[%d]: %v", i, err)
		}

		// Check: TaxableAmount * Percent / 100 = TaxAmount.
		// taxableAmount is in cents (1000.00 -> 100000), percent is an integer (20).
		expectedTaxAmount := (taxableAmount * int64(percent)) / 100 // Divide by 100 to adjust for cents.
		if expectedTaxAmount != taxAmount {
			return 0, fmt.Errorf("incorrect tax calculation in TaxSubtotal[%d]: %d * %d / 100 != %d",
				i, taxableAmount, percent, taxAmount)
		}

		sumTaxSubtotal += taxAmount
	}

	// Check: sum of TaxAmount from TaxSubtotal equals total TaxAmount.
	if sumTaxSubtotal != totalTaxAmount {
		return 0, fmt.Errorf("total tax amount (%d) does not match TaxSubtotal sum (%d)",
			totalTaxAmount, sumTaxSubtotal)
	}

	// Calculate total amount: LineExtensionAmount + TaxAmount.
	totalAmount := lineExtAmount + totalTaxAmount

	return totalAmount, nil
}

// formatInt64ToString converts int64 back to a string with two decimal places.
func formatInt64ToString(value int64) string {
	// Create a big.Float from the int64 value.
	f := new(big.Float).SetInt64(value)

	// Divide by 100 to convert from cents to dollars with decimals.
	hundred := big.NewFloat(100)
	f.Quo(f, hundred)

	// Format the result with exactly two decimal places.
	return f.Text('f', 2)
}

func main() {
	// Example data from the task.
	invoice := Invoice{
		LineExtensionAmount: "1500.00",
		TaxTotal: TaxTotal{
			TaxAmount: "350.00",
			TaxSubtotal: []TaxSubtotal{
				{
					TaxableAmount: "1000.00",
					TaxAmount:     "200.00",
					TaxCategory: TaxCategory{
						ID:      "S",
						Percent: "20",
					},
				},
				{
					TaxableAmount: "500.00",
					TaxAmount:     "50.00",
					TaxCategory: TaxCategory{
						ID:      "R",
						Percent: "10",
					},
				},
			},
		},
	}

	// Call the function.
	total, err := CalculateInvoiceTotal(invoice)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Println("Total invoice amount: 0.00")
	} else {
		fmt.Printf("Total invoice amount: %s\n", formatInt64ToString(total))
	}
}
