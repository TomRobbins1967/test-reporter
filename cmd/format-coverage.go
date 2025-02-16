package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/TomRobbins1967/test-reporter/formatters"
	"github.com/TomRobbins1967/test-reporter/formatters/clover"
	"github.com/TomRobbins1967/test-reporter/formatters/cobertura"
	"github.com/TomRobbins1967/test-reporter/formatters/coveragepy"
	"github.com/TomRobbins1967/test-reporter/formatters/excoveralls"
	"github.com/TomRobbins1967/test-reporter/formatters/gcov"
	"github.com/TomRobbins1967/test-reporter/formatters/gocov"
	"github.com/TomRobbins1967/test-reporter/formatters/jacoco"
	"github.com/TomRobbins1967/test-reporter/formatters/lcov"
	"github.com/TomRobbins1967/test-reporter/formatters/lcovjson"
	"github.com/TomRobbins1967/test-reporter/formatters/simplecov"
	"github.com/TomRobbins1967/test-reporter/formatters/xccov"
	"github.com/gobuffalo/envy"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type CoverageFormatter struct {
	CoveragePath string
	In           formatters.Formatter
	InputType    string
	Output       string
	Prefix       string
	AddPrefix    string
	writer       io.Writer
}

var formatOptions = CoverageFormatter{}

// a prioritized list of the formatters to use
var formatterList = []string{"clover", "cobertura", "coverage.py", "excoveralls", "gcov", "gocov", "jacoco", "lcov", "lcov-json", "simplecov", "xccov"}

// a map of the formatters to use
var formatterMap = map[string]formatters.Formatter{
	"clover":      &clover.Formatter{},
	"cobertura":   &cobertura.Formatter{},
	"coverage.py": &coveragepy.Formatter{},
	"excoveralls": &excoveralls.Formatter{},
	"gcov":        &gcov.Formatter{},
	"gocov":       &gocov.Formatter{},
	"jacoco":      &jacoco.Formatter{},
	"lcov":        &lcov.Formatter{},
	"lcov-json":   &lcovjson.Formatter{},
	"simplecov":   &simplecov.Formatter{},
	"xccov":       &xccov.Formatter{},
}

// formatCoverageCmd represents the format command
var formatCoverageCmd = &cobra.Command{
	Use:   "format-coverage [coverage file]",
	Short: "Locate, parse, and re-format supported coverage sources.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			logrus.Debugf("coverage path %s", args[0])
			formatOptions.CoveragePath = args[0]
			if formatOptions.InputType == "" {
				return errors.WithStack(errors.Errorf("please specify the format of the coverage file \"%s\" using the --input-type flag", formatOptions.CoveragePath))
			}
		}
		return runFormatter(formatOptions)
	},
}

func runFormatter(formatOptions CoverageFormatter) error {
	envy.Set("PREFIX", formatOptions.Prefix)
	envy.Set("ADD_PREFIX", formatOptions.AddPrefix)

	// if a type is specified use that
	if formatOptions.InputType != "" {
		if f, ok := formatterMap[formatOptions.InputType]; ok {
			logrus.Debugf("using formatter %s", formatOptions.InputType)
			_, err := f.Search(formatOptions.CoveragePath)
			if err != nil {
				logrus.Errorf("could not find coverage file %s\n%s", formatOptions.CoveragePath, err)
				return errors.WithStack(err)
			}
			formatOptions.In = f
		} else {
			return errors.WithStack(errors.Errorf("could not find a formatter of type %s", formatOptions.InputType))
		}
	} else {
		logrus.Debug("searching for a formatter to use")
		// else start searching for files:
		for n, f := range formatterMap {
			logrus.Debugf("checking %s formatter", n)
			if p, err := f.Search(); err == nil {
				logrus.Debugf("found file %s for %s formatter", p, n)
				formatOptions.In = f
				break
			}
		}
	}

	if formatOptions.In == nil {
		return errors.WithStack(errors.Errorf("could not find any viable formatter. available formatters: %s", strings.Join(formatterList, ", ")))
	}

	return formatOptions.Save()
}

func (f CoverageFormatter) Save() error {
	rep, err := f.In.Format()
	if err != nil {
		return errors.WithStack(err)
	}

	if len(rep.SourceFiles) == 0 {
		return errors.WithStack(errors.New("could not find coverage info for source files"))
	}

	if f.writer == nil {
		f.writer, err = writer(formatOptions.Output)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	err = rep.Save(f.writer)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func init() {
	pwd, _ := os.Getwd()
	formatCoverageCmd.Flags().StringVarP(&formatOptions.Prefix, "prefix", "p", pwd, "the root directory where the coverage analysis was performed")
	formatCoverageCmd.Flags().StringVar(&formatOptions.AddPrefix, "add-prefix", "", "add this prefix to file paths")
	formatCoverageCmd.Flags().StringVarP(&formatOptions.Output, "output", "o", ccDefaultCoveragePath, "output path")
	formatCoverageCmd.Flags().StringVarP(&formatOptions.InputType, "input-type", "t", "", fmt.Sprintf("type of input source to use [%s]", strings.Join(formatterList, ", ")))
	RootCmd.AddCommand(formatCoverageCmd)
}
