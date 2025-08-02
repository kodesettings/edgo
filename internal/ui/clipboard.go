package ui

import (
	. "github.com/vipmax/edgo/internal/operations"
	. "github.com/vipmax/edgo/internal/utils"
	"github.com/atotto/clipboard"
	"strings"
)

func (e *Editor) OnCopy() {
	selectionString := e.Selection.GetSelectionString(e.Lines)
	clipboard.WriteAll(selectionString)
}

func (e *Editor) OnPaste() {

	if e.Selection.IsSelectionNonEmpty() {
		e.Cut(false)
	}

	text, _ := clipboard.ReadAll()
	lines := strings.Split(text, "\n")

	if len(lines) == 0 { return }

	if len(lines) == 1 { // single Line paste
		e.InsertString(e.Row, e.Col, lines[0])
	}

	if len(lines) > 1 { // multiple Line paste
		e.InsertLines(e.Row, e.Col, lines)
	}
	
	e.Update = true
	e.UpdateNeeded()
}

func (e *Editor) Cut(isCopySelected bool) {
	e.Focus()

	if len(e.Lines) < 1 {
		e.Lines[0].Buf = []rune{};
		e.Row, e.Col = 0, 0
		return
	}

	var ops = EditOperation{}

	if isCopySelected {
		selectionString := e.Selection.GetSelectionString(e.Lines)
		clipboard.WriteAll(selectionString)
	}

	ops = append(ops, Operation{MoveCursor, ' ', e.Row, e.Col})
	selectedIndices := e.Selection.GetSelectedIndices(e.Lines)

	firstCol := e.Col
	firstRow := e.Row

	// Sort selectedIndices in reverse order to delete characters from the end
	index := len(selectedIndices) - 1
repeat:
	indices := selectedIndices[index]
	xd := indices[0]
	yd := indices[1]
	e.Col, e.Row = xd, yd

	if len(e.Lines[yd].Buf) > 0 {
		// Delete the character at index (x, j)
		ch := e.Lines[yd].Buf[xd]
		ops = append(ops, Operation{Delete, ch, yd, xd})
		e.Lines[yd].Buf = append(e.Lines[yd].Buf[:xd], e.Lines[yd].Buf[xd+1:]...)
	}

	if len(e.Lines[yd].Buf) == 0 { // delete Line
		if e.Row == 0 {
			ops = append(ops, Operation{DeleteLine, '\n', 0, 0})
		} else {
			ops = append(ops, Operation{DeleteLine, '\n', e.Row -1, len(e.Lines[e.Row-1].Buf)})
		}

		e.Lines = append(e.Lines[:yd], e.Lines[yd+1:]...)
	}

	if index > 0 { index--; goto repeat; }

	code := ConvertLinesToString(e.Lines)
	e.treeSitterHighlighter.RemoveCharsEdit(&code, firstCol, firstRow, e.Col, e.Row)

	if len(e.Lines) == 0 {
		e.Lines = make([]Line, 1)
	}

	if e.Row >= len(e.Lines)  {
		e.Row = len(e.Lines) - 1
		if e.Col >= len(e.Lines[e.Row].Buf) { e.Col = len(e.Lines[e.Row].Buf) - 1 }
	}

	if e.Row < 0 { e.Row = 0 }
	if e.Col < 0 { e.Col = 0 }

	e.UpdateColors()
	e.Undo = append(e.Undo, ops)
	e.Selection.CleanSelection()
	e.Update = true
	e.UpdateLsp(false, ConvertLinesToString(e.Lines))
	e.IsContentChanged = true
	e.UpdateNeeded() // optimize
}

func (e *Editor) Duplicate() {
	e.Focus()

	if len(e.Lines) == 0 { return }

	if e.Selection.Ssx == -1 && e.Selection.Ssy == -1 ||
		e.Selection.Ssx == e.Selection.Sex && e.Selection.Ssy == e.Selection.Sey {
		var ops = EditOperation{}
		ops = append(ops, Operation{MoveCursor, ' ', e.Row, e.Col})
		ops = append(ops, Operation{Enter, '\n', e.Row, len(e.Lines[e.Row].Buf)})

		duplicatedSlice := make([]rune, len(e.Lines[e.Row].Buf))
		copy(duplicatedSlice, e.Lines[e.Row].Buf)
		for i, ch := range duplicatedSlice {
			ops = append(ops, Operation{Insert, ch, e.Row, i})
		}
		e.Row++
		e.Lines = InsertTo(e.Lines, e.Row, Line{duplicatedSlice})

		e.UpdateColors()
		e.Undo = append(e.Undo, ops)
		e.Update = true
		e.IsContentChanged = true
		e.FindTests()
	} else {
		selection := e.Selection.GetSelectionString(e.Lines)
		if len(selection) == 0 { return }
		lines := strings.Split(selection, "\n")

		if len(lines) == 0 { return }

		if len(lines) == 1 { // single Line
			lines[0] = " " + lines[0]// add space before
			e.InsertString(e.Row, e.Col, lines[0])
		}

		if len(lines) > 1 { // multiple Line
			e.InsertLines(e.Row, e.Col, lines)
		}
		e.UpdateColors()
		e.Selection.CleanSelection()
		e.UpdateNeeded()
	}

}

