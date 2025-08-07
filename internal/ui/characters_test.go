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

func TestMaybeAddPair(t *testing.T) {
	e.Lines = []Line{
		Line{Buf: []rune("")},
	}

	e.Row = 0
	e.Col = 0

	e.MaybeAddPair('[')
	assert.Equal(t, "]", e.code, "maybeaddpair error")
}
