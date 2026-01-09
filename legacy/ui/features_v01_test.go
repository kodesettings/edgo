package ui

import (
	"testing"
	"github.com/zyedidia/rope"
	assert "github.com/stretchr/testify/assert"
)

func TestOnCommentLine(t *testing.T) {
	e.code = rope.New([]byte("this is a sample text"))

	apply_highlighter(e.code.Value(), "", "")

	e.langConf.Comment = "//"
	e.OnCommentLine()

	expected := "//this is a sample text"
	assert.Equal(t, expected, string(e.code.Value()), "comment line mismatch")

	e.OnUndo()
	actual := "this is a sample text"
	assert.Equal(t, actual, string(e.code.Value()), "undo comment line mismatch")

	e.OnRedo()
	assert.Equal(t, expected, string(e.code.Value()), "redo comment line mismatch")
}

func TestOnSwapLinesUp(t *testing.T) {
	e.code = rope.New([]byte("first line of text\nsecond line of text\nthird line of text\n"))

	e.Row = 1
	e.Col = 0

	apply_highlighter(e.code.Value(), "", "")
	e.OnSwapLinesUp()

	expected := "second line of text\nfirst line of text\nthird line of text\n"
	assert.Equal(t, expected, string(e.code.Value()), "swaplinesup mismatch")

	e.OnUndo()
	actual := "first line of text\nsecond line of text\nthird line of text\n"
	assert.Equal(t, actual, string(e.code.Value()), "undo swaplinesup mismatch")

	e.OnRedo()
	assert.Equal(t, expected, string(e.code.Value()), "redo swaplinesup mismatch")
}

func TestOnSwapLinesDown(t *testing.T) {
	e.code = rope.New([]byte("first line of text\nsecond line of text\nthird line of text\n"))

	e.Row = 1
	e.Col = 0

	apply_highlighter(e.code.Value(), "", "")
	e.OnSwapLinesDown()

	expected := "first line of text\nthird line of text\nsecond line of text\n"
	assert.Equal(t, expected, string(e.code.Value()), "swaplinesup mismatch")

	e.OnUndo()
	actual := "first line of text\nsecond line of text\nthird line of text\n"
	assert.Equal(t, actual, string(e.code.Value()), "undo swaplinesup mismatch")

	e.OnRedo()
	assert.Equal(t, expected, string(e.code.Value()), "redo swaplinesup mismatch")
}
