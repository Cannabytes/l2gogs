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

// CalcInt Общее количество всех чисел разом
func CalcInt(args ...int) int {
	total := 0
	for _, v := range args {
		total += v
	}
	return total
}
