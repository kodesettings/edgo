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

func TestCutAction(t *testing.T) {
	e.Lines = []Line{
		Line{Buf: []rune("this is a sample text")},
		Line{Buf: []rune("and another line of text to cut")},
		Line{Buf: []rune("one more line of text")},
	}

	e.Selection.Ssx = 0
	e.Selection.Ssy = 1
	e.Selection.Sex = 11
	e.Selection.Sey = 2

	got := e.Selection.GetSelectionString(e.Lines)
	expected := "and another line of text to cut\none more li"
	assert.Equal(t, expected, got, "selection of string error")

	apply_highlighter(e.Lines, "", "")

	e.Cut(true)

	e.code = ConvertLinesToString(e.Lines)
	assert.Equal(t, "this is a sample text\nne of text", e.code, "cut lines mismatch")
}

func TestDuplicateAction(t *testing.T) {
	e.Lines = []Line{
		Line{Buf: []rune("this is a sample text")},
		Line{Buf: []rune("and another line of text to cut")},
		Line{Buf: []rune("one more line of text")},
	}

	e.Selection.Ssx = 15
	e.Selection.Ssy = 0
	e.Selection.Sex = 0
	e.Selection.Sey = 0

	got := e.Selection.GetSelectionString(e.Lines)
	expected := "this is a sampl"
	assert.Equal(t, expected, got, "selection of string error")

	apply_highlighter(e.Lines, "", "")

	e.Duplicate()

	e.code = ConvertLinesToString(e.Lines)
	expected = "this is a sample text\nthis is a sampland another line of text to cut\none more line of text"
	assert.Equal(t, expected, e.code, "duplicate lines mismatch")
}
