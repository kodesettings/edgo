package highlighter

import (
	"context"
	sitter "github.com/tree-sitter/go-tree-sitter"
	javascript "github.com/tree-sitter/tree-sitter-javascript/bindings/go"
	. "github.com/vipmax/edgo/internal/highlighter/langs"
	"testing"
	assert "github.com/stretchr/testify/assert"
)

func TestTreeSitterQueries(t *testing.T) {
	sourceCode := `
function hello() { 
	// comment line 
	console.log('hello') 
	if (true) { console.log('true') }
	return "value"
}
`

	// Setup new parser
	language := sitter.NewLanguage(javascript.Language())
	parser := sitter.NewParser()
	defer parser.Close()
	parser.SetLanguage(language)

	// Parse source code
	tree := parser.ParseCtx(context.Background(), []byte(sourceCode), nil)
	defer tree.Close()
	query, _ := sitter.NewQuery(language, MatchQueryLang("javascript"))

	// Execute the query
	queryCursor := sitter.NewQueryCursor()
	matches := queryCursor.Matches(query, tree.RootNode(), []byte(sourceCode))

	for match := matches.Next(); match != nil; match = matches.Next() {
		for _, capture := range match.Captures {
			node := capture.Node;
			nodename := query.CaptureNames()[capture.Index]
			content := node.Utf8Text([]byte(sourceCode))
			assert.NotNil(t, capture.Node.StartPosition(), "node startpoint is nil")
			assert.NotNil(t, capture.Node.EndPosition(), "node endpoint is nil")
			assert.NotNil(t, nodename, "name is nil")
			assert.NotNil(t, node.Kind(), "node type is nil")
			assert.NotNil(t, content, "content is nil")
		}
	}
}
