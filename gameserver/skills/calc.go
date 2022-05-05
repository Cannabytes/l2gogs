package skills

func CapMath(mainVal int, twoVal int, cap string) int {
	defaultVal := 99999
	switch cap {
	case "per": //Добавить N процентов
		return mainVal * twoVal / 100
	}
	return defaultVal
}
