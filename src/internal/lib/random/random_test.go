package random

import (
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
)

func TestString(t *testing.T) {
	n := rand.Intn(100)
	str1, err := String(n)
	require.NoError(t, err)
	str2, err := String(n)
	require.NoError(t, err)
	if str1 == str2 {
		t.Fatalf("duplicate result")
	}

	str1, err = String(0)
	require.NoError(t, err)
	str2, err = String(0)
	require.NoError(t, err)

	if str1 != str2 {
		t.Fatal("\"\" != \"\"")
	}
}
