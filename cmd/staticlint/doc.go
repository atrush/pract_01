/**
Static tests for project. Includes:
- golang.org/x/tools/go/analysis tests
- staticchk.io tesst
	- Staticchecks all "SA****"
	- Codeimplifications
		- "S1009" - Omit redundant nil check on slice.
		- "S1028" - Simplify error construction with fmt.Errorf.
		- "S1011" - Use a single append to concatenate two slices.
		- "S1031" - Omit redundant nil check around loop.
	- Stylechecks
		- "ST1019" - Importing the same package multiple times.
		- "ST1013" - Should use constants for HTTP error codes, not magic numbers.
		- "ST1012" - Poorly chosen name for error variable.
		- "ST1008" - A function’s error value should be its last return value.
- github.com/go-critic/go-critic tests
- custom check for os.Exit() in main function

To run tests use
go run cmd/staticlint/main.go ./...

*/
package main
