package simplecov

import (
	"testing"

	"github.com/codeclimate/test-reporter/env"
	"github.com/stretchr/testify/require"
)

func Test_Format(t *testing.T) {
	gb := env.GitBlob
	defer func() { env.GitBlob = gb }()
	env.GitBlob = func(s string) (string, error) {
		return s, nil
	}

	r := require.New(t)

	rb := Formatter{
		Path: "./simplecov-example.json",
	}
	err := rb.Parse()
	r.NoError(err)

	rep, err := rb.Format()
	r.NoError(err)

	r.InDelta(97.95, rep.CoveredPercent, 1)

	r.Len(rep.SourceFiles, len(rb.Tests[0].SourceFiles))

	sf := rb.Tests[0].SourceFiles[0]
	r.InDelta(100, sf.CoveredPercent, 1)

	lc := rep.LineCounts
	r.Equal(lc.Covered, 56)
	r.Equal(lc.Missed, 1)
	r.Equal(lc.Total, 57)
}