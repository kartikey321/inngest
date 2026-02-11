package devserver

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseCSVSet(t *testing.T) {
	actual := parseCSVSet(" fn-a,fn-b , ,fn-c,fn-a ")
	require.Len(t, actual, 3)
	require.Contains(t, actual, "fn-a")
	require.Contains(t, actual, "fn-b")
	require.Contains(t, actual, "fn-c")
}
