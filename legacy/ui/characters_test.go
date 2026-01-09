package ui

import (
	"testing"
	"github.com/zyedidia/rope"
	assert "github.com/stretchr/testify/assert"
)

func TestAddChar(t *testing.T) {
	e.code = rope.New([]byte("this is a sample text"))

	e.Row = 0
	e.Col = 11

	apply_highlighter(e.code.Value(), "", "")
	e.AddChar('g')

	text := "this is a sample text"
	expected := "this is a sgample text"
	assert.Equal(t, expected, string(e.code.Value()), "add char error")

	e.OnUndo()
	assert.Equal(t, text, string(e.code.Value()), "undo add char mismatch")

	e.OnRedo()
	assert.Equal(t, expected, string(e.code.Value()), "redo add char mismatch")
}

func TestInsertCharacter(t *testing.T) {
	e.code = rope.New([]byte("this is a sample text"))

	e.Row = 0
	e.Col = 2

	apply_highlighter(e.code.Value(), "", "")
	e.InsertCharacter(0, 2, 'f')

	text := "this is a sample text"
	expected := "thfis is a sample text"
	assert.Equal(t, expected, string(e.code.Value()), "insert character error")

	e.OnUndo()
	assert.Equal(t, text, string(e.code.Value()), "undo insert character mismatch")

	e.OnRedo()
	assert.Equal(t, expected, string(e.code.Value()), "redo insert character mismatch")
}

func TestInsertString(t *testing.T) {
	e.code = rope.New([]byte("this is a sample text"))

	e.Row = 0
	e.Col = 1

	apply_highlighter(e.code.Value(), "", "")
	e.InsertString(0, 1, "some")

	text := "this is a sample text"
	expected := "thsomeis is a sample text"
	assert.Equal(t, expected, string(e.code.Value()), "insert string error")

	e.OnUndo()
	assert.Equal(t, text, string(e.code.Value()), "undo insert string mismatch")

	e.OnRedo()
	assert.Equal(t, expected, string(e.code.Value()), "redo insert string mismatch")
}

func TestInsertLines(t *testing.T) {
	e.code = rope.New([]byte("this is a sample text"))

	e.Row = 0
	e.Col = 10

	apply_highlighter(e.code.Value(), "", "")
	e.InsertString(0, 10, "new line\nanother line\n")

	text := "this is a sample text"
	expected := "this is a snew line\nanother line\nample text"
	assert.Equal(t, expected, string(e.code.Value()), "insert lines error")

	e.OnUndo()
	assert.Equal(t, text, string(e.code.Value()), "undo insert lines mismatch")

	e.OnRedo()
	assert.Equal(t, expected, string(e.code.Value()), "redo insert lines mismatch")
}

func TestDeleteCharacter(t *testing.T) {
	e.code = rope.New([]byte("this is a sample text"))

	e.Row = 0
	e.Col = 2

	apply_highlighter(e.code.Value(), "", "")
	e.DeleteCharacter(0, 2)

	text := "this is a sample text"
	expected := "ths is a sample text"
	assert.Equal(t, expected, string(e.code.Value()), "delete character error")

	e.OnUndo()
	assert.Equal(t, text, string(e.code.Value()), "undo delete character mismatch")

	e.OnRedo()
	assert.Equal(t, expected, string(e.code.Value()), "redo delete character mismatch")
}

func TestReplaceString(t *testing.T) {
	e.code = rope.New([]byte("this is a sample text"))

	apply_highlighter(e.code.Value(), "", "")
	e.ReplaceString(0, 3, 11, "asd")

	text := "this is a sample text"
	expected := "thisasdmple text"
	assert.Equal(t, expected, string(e.code.Value()), "delete string error")

	e.OnUndo()
	assert.Equal(t, text, string(e.code.Value()), "undo delete string mismatch")

	e.OnRedo()
	assert.Equal(t, expected, string(e.code.Value()), "redo delete string mismatch")
}

func TestDeleteLine(t *testing.T) {
	e.code = rope.New([]byte("this is a sample text\n"))

	text := "this is a sample text\n"
	expected := "this is a sample text"

	apply_highlighter(e.code.Value(), "", "")
	e.DeleteCharacter(0, len(text) - 1)

	e.OnUndo()
	assert.Equal(t, text, string(e.code.Value()), "undo delete line mismatch")

	e.OnRedo()
	assert.Equal(t, expected, string(e.code.Value()), "redo delete line mismatch")
}

func TestInsertEnter(t *testing.T) {
	e.code = rope.New([]byte("this is a sample text\n"))

	text := "this is a sample text\n"
	expected := "this is a sample text\n\n"

	apply_highlighter(e.code.Value(), "", "")
	e.InsertCharacter(0, len(text) - 1, '\n')

	e.OnUndo()
	assert.Equal(t, text, string(e.code.Value()), "undo insert enter mismatch")

	e.OnRedo()
	assert.Equal(t, expected, string(e.code.Value()), "redo insert enter mismatch")
}

func TestShiftWithTabsToRight(t *testing.T) {
	e.code = rope.New([]byte("this is a sample text\none more line"))

	e.Selection.Ssx = 0
	e.Selection.Ssy = 0
	e.Selection.Sex = 0
	e.Selection.Sey = 2

	got := e.Selection.GetSelectionString(e.code.Value())
	text := "this is a sample text\none more line"
	assert.Equal(t, text, got, "selection of string error")

	apply_highlighter(e.code.Value(), "", "")

	selectedLines := e.Selection.GetSelectedLines(e.code.Value())
	e.ShiftWithTabsToRight(0, 0, selectedLines)

	expected := "\tthis is a sample text\t\none more line"
	assert.Equal(t, expected, string(e.code.Value()), "shift with tabs error")

	e.OnUndo()
	assert.Equal(t, text, string(e.code.Value()), "undo shift with tabs mismatch")

	e.OnRedo()
	assert.Equal(t, expected, string(e.code.Value()), "redo shift with tabs mismatch")
}

func TestMaybeAddPair(t *testing.T) {
	e.code = rope.New([]byte("test["))

	val, found := e.MaybeAddPair(0, 4, '[')

	assert.Equal(t, ']', val, "maybeaddpair error")
	assert.Equal(t, true, found, "maybeaddpair error")
}
