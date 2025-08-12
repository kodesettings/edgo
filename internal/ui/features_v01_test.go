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
	expected = "this is a sample text"
	assert.Equal(t, expected, e.code, "undo comment line mismatch")

	e.OnRedo()
	expected = "//this is a sample text"
	assert.Equal(t, expected, e.code, "redo comment line mismatch")
}
