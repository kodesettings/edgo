package ui

import (
	. "github.com/vipmax/edgo/internal/operations"
	. "github.com/vipmax/edgo/internal/utils"
)

func (e *Editor) OnCommentLine() {
	e.Focus()

	var found bool
	pos := LineOffset(e.code.Value(), e.Row)
	ch_1 := e.code.At(pos)
	ch_2 := e.code.At(pos + 1)

	if len(e.langConf.Comment) == 1 && ch_1 == e.langConf.Comment[0] {
		e.code.Remove(pos, pos + 1)
		e.Undo = append(e.Undo, EditOperation{
			{MoveCursor, rune(ch_1), e.Row, pos+1},
			{Delete, rune(ch_1), e.Row, pos},
		})
		found = true
	} else if len(e.langConf.Comment) == 2 && ch_1 == e.langConf.Comment[0] && ch_2 == e.langConf.Comment[1] {
		e.code.Remove(pos, pos + 1)
		e.code.Remove(pos, pos + 1)
		e.Undo = append(e.Undo, EditOperation{
			{MoveCursor, rune(ch_1), e.Row, pos + 1},
			{Delete, rune(ch_1), e.Row, pos},
			{MoveCursor, rune(ch_2), e.Row, pos + 1},
			{Delete, rune(ch_2), e.Row, pos},
		})
		found = true
	}

	tabs := CountTabs(e.code.Value(), e.Row + e.Col)
	spaces := CountSpaces(e.code.Value(), e.Row + e.Col)

	from := tabs
	ops := EditOperation{}

	if found { goto exit }
	if tabs == 0 && spaces != 0 { from = spaces }

	e.code.Insert(from, []byte(e.langConf.Comment))
	for _, ch := range e.langConf.Comment {
		ops = append(ops, Operation{Insert, ch, e.Row, from})
	}
	e.Undo = append(e.Undo, ops)
exit:
	e.treeSitterHighlighter.UpdateCharsEdit(e.code.Value(), pos, 0, pos + 1, 0)
	e.UpdateLsp(false, string(e.code.Value()))
	e.set_update_parameters(true)
}

func (e *Editor) OnSwapLinesUp() {
	e.Focus()

	if e.Row == 0 { return }

	from := LineOffset(e.code.Value(), e.Row - 1)
	to := LineOffset(e.code.Value(), e.Row)
	offset := LineOffset(e.code.Value(), e.Row + 1)

	line_1 := make([]byte, to - from)
	copy(line_1, e.code.Slice(from, to + 1))

	line_2 := make([]byte, offset - to)
	copy(line_2, e.code.Slice(to + 1, offset + 1))

	// TOOD: record undo/redo operations here

	e.code.Remove(to, offset)  // remove line_2 from current position
	e.code.Insert(from, line_2) // add line_2 to top
	offset = LineOffset(e.code.Value(), e.Row + 1)

	e.Row--
	e.Undo = append(e.Undo, EditOperation{})
	e.treeSitterHighlighter.UpdateCharsEdit(e.code.Value(), from, 0, offset, 0)
	e.UpdateLsp(false, string(e.code.Value()))
	e.set_update_parameters(true)
}

func (e *Editor) OnSwapLinesDown() {
	e.Focus()

	if e.Row < 1 { return }

	from := LineOffset(e.code.Value(), e.Row)
	to := LineOffset(e.code.Value(), e.Row + 1)
	offset := LineOffset(e.code.Value(), e.Row + 2)

	line_1 := make([]byte, to - from)
	copy(line_1, e.code.Slice(from, to + 1))

	line_2 := make([]byte, offset - to)
	copy(line_2, e.code.Slice(to + 1, offset))

	// TOOD: record undo/redo operations here

	e.code.Remove(from, to)  // remove line_1 from current position
	offset = LineOffset(e.code.Value(), e.Row + 1)
	e.code.Insert(offset, line_1) // add line_1 to bottom

	e.Row--
	e.Undo = append(e.Undo, EditOperation{})
	e.treeSitterHighlighter.UpdateCharsEdit(e.code.Value(), from, 0, offset, 0)
	e.UpdateLsp(false, string(e.code.Value()))
	e.set_update_parameters(true)
}
