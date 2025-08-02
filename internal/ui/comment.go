package ui

import (
	. "github.com/vipmax/edgo/internal/operations"
	. "github.com/vipmax/edgo/internal/utils"
)

func (e *Editor) OnCommentLine() {
	e.Focus()

	found := false
	index := 0
repeat:
	if len(e.Lines[e.Row].Buf) == 0 { return } // exit this function if useless
	ch := e.Lines[e.Row].Buf[index]

	if len(e.langConf.Comment) == 1 && ch == rune(e.langConf.Comment[0]) {
		// found 1 char comment, uncomment
		e.Col = index
		e.Undo = append(e.Undo, EditOperation{
			{MoveCursor, ch, e.Row, index+1},
			{Delete, ch, e.Row, index},
		})
		e.Lines[e.Row].Buf = Remove(e.Lines[e.Row].Buf, index)

		code := ConvertLinesToString(e.Lines)
		e.treeSitterHighlighter.RemoveCharEdit(&code, e.Row, index, ch)

		found = true
	} else if len(e.langConf.Comment) == 2 && ch == rune(e.langConf.Comment[0]) &&
			e.Lines[e.Row].Buf[index+1] == rune(e.langConf.Comment[1]) {
		// found 2 char comment, uncomment
		e.Col = index
		e.Undo = append(e.Undo, EditOperation{
			{MoveCursor, ch, e.Row, index+1},
			{Delete, ch, e.Row, index},
			{MoveCursor, e.Lines[e.Row].Buf[index+1], e.Row, index+1},
			{Delete, ch, e.Row, index},
		})
		e.Lines[e.Row].Buf = Remove(e.Lines[e.Row].Buf, index)
		e.Lines[e.Row].Buf = Remove(e.Lines[e.Row].Buf, index)

		code := ConvertLinesToString(e.Lines)
		e.treeSitterHighlighter.RemoveCharEdit(&code, e.Row, index, ch)
		e.treeSitterHighlighter.RemoveCharEdit(&code, e.Row, index, ch)

		found = true
	}

	if index < len(e.Lines[e.Row].Buf) - 1 { index++; goto repeat; }

	if found {
		if e.Col < 0 { e.Col = 0 }
		e.OnDown(false)
		e.Update = true
		e.IsContentChanged = true
		return
	}

	tabs := CountTabs(e.Lines[e.Row].Buf, e.Col)
	spaces := CountSpaces(e.Lines[e.Row].Buf, e.Col)

	from := tabs
	if tabs == 0 && spaces != 0 { from = spaces }

	e.Col = from
	ops := EditOperation{}
	for _, ch := range e.langConf.Comment {
		e.Lines[e.Row].Buf = InsertTo(e.Lines[e.Row].Buf, from, ch)
		code := ConvertLinesToString(e.Lines)
		e.treeSitterHighlighter.AddCharEdit(&code, e.Row, from, ch)
		ops = append(ops, Operation{Insert, ch, e.Row, from})
	}

	e.Undo = append(e.Undo, ops)
	if e.Col < 0 { e.Col = 0 }
	e.OnDown(false)
	e.Update = true
	e.UpdateLsp(false, ConvertLinesToString(e.Lines))
	e.IsContentChanged = true
}

