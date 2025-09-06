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
		if len(selectionString) == 0 { return } // skipping function if selection is empty
		clipboard.WriteAll(selectionString)
	}

	selectedIndices := e.Selection.GetSelectedIndices(e.code.Value())
	lastElement := len(selectedIndices) - 1

	exd := selectedIndices[lastElement][0]
	eyd := selectedIndices[lastElement][1]
	sxd := selectedIndices[0][0]
	syd := selectedIndices[0][1]

	e.Row = syd
	e.Col = sxd

	sxd += LineOffset(e.code.Value(), syd) + 1
	exd += LineOffset(e.code.Value(), eyd) + 2

	text := make([]byte, exd - sxd)
	copy(text, e.code.Slice(sxd, exd))

	e.code.Remove(sxd, exd)
	e.treeSitterHighlighter.RemoveTextEdit(e.code.Value(), sxd, len(text))

	e.Undo = append(e.Undo, EditOperation{{Delete, text, sxd, CursorMove{e.Row, e.Col}}})
	e.Selection.CleanSelection()
	e.UpdateLsp(false, string(e.code.Value()))
	e.set_update_parameters(true)
}

func (e *Editor) Duplicate() {
	e.Focus()

	if e.code.Len() == 0 { return }
	syd := LineOffset(e.code.Value(), e.Row) + 1
	eyd := LineOffset(e.code.Value(), e.Row + 1) + 1
	if e.Row == 0 { syd-- } // this is required for first line only
	duplicatedSlice := e.code.Slice(syd, eyd)

	e.code.Insert(eyd, duplicatedSlice)
	e.Row++

	eyd = LineOffset(e.code.Value(), e.Row) + 1

	e.treeSitterHighlighter.InsertTextEdit(e.code.Value(), syd, len(duplicatedSlice))
	e.Undo = append(e.Undo, EditOperation{{Insert, duplicatedSlice, eyd, CursorMove{e.Row, e.Col}}})
	e.UpdateLsp(false, string(e.code.Value()))
	e.set_update_parameters(true)
}

func (e *Editor) OnUndo() {
	if len(e.Undo) == 0 { e.UpdateLsp(true, string(e.code.Value())); return }

	lastOperation := e.Undo[len(e.Undo)-1]
	e.Undo = e.Undo[:len(e.Undo)-1]
	e.Focus()

	index := len(lastOperation) - 1
	var o Operation
undo:
	if len(lastOperation) > 0 { o = lastOperation[index] } else { goto exit }

	switch o.Action {
	case Insert, AddChar:
		e.code.Remove(o.Offset, o.Offset + len(o.Text))
		e.treeSitterHighlighter.RemoveTextEdit(e.code.Value(), o.Offset, len(o.Text))
		e.Row = o.Cursor.Line; e.Col = o.Cursor.Column
	break;
	case Delete:
		e.code.Insert(o.Offset, o.Text)
		e.treeSitterHighlighter.InsertTextEdit(e.code.Value(), o.Offset, len(o.Text))
		e.Row = o.Cursor.Line; e.Col = o.Cursor.Column
	break;
	}

	if index > 0 { index--; goto undo; }
exit:
	e.Redo = append(e.Redo, lastOperation)
	e.UpdateLsp(false, string(e.code.Value()))
	e.set_update_parameters(true)
}

func (e *Editor) OnRedo() {
	if len(e.Redo) == 0 { return }

	lastRedoOperation := e.Redo[len(e.Redo)-1]
	e.Redo = e.Redo[:len(e.Redo)-1]

	index, offset := 0, 0
	var o Operation
redo:
	if len(lastRedoOperation) > 0 { o = lastRedoOperation[index] } else { goto exit }

	switch o.Action {
	case AddChar:
		offset++
		fallthrough
	case Insert:
		e.code.Insert(o.Offset, o.Text)
		e.treeSitterHighlighter.InsertTextEdit(e.code.Value(), o.Offset, len(o.Text))
		e.Row = o.Cursor.Line; e.Col = o.Cursor.Column + offset
	break;
	case Delete:
		e.code.Remove(o.Offset, o.Offset + len(o.Text))
		e.treeSitterHighlighter.RemoveTextEdit(e.code.Value(), o.Offset, len(o.Text))
		e.Row = o.Cursor.Line; e.Col = o.Cursor.Column
	break;
	}

	if index < len(lastRedoOperation) - 1 { index++; goto redo; }
exit:
	e.Undo = append(e.Undo, lastRedoOperation)
	e.UpdateLsp(false, string(e.code.Value()))
	e.set_update_parameters(true)
}
