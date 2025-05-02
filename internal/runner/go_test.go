package runner

import (
	"context"
	"fmt"
	. "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
	"testing"
	assert "github.com/stretchr/testify/assert"
)

func TestGoFindTest(t *testing.T) {
	filename := "main.go"
	lang := "go"

	code := `
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
			Filename: filename,
			Line:     8,
		},
	}

	codeBytes := []byte(code)

	run := GoRun{}
	query := run.Query()

	language := golang.GetLanguage()
	q, _ := NewQuery([]byte(query), language)

	testFinder := RunQueryFinder{Query: q, Lang: lang}
	node, _ := ParseCtx(context.Background(), codeBytes, language)

	tests := run.Find(&testFinder, node, filename, codeBytes)
	fmt.Println(tests)

	assert.NotNil(t, tests, "tests can't be nil in this case")
	assert.Equal(t, len(tests), len(expectedTest), "tests must be same size %d %d", len(tests), len(expectedTest))

	for line, expected := range expectedTest {
		actual, found := tests[line]
		assert.Equal(t, found, true, "expected test on line %d, but not found", line)
		assert.Equal(t, actual, expected, "expected test on line %d to be %v, but got %v", line, expected, actual)
	}
}
