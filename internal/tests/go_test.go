package tests

import (
	"context"
	"fmt"
	. "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
	"testing"
	assert "github.com/stretchr/testify/assert"
)

func TestGoFindTest(t *testing.T) {
	lang := "go"

	code := `
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

	codeBytes := []byte(code)

	test := GoTest{}
	query := test.TestQuery()

	language := golang.GetLanguage()
	q, _ := NewQuery([]byte(query), language)

	testFinder := TestFinder{TestQuery: q, Lang: lang}
	node, _ := ParseCtx(context.Background(), codeBytes, language)

	tests := test.Find(&testFinder, node, "example_test.go", codeBytes)
	fmt.Println(tests)

	assert.NotNil(t, tests, "tests can't be nil in this case")
	assert.Equal(t, len(tests), len(expectedTest), "tests must be same size %d %d", len(tests), len(expectedTest))

	for line, expected := range expectedTest {
		actual, found := tests[line]
		assert.Equal(t, found, true, "expected test on line %d, but not found", line)
		assert.Equal(t, actual, expected, "expected test on line %d to be %v, but got %v", line, expected, actual)
	}
}
