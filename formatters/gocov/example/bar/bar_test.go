package bar_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/TomRobbins1967/test-reporter/formatters/gocov/example/bar"
)

func Test_String(t *testing.T) {
	bar := bar.New()
	require.Equal(t, "bar", bar.String())
}
