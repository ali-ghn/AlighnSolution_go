package cryptography

import (
	"fmt"
	"testing"
)

func TestHash(t *testing.T) {
	data := "Hello world"
	hashedData := Hash([]byte(data))
	correctHash := "64ec88ca00b268e5ba1a35678a1b5316d212f4f366b2477232534a8aeca37f3c"
	if fmt.Sprintf("%x", hashedData) != correctHash {
		t.Errorf("expected %v, got %v", correctHash, fmt.Sprintf("%x", hashedData))
	}
}
