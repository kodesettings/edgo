package highlighter

import (
	"context"
	"github.com/vipmax/edgo/internal/utils"
	"fmt"
	sitter "github.com/tree-sitter/go-tree-sitter"
	golang "github.com/tree-sitter/tree-sitter-go/bindings/go"
	javascript "github.com/tree-sitter/tree-sitter-javascript/bindings/go"
	python "github.com/tree-sitter/tree-sitter-python/bindings/go"
	"testing"
	"time"
	assert "github.com/stretchr/testify/assert"
)

func validation_helper(t *testing.T, node *sitter.Node, code []byte) {
	// Process the current node here.
	assert.NotNil(t, node.StartPosition(), "node startpoint is nil")
	assert.NotNil(t, node.EndPosition(), "node endpoint is nil")
	assert.NotNil(t, node.Kind(), "node type is nil")
	assert.NotNil(t, node.Utf8Text(code), "content is nil")

	// Visit each child node recursively.
	childCount := int(node.NamedChildCount())
	for i := 0; i < childCount; i++ {
		child := node.NamedChild(uint(i))
		validation_helper(t, child, code) // recursive validation
	}
}

func TestTreeSitter(t *testing.T) {
	sourceCode, _ := utils.ReadFileToString("tree-sitter_test.go")
	start := time.Now()

	// Setup new parser
	language := sitter.NewLanguage(golang.Language())
	parser := sitter.NewParser()
	defer parser.Close()
	parser.SetLanguage(language)

	// Parse source code
	tree := parser.ParseCtx(context.Background(), []byte(sourceCode), nil)
	defer tree.Close()

	assert.NotNil(t, tree, "error occurred in tree sitter test")
	fmt.Println("parsed, elapsed", time.Since(start))
}

func TestTreeSitterGo(t *testing.T) {
	sourceCode := `
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
	start := time.Now()

	// Setup new parser
	language := sitter.NewLanguage(golang.Language())
	parser := sitter.NewParser()
	defer parser.Close()
	parser.SetLanguage(language)

	// Parse source code
	tree := parser.ParseCtx(context.Background(), []byte(sourceCode), nil)
	defer tree.Close()

	assert.NotNil(t, tree, "error occurred in tree sitter go")
	fmt.Println("parsed, elapsed", time.Since(start))
}

func TestTreeSitterPython(t *testing.T) {
	sourceCode := `
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
	start := time.Now()

	// Setup new parser
	language := sitter.NewLanguage(python.Language())
	parser := sitter.NewParser()
	defer parser.Close()
	parser.SetLanguage(language)

	// Parse source code
	tree := parser.ParseCtx(context.Background(), []byte(sourceCode), nil)
	defer tree.Close()

	assert.NotNil(t, tree, "error occurred in tree sitter python")
	fmt.Println("parsed, elapsed", time.Since(start))
}

func TestTreeSitterJs(t *testing.T) {
	sourceCode := []byte(`
function hello() { 
	console.log('hello') 
}
`)
	start := time.Now()

	// Setup new parser
	language := sitter.NewLanguage(javascript.Language())
	parser := sitter.NewParser()
	defer parser.Close()
	parser.SetLanguage(language)

	// Parse source code
	tree := parser.ParseCtx(context.Background(), []byte(sourceCode), nil)
	defer tree.Close()

	assert.NotNil(t, tree, "error occurred in tree sitter javascript")
	fmt.Println("parsed, elapsed", time.Since(start))
	validation_helper(t, tree.RootNode(), sourceCode)
}

