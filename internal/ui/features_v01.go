package ui

import (
	. "github.com/vipmax/edgo/internal/operations"
	. "github.com/vipmax/edgo/internal/utils"
)

func (e *Editor) OnCommentLine() {
	e.Focus()

	var found bool
	offset := LineOffset(e.code.Value(), e.Row)
	ch_1 := e.code.At(offset)
	ch_2 := e.code.At(offset + 1)

	if len(e.langConf.Comment) == 1 && ch_1 == e.langConf.Comment[0] {
		e.code.Remove(offset, offset + 1)
		e.Undo = append(e.Undo, EditOperation{{Delete, []byte{ch_1}, offset, CursorMove{offset, 0}}})
		found = true
	} else if len(e.langConf.Comment) == 2 && ch_1 == e.langConf.Comment[0] && ch_2 == e.langConf.Comment[1] {
		e.code.Remove(offset, offset + 2)
		e.Undo = append(e.Undo, EditOperation{{Delete, []byte{ch_1, ch_2}, offset, CursorMove{offset, 0}}})
		found = true
	}

	tabs := CountTabs(e.code.Value(), offset)
	spaces := CountSpaces(e.code.Value(), offset)
	from := tabs

	if found { goto exit }
	if tabs == 0 && spaces != 0 { from = spaces }

	e.code.Insert(from, []byte(e.langConf.Comment))
	e.Undo = append(e.Undo, EditOperation{{Insert, []byte(e.langConf.Comment), from, CursorMove{e.Row, from}}})
exit:
	e.treeSitterHighlighter.UpdateCharsEdit(e.code.Value(), offset, 0, offset + 1, 0)
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

	e.Undo = append(e.Undo, EditOperation{
		{Delete, line_2, to + 1, CursorMove{e.Row, 0}},
		{Insert, line_2, from, CursorMove{e.Row, 0}},
	})

	e.code.Remove(to, offset) // remove line_2 from current position
	e.code.Insert(from, line_2) // add line_2 to top
	offset = LineOffset(e.code.Value(), e.Row + 1)

	e.Row--
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

	e.Undo = append(e.Undo, EditOperation{
		{Delete, line_1, from, CursorMove{e.Row, 0}},
		{Insert, line_1, to - 1, CursorMove{e.Row, 0}},
	})

	e.code.Remove(from, to) // remove line_1 from current position
	offset = LineOffset(e.code.Value(), e.Row + 1)
	e.code.Insert(offset, line_1) // add line_1 to bottom

	e.Row--
	e.treeSitterHighlighter.UpdateCharsEdit(e.code.Value(), from, 0, offset, 0)
	e.UpdateLsp(false, string(e.code.Value()))
	e.set_update_parameters(true)
}
