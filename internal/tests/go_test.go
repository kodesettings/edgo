package tests

import (
	"context"
	"testing"
	sitter "github.com/tree-sitter/go-tree-sitter"
	golang "github.com/tree-sitter/tree-sitter-go/bindings/go"
	assert "github.com/stretchr/testify/assert"
)

func TestGoFindTest(t *testing.T) {
	lang := "go"

	sourceCode := `
import (
	"testing"
)

func simple(t *testing.T) { }

func Test1(t *testing.T) {
	if 1 != 2 {
		t.Errorf("Expected ")
	}
}

func Test2(t *testing.T) {
	if 1 != 2 { t.Errorf("Expected ") }
}
`
	expectedTest := map[int]TestData{
		7: {
			Name:     "Test1",
			Filename: "example_test.go",
			Line:     7,
		},
		13: {
			Name:     "Test2",
			Filename: "example_test.go",
			Line:     13,
		},
	}

	goTest := GoTest{}
	query := goTest.TestQuery()

	language := sitter.NewLanguage(golang.Language())
	q, _ := sitter.NewQuery(language, query)

	testFinder := TestFinder{TestQuery: q, Lang: lang}

	parser := sitter.NewParser()
	defer parser.Close()
	parser.SetLanguage(language)

	tree := parser.ParseCtx(context.Background(), []byte(sourceCode), nil)
	tests := goTest.Find(&testFinder, tree.RootNode(), "example_test.go", []byte(sourceCode))

	assert.NotNil(t, tests, "tests can't be nil in this case")
	assert.Equal(t, len(tests), len(expectedTest), "tests must be same size %d %d", len(tests), len(expectedTest))

	for line, expected := range expectedTest {
		actual, found := tests[line]
		assert.Equal(t, found, true, "expected test on line %d, but not found", line)
		assert.Equal(t, actual, expected, "expected test on line %d to be %v, but got %v", line, expected, actual)
	}
}