func TestTreeSitterJsEdit(t *testing.T) {
	// Setup new parser
	language := sitter.NewLanguage(javascript.Language())
	parser := sitter.NewParser()
	defer parser.Close()
	parser.SetLanguage(language)

	start := time.Now()

	sourceCode := []byte(`function hello() { console.log('hello') }`)
	oldEndIndex := uint(len(sourceCode))

	// Parse source code
	tree := parser.ParseCtx(context.Background(), []byte(sourceCode), nil)
	defer tree.Close()
	assert.NotNil(t, tree, "error occurred in tree sitter js edit")

    elapsedFirst := time.Since(start)
	fmt.Println("parsed, elapsed", time.Since(start))
	validation_helper(t, tree.RootNode(), sourceCode)

	fmt.Println("Edit input")

	sourceCode = []byte(`function hello2() { console.log('hello') }`)
	newEndIndex := uint(len(sourceCode))

	tree.Edit(&sitter.InputEdit{
		StartByte:  14,
		OldEndByte: oldEndIndex,
		NewEndByte: newEndIndex,
		StartPosition: sitter.Point{Row: 0, Column: 14 },
		OldEndPosition: sitter.Point{Row: 0, Column: 14},
		NewEndPosition: sitter.Point{Row: 0, Column: 15},
	})

	start = time.Now()
	tree = parser.ParseCtx(context.Background(), []byte(sourceCode), tree)
	assert.NotNil(t, tree, "error occurred in tree sitter js edit")

	elapsedSecond := time.Since(start)
	fmt.Println("parsed again, elapsed", elapsedSecond)
	speedup := float64(elapsedFirst) / float64(elapsedSecond)
	fmt.Printf("Speedup factor: %.2f\n", speedup)

	validation_helper(t, tree.RootNode(), sourceCode)
}

func TestTreeSitterJsEditDelete(t *testing.T) {
	// Setup new parser
	language := sitter.NewLanguage(javascript.Language())
	parser := sitter.NewParser()
	defer parser.Close()
	parser.SetLanguage(language)

	start := time.Now()

	sourceCode := []byte(`function hello() { console.log('hello') }`)
	oldEndIndex := uint(len(sourceCode))

	// Parse source code
	tree := parser.ParseCtx(context.Background(), []byte(sourceCode), nil)
	defer tree.Close()
	assert.NotNil(t, tree, "error occurred in tree sitter in js edit delete")

    elapsedFirst := time.Since(start)
	fmt.Println("parsed, elapsed", time.Since(start))
	validation_helper(t, tree.RootNode(), sourceCode)

	fmt.Println("Edit input")

	sourceCode = []byte(`function hel() { console.log('hello') }`)
	newEndIndex := uint(len(sourceCode))

	tree.Edit(&sitter.InputEdit{
		StartByte:  14,
		OldEndByte: oldEndIndex,
		NewEndByte: newEndIndex,
		StartPosition: sitter.Point{Row: 0, Column: 14 },
		OldEndPosition: sitter.Point{Row: 0, Column: oldEndIndex},
		NewEndPosition: sitter.Point{Row: 0, Column: newEndIndex},
	})

	start = time.Now()
	tree = parser.ParseCtx(context.Background(), []byte(sourceCode), tree)
	assert.NotNil(t, tree, "error occurred in tree sitter js edit")

	elapsedSecond := time.Since(start)
	fmt.Println("parsed again, elapsed", elapsedSecond)
	speedup := float64(elapsedFirst) / float64(elapsedSecond)
	fmt.Printf("Speedup factor: %.2f\n", speedup)

	validation_helper(t, tree.RootNode(), sourceCode)
}

