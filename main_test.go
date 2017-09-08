package main

import (
	"go/parser"
	"go/token"
	"log"
	"os"
	"testing"
	//"runtime"

	"github.com/bouk/monkey"
	"github.com/fatih/set"
	"github.com/stretchr/testify/assert"
)

func TestGetDeclaredNamesWithSimpleFunctions(t *testing.T) {
	t.Parallel()

	in, err := parser.ParseFile(token.NewFileSet(), "example_packages/simple/main.go", nil, parser.AllErrors)
	if err != nil {
		t.Logf("failing because ParseFile returned error: %v", err)
		t.FailNow()
	}

	expectedDeclarations := []string{"A", "B", "C", "outer"}
	expected := set.New()
	for _, x := range expectedDeclarations {
		expected.Add(x)
	}

	actual := set.New()

	getDeclaredNames(in, actual)

	assert.Equal(t, expected, actual, "expected output did not match actual output")
}

func TestGetDeclaredNamesWithStructMethods(t *testing.T) {
	t.Parallel()

	in, err := parser.ParseFile(token.NewFileSet(), "example_packages/methods/main.go", nil, parser.AllErrors)
	if err != nil {
		t.Logf("failing because ParseFile returned error: %v", err)
		t.FailNow()
	}

	expectedDeclarations := []string{"Example.A", "Example.B", "Example.C", "outer"}
	expected := set.New()
	for _, x := range expectedDeclarations {
		expected.Add(x)
	}

	actual := set.New()
	getDeclaredNames(in, actual)

	assert.Equal(t, expected, actual, "expected output did not match actual output")
}

func TestGetCalledNamesWithStructMethods(t *testing.T) {
	t.Parallel()

	in, err := parser.ParseFile(token.NewFileSet(), "example_packages/methods/main_test.go", nil, parser.AllErrors)
	if err != nil {
		t.Logf("failing because ParseFile returned error: %v", err)
		t.FailNow()
	}

	expectedDeclarations := []string{".Parallel", "Example.A", "Example.C", "outer"}
	expected := set.New()
	for _, x := range expectedDeclarations {
		expected.Add(x)
	}

	actual := set.New()
	getCalledNames(in, actual)

	assert.Equal(t, expected, actual, "expected output did not match actual output")
}

func TestMainOperatesAsExpected(t *testing.T) {
	originalArgs := os.Args
	os.Args = []string{
		originalArgs[0],
		"--package=github.com/verygoodsoftwarenotvirus/tarp/example_packages/simple",
	}

	main()
	os.Args = originalArgs
}

func TestMainFailsWhenPackageIsNonexistent(t *testing.T) {
	originalArgs := os.Args
	os.Args = []string{
		originalArgs[0],
		"--package=github.com/nosuchrealusername/absolutelynosuchpackage",
		"--fail-on-extras",
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
	assert.True(t, fatalfCalled, "main should call log.Fatalf() when --fail-on-extras is passed in and extras are found")

	os.Args = originalArgs
	monkey.Unpatch(log.Fatalf)
}

func TestMainFailsWhenCodeAnalyzedIsInvalid(t *testing.T) {
	originalArgs := os.Args
	os.Args = []string{
		originalArgs[0],
		"--package=github.com/verygoodsoftwarenotvirus/tarp/example_packages/invalid",
		"--fail-on-extras",
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

	assert.True(t, fatalCalled, "main should call log.Fatal() when --fail-on-extras is passed in and extras are found")

	os.Args = originalArgs
	monkey.Unpatch(log.Fatal)
}

func TestMainFailsArgumentsAreInvalid(t *testing.T) {
	originalArgs := os.Args
	os.Args = []string{originalArgs[0]}

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
	assert.True(t, fatalCalled, "main should call log.Fatal when --fail-on-extras is passed in and extras are found")
	os.Args = originalArgs
	monkey.Unpatch(log.Fatal)
}

func TestSimplePackageFailsWhenArgsInstructItTo(t *testing.T) {
	originalArgs := os.Args
	os.Args = []string{
		originalArgs[0],
		"--package=github.com/verygoodsoftwarenotvirus/tarp/example_packages/simple",
		"--fail-on-extras",
	}

	var fatalCalled bool
	monkey.Patch(log.Fatal, func(...interface{}) {
		fatalCalled = true
	})

	main()
	assert.True(t, fatalCalled, "main should call log.Fatal() when --fail-on-extras is passed in and extras are found")
	os.Args = originalArgs
	monkey.Unpatch(log.Fatal)
}
