package pythonstart_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnitPythonStart(t *testing.T) {
	suite := spec.New("python-start", spec.Report(report.Terminal{}), spec.Parallel())
	suite("Build", testBuild)
	suite("Detect", testDetect)
	suite.Run(t)
}