func (e *Editor) OnUndo() {
	if len(e.Undo) == 0 { e.UpdateLsp(true, ConvertLinesToString(e.Lines)); return }

	lastOperation := e.Undo[len(e.Undo)-1]
	e.Undo = e.Undo[:len(e.Undo)-1]
	e.Focus()

	index := len(lastOperation) - 1
repeat:
	o := lastOperation[index]

	if o.Action == Insert {
		e.Row = o.Line; e.Col = o.Column
		e.Lines[e.Row].Buf = append(e.Lines[e.Row].Buf[:e.Col], e.Lines[e.Row].Buf[e.Col+1:]...)
	} else if o.Action == Delete {
		e.Row = o.Line; e.Col = o.Column
		e.Lines[e.Row].Buf = InsertTo(e.Lines[e.Row].Buf, e.Col, o.Char)
	} else if o.Action == Enter {
		// Merge lines
		e.Lines[o.Line].Buf = append(e.Lines[o.Line].Buf, e.Lines[o.Line+1].Buf...)
		e.Lines = append(e.Lines[:o.Line+1], e.Lines[o.Line+2:]...)
		e.Row = o.Line; e.Col = o.Column
	} else if o.Action == DeleteLine {
		// Insert enter
		e.Row = o.Line; e.Col = o.Column
		after := e.Lines[e.Row].Buf[e.Col:]
		before := e.Lines[e.Row].Buf[:e.Col]
		e.Lines[e.Row].Buf = before
		e.Row++; e.Col = 0
		newline := append([]rune{}, after...)
		e.Lines = InsertTo(e.Lines, e.Row, Line{newline})
	} else if o.Action == MoveCursor {
		e.Row = o.Line; e.Col = o.Column
	} else {
		e.OnCursorChanged()
	}

	if index > 0 { index--; goto repeat; }

	e.UpdateColors()
	e.UpdateLsp(false, ConvertLinesToString(e.Lines))
	e.Redo = append(e.Redo, lastOperation)
	e.UpdateNeeded()
}

func (e *Editor) OnRedo() {
	if len(e.Redo) == 0 { return }

	lastRedoOperation := e.Redo[len(e.Redo)-1]
	e.Redo = e.Redo[:len(e.Redo)-1]

	index := 0
repeat:
	o := lastRedoOperation[index]

	if o.Action == Insert {
		e.Row = o.Line; e.Col = o.Column
		e.Lines[e.Row].Buf = InsertTo(e.Lines[e.Row].Buf, e.Col, o.Char)
		e.Col++
	} else if o.Action == Delete {
		e.Row = o.Line; e.Col = o.Column
		e.Lines[e.Row].Buf = append(e.Lines[e.Row].Buf[:e.Col], e.Lines[e.Row].Buf[e.Col+1:]...)
	} else if o.Action == Enter {
		e.Row = o.Line; e.Col = o.Column
		after := e.Lines[e.Row].Buf[e.Col:]
		before := e.Lines[e.Row].Buf[:e.Col]
		e.Lines[e.Row].Buf = before
		e.Row++; e.Col = 0
		newline := append([]rune{}, after...)
		e.Lines = InsertTo(e.Lines, e.Row, Line{newline})
	} else if o.Action == DeleteLine {
		// Merge lines
		e.Lines[o.Line].Buf = append(e.Lines[o.Line].Buf, e.Lines[o.Line+1].Buf...)
		e.Lines = append(e.Lines[:o.Line+1], e.Lines[o.Line+2:]...)
		e.Row = o.Line; e.Col = o.Column
	} else if o.Action == MoveCursor {
		e.Row = o.Line; e.Col = o.Column
	}

	if index < len(lastRedoOperation) - 1 { index++; goto repeat; }

	e.UpdateColors()
	e.UpdateLsp(false, ConvertLinesToString(e.Lines))
	e.Undo = append(e.Undo, lastRedoOperation)
	e.UpdateNeeded()
}
