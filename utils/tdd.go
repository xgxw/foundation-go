package utils

import "os"

func SkipTest() bool {
	if debug := os.Getenv("debug"); debug == "debug" {
		return false
	}
	return true
}
