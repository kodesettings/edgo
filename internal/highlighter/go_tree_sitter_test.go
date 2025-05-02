package highlighter

import (
	"context"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/javascript"
	"testing"
	assert "github.com/stretchr/testify/assert"
)

func TestTreeSitterQueries(t *testing.T) {
	code := []byte(`
function hello() { 
	// comment line 
	console.log('hello') 
	if (true) { console.log('true') }
	return "value"
}
`)

	query := `

[
  "async"
  "debugger"
  "delete"
  "extends"
  "from"
  "get"
  "new"
  "set"
  "target"
  "typeof"
  "instanceof"
  "void"
  "with"
] @keyword

[
  "of"
  "as"
  "in"
] @keyword.operator

[
  "function"
] @keyword.function

[
  "class"
  "let"
  "var"
] @keyword.storage.type

[
  "const"
  "static"
] @keyword.storage.modifier

[
  "default"
  "yield"
  "finally"
  "do"
  "await"
] @keyword.control

[
  "if"
  "else"
  "switch"
  "case"
  "while"
] @keyword.control.conditional

[
  "for"
] @keyword.control.repeat

[
  "import"
  "export"
] @keyword.control.import 

[
  "return"
  "break"
  "continue"
] @keyword.control.return

[
  "throw"
  "try"
  "catch"
] @keyword.control.exception
`

	// Parse source code
	lang := javascript.GetLanguage()
	n, _ := sitter.ParseCtx(context.Background(), code, lang)

	// Execute the query
	q, _ := sitter.NewQuery([]byte(query), lang)
	qc := sitter.NewQueryCursor()
	qc.Exec(q, n)

	for {
		m, ok := qc.NextMatch()
		if !ok { break }
		m = qc.FilterPredicates(m, code)
		for _, c := range m.Captures {
			name := q.CaptureNameForId(c.Index)
			content := c.Node.Content(code)
			assert.NotNil(t, c.Node.StartPoint(), "node startpoint is nil")
			assert.NotNil(t, c.Node.EndPoint(), "node endpoint is nil")
			assert.NotNil(t, name, "name is nil")
			assert.NotNil(t, c.Node.Type(), "node type is nil")
			assert.NotNil(t, content, "content is nil")
		}
	}
}
