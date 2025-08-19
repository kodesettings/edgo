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

	if len(e.Redo) > 0 { e.Redo = []EditOperation{} }

	e.Update = true
	e.UpdateLsp(false, string(e.code.Value()))
	e.IsContentChanged = true
	e.FindTests()
}

func (e *Editor) InsertCharacter(line, pos int, ch rune) {
	pos += LineOffset(e.code.Value(), line)
	e.code.Insert(pos, []byte(string(ch)))
	e.Undo = append(e.Undo, EditOperation{{Insert, ch, line, pos}})
	e.treeSitterHighlighter.AddCharEdit(e.code.Value(), line, pos, ch)
}

func (e *Editor) InsertString(line, pos int, linestring string) {
	// Convert the string to insert to a slice of runes
	l := RemoveLeadingTabsSpaces(linestring)
	pos += LineOffset(e.code.Value(), line)

	// Record the operation on the undo stack. Note that we're creating a new EditOperation
	// and adding all the Operations to it
	var ops = EditOperation{}
	for _, ch := range l {
		e.code.Insert(pos, []byte(string(ch)))
		ops = append(ops, Operation{Insert, ch, line, pos})
		pos++
	}

	e.Col = pos
	e.Undo = append(e.Undo, ops)
}

func (e *Editor) DeleteCharacter(line, pos int) {
	pos += LineOffset(e.code.Value(), line)
	ch := rune(e.code.At(pos))
	e.Undo = append(e.Undo, EditOperation{
		{MoveCursor, ch, line, pos + 1},
		{Delete, ch, line, pos},
	})

	e.code.Remove(pos, pos + 1)
	e.treeSitterHighlighter.RemoveCharEdit(e.code.Value(), line, pos, ch)
}

func (e *Editor) ShiftWithTabsToRight(line, pos int, selectedLines []int) {
	e.Selection.Ssx = 0

	var ops = EditOperation{}
	for _, linenumber := range selectedLines {
		line = LineOffset(e.code.Value(), linenumber)
		e.code.Insert(line, []byte(string('\t')))
		ops = append(ops, Operation{Insert, '\t', line, 0})
	}

	e.Selection.Sex = pos
	e.Undo = append(e.Undo, ops)
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
