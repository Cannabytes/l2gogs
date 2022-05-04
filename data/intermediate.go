// Промежуточный пакет, чтоб не повторять код
package data

import (
	"l2gogameserver/data/logger"
	"strconv"
)

func StrToInt(value string) int {
	nvalue, err := strconv.Atoi(value)
	if err != nil {
		logger.Info.Panicln(err)
	}
	return nvalue
}

func StrToInt64(value string) int64 {
	nvalue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		logger.Info.Panicln(err)
	}
	return nvalue
}

func StrToInt32(value string) int32 {
	nvalue, err := strconv.ParseInt(value, 16, 32)
	if err != nil {
		logger.Info.Panicln(err)
	}
	return int32(nvalue)
}

// CalcInt Общее количество всех чисел разом
func CalcInt(args ...int) int {
	total := 0
	for _, v := range args {
		total += v
	}
	return total
}
