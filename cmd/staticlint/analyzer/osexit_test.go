package analyzer

import (
	"golang.org/x/tools/go/analysis/analysistest"
	"testing"
)

func TestAnalyzerOsExit(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), AnalyzerOsExit)
}
