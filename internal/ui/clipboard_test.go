package ui

import (
	"testing"
	. "github.com/vipmax/edgo/internal/highlighter"
	. "github.com/vipmax/edgo/internal/utils"
	assert "github.com/stretchr/testify/assert"
)

var e *Editor = new(Editor)

func apply_highlighter(lines []Line, theme string, lang string) {
	e.code = ConvertLinesToString(lines)
	e.treeSitterHighlighter = NewTreeSitter()
	e.treeSitterHighlighter.SetTheme(theme)
	e.treeSitterHighlighter.SetLang(lang)
	e.treeSitterHighlighter.Parse(&e.code)
}

func TestCutActionPosition(t *testing.T) {
	e.Lines = []Line{
		Line{Buf: []rune("this is a sample text")},
		Line{Buf: []rune("and another line of text to cut")},
		Line{Buf: []rune("one more line of text")},
	}

	text := ConvertLinesToString(e.Lines)

	e.Selection.Ssx = 0
	e.Selection.Ssy = 1
	e.Selection.Sex = 11
	e.Selection.Sey = 2

	got := e.Selection.GetSelectionString(e.Lines)
	expected := "and another line of text to cut\none more li"
	assert.Equal(t, expected, got, "selection of string error")

	apply_highlighter(e.Lines, "", "")

	e.Cut(true)

	expected = "this is a sample text\nne of text"
	assert.Equal(t, expected, e.code, "cut lines mismatch")

	e.OnUndo()
	assert.Equal(t, text, e.code, "undo cut lines mismatch")

	e.OnRedo()
	assert.Equal(t, expected, e.code, "redo cut lines mismatch")
}

func TestCutActionLinesOnly(t *testing.T) {
	e.Lines = []Line{
		Line{Buf: []rune("this is a sample text")},
		Line{Buf: []rune("and another line of text to cut")},
		Line{Buf: []rune("one more line of text")},
	}

	text := ConvertLinesToString(e.Lines)

	e.Selection.Ssx = 0
	e.Selection.Ssy = 1
	e.Selection.Sex = 0
	e.Selection.Sey = 3

	got := e.Selection.GetSelectionString(e.Lines)
	expected := "and another line of text to cut\none more line of text"
	assert.Equal(t, expected, got, "selection of string error")

	apply_highlighter(e.Lines, "", "")

	e.Cut(true)

	expected = "this is a sample text"
	assert.Equal(t, expected, e.code, "cut lines mismatch")

	e.OnUndo()
	assert.Equal(t, text, e.code, "undo cut lines mismatch")

	e.OnRedo()
	assert.Equal(t, expected, e.code, "redo cut lines mismatch")
}

func TestDuplicateAction(t *testing.T) {
	e.Lines = []Line{
		Line{Buf: []rune("this is a sample text")},
		Line{Buf: []rune("and another line of text to cut")},
		Line{Buf: []rune("one more line of text")},
	}

	text := ConvertLinesToString(e.Lines)

	e.Selection.Ssx = 15
	e.Selection.Ssy = 0
	e.Selection.Sex = 0
	e.Selection.Sey = 0

	got := e.Selection.GetSelectionString(e.Lines)
	expected := "this is a sampl"
	assert.Equal(t, expected, got, "selection of string error")

	apply_highlighter(e.Lines, "", "")

	e.Duplicate()

	expected = "this is a sample text\nthis is a sampland another line of text to cut\none more line of text"
	assert.Equal(t, expected, e.code, "duplicate lines mismatch")

	e.OnUndo()
	assert.Equal(t, text, e.code, "undo duplicate lines mismatch")

	e.OnRedo()
	assert.Equal(t, expected, e.code, "redo duplicate lines mismatch")
}

func TestUndoRedoStack(t *testing.T) {
	e.Lines = []Line{Line{Buf: []rune{}}}

	e.Row = 0
	e.Col = 0

	apply_highlighter(e.Lines, "", "")

	e.AddChar('a')
	e.AddChar('b')
	e.AddChar('c')
	e.AddChar('d')

	expected := "abcd"
	assert.Equal(t, expected, e.code, "added characters error")

	e.OnUndo()
	assert.Equal(t, "abc", e.code, "undo stack mismatch")
	e.OnUndo()
	assert.Equal(t, "ab", e.code, "undo stack mismatch")
	e.OnUndo()
	e.OnUndo()
	e.OnUndo()
	assert.Equal(t, "", e.code, "undo stack mismatch")

	e.OnUndo()
	e.OnRedo()
	assert.Equal(t, "a", e.code, "redo stack mismatch")

	e.OnRedo()
	e.OnRedo()
	e.OnRedo()
	e.OnRedo()
	assert.Equal(t, "abcd", e.code, "redo stack mismatch")
}
