package utils

func RemoveEmptyStr(a []string) []string {
	var result []string
	for _, s := range a {
		if s != "" {
			result = append(result, s)
		}
	}
	return result
}
