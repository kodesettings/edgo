package ui

import (
	"testing"
	. "github.com/vipmax/edgo/internal/highlighter"
	. "github.com/vipmax/edgo/internal/operations"
	"github.com/zyedidia/rope"
	assert "github.com/stretchr/testify/assert"
)

var e *Editor = new(Editor)

func apply_highlighter(code []byte, theme string, lang string) {
	e.treeSitterHighlighter = NewTreeSitter()
	e.treeSitterHighlighter.SetTheme(theme)
	e.treeSitterHighlighter.SetLang(lang)
	e.treeSitterHighlighter.Parse(code)
}

func TestCutActionPosition(t *testing.T) {
	e.code = rope.New([]byte("this is a sample text\nand another line of text to cut\none more line of text\n"))

	e.Selection.Ssx = 0
	e.Selection.Ssy = 1
	e.Selection.Sex = 11
	e.Selection.Sey = 2

	got := e.Selection.GetSelectionString(e.code.Value())
	expected := "and another line of text to cut\none more li"
	assert.Equal(t, expected, got, "selection of string error")

	text := string(e.code.Value())
	apply_highlighter(e.code.Value(), "", "")
	e.Cut(true)

	expected = "this is a sample text\ne of text\n"
	assert.Equal(t, expected, string(e.code.Value()), "cut lines mismatch")

	e.OnUndo()
	assert.Equal(t, text, string(e.code.Value()), "undo cut lines mismatch")

	e.OnRedo()
	assert.Equal(t, expected, string(e.code.Value()), "redo cut lines mismatch")
}

func TestCutActionLinesOnly(t *testing.T) {
	e.code = rope.New([]byte("this is a sample text\nand another line of text to cut\none more line of text\n"))

	e.Selection.Ssx = 0
	e.Selection.Ssy = 1
	e.Selection.Sex = 0
	e.Selection.Sey = 3

	got := e.Selection.GetSelectionString(e.code.Value())
	expected := "and another line of text to cut\none more line of text"
	assert.Equal(t, expected, got, "selection of string error")

	text := string(e.code.Value())
	apply_highlighter(e.code.Value(), "", "")
	e.Cut(true)

	expected = "this is a sample text\n"
	assert.Equal(t, expected, string(e.code.Value()), "cut lines mismatch")

	e.OnUndo()
	assert.Equal(t, text, string(e.code.Value()), "undo cut lines mismatch")

	e.OnRedo()
	assert.Equal(t, expected, string(e.code.Value()), "redo cut lines mismatch")
}

func TestDuplicateAction(t *testing.T) {
	e.Row = 0
	e.Col = 0

	e.code = rope.New([]byte("this is a sample text\nand another line of text to cut\none more line of text"))

	text := string(e.code.Value())
	apply_highlighter(e.code.Value(), "", "")
	e.Duplicate()

	expected := "this is a sample text\nthis is a sample text\nand another line of text to cut\none more line of text"
	assert.Equal(t, expected, string(e.code.Value()), "duplicate lines mismatch")

	e.OnUndo()
	assert.Equal(t, text, string(e.code.Value()), "undo duplicate lines mismatch")

	e.OnRedo()
	assert.Equal(t, expected, string(e.code.Value()), "redo duplicate lines mismatch")
}

func TestUndoRedoStack(t *testing.T) {
	e.Row = 0
	e.Col = 0

	e.code = rope.New([]byte(""))
	apply_highlighter(e.code.Value(), "", "")

	e.Undo = make([]EditOperation, 0)
	e.Redo = make([]EditOperation, 0)

	e.AddChar('a')
	e.AddChar('b')
	e.AddChar('c')
	e.AddChar('d')

	expected := "abcd"
	assert.Equal(t, expected, string(e.code.Value()), "added characters error")

	e.OnUndo()
	assert.Equal(t, "abc", string(e.code.Value()), "undo stack mismatch")
	e.OnUndo()
	assert.Equal(t, "ab", string(e.code.Value()), "undo stack mismatch")
	e.OnUndo()
	e.OnUndo()
	e.OnUndo()
	assert.Equal(t, "", string(e.code.Value()), "undo stack mismatch")

	e.OnUndo()
	e.OnRedo()
	assert.Equal(t, "a", string(e.code.Value()), "redo stack mismatch")

	e.OnRedo()
	e.OnRedo()
	e.OnRedo()
	e.OnRedo()
	assert.Equal(t, "abcd", string(e.code.Value()), "redo stack mismatch")
}
