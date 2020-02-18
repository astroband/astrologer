package support

func Difference(a, b []int) (diff []int) {
	m := make(map[int]bool)

	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return
}

func Unique(arr []int) []int {
	keys := make(map[int]bool)
	res := []int{}
	for _, entry := range arr {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			res = append(res, entry)
		}
	}
	return res
}

//Creates (from, to) range slice
// e.g. for 4, 7 returns [5, 6]
func MakeRange(from, to int) []int {
	a := make([]int, to-from-1)
	for i := range a {
		a[i] = from + 1 + i
	}
	return a
}
