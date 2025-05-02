package highlighter

import (
	"context"
	"edgo/internal/utils"
	"fmt"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/smacker/go-tree-sitter/javascript"
	"github.com/smacker/go-tree-sitter/python"
	"testing"
	"time"
	assert "github.com/stretchr/testify/assert"
)

func validation_helper(t *testing.T, node *sitter.Node, code []byte) {
	// Process the current node here.
	assert.NotNil(t, node.StartPoint(), "node startpoint is nil")
	assert.NotNil(t, node.EndPoint(), "node endpoint is nil")
	assert.NotNil(t, node.Type(), "node type is nil")
	assert.NotNil(t, node.Content(code), "content is nil")

	// Visit each child node recursively.
	childCount := int(node.NamedChildCount())
	for i := 0; i < childCount; i++ {
		child := node.NamedChild(i)
		validation_helper(t, child, code) // recursive validation
	}
}

func TestTreeSitter(t *testing.T) {
	//code, _ := utils.ReadFileToString("tree-sitter_test.go")
	code, _ := utils.ReadFileToString("../editor/editor.go")

	sourceCode := []byte(code)
	//lang := javascript.GetLanguage()
	lang := golang.GetLanguage()
	start := time.Now()
	n, err := sitter.ParseCtx(context.Background(), sourceCode, lang)

	assert.Nil(t, err, "error occurred in tree sitter test")
	assert.NotNil(t, n, "error occurred in tree sitter test")

	fmt.Println("parsed, elapsed", time.Since(start))
}

func TestTreeSitterGo(t *testing.T) {
	//code, _ := utils.ReadFileToString("tree-sitter_test.go")
	code := `
package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	var count = 0
	for i := 0; i <= 100000000; i++ {
		count += i
		fmt.Println(count)
		time.Sleep(time.Millisecond * 10)
	}
	fmt.Println(count, "elapsed", time.Since(start))
}
`

	sourceCode := []byte(code)
	//lang := javascript.GetLanguage()
	lang := golang.GetLanguage()
	start := time.Now()
	n, err := sitter.ParseCtx(context.Background(), sourceCode, lang)

	assert.Nil(t, err, "error occurred in tree sitter go")
	assert.NotNil(t, n, "error occurred in tree sitter go")

	fmt.Println("parsed, elapsed", time.Since(start))
}

func TestTreeSitterPython(t *testing.T) {
	//code, _ := utils.ReadFileToString("tree-sitter_test.go")
	code := `
import time

print("starting")
start_time = time.time()

for i in range(100000):
    print(i)
    time.sleep(0.01)

print("done")

elapsed_time = time.time() - start_time
print("Elapsed time:", elapsed_time, q"seconds")

`

	sourceCode := []byte(code)
	lang := python.GetLanguage()
	start := time.Now()
	n, err := sitter.ParseCtx(context.Background(), sourceCode, lang)

	assert.Nil(t, err, "error occurred in tree sitter python")
	assert.NotNil(t, n, "error occurred in tree sitter python")

	fmt.Println("parsed, elapsed", time.Since(start))
}

func TestTreeSitterJs(t *testing.T) {
	code := []byte(`
function hello() { 
	console.log('hello') 
}
`)
	parser := sitter.NewParser()
	parser.SetLanguage(javascript.GetLanguage())

	start := time.Now()
	tree, err := parser.ParseCtx(context.Background(),nil, code)

	assert.Nil(t, err, "error occurred in tree sitter js")
	assert.NotNil(t, tree, "error occurred in tree sitter js")

	fmt.Println("parsed, elapsed", time.Since(start))
	validation_helper(t, tree.RootNode(), code)
}

