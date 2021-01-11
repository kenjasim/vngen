package hash

import (
	"hash/fnv"
)

// Hash - Return unencrypted hash of input string
func Hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
