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

	expected := "this is a sgample text"
	assert.Equal(t, expected, e.code, "add char error")
}

func TestInsertCharacter(t *testing.T) {
	e.Lines = []Line{
		Line{Buf: []rune("this is a sample text")},
	}

	e.Row = 0
	e.Col = 11

	e.InsertCharacter(0, 2, 'f')

	expected := "thfis is a sample text"
	assert.Equal(t, expected, e.code, "insert character error")
}

func TestInsertString(t *testing.T) {
	e.Lines = []Line{
		Line{Buf: []rune("this is a sample text")},
	}

	e.Row = 0
	e.Col = 11

	e.InsertString(0, 2, "some")

	e.code = ConvertLinesToString(e.Lines)
	expected := "thsomeis is a sample text"
	assert.Equal(t, expected, e.code, "insert string error")
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
	expected := "this is a snew line\nanother line\nthis is a sample text"
	assert.Equal(t, expected, e.code, "insert lines error")
}

func TestDeleteCharacter(t *testing.T) {
	e.Lines = []Line{
		Line{Buf: []rune("this is a sample text")},
	}

	e.Row = 0
	e.Col = 11

	e.DeleteCharacter(0, 2)

	expected := "ths is a sample text"
	assert.Equal(t, expected, e.code, "delete character error")
}

func TestDeleteLine(t *testing.T) {
	e.Lines = []Line{
		Line{Buf: []rune("this is a sample text")},
		Line{Buf: []rune("\n")},
	}

	assert.Equal(t, 2, len(e.Lines), "delete line error")
	e.DeleteLine(1, 0)
	assert.Equal(t, 1, len(e.Lines), "delete line error")
}

func TestInsertEnter(t *testing.T) {
	e.Lines = []Line{
		Line{Buf: []rune("this is a sample text")},
		Line{Buf: []rune("\n")},
	}

	assert.Equal(t, 2, len(e.Lines), "insert enter error")
	e.InsertEnter(1, 0)
	assert.Equal(t, 3, len(e.Lines), "insert enter error")
}

func TestShiftWithTabsToRight(t *testing.T) {
	e.Lines = []Line{
		Line{Buf: []rune("this is a sample text")},
		Line{Buf: []rune("one more line")},
	}

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
