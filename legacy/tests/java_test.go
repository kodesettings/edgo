package tests

import (
	"context"
	"testing"
	sitter "github.com/tree-sitter/go-tree-sitter"
	java "github.com/tree-sitter/tree-sitter-java/bindings/go"
	assert "github.com/stretchr/testify/assert"
)

func TestJavaFindTest(t *testing.T) {
	lang := "java"

	sourceCode := `
import static org.junit.jupiter.api.Assertions.assertEquals;

import calc.Calculator;
import org.junit.jupiter.api.Test;


class MyFirstJUnitJupiterTests {

    private final Calculator calculator = new Calculator();

    @Test
    void addition() {
        System.out.println("test");
        assertEquals(2, calculator.add(1, 1));
    }
	
	void simpleFunction() {
        System.out.println("hi");
	}
}
`
	expectedTest := map[int]TestData{
		12: {
			Name:     `addition`,
			Filename: "test.java",
			Line:     12,
		},

	}

	javaTest := JavaTest{}
	query := javaTest.TestQuery()

	language := sitter.NewLanguage(java.Language())
	q, _ := sitter.NewQuery(language, query)

	testFinder := TestFinder{TestQuery: q, Lang: lang}

	parser := sitter.NewParser()
	defer parser.Close()
	parser.SetLanguage(language)

	tree := parser.ParseCtx(context.Background(), []byte(sourceCode), nil)
	tests := javaTest.Find(&testFinder, tree.RootNode(), "test.java", []byte(sourceCode))

	assert.NotNil(t, tests, "tests can't be nil in this case")
	assert.Equal(t, len(expectedTest), len(tests), "tests must be same size")

	for line, expected := range expectedTest {
		actual, found := tests[line]
		assert.Equal(t, true, found, "expected test on line %d, but not found", line)
		assert.Equal(t, expected, actual, "expected test on line %d to be %v, but got %v", line, expected, actual)
	}
}
