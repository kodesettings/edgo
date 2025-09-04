package ui

import (
	. "github.com/vipmax/edgo/internal/operations"
	. "github.com/vipmax/edgo/internal/utils"
)

func (e *Editor) OnCommentLine() {
	e.Focus()

	var found bool
	offset := LineOffset(e.code.Value(), e.Row)
	if e.Row > 0 { offset++ } // due to newline character
	ch_1 := e.code.At(offset)
	ch_2 := e.code.At(offset + 1)

	if len(e.langConf.Comment) == 1 && ch_1 == e.langConf.Comment[0] {
		e.code.Remove(offset, offset + 1)
		e.treeSitterHighlighter.RemoveTextEdit(e.code.Value(), offset, 1)
		e.Undo = append(e.Undo, EditOperation{{Delete, []byte{ch_1}, offset, CursorMove{offset, 0}}})
		found = true
	} else if len(e.langConf.Comment) == 2 && ch_1 == e.langConf.Comment[0] && ch_2 == e.langConf.Comment[1] {
		e.code.Remove(offset, offset + 2)
		e.treeSitterHighlighter.RemoveTextEdit(e.code.Value(), offset, 2)
		e.Undo = append(e.Undo, EditOperation{{Delete, []byte{ch_1, ch_2}, offset, CursorMove{offset, 0}}})
		found = true
	}

	if found { goto exit }

	e.code.Insert(offset, []byte(e.langConf.Comment))
	e.treeSitterHighlighter.InsertTextEdit(e.code.Value(), offset, len(e.langConf.Comment))
	e.Undo = append(e.Undo, EditOperation{{Insert, []byte(e.langConf.Comment), offset, CursorMove{e.Row, 0}}})
exit:
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
	text_len := len(append(line_1, line_2...))
	e.treeSitterHighlighter.InsertTextEdit(e.code.Value(), from, text_len)
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
	text_len := len(append(line_1, line_2...))
	e.treeSitterHighlighter.InsertTextEdit(e.code.Value(), from, text_len)
	e.UpdateLsp(false, string(e.code.Value()))
	e.set_update_parameters(true)
}
