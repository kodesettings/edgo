package runner

import (
	. "github.com/tree-sitter/go-tree-sitter"
)

type Run interface {
	Query() string
	Find(tfinder *RunQueryFinder, root *Node, filename string, code []byte) map[int]RunData
	Run(data RunData) []string
}

func GetRunnerByLang(lang string, filepath string) Run {
	switch lang {
	case "go": return &GoRun{}
	default:
	}

	return nil
}

type RunData struct {
	Name string
	Filename string
	Line int
}

type RunQueryFinder struct {
	Query *Query
	Lang  string
}

func (this *RunQueryFinder) Find(root *Node, filename string, code []byte) map[int]RunData {
	results := make(map[int]RunData)

	queryCursor := NewQueryCursor()
	matches := queryCursor.Matches(this.Query, root, code)

	for match := matches.Next(); match != nil; match = matches.Next() {
		for _, capture := range match.Captures {
			node := capture.Node;
			nodename := this.Query.CaptureNames()[capture.Index]
			content := node.Utf8Text(code)
			isTestFound := nodename == "test-name"
			if isTestFound {
				line := int(node.StartPosition().Row)
				results[line] = RunData{
					Name: content,
					Filename: filename,
					Line: line,
				}
			}
		}
	}

	return results
}
