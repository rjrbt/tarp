// +build !unit

package main

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/bouk/monkey"
	"github.com/stretchr/testify/assert"
)

func TestMainFunc(t *testing.T) {
	originalArgs := os.Args

	optimal := func(t *testing.T) {
		os.Args = []string{
			originalArgs[0],
			"analyze",
			fmt.Sprintf("--package=%s", buildExamplePackagePath(t, "simple", false)),
		}

		main()
		os.Args = originalArgs
	}

	nonexistentPackage := func(t *testing.T) {
		os.Args = []string{
			originalArgs[0],
			"analyze",
			fmt.Sprintf("--package=%s", buildExamplePackagePath(t, "absolutelynosuchpackage", false)),
			"--fail-on-found",
		}

		defer func() {
			if r := recover(); r != nil {
				// recovered from our monkey patched log.Fatalf
				assert.True(t, true)
			}
		}()

		var fatalfCalled bool
		monkey.Patch(log.Fatalf, func(string, ...interface{}) {
			fatalfCalled = true
			panic("log.Fatalf")
		})

		main()
		assert.True(t, fatalfCalled, "main should call log.Fatalf() when --fail-on-found is passed in and extras are found")

		os.Args = originalArgs
		monkey.Unpatch(log.Fatalf)
	}

	emptyPackage := func(t *testing.T) {
		os.Args = []string{
			originalArgs[0],
			"analyze",
			fmt.Sprintf("--package=%s", buildExamplePackagePath(t, "no_go_files", false)),
			"--fail-on-found",
		}

		defer func() {
			if r := recover(); r != nil {
				// recovered from our monkey patched log.Fatalf
				assert.True(t, true)
			}
		}()

		var fatalfCalled bool
		monkey.Patch(log.Fatalf, func(string, ...interface{}) {
			fatalfCalled = true
			panic("log.Fatalf")
		})

		main()
		assert.True(t, fatalfCalled, "main should call log.Fatalf() when --fail-on-found is passed in and extras are found")

		os.Args = originalArgs
		monkey.Unpatch(log.Fatalf)
	}

	invalidCodeTest := func(t *testing.T) {
		os.Args = []string{
			originalArgs[0],
			"analyze",
			fmt.Sprintf("--package=%s", buildExamplePackagePath(t, "invalid", false)),
			"--fail-on-found",
		}

		invalidCodePath := buildExamplePackagePath(t, "invalid", true)
		err := os.MkdirAll(invalidCodePath, os.ModePerm)
		if err != nil {
			t.Log("error encountered creating temp path for invalid code test")
			t.FailNow()
		}

		f, err := os.Create(fmt.Sprintf("%s/main.go", invalidCodePath))
		if err != nil {
			t.Log("error encountered creating temp file for invalid code test")
			t.FailNow()
		}
		invalidCode := `
		package invalid

		import (
			"log"

		funk main() {
			return x
		)`
		fmt.Fprint(f, invalidCode)

		defer func() {
			if r := recover(); r != nil {
				// recovered from our monkey patched log.Fatal
				err = os.RemoveAll(invalidCodePath)
				if err != nil {
					t.Logf("error encountered deleting temp directory: %v", err)
					t.FailNow()
				}
			}
		}()

		var fatalCalled bool
		monkey.Patch(log.Fatal, func(...interface{}) {
			fatalCalled = true
			panic("log.Fatal")
		})

		main()

		assert.True(t, fatalCalled, "main should call log.Fatal() when --fail-on-found is passed in and extras are found")

		os.Args = originalArgs
		monkey.Unpatch(log.Fatal)
	}

	invalidArguments := func(t *testing.T) {
		os.Args = []string{
			originalArgs[0],
			"analyze",
		}
		defer func() {
			if r := recover(); r != nil {
				// recovered from our monkey patched log.Fatal
				assert.True(t, true)
			}
		}()

		var fatalCalled bool
		monkey.Patch(log.Fatal, func(...interface{}) {
			fatalCalled = true
			panic("log.Fatal")
		})

		main()
		assert.True(t, fatalCalled, "main should call log.Fatal when invalid arguments are passed to analyze")
		os.Args = originalArgs
		monkey.Unpatch(log.Fatal)
	}

	failsWhenInstructed := func(t *testing.T) {
		os.Args = []string{
			originalArgs[0],
			"analyze",
			fmt.Sprintf("--package=%s", buildExamplePackagePath(t, "simple", false)),
			"--fail-on-found",
		}

		var fatalCalled bool
		monkey.Patch(log.Fatal, func(...interface{}) {
			fatalCalled = true
		})

		main()
		assert.True(t, fatalCalled, "main should call log.Fatal() when --fail-on-found is passed in and extras are found")
		os.Args = originalArgs
		monkey.Unpatch(log.Fatal)
	}

	subtests := []subtest{
		{
			Message: "optimal",
			Test:    optimal,
		},
		{
			Message: "nonexistent package",
			Test:    nonexistentPackage,
		},
		{
			Message: "empty package",
			Test:    emptyPackage,
		},
		{
			Message: "invalid code",
			Test:    invalidCodeTest,
		},
		{
			Message: "invalid args",
			Test:    invalidArguments,
		},
		{
			Message: "fails with --fail-on-found",
			Test:    failsWhenInstructed,
		},
	}
	runSubtestSuite(t, subtests)
}
