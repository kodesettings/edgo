package tests

import (
	. "github.com/tree-sitter/go-tree-sitter"
	"strings"
)

type Test interface {
	TestQuery() string
	Find(tfinder *TestFinder, root *Node, filename string, code []byte) map[int]TestData
	Run(test TestData) []string
}

type TestData struct {
	Name string
	Filename string
	Line int
}

type TestFinder struct {
	TestQuery *Query
	Lang      string
}


func (this *TestFinder) Find(root *Node, filename string, code []byte) map[int]TestData {
	tests := make(map[int]TestData)

	queryCursor := NewQueryCursor()
	matches := queryCursor.Matches(this.TestQuery, root, code)

	for match := matches.Next(); match != nil; match = matches.Next() {
		for _, capture := range match.Captures {
			node := capture.Node;
			nodename := this.TestQuery.CaptureNames()[capture.Index]
			content := node.Utf8Text(code)
			isTestFound := nodename == "test-name"
			if isTestFound {
				line := int(node.StartPosition().Row)
				tests[line] = TestData{
					Name: content,
					Filename: filename,
					Line: line,
				}
			}
		}
	}

	return tests
}

func GetTestByLang(lang string, filepath string) Test {
	switch lang {
	case "go":
		if !strings.HasSuffix(filepath, "test.go") { return nil }
		return &GoTest{}

	case "python": return &PythonTest{}
	case "javascript": return &JavascriptTest{}
	case "java":return &JavaTest{}
	default:
	}

	return nil
}
