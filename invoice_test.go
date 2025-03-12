package main

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestParseToInt64 tests the parseToInt64 function with various inputs.
func TestParseToInt64(t *testing.T) {
	tests := []struct {
		input   string
		want    int64
		wantErr bool
		errMsg  string
	}{
		{"1500.00", 150000, false, ""},            // Valid input.
		{"0.00", 0, false, ""},                    // Zero value.
		{"123.45", 12345, false, ""},             // Non-round number.
		{"-1500.00", -150000, false, ""},         // Negative value.
		{"abc", 0, true, "invalid value: abc"},   // Invalid string.
		{"", 0, true, "invalid value: "},         // Empty string.
		{"1.234", 123, false, ""},                // More than two decimals (truncated).
		{"999999999999.99", 99999999999999, false, ""}, // Large number.
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := parseToInt64(tt.input)
			if tt.wantErr {
				require.Error(t, err, "expected an error")
				require.Equal(t, tt.errMsg, err.Error(), "error message mismatch")
			} else {
				require.NoError(t, err, "unexpected error")
			}
			require.Equal(t, tt.want, got, "parseToInt64 result mismatch")
		})
	}
}

// TestFormatInt64ToString tests the formatInt64ToString function with various inputs.
func TestFormatInt64ToString(t *testing.T) {
	tests := []struct {
		input int64
		want  string
	}{
		{150000, "1500.00"},    // Standard positive value.
		{0, "0.00"},            // Zero.
		{12345, "123.45"},      // Non-round number.
		{-150000, "-1500.00"},  // Negative value.
		{50, "0.50"},           // Small positive value.
		{-50, "-0.50"},         // Small negative value.
		{99999999999999, "999999999999.99"}, // Large number.
	}

	for _, tt := range tests {
		t.Run(strconv.FormatInt(tt.input, 10), func(t *testing.T) {
			got := formatInt64ToString(tt.input)
			require.Equal(t, tt.want, got, "formatInt64ToString result mismatch")
		})
	}
}

