package utils

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func RandN(n int, seed int64) string {
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	result := make([]string, n)
	for i := 0; i < n; i++ {
		result[i] = strconv.FormatInt(rand.Int63n(10), 10)
	}
	return strings.Join(result, "")
}