func TestTreeSitterJsEdit(t *testing.T) {
	parser := sitter.NewParser()
	parser.SetLanguage(javascript.GetLanguage())

	code := []byte(`function hello() { console.log('hello') }`)
	oldEndIndex := uint32(len(code))

	start := time.Now()

	tree, err := parser.ParseCtx(context.Background(),nil, code)
	assert.Nil(t, err, "error occurred in tree sitter js edit")
	assert.NotNil(t, tree, "error occurred in tree sitter js edit")

	elapsedFirst := time.Since(start)
	fmt.Println("parsed, elapsed", elapsedFirst)
	validation_helper(t, tree.RootNode(), code)

	fmt.Println("Edit input")

	code = []byte(`function hello2() { console.log('hello') }`)
	newEndIndex := uint32(len(code))

	tree.Edit(sitter.EditInput{
		StartIndex:  14,
		OldEndIndex: oldEndIndex,
		NewEndIndex: newEndIndex,
		StartPoint: sitter.Point{Row: 0, Column: 14 },
		OldEndPoint: sitter.Point{Row: 0, Column: 14},
		NewEndPoint: sitter.Point{Row: 0, Column: 15},
	})

	start = time.Now()
	tree, err = parser.ParseCtx(context.Background(), tree, code)
	assert.Nil(t, err, "error occurred in tree sitter js edit")
	assert.NotNil(t, tree, "error occurred in tree sitter js edit")

	elapsedSecond := time.Since(start)
	fmt.Println("parsed again, elapsed", elapsedSecond)
	speedup := float64(elapsedFirst) / float64(elapsedSecond)
	fmt.Printf("Speedup factor: %.2f\n", speedup)

	validation_helper(t, tree.RootNode(), code)
}

func TestTreeSitterJsEditDelete(t *testing.T) {
	parser := sitter.NewParser()
	parser.SetLanguage(javascript.GetLanguage())

	code := []byte(`function hello() { console.log('hello') }`)
	oldEndIndex := uint32(len(code))

	start := time.Now()

	tree, err := parser.ParseCtx(context.Background(),nil, code)
	assert.Nil(t, err, "error occurred in tree sitter js edit delete")
	assert.NotNil(t, tree, "error occurred in tree sitter js edit delete")

	elapsedFirst := time.Since(start)
	fmt.Println("parsed, elapsed", elapsedFirst)
	validation_helper(t, tree.RootNode(), code)

	fmt.Println("Edit input")

	code = []byte(`function hel() { console.log('hello') }`)
	newEndIndex := uint32(len(code))

	tree.Edit(sitter.EditInput{
		StartIndex:  14,
		OldEndIndex: oldEndIndex,
		NewEndIndex: newEndIndex,
		StartPoint: sitter.Point{Row: 0, Column: 14 },
		OldEndPoint: sitter.Point{Row: 0, Column: oldEndIndex},
		NewEndPoint: sitter.Point{Row: 0, Column: newEndIndex},
	})

	start = time.Now()
	tree, err = parser.ParseCtx(context.Background(), tree, code)
	assert.Nil(t, err, "error occurred in tree sitter js edit delete")
	assert.NotNil(t, tree, "error occurred in tree sitter js edit delete")

	elapsedSecond := time.Since(start)
	fmt.Println("parsed again, elapsed", elapsedSecond)
	speedup := float64(elapsedFirst) / float64(elapsedSecond)
	fmt.Printf("Speedup factor: %.2f\n", speedup)

	validation_helper(t, tree.RootNode(), code)
}

func TestTreeSitterJsEditDeleteMultiple(t *testing.T) {
	parser := sitter.NewParser()
	parser.SetLanguage(javascript.GetLanguage())

	code := []byte(`function hello() { console.log('hello') }`)
	oldEndIndex := uint32(len(code))

	start := time.Now()

	tree, err := parser.ParseCtx(context.Background(),nil, code)
	assert.Nil(t, err, "error occurred in tree sitter js edit delete multiple")
	assert.NotNil(t, tree, "error occurred in tree sitter js edit delete multiple")

	elapsedFirst := time.Since(start)
	fmt.Println("parsed, elapsed", elapsedFirst)
	validation_helper(t, tree.RootNode(), code)

	fmt.Println("Edit input")

	code = []byte(`function hel() { console.log('hello') }`)
	newEndIndex := uint32(len(code))

	tree.Edit(sitter.EditInput{
		StartIndex:  14, OldEndIndex: oldEndIndex, NewEndIndex: newEndIndex,
		StartPoint: sitter.Point{Row: 0, Column: 14 },
		OldEndPoint: sitter.Point{Row: 0, Column: oldEndIndex},
		NewEndPoint: sitter.Point{Row: 0, Column: newEndIndex},
	})

	start = time.Now()
	tree, err = parser.ParseCtx(context.Background(), tree, code)
	if err != nil { fmt.Println(err)}
	elapsedSecond := time.Since(start)
	fmt.Println("parsed again, elapsed", elapsedSecond)
	speedup := float64(elapsedFirst) / float64(elapsedSecond)
	fmt.Printf("Speedup factor: %.2f\n", speedup)
	validation_helper(t, tree.RootNode(), code)

	fmt.Println("Edit input")

	code = []byte(`function h() { console.log('hello') }`)
	newEndIndex = uint32(len(code))

	tree.Edit(sitter.EditInput{
		StartIndex:  10, OldEndIndex: oldEndIndex, NewEndIndex: newEndIndex,
		StartPoint: sitter.Point{Row: 0, Column: 10 },
		OldEndPoint: sitter.Point{Row: 0, Column: oldEndIndex},
		NewEndPoint: sitter.Point{Row: 0, Column: newEndIndex},
	})

	start = time.Now()
	tree, err = parser.ParseCtx(context.Background(), tree, code)
	assert.Nil(t, err, "error occurred in tree sitter js edit delete multiple")
	assert.NotNil(t, tree, "error occurred in tree sitter js edit delete multiple")

	elapsedSecond = time.Since(start)
	fmt.Println("parsed again, elapsed", elapsedSecond)
	speedup = float64(elapsedFirst) / float64(elapsedSecond)
	fmt.Printf("Speedup factor: %.2f\n", speedup)
	validation_helper(t, tree.RootNode(), code)
}

