package ui

import (
	. "github.com/vipmax/edgo/internal/operations"
	. "github.com/vipmax/edgo/internal/utils"
	"github.com/atotto/clipboard"
)

// this function sets the required parameters to validate updates
func (e *Editor) set_update_parameters(changed bool) {
	e.Update = true
	e.IsContentChanged = changed
	e.FindTests()
}

func (e *Editor) OnCopy() {
	selectionString := e.Selection.GetSelectionString(e.code.Value())
	clipboard.WriteAll(selectionString)
}

func (e *Editor) OnPaste() {

	if e.Selection.IsSelectionNonEmpty() {
		e.Cut(false)
	}

	text, _ := clipboard.ReadAll()

	if len(text) == 0 { return }
	e.InsertString(e.Row, e.Col, text)
	e.set_update_parameters(true)
}

func (e *Editor) Cut(isCopySelected bool) {
	e.Focus()

	if (e.code.Len() < 0) { e.Row, e.Col = 0, 0; return }

	if isCopySelected {
		selectionString := e.Selection.GetSelectionString(e.code.Value())
		clipboard.WriteAll(selectionString)
	}

	selectedIndices := e.Selection.GetSelectedIndices(e.code.Value())
	lastElement := len(selectedIndices) - 1

	exd := selectedIndices[lastElement][0]
	eyd := selectedIndices[lastElement][1]
	sxd := selectedIndices[0][0]
	syd := selectedIndices[0][1]

	sxd += LineOffset(e.code.Value(), syd)
	exd += LineOffset(e.code.Value(), eyd)
	e.code.Remove(sxd, exd + 2)

	e.Col, e.Row = exd, eyd

	e.treeSitterHighlighter.UpdateCharsEdit(e.code.Value(), syd, sxd, eyd, exd)
	e.Undo = append(e.Undo, EditOperation{}) // TODO: refactor this
	e.Selection.CleanSelection()

	e.UpdateLsp(false, string(e.code.Value()))
	e.set_update_parameters(true)
}

func (e *Editor) Duplicate() {
	e.Focus()

	if e.code.Len() == 0 { return }
	syd := LineOffset(e.code.Value(), e.Row)
	eyd := LineOffset(e.code.Value(), e.Row + 1)
	duplicatedSlice := e.code.Slice(syd, eyd)

	e.code.Insert(eyd, []byte("\n"))
	e.code.Insert(eyd + 1, duplicatedSlice)
	e.Row++

	eyd = LineOffset(e.code.Value(), e.Row)
	e.treeSitterHighlighter.UpdateCharsEdit(e.code.Value(), syd, 0, eyd, 0)
	e.Undo = append(e.Undo, EditOperation{}) // TODO: refactor this

	e.UpdateLsp(false, string(e.code.Value()))
	e.set_update_parameters(true)
}

func (e *Editor) OnUndo() {
	if len(e.Undo) == 0 { e.UpdateLsp(true, string(e.code.Value())); return }

	lastOperation := e.Undo[len(e.Undo)-1]
	e.Undo = e.Undo[:len(e.Undo)-1]
	e.Focus()

	fromCol := e.Col
	fromRow := e.Row

	index := len(lastOperation) - 1
undo:
	o := lastOperation[index]

	switch o.Action {
	case Insert, Enter:
		e.code.Remove(o.Line + o.Column, o.Line + o.Column + 1)
	break;
	case Delete:
		e.code.Insert(o.Line + o.Column, []byte(string(o.Char)))
	break;
	case MoveCursor:
		e.Row = o.Line; e.Col = o.Column
	break;
	}

	if index > 0 { index--; goto undo; }

	e.treeSitterHighlighter.UpdateCharsEdit(e.code.Value(), fromRow, fromCol, e.Row, e.Col)

	e.Redo = append(e.Redo, lastOperation)
	e.UpdateLsp(false, string(e.code.Value()))
	e.set_update_parameters(true)
}

func (e *Editor) OnRedo() {
	if len(e.Redo) == 0 { return }

	lastRedoOperation := e.Redo[len(e.Redo)-1]
	e.Redo = e.Redo[:len(e.Redo)-1]

	fromCol := e.Col
	fromRow := e.Row

	index := 0
redo:
	o := lastRedoOperation[index]

	switch o.Action {
	case Insert, Enter:
		e.code.Insert(o.Line + o.Column, []byte(string(o.Char)))
	break;
	case Delete:
		e.code.Remove(o.Line + o.Column, o.Line + o.Column + 1)
	break;
	case MoveCursor:
		e.Row = o.Line; e.Col = o.Column
	break;
	}

	if index < len(lastRedoOperation) - 1 { index++; goto redo; }

	e.treeSitterHighlighter.UpdateCharsEdit(e.code.Value(), fromRow, fromCol, e.Row, e.Col)

	e.Undo = append(e.Undo, lastRedoOperation)
	e.UpdateLsp(false, string(e.code.Value()))
	e.set_update_parameters(true)
}
