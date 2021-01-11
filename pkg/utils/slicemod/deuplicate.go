package slicemod

import (
	hash "nenvoy.com/pkg/utils/hash"
)

func Deduplicate(stringSlice []string) (uniqueSlice []string) {
	a := make(map[uint32]int, 0)
	for _, s := range stringSlice {
		hash := hash.Hash(s)
		if _, ok := a[hash]; !ok {
			a[hash] = 1
			uniqueSlice = append(uniqueSlice, s)
		}
	}
	return uniqueSlice
}
