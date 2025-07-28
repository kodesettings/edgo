package tests

import (
	"context"
	"testing"
	sitter "github.com/tree-sitter/go-tree-sitter"
	javascript "github.com/tree-sitter/tree-sitter-javascript/bindings/go"
	assert "github.com/stretchr/testify/assert"
)

func TestJavascriptFindTest(t *testing.T) {
	lang := "javascript"

	sourceCode := `
function sum(a, b) {
    return a + b;
}


describe("math tests", () => {

    it("positive", () => {
        expect(sum(1, 1)).toBe(2);
    });

    it("negative", () => {
        expect(sum(-1, -1)).toBe(-2);
    });
    
    it("failed", () => {
        expect(sum(-1, -1)).toBe(2);
    });
   
});
`
	expectedTest := map[int]TestData{
		6: {
			Name:     `"math tests"`,
			Filename: "test.js",
			Line:     6,
		},
		8: {
			Name:     `"positive"`,
			Filename: "test.js",
			Line:     8,
		},
		12: {
			Name:     `"negative"`,
			Filename: "test.js",
			Line:     12,
		},
		16: {
			Name:     `"failed"`,
			Filename: "test.js",
			Line:     16,
		},
	}

	javascriptTest := JavascriptTest{}
	query := javascriptTest.TestQuery()

	language := sitter.NewLanguage(javascript.Language())
	q, _ := sitter.NewQuery(language, query)

	testFinder := TestFinder{TestQuery: q, Lang: lang}

	parser := sitter.NewParser()
	defer parser.Close()
	parser.SetLanguage(language)

	tree := parser.ParseCtx(context.Background(), []byte(sourceCode), nil)
	tests := javascriptTest.Find(&testFinder, tree.RootNode(), "test.js", []byte(sourceCode))

	assert.NotNil(t, tests, "tests can't be nil in this case")
	assert.Equal(t, len(tests), len(expectedTest), "tests must be same size %d %d", len(tests), len(expectedTest))

	for line, expected := range expectedTest {
		actual, found := tests[line]
		assert.Equal(t, found, true, "expected test on line %d, but not found", line)
		assert.Equal(t, actual, expected, "expected test on line %d to be %v, but got %v", line, expected, actual)
	}
}
