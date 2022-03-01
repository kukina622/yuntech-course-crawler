package utils

func RemoveIndex(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func FindIndex(s []string, item string) int {
	for i, j := range s {
		if j == item {
			return i
		}
	}
	return -1
}