func TestTreeSitterJsEditDeleteMultiple(t *testing.T) {
	// Setup new parser
	language := sitter.NewLanguage(javascript.Language())
	parser := sitter.NewParser()
	defer parser.Close()
	parser.SetLanguage(language)

	start := time.Now()

	sourceCode := []byte(`function hello() { console.log('hello') }`)
	oldEndIndex := uint(len(sourceCode))

	// Parse source code
	tree := parser.ParseCtx(context.Background(), []byte(sourceCode), nil)
	defer tree.Close()
	assert.NotNil(t, tree, "error occurred in tree sitter in js edit delete multiple")

    elapsedFirst := time.Since(start)
	fmt.Println("parsed, elapsed", time.Since(start))
	validation_helper(t, tree.RootNode(), sourceCode)

	fmt.Println("Edit input")

	sourceCode = []byte(`function hel() { console.log('hello') }`)
	newEndIndex := uint(len(sourceCode))

	tree.Edit(&sitter.InputEdit{
		StartByte:  14,
		OldEndByte: oldEndIndex,
		NewEndByte: newEndIndex,
		StartPosition: sitter.Point{Row: 0, Column: 14 },
		OldEndPosition: sitter.Point{Row: 0, Column: oldEndIndex},
		NewEndPosition: sitter.Point{Row: 0, Column: newEndIndex},
	})

	start = time.Now()
	tree = parser.ParseCtx(context.Background(), []byte(sourceCode), tree)
	assert.NotNil(t, tree, "error occurred in tree sitter js edit")

	elapsedSecond := time.Since(start)
	fmt.Println("parsed again, elapsed", elapsedSecond)
	speedup := float64(elapsedFirst) / float64(elapsedSecond)
	fmt.Printf("Speedup factor: %.2f\n", speedup)

	validation_helper(t, tree.RootNode(), sourceCode)

	fmt.Println("Edit input")

	sourceCode = []byte(`function h() { console.log('hello') }`)
	newEndIndex = uint(len(sourceCode))

	tree.Edit(&sitter.InputEdit{
		StartByte:  10,
		OldEndByte: oldEndIndex,
		NewEndByte: newEndIndex,
		StartPosition: sitter.Point{Row: 0, Column: 10 },
		OldEndPosition: sitter.Point{Row: 0, Column: oldEndIndex},
		NewEndPosition: sitter.Point{Row: 0, Column: newEndIndex},
	})

	start = time.Now()
	tree = parser.ParseCtx(context.Background(), []byte(sourceCode), tree)
	assert.NotNil(t, tree, "error occurred in tree sitter js edit")

	elapsedSecond = time.Since(start)
	fmt.Println("parsed again, elapsed", elapsedSecond)
	speedup = float64(elapsedFirst) / float64(elapsedSecond)
	fmt.Printf("Speedup factor: %.2f\n", speedup)

	validation_helper(t, tree.RootNode(), sourceCode)
}

func TestTreeSitterJsEditEnter(t *testing.T) {
	// Setup new parser
	language := sitter.NewLanguage(javascript.Language())
	parser := sitter.NewParser()
	defer parser.Close()
	parser.SetLanguage(language)

	start := time.Now()

	sourceCode := []byte(`function hello() { console.log('hello') }`)
	oldEndIndex := uint(len(sourceCode))

	// Parse source code
	tree := parser.ParseCtx(context.Background(), []byte(sourceCode), nil)
	defer tree.Close()
	assert.NotNil(t, tree, "error occurred in tree sitter in js edit enter")

    elapsedFirst := time.Since(start)
	fmt.Println("parsed, elapsed", time.Since(start))
	validation_helper(t, tree.RootNode(), sourceCode)

	fmt.Println("Edit input")

	sourceCode = []byte(`function hello() 
{ console.log('hello') }`)
	newEndIndex := uint(len(sourceCode))

	tree.Edit(&sitter.InputEdit{
		StartByte:  19,
		OldEndByte: oldEndIndex,
		NewEndByte: newEndIndex,
		StartPosition: sitter.Point{Row: 1, Column: 0 },
		OldEndPosition: sitter.Point{Row: 0, Column: oldEndIndex},
		NewEndPosition: sitter.Point{Row: 1, Column: newEndIndex},
	})

	start = time.Now()
	tree = parser.ParseCtx(context.Background(), []byte(sourceCode), tree)
	assert.NotNil(t, tree, "error occurred in tree sitter js edit")

	elapsedSecond := time.Since(start)
	fmt.Println("parsed again, elapsed", elapsedSecond)
	speedup := float64(elapsedFirst) / float64(elapsedSecond)
	fmt.Printf("Speedup factor: %.2f\n", speedup)

	validation_helper(t, tree.RootNode(), sourceCode)
}
