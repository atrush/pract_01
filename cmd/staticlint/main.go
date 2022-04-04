package main

import (
	custom "github.com/atrush/pract_01.git/cmd/staticlint/analyzer"
	"github.com/go-critic/go-critic/checkers/analyzer"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

func main() {
	mychecks := []*analysis.Analyzer{}

	// append checks from analysis
	mychecks = append(mychecks,
		asmdecl.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		bools.Analyzer,
		buildtag.Analyzer,
		cgocall.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		errorsas.Analyzer,
		framepointer.Analyzer,
		httpresponse.Analyzer,
		ifaceassert.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		printf.Analyzer,
		shift.Analyzer,
		sigchanyzer.Analyzer,
		stdmethods.Analyzer,
		stringintconv.Analyzer,
		structtag.Analyzer,
		tests.Analyzer,
		testinggoroutine.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
	)

	//  append 'SA' checks from analysis staticcheck.
	for _, ch := range staticcheck.Analyzers {
		mychecks = append(mychecks, ch)
	}

	mychecks = append(mychecks,
		//  append some 'S1' from analysis simple
		simple.Analyzers["S1009"], // Omit redundant nil check on slice.
		simple.Analyzers["S1011"], // Use a single append to concatenate two slices.
		simple.Analyzers["S1028"], // Simplify error construction with fmt.Errorf.
		simple.Analyzers["S1031"], // Omit redundant nil check around loop.

		//  append some 'ST1' from analysis stylecheck
		stylecheck.Analyzers["ST1019"], // Importing the same package multiple times.
		stylecheck.Analyzers["ST1013"], // Should use constants for HTTP error codes, not magic numbers.
		stylecheck.Analyzers["ST1012"], // Poorly chosen name for error variable.
		stylecheck.Analyzers["ST1008"], // A function’s error value should be its last return value.
	)

	//  append go-critic analysers
	mychecks = append(mychecks, analyzer.Analyzer)

	//  append custom check for os.Exit in main func
	mychecks = append(mychecks, custom.AnalyzerOsExit)
	multichecker.Main(mychecks...)
}
