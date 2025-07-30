package ui

import (
	. "github.com/vipmax/edgo/internal/operations"
	. "github.com/vipmax/edgo/internal/utils"
	"github.com/atotto/clipboard"
	"strings"
)

func (e *Editor) OnCopy() {
	selectionString := e.Selection.GetSelectionString(e.Content)
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

	if len(e.Content) < 1 {
		e.Content[0] = []rune{};
		e.Row, e.Col = 0, 0
		return
	}

	var ops = EditOperation{}

	if isCopySelected {
		selectionString := e.Selection.GetSelectionString(e.Content)
		clipboard.WriteAll(selectionString)
	}

	ops = append(ops, Operation{MoveCursor, ' ', e.Row, e.Col})
	selectedIndices := e.Selection.GetSelectedIndices(e.Content)

	firstCol := e.Col
	firstRow := e.Row

	// Sort selectedIndices in reverse order to delete characters from the end
	index := len(selectedIndices) - 1
repeat:
	indices := selectedIndices[index]
	xd := indices[0]
	yd := indices[1]
	e.Col, e.Row = xd, yd

	if len(e.Content[yd]) > 0 {
		// Delete the character at index (x, j)
		ch := e.Content[yd][xd]
		ops = append(ops, Operation{Delete, ch, yd, xd})
		e.Content[yd] = append(e.Content[yd][:xd], e.Content[yd][xd+1:]...)
	}

	if len(e.Content[yd]) == 0 { // delete Line
		if e.Row == 0 {
			ops = append(ops, Operation{DeleteLine, '\n', 0, 0})
		} else {
			ops = append(ops, Operation{DeleteLine, '\n', e.Row -1, len(e.Content[e.Row-1])})
		}

		e.Content = append(e.Content[:yd], e.Content[yd+1:]...)
	}

	if index > 0 { index--; goto repeat; }

	code := ConvertContentToString(e.Content)
	e.treeSitterHighlighter.RemoveCharsEdit(&code, firstCol, firstRow, e.Col, e.Row)

	if len(e.Content) == 0 {
		e.Content = make([][]rune, 1)
	}

	if e.Row >= len(e.Content)  {
		e.Row = len(e.Content) - 1
		if e.Col >= len(e.Content[e.Row]) { e.Col = len(e.Content[e.Row]) - 1 }
	}

	if e.Row < 0 { e.Row = 0 }
	if e.Col < 0 { e.Col = 0 }

	e.UpdateColors()
	e.Undo = append(e.Undo, ops)
	e.Selection.CleanSelection()
	e.Update = true
	e.UpdateLsp(false, ConvertContentToString(e.Content))
	e.IsContentChanged = true
	e.UpdateNeeded() // optimize
}

func (e *Editor) Duplicate() {
	e.Focus()

	if len(e.Content) == 0 { return }

	if e.Selection.Ssx == -1 && e.Selection.Ssy == -1 ||
		e.Selection.Ssx == e.Selection.Sex && e.Selection.Ssy == e.Selection.Sey {
		var ops = EditOperation{}
		ops = append(ops, Operation{MoveCursor, ' ', e.Row, e.Col})
		ops = append(ops, Operation{Enter, '\n', e.Row, len(e.Content[e.Row])})

		duplicatedSlice := make([]rune, len(e.Content[e.Row]))
		copy(duplicatedSlice, e.Content[e.Row])
		for i, ch := range duplicatedSlice {
			ops = append(ops, Operation{Insert, ch, e.Row, i})
		}
		e.Row++
		e.Content = InsertTo(e.Content, e.Row, duplicatedSlice)

		e.UpdateColors()
		e.Undo = append(e.Undo, ops)
		e.Update = true
		e.IsContentChanged = true
		e.FindTests()
	} else {
		selection := e.Selection.GetSelectionString(e.Content)
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
	if len(e.Undo) == 0 { e.UpdateLsp(true, ConvertContentToString(e.Content)); return }

	lastOperation := e.Undo[len(e.Undo)-1]
	e.Undo = e.Undo[:len(e.Undo)-1]
	e.Focus()

	index := len(lastOperation) - 1
repeat:
	o := lastOperation[index]

	if o.Action == Insert {
		e.Row = o.Line; e.Col = o.Column
		e.Content[e.Row] = append(e.Content[e.Row][:e.Col], e.Content[e.Row][e.Col+1:]...)
	} else if o.Action == Delete {
		e.Row = o.Line; e.Col = o.Column
		e.Content[e.Row] = InsertTo(e.Content[e.Row], e.Col, o.Char)
	} else if o.Action == Enter {
		// Merge lines
		e.Content[o.Line] = append(e.Content[o.Line], e.Content[o.Line+1]...)
		e.Content = append(e.Content[:o.Line+1], e.Content[o.Line+2:]...)
		e.Row = o.Line; e.Col = o.Column
	} else if o.Action == DeleteLine {
		// Insert enter
		e.Row = o.Line; e.Col = o.Column
		after := e.Content[e.Row][e.Col:]
		before := e.Content[e.Row][:e.Col]
		e.Content[e.Row] = before
		e.Row++; e.Col = 0
		newline := append([]rune{}, after...)
		e.Content = InsertTo(e.Content, e.Row, newline)
	} else if o.Action == MoveCursor {
		e.Row = o.Line; e.Col = o.Column
	} else {
		e.OnCursorChanged()
	}

	if index > 0 { index--; goto repeat; }

	e.UpdateColors()
	e.UpdateLsp(false, ConvertContentToString(e.Content))
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
		e.Content[e.Row] = InsertTo(e.Content[e.Row], e.Col, o.Char)
		e.Col++
	} else if o.Action == Delete {
		e.Row = o.Line; e.Col = o.Column
		e.Content[e.Row] = append(e.Content[e.Row][:e.Col], e.Content[e.Row][e.Col+1:]...)
	} else if o.Action == Enter {
		e.Row = o.Line; e.Col = o.Column
		after := e.Content[e.Row][e.Col:]
		before := e.Content[e.Row][:e.Col]
		e.Content[e.Row] = before
		e.Row++; e.Col = 0
		newline := append([]rune{}, after...)
		e.Content = InsertTo(e.Content, e.Row, newline)
	} else if o.Action == DeleteLine {
		// Merge lines
		e.Content[o.Line] = append(e.Content[o.Line], e.Content[o.Line+1]...)
		e.Content = append(e.Content[:o.Line+1], e.Content[o.Line+2:]...)
		e.Row = o.Line; e.Col = o.Column
	} else if o.Action == MoveCursor {
		e.Row = o.Line; e.Col = o.Column
	}

	if index < len(lastRedoOperation) - 1 { index++; goto repeat; }

	e.UpdateColors()
	e.UpdateLsp(false, ConvertContentToString(e.Content))
	e.Undo = append(e.Undo, lastRedoOperation)
	e.UpdateNeeded()
}
