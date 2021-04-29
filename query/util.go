package query

func Zip(v1 []string, v2 []string) [][]string {
	values := make([][]string, 0, len(v1))
	for i, value := range v1 {
		values = append(values, []string{value, v2[i]})
	}
	return values
}