func TestTreeSitterJsEditEnter(t *testing.T) {
	parser := sitter.NewParser()
	parser.SetLanguage(javascript.GetLanguage())

	code := []byte(`function hello() { console.log('hello') }`)
	oldEndIndex := uint32(len(code))

	start := time.Now()

	tree, err := parser.ParseCtx(context.Background(),nil, code)
	assert.Nil(t, err, "error occurred in tree sitter js edit delete enter")
	assert.NotNil(t, tree, "error occurred in tree sitter js edit delete enter")

	elapsedFirst := time.Since(start)
	fmt.Println("parsed, elapsed", elapsedFirst)
	validation_helper(t, tree.RootNode(), code)

	fmt.Println("Edit input")

	code = []byte(`function hello() { 
console.log('hello')}`)
	newEndIndex := uint32(len(code))

	tree.Edit(sitter.EditInput{
		StartIndex:  19,
		OldEndIndex: oldEndIndex,
		NewEndIndex: newEndIndex,
		StartPoint: sitter.Point{Row: 1, Column: 0 },
		OldEndPoint: sitter.Point{Row: 0, Column: oldEndIndex},
		NewEndPoint: sitter.Point{Row: 1, Column: newEndIndex},
	})

	start = time.Now()
	tree, err = parser.ParseCtx(context.Background(), tree, code)
	assert.Nil(t, err, "error occurred in tree sitter js edit delete enter")
	assert.NotNil(t, tree, "error occurred in tree sitter js edit delete enter")

	elapsedSecond := time.Since(start)
	fmt.Println("parsed again, elapsed", elapsedSecond)
	speedup := float64(elapsedFirst) / float64(elapsedSecond)
	fmt.Printf("Speedup factor: %.2f\n", speedup)

	validation_helper(t, tree.RootNode(), code)
}

func TestTreeSitterQuery(t *testing.T) {
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
  "function"
  "if"
  "return"
] @keyword

(comment) @comment
`
	lang := javascript.GetLanguage()
	start := time.Now()
	q, _ := sitter.NewQuery([]byte(query), lang)
	qc := sitter.NewQueryCursor()

	n, _ := sitter.ParseCtx(context.Background(), code, lang)
	fmt.Println("parsed , elapsed", time.Since(start))

	// Execute the query

	qc.Exec(q, n)
	fmt.Println("query exec, elapsed", time.Since(start))

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

func TestTreeSitterQuery2(t *testing.T) {
	sourceCode := []byte(`
package main
import "fmt"
func main() {
	// comment line 
	fmt.Println("Hello, world!")
}
`)

	// Query with predicates
	query := `
(identifier) @keyword
`
	// Parse source code
	lang := golang.GetLanguage()
	n, _ := sitter.ParseCtx(context.Background(), sourceCode, lang)

	// Execute the query
	q, _ := sitter.NewQuery([]byte(query), lang)
	qc := sitter.NewQueryCursor()
	qc.Exec(q, n)

	// Iterate over query results
	for {
		m, ok := qc.NextMatch()
		if !ok { break }
		// Apply predicates filtering
		m = qc.FilterPredicates(m, sourceCode)
		for _, c := range m.Captures {
			assert.NotNil(t, c.Node.StartPoint(), "node startpoint is nil")
			assert.NotNil(t, c.Node.EndPoint(), "node endpoint is nil")
			assert.NotNil(t, c.Node.Type(), "node type is nil")
			assert.NotNil(t, c.Node.Content(sourceCode), "content is nil")
		}
	}
}
