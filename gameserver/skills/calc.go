package skills

func CapMath(mainVal, twoVal float64, cap string) int {
	defaultVal := 99999
	switch cap {
	case "per": //Добавить N процентов
		return int(mainVal * twoVal / 100)
	}
	return defaultVal
}