// TestCalculateInvoiceTotal tests the CalculateInvoiceTotal function with various scenarios.
func TestCalculateInvoiceTotal(t *testing.T) {
	tests := []struct {
		name    string
		invoice Invoice
		want    int64
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid invoice",
			invoice: Invoice{
				LineExtensionAmount: "1500.00",
				TaxTotal: TaxTotal{
					TaxAmount: "250.00",
					TaxSubtotal: []TaxSubtotal{
						{TaxableAmount: "1000.00", TaxAmount: "200.00", TaxCategory: TaxCategory{ID: "S", Percent: "20"}},
						{TaxableAmount: "500.00", TaxAmount: "50.00", TaxCategory: TaxCategory{ID: "R", Percent: "10"}},
					},
				},
			},
			want:    175000, // 1500.00 + 250.00 = 1750.00 (in cents).
			wantErr: false,
			errMsg:  "",
		},
		{
			name: "Invalid LineExtensionAmount",
			invoice: Invoice{
				LineExtensionAmount: "abc",
				TaxTotal: TaxTotal{
					TaxAmount: "250.00",
					TaxSubtotal: []TaxSubtotal{
						{TaxableAmount: "1000.00", TaxAmount: "200.00", TaxCategory: TaxCategory{ID: "S", Percent: "20"}},
					},
				},
			},
			want:    0,
			wantErr: true,
			errMsg:  "invalid LineExtensionAmount: invalid value: abc",
		},
		{
			name: "Invalid TaxAmount in TaxTotal",
			invoice: Invoice{
				LineExtensionAmount: "1500.00",
				TaxTotal: TaxTotal{
					TaxAmount: "xyz",
					TaxSubtotal: []TaxSubtotal{
						{TaxableAmount: "1000.00", TaxAmount: "200.00", TaxCategory: TaxCategory{ID: "S", Percent: "20"}},
					},
				},
			},
			want:    0,
			wantErr: true,
			errMsg:  "invalid TaxAmount in TaxTotal: invalid value: xyz",
		},
		{
			name: "Invalid TaxableAmount in TaxSubtotal",
			invoice: Invoice{
				LineExtensionAmount: "1500.00",
				TaxTotal: TaxTotal{
					TaxAmount: "200.00",
					TaxSubtotal: []TaxSubtotal{
						{TaxableAmount: "abc", TaxAmount: "200.00", TaxCategory: TaxCategory{ID: "S", Percent: "20"}},
					},
				},
			},
			want:    0,
			wantErr: true,
			errMsg:  "invalid TaxableAmount in TaxSubtotal[0]: invalid value: abc",
		},
		{
			name: "Invalid TaxAmount in TaxSubtotal",
			invoice: Invoice{
				LineExtensionAmount: "1500.00",
				TaxTotal: TaxTotal{
					TaxAmount: "200.00",
					TaxSubtotal: []TaxSubtotal{
						{TaxableAmount: "1000.00", TaxAmount: "def", TaxCategory: TaxCategory{ID: "S", Percent: "20"}},
					},
				},
			},
			want:    0,
			wantErr: true,
			errMsg:  "invalid TaxAmount in TaxSubtotal[0]: invalid value: def",
		},
		{
			name: "Invalid Percent in TaxSubtotal",
			invoice: Invoice{
				LineExtensionAmount: "1500.00",
				TaxTotal: TaxTotal{
					TaxAmount: "200.00",
					TaxSubtotal: []TaxSubtotal{
						{TaxableAmount: "1000.00", TaxAmount: "200.00", TaxCategory: TaxCategory{ID: "S", Percent: "gh"}},
					},
				},
			},
			want:    0,
			wantErr: true,
			errMsg:  "invalid Percent in TaxSubtotal[0]: strconv.Atoi: parsing \"gh\": invalid syntax",
		},
		{
			name: "Incorrect tax calculation",
			invoice: Invoice{
				LineExtensionAmount: "1500.00",
				TaxTotal: TaxTotal{
					TaxAmount: "200.00",
					TaxSubtotal: []TaxSubtotal{
						{TaxableAmount: "1000.00", TaxAmount: "250.00", TaxCategory: TaxCategory{ID: "S", Percent: "20"}},
					},
				},
			},
			want:    0,
			wantErr: true,
			errMsg:  "incorrect tax calculation in TaxSubtotal[0]: 100000 * 20 / 100 != 25000",
		},
		{
			name: "Tax sum mismatch",
			invoice: Invoice{
				LineExtensionAmount: "1500.00",
				TaxTotal: TaxTotal{
					TaxAmount: "350.00",
					TaxSubtotal: []TaxSubtotal{
						{TaxableAmount: "1000.00", TaxAmount: "200.00", TaxCategory: TaxCategory{ID: "S", Percent: "20"}},
						{TaxableAmount: "500.00", TaxAmount: "50.00", TaxCategory: TaxCategory{ID: "R", Percent: "10"}},
					},
				},
			},
			want:    0,
			wantErr: true,
			errMsg:  "total tax amount (35000) does not match TaxSubtotal sum (25000)",
		},
		{
			name: "Negative values",
			invoice: Invoice{
				LineExtensionAmount: "-1500.00",
				TaxTotal: TaxTotal{
					TaxAmount: "-250.00",
					TaxSubtotal: []TaxSubtotal{
						{TaxableAmount: "-1000.00", TaxAmount: "-200.00", TaxCategory: TaxCategory{ID: "S", Percent: "20"}},
						{TaxableAmount: "-500.00", TaxAmount: "-50.00", TaxCategory: TaxCategory{ID: "R", Percent: "10"}},
					},
				},
			},
			want:    -175000, // -1500.00 + (-250.00) = -1750.00 (in cents).
			wantErr: false,
			errMsg:  "",
		},
		{
			name: "Zero values",
			invoice: Invoice{
				LineExtensionAmount: "0.00",
				TaxTotal: TaxTotal{
					TaxAmount: "0.00",
					TaxSubtotal: []TaxSubtotal{
						{TaxableAmount: "0.00", TaxAmount: "0.00", TaxCategory: TaxCategory{ID: "S", Percent: "20"}},
					},
				},
			},
			want:    0,
			wantErr: false,
			errMsg:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateInvoiceTotal(tt.invoice)
			if tt.wantErr {
				require.Error(t, err, "expected an error")
				require.Equal(t, tt.errMsg, err.Error(), "error message mismatch")
			} else {
				require.NoError(t, err, "unexpected error")
			}
			require.Equal(t, tt.want, got, "CalculateInvoiceTotal result mismatch")
		})
	}
}
