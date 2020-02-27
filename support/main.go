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
func MakeRangeGtLt(gt, lt int) []int {
	a := make([]int, lt-gt-1)
	for i := range a {
		a[i] = gt + 1 + i
	}
	return a
}

//Creates (from, to] range slice
// e.g. for 4, 7 returns [5, 6, 7]
func MakeRangeGtLte(gt, lte int) []int {
	a := make([]int, lte-gt)
	for i := range a {
		a[i] = gt + 1 + i
	}
	return a
}

//Creates [from, to) range slice
// e.g. for 4, 7 returns [4, 5, 6]
func MakeRangeGteLt(gte, lt int) []int {
	a := make([]int, lt-gte)
	for i := range a {
		a[i] = gte + i
	}
	return a
}
