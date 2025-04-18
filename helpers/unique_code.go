package helpers

import (
	"strconv"

	"github.com/mqdvi-dp/go-common/convert"
)

func CountChangedDigits(amount int, numChange string) int {
	count := 3
	strAmount := strconv.Itoa(amount)

	numChangeFloat := convert.StringToFloat(numChange)
	if numChangeFloat < 0 {
		return 0
	}

	amountFloat := convert.StringToFloat(strAmount)
	if amountFloat < 0 {
		return 0
	}
	total := numChangeFloat + amountFloat

	strTotal := strconv.Itoa(int(total))
	if len(strTotal) > len(strAmount) {
		return len(strTotal) + countDotsBasedOnDigits(len(strTotal))
	}

	newTotal := strTotal[:len(strTotal)-3]
	for i := 0; i < len(newTotal); i++ {
		if newTotal[i] != strAmount[i] {
			count++
		}
	}

	return count + countDotsBasedOnDigits(count)
}

func countDotsBasedOnDigits(digits int) int {
	// Jumlah titik ditambahkan setiap 3 digit dari kanan
	if digits <= 3 {
		return 0
	}

	// Jumlah titik yang diperlukan
	dotsCount := (digits - 1) / 3

	return dotsCount
}
