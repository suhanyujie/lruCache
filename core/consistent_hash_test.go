package core

import (
	"strconv"
	"testing"
)

func TestHash1(t *testing.T) {
	hash := NewMap(4, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})
	hash.Add("5", "6", "8")
	hash.Add("9")
	testCases := map[string]string{
		"5": "5",
		"15": "5",
		"29": "9",
	}
	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}
}
