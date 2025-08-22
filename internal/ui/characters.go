package ui

import (
	. "github.com/vipmax/edgo/internal/operations"
	. "github.com/vipmax/edgo/internal/utils"
)

func (e *Editor) AddChar(ch rune) {
	if len(e.Selection.GetSelectionString(e.code.Value())) != 0 { e.Cut(false) }

	e.Focus()
	e.InsertCharacter(e.Row, e.Col, ch)
	e.Col++

	val, found := e.MaybeAddPair(e.Row, e.Col, ch)
	if found {
		e.InsertCharacter(e.Row, e.Col, val)
	}

	e.Update = true
	e.UpdateLsp(false, string(e.code.Value()))
	e.IsContentChanged = true
	e.FindTests()
}

func (e *Editor) InsertCharacter(line, pos int, ch rune) {
	offset := LineOffset(e.code.Value(), line) + pos
	e.code.Insert(offset, []byte(string(ch)))
	e.treeSitterHighlighter.InsertTextEdit(e.code.Value(), offset, len(string(ch)))
	e.Undo = append(e.Undo, EditOperation{{Insert, []byte(string(ch)), offset, CursorMove{line, pos}}})
	if len(e.Redo) > 0 { e.Redo = make([]EditOperation, 0) }
}

func (e *Editor) InsertString(line, pos int, linestring string) {
	// modify string to remove tabs and space from beginning
	l := RemoveLeadingTabsSpaces(linestring)
	offset := LineOffset(e.code.Value(), line) + pos

	// record the operation on the undo stack. Note that we're creating a new EditOperation
	// and adding all the Operations to it
	e.Undo = append(e.Undo, EditOperation{{Insert, []byte(l), offset, CursorMove{line, pos}}})
	if len(e.Redo) > 0 { e.Redo = make([]EditOperation, 0) }

	e.code.Insert(offset, []byte(l))
	e.Col++
}

func (e *Editor) DeleteCharacter(line, pos int) {
	offset := LineOffset(e.code.Value(), line) + pos
	ch := e.code.At(offset)
	e.code.Remove(offset, offset + 1)
	e.treeSitterHighlighter.RemoveTextEdit(e.code.Value(), offset, len(string(ch)))
	e.Undo = append(e.Undo, EditOperation{{Delete, []byte{ch}, offset, CursorMove{line, pos + 1}}})
	if len(e.Redo) > 0 { e.Redo = make([]EditOperation, 0) }
}

func (e *Editor) ShiftWithTabsToRight(line, pos int, selectedLines []int) {
	e.Selection.Ssx = 0

	var ops = EditOperation{}
	for _, linenumber := range selectedLines {
		offset := LineOffset(e.code.Value(), linenumber)
		e.code.Insert(offset, []byte("\t"))
		ops = append(ops, Operation{Insert, []byte("\t"), offset, CursorMove{line, pos}})
	}

	e.Selection.Sex = pos
	e.Undo = append(e.Undo, ops)
	if len(e.Redo) > 0 { e.Redo = make([]EditOperation, 0) }
}

func (e *Editor) MaybeAddPair(line, pos int, ch rune) (rune, bool) {
	pairMap := map[rune]rune{
		'(': ')', '{': '}', '[': ']',
		'"': '"', '\'': '\'', '`': '`',
	}

	if e.code.Len() == 0 { return rune('\000'), false }

	offset_start := LineOffset(e.code.Value(), line)
	offset_end := LineOffset(e.code.Value(), line + 1)
	line_str := e.code.Slice(offset_start, offset_end)
	current_char := e.code.At(offset_start + pos - 1)

	if closeChar, found := pairMap[ch]; found {
		noMoreChars := pos >= len(line_str) - 1
		isSpaceNext := pos < len(line_str) - 1 && current_char == ' '
		isStringAndClosedBracketNext := closeChar == '"' && pos < len(line_str) - 1 && current_char == ')'

		if noMoreChars || isSpaceNext || isStringAndClosedBracketNext {
			return closeChar, found
		}
	}

	return rune('\000'), false
}
