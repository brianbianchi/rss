package data

import "testing"

func TestGenerateRandomID(t *testing.T) {
	n := 6
	result := generateRandomID(n)
	if len(result) != n {
		t.Fatalf(`RandSeq(%d) returned %s, which is not %d characters long.`, n, result, n)
	}
	n = 10
	result = generateRandomID(n)
	if len(result) != n {
		t.Fatalf(`RandSeq(%d) returned %s, which is not %d characters long.`, n, result, n)
	}
}
