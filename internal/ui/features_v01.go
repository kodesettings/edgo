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

	e.code = ConvertLinesToString(e.Lines)

	e.Undo = append(e.Undo, ops)
	if e.Col < 0 { e.Col = 0 }
	e.OnDown(false)
	e.Update = true
	e.UpdateLsp(false, e.code)
	e.IsContentChanged = true
}

func (e *Editor) OnSwapLinesUp() {
	e.Focus()

	if e.Row == 0 { return }
	var ops = EditOperation{}
	ops = append(ops, Operation{MoveCursor, ' ', e.Row, e.Col})

	line1 := e.Lines[e.Row].Buf; line2 := e.Lines[e.Row-1].Buf

	for i := len(line1)-1; i >= 0; i-- { ops = append(ops, Operation{Delete, line1[i], e.Row, i}) }
	for i := len(line2)-1; i >= 0; i-- { ops = append(ops, Operation{Delete, line2[i], e.Row -1, i}) }
	for i, ch := range line1 { ops = append(ops, Operation{Insert, ch, e.Row -1, i}) }
	for i, ch := range line2 { ops = append(ops, Operation{Insert, ch, e.Row, i}) }

	e.Lines[e.Row].Buf = line2; e.Lines[e.Row-1].Buf = line1 // swap
	e.Row--

	e.code = ConvertLinesToString(e.Lines)

	e.Undo = append(e.Undo, ops)
	e.Selection.CleanSelection()
	e.Update = true
	e.IsContentChanged = true
	e.FindTests()
}

func (e *Editor) OnSwapLinesDown() {
	e.Focus()

	if e.Row+1 == len(e.Lines) { return }

	var ops = EditOperation{}
	ops = append(ops, Operation{MoveCursor, ' ', e.Row, e.Col})

	line1 := e.Lines[e.Row].Buf; line2 := e.Lines[e.Row+1].Buf

	for i := len(line1)-1; i >= 0; i-- { ops = append(ops, Operation{Delete, line1[i], e.Row, i}) }
	for i := len(line2)-1; i >= 0; i-- { ops = append(ops, Operation{Delete, line2[i], e.Row +1, i}) }
	for i, ch := range line1 { ops = append(ops, Operation{Insert, ch, e.Row +1, i}) }
	for i, ch := range line2 { ops = append(ops, Operation{Insert, ch, e.Row, i}) }

	e.Lines[e.Row].Buf = line2; e.Lines[e.Row+1].Buf = line1 // swap
	e.Row++

	e.code = ConvertLinesToString(e.Lines)

	e.Undo = append(e.Undo, ops)
	e.Selection.CleanSelection()
	e.Update = true
	e.IsContentChanged = true
	e.FindTests()
}
