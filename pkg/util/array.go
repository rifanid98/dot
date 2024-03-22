package util

import (
	"golang.org/x/exp/rand"
	"time"
)

func ArrayOfN(n int) []int {
	var arr []int
	for i := 1; i <= n; i++ {
		arr = append(arr, i)
	}
	return arr
}

func CompareAndRemoveSameOne(s []int, r []int) []int {
	var list []int

	m := make(map[int]bool)
	for _, i := range r {
		m[i] = true
	}

	for _, i := range s {
		if _, value := m[i]; !value {
			list = append(list, i)
		}
	}
	return list
}

func RandomIndexFromIntArr(p []int) int {
	rand.Seed(uint64(time.Now().UnixNano()))
	a := rand.Intn(len(p))

	return p[a]
}
