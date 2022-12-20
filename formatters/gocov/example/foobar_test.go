package example_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/TomRobbins1967/test-reporter/formatters/gocov/example/bar"
	"github.com/TomRobbins1967/test-reporter/formatters/gocov/example/foo"
)

func Test_String(t *testing.T) {
	foo := foo.New(bar.New())
	require.Equal(t, "foo bar", foo.String())
}
