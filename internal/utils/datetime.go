package utils

import (
	"fmt"
	"strconv"
	"time"
)

func CurrentYearMonthDay() (string, string, string) {
	year, month, day := time.Now().Date()
	return strconv.Itoa(year), fmt.Sprintf("%02d", int(month)), fmt.Sprintf("%02d", day)
}
