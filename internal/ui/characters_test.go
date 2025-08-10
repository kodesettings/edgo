package ui

import (
	"testing"
	. "github.com/vipmax/edgo/internal/utils"
	assert "github.com/stretchr/testify/assert"
)

func TestAddChar(t *testing.T) {
	e.Lines = []Line{
		Line{Buf: []rune("this is a sample text")},
	}

	e.Row = 0
	e.Col = 11

	apply_highlighter(e.Lines, "", "")

	e.AddChar('g')

	text := "this is a sample text"
	expected := "this is a sgample text"
	assert.Equal(t, expected, e.code, "add char error")

	e.OnUndo()
	assert.Equal(t, text, e.code, "undo add char mismatch")

	e.OnRedo()
	assert.Equal(t, expected, e.code, "redo add char mismatch")
}

func TestInsertCharacter(t *testing.T) {
	e.Lines = []Line{
		Line{Buf: []rune("this is a sample text")},
	}

	e.Row = 0
	e.Col = 11

	e.InsertCharacter(0, 2, 'f')

	text := "this is a sample text"
	expected := "thfis is a sample text"
	assert.Equal(t, expected, e.code, "insert character error")

	e.OnUndo()
	assert.Equal(t, text, e.code, "undo insert character mismatch")

	e.OnRedo()
	assert.Equal(t, expected, e.code, "redo insert character mismatch")
}

func TestInsertString(t *testing.T) {
	e.Lines = []Line{
		Line{Buf: []rune("this is a sample text")},
	}

	e.Row = 0
	e.Col = 11

	e.InsertString(0, 2, "some")

	text := "this is a sample text"
	expected := "thsomeis is a sample text"
	assert.Equal(t, expected, e.code, "insert string error")

	e.OnUndo()
	assert.Equal(t, text, e.code, "undo insert string mismatch")

	e.OnRedo()
	assert.Equal(t, expected, e.code, "redo insert string mismatch")
}

func TestInsertLines(t *testing.T) {
	e.Lines = []Line{
		Line{Buf: []rune("this is a sample text")},
	}

	e.Row = 0
	e.Col = 11

	line := []string{"new line", "another line"}
	e.InsertLines(0, 2, line)

	e.code = ConvertLinesToString(e.Lines)

	text := "this is a sample text"
	expected := "this is a snew line\nanother line\nthis is a sample text"
	assert.Equal(t, expected, e.code, "insert lines error")

	e.OnUndo()
	assert.Equal(t, text, e.code, "undo insert lines mismatch")

	e.OnRedo()
	assert.Equal(t, expected, e.code, "redo insert lines mismatch")
}

func TestDeleteCharacter(t *testing.T) {
	e.Lines = []Line{
		Line{Buf: []rune("this is a sample text")},
	}

	e.Row = 0
	e.Col = 11

	e.DeleteCharacter(0, 2)

	text := "this is a sample text"
	expected := "ths is a sample text"
	assert.Equal(t, expected, e.code, "delete character error")

	e.OnUndo()
	assert.Equal(t, text, e.code, "undo delete character mismatch")

	e.OnRedo()
	assert.Equal(t, expected, e.code, "redo delete character mismatch")
}

func TestDeleteLine(t *testing.T) {
	e.Lines = []Line{
		Line{Buf: []rune("this is a sample text")},
		Line{Buf: []rune("\n")},
	}

	text := "this is a sample text\n"
	expected := "this is a sample text"

	assert.Equal(t, 2, len(e.Lines), "delete line error")
	e.DeleteLine(1, 0)
	assert.Equal(t, 1, len(e.Lines), "delete line error")

	e.OnUndo()
	assert.Equal(t, text, e.code, "undo delete line mismatch")

	e.OnRedo()
	assert.Equal(t, expected, e.code, "redo delete line mismatch")
}

func TestInsertEnter(t *testing.T) {
	e.Lines = []Line{
		Line{Buf: []rune("this is a sample text")},
		Line{Buf: []rune("\n")},
	}

	text := "this is a sample text\n"
	expected := "this is a sample text\n\n"

	assert.Equal(t, 2, len(e.Lines), "insert enter error")
	e.InsertEnter(1, 0)
	assert.Equal(t, 3, len(e.Lines), "insert enter error")

	e.OnUndo()
	assert.Equal(t, text, e.code, "undo insert enter mismatch")

	e.OnRedo()
	assert.Equal(t, expected, e.code, "redo insert enter mismatch")
}

func TestShiftWithTabsToRight(t *testing.T) {
	e.Lines = []Line{
		Line{Buf: []rune("this is a sample text")},
		Line{Buf: []rune("one more line")},
	}

	text := ConvertLinesToString(e.Lines)

	e.Selection.Ssx = 0
	e.Selection.Ssy = 0
	e.Selection.Sex = 0
	e.Selection.Sey = 2

	got := e.Selection.GetSelectionString(e.Lines)
	expected := "this is a sample text\none more line"
	assert.Equal(t, expected, got, "selection of string error")

	selectedLines := e.Selection.GetSelectedLines(e.Lines)
	e.ShiftWithTabsToRight(0, 0, selectedLines)

	e.code = ConvertLinesToString(e.Lines)
	expected = "\tthis is a sample text\n\tone more line"
	assert.Equal(t, expected, e.code, "shift with tabs error")

	e.OnUndo()
	assert.Equal(t, text, e.code, "undo shift with tabs mismatch")

	e.OnRedo()
	assert.Equal(t, expected, e.code, "redo shift with tabs mismatch")
}

func TestMaybeAddPair(t *testing.T) {
	e.Lines = []Line{
		Line{Buf: []rune("")},
	}

	e.Row = 0
	e.Col = 0

	e.MaybeAddPair('[')
	assert.Equal(t, "]", e.code, "maybeaddpair error")
}
