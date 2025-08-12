package ui

import (
	"testing"
	. "github.com/vipmax/edgo/internal/utils"
	assert "github.com/stretchr/testify/assert"
)

func TestOnCommentLine(t *testing.T) {
	e.Lines = []Line{
		Line{Buf: []rune("this is a sample text")},
	}

	apply_highlighter(e.Lines, "", "")

	e.langConf.Comment = "//"
	e.OnCommentLine()

	expected := "//this is a sample text"
	assert.Equal(t, expected, e.code, "comment line mismatch")

	e.OnUndo()
	actual := "this is a sample text"
	assert.Equal(t, actual, e.code, "undo comment line mismatch")

	e.OnRedo()
	assert.Equal(t, expected, e.code, "redo comment line mismatch")
}

func TestOnSwapLinesUp(t *testing.T) {
	e.Lines = []Line{
		Line{Buf: []rune("first line of text")},
		Line{Buf: []rune("second line of text")},
		Line{Buf: []rune("third line of text")},
	}

	e.Row = 1
	e.Col = 0

	apply_highlighter(e.Lines, "", "")

	e.OnSwapLinesUp()

	expected := "second line of text\nfirst line of text\nthird line of text"
	assert.Equal(t, expected, e.code, "swaplinesup mismatch")

	e.OnUndo()
	actual := "first line of text\nsecond line of text\nthird line of text"
	assert.Equal(t, actual, e.code, "undo swaplinesup mismatch")

	e.OnRedo()
	assert.Equal(t, expected, e.code, "redo swaplinesup mismatch")
}

func TestOnSwapLinesDown(t *testing.T) {
	e.Lines = []Line{
		Line{Buf: []rune("first line of text")},
		Line{Buf: []rune("second line of text")},
		Line{Buf: []rune("third line of text")},
	}

	e.Row = 1
	e.Col = 0

	apply_highlighter(e.Lines, "", "")

	e.OnSwapLinesDown()

	expected := "first line of text\nthird line of text\nsecond line of text"
	assert.Equal(t, expected, e.code, "swaplinesup mismatch")

	e.OnUndo()
	actual := "first line of text\nsecond line of text\nthird line of text"
	assert.Equal(t, actual, e.code, "undo swaplinesup mismatch")

	e.OnRedo()
	assert.Equal(t, expected, e.code, "redo swaplinesup mismatch")
}
