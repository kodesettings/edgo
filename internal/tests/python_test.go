package tests

import (
	"context"
	"testing"
	sitter "github.com/tree-sitter/go-tree-sitter"
	python "github.com/tree-sitter/tree-sitter-python/bindings/go"
	assert "github.com/stretchr/testify/assert"
)

func TestPythonFindTest(t *testing.T) {
	lang := "python"

	sourceCode := `import pytest
from random import randint

class TestYo:
    def mess(self, value):
        return randint(0, value)

    def test_pass(self):
        assert 1 == 1

    def test_fail_sometimes(self):
        assert 1 == self.mess(1)

`
	expectedTest := map[int]TestData{
		3: {
			Name:     "TestYo",
			Filename: "test_yo.py",
			Line:     3,
		},
		7: {
			Name:     "test_pass",
			Filename: "test_yo.py",
			Line:     7,
		},
		10: {
			Name:     "test_fail_sometimes",
			Filename: "test_yo.py",
			Line:     10,
		},
	}

	pythonTest := PythonTest{}
	query := pythonTest.TestQuery()

	language := sitter.NewLanguage(python.Language())
	q, _ := sitter.NewQuery(language, query)

	testFinder := TestFinder{TestQuery: q, Lang: lang}

	parser := sitter.NewParser()
	defer parser.Close()
	parser.SetLanguage(language)

	tree := parser.ParseCtx(context.Background(), []byte(sourceCode), nil)
	tests := pythonTest.Find(&testFinder, tree.RootNode(), "test_yo.py", []byte(sourceCode))

	assert.NotNil(t, tests, "tests can't be nil in this case")
	assert.Equal(t, len(tests), len(expectedTest), "tests must be same size %d %d", len(tests), len(expectedTest))

	for line, expected := range expectedTest {
		actual, found := tests[line]
		assert.Equal(t, found, true, "expected test on line %d, but not found", line)
		assert.Equal(t, actual, expected, "expected test on line %d to be %v, but got %v", line, expected, actual)
	}
}
