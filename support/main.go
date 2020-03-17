package support

import (
	"fmt"
)

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

//Creates [from, to] range slice
// e.g. for 4, 7 returns [4, 5, 6, 7]
func MakeRangeGteLte(gte, lte int) []int {
	a := make([]int, lte-gte+1)
	for i := range a {
		a[i] = gte + i
	}
	return a
}

func ByteCountBinary(b int) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}
