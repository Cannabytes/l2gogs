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
