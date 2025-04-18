package convert

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	"golang.org/x/exp/constraints"
)

// FormatCurrency is a type that represents the currency format
type FormatCurrency string

const (
	IDR FormatCurrency = "IDR"
	USD FormatCurrency = "USD"
	CNY FormatCurrency = "CNY"
	RUB FormatCurrency = "RUB"
)

// String is a method that returns the string representation of the FormatCurrency
func (f FormatCurrency) String() string {
	return string(f)
}

// formatFunc is a function that formats the currency
type formatFunc func(float64) string

// getMapFormatCurrency is a function that returns a map of currency format
// this is used to get the format function based on the currency
func getMapFormatCurrency() map[FormatCurrency]formatFunc {
	return map[FormatCurrency]formatFunc{
		IDR: formatIdr,
		USD: formatUsd,
		CNY: formatCny,
		RUB: formatRub,
	}
}

// GetFormatCurrencyFromString is a function that returns the FormatCurrency from the string
func GetFormatCurrencyFromString(val string) FormatCurrency {
	return FormatCurrency(val)
}

// Number is a type that represents the number type
type Number interface {
	constraints.Integer | constraints.Float
}

// formatIdr is a function that formats the currency to IDR
func formatIdr(val float64) string {
	if val < 0 {
		return fmt.Sprintf("-Rp%s", strings.ReplaceAll(humanize.Commaf(val*-1), ",", "."))
	}
	return fmt.Sprintf("Rp%s", strings.ReplaceAll(humanize.Commaf(val), ",", "."))
}

// formatUsd is a function that formats the currency to USD
func formatUsd(val float64) string {
	if val < 0 {
		return fmt.Sprintf("-US$%s", humanize.Commaf(val*-1))
	}
	return fmt.Sprintf("US$%s", humanize.Commaf(val))
}

// formatCny is a function that formats the currency to CNY
func formatCny(val float64) string {
	if val < 0 {
		return fmt.Sprintf("-¥%s", humanize.Commaf(val*-1))
	}
	return fmt.Sprintf("¥%s", humanize.Commaf(val))
}

// formatRub is a function that formats the currency to RUB
func formatRub(val float64) string {
	if val < 0 {
		return fmt.Sprintf("-₽%s", humanize.Commaf(val*-1))
	}
	return fmt.Sprintf("₽%s", humanize.Commaf(val))
}

// ToCurrency is a function that converts the number to currency
func ToCurrency[T Number](val T, formats ...FormatCurrency) string {
	v := float64(val)

	// default format is IDR
	format := IDR
	if len(formats) > 0 {
		format = formats[0]
	}

	// when data is found in map, we should use that
	// otherwise we will use idr
	fm, ok := getMapFormatCurrency()[format]
	if ok {
		return fm(v)
	}

	return formatIdr(v)
}

// StringToCurrency is a function that converts the string to currency
func StringToCurrency(val string, formats ...FormatCurrency) (string, error) {
	v, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return "", err
	}

	format := IDR
	if len(formats) > 0 {
		format = formats[0]
	}

	// when data is found in map, we should use that
	// otherwise we will use idr
	fm, ok := getMapFormatCurrency()[format]
	if ok {
		return fm(v), nil
	}

	return formatIdr(v), nil
}
