package runner

import (
	"context"
	"testing"
	sitter "github.com/tree-sitter/go-tree-sitter"
	golang "github.com/tree-sitter/tree-sitter-go/bindings/go"
	assert "github.com/stretchr/testify/assert"
)

func TestGoFindTest(t *testing.T) {
	fileName := "main.go"
	lang := "go"

	sourceCode := `
package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	var count = 0
	for i := 0; i <= 10000000; i++ {
		count += i
	}
	
	fmt.Println(count, "elapsed", time.Since(start))
}

func main2() int { return 1 }
`

	expectedTest := map[int]RunData{
		8: {
			Name:     "main",
			Filename: fileName,
			Line:     8,
		},
	}

	run := GoRun{}
	query := run.Query()

	language := sitter.NewLanguage(golang.Language())
	q, _ := sitter.NewQuery(language, query)

	testFinder := RunQueryFinder{Query: q, Lang: lang}

	parser := sitter.NewParser()
	defer parser.Close()
	parser.SetLanguage(language)

	tree := parser.ParseCtx(context.Background(), []byte(sourceCode), nil)

	tests := run.Find(&testFinder, tree.RootNode(), fileName, []byte(sourceCode))

	assert.NotNil(t, tests, "tests can't be nil in this case")
	assert.Equal(t, len(tests), len(expectedTest), "tests must be same size %d %d", len(tests), len(expectedTest))

	for line, expected := range expectedTest {
		actual, found := tests[line]
		assert.Equal(t, found, true, "expected test on line %d, but not found", line)
		assert.Equal(t, actual, expected, "expected test on line %d to be %v, but got %v", line, expected, actual)
	}
}
