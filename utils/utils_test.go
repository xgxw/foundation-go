package utils

import "testing"

func TestRandN(t *testing.T) {
	n := 10
	s := RandN(n, 0)
	t.Logf("test randN. n:%d, s: %s", n, s)
	if len(s) != n {
		t.FailNow()
	}
}
