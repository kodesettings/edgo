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

	// adding one character offset after first line
	// this is required due to the newline characters in the string buffer
	if line > 0 { offset += 1 }

	e.code.Insert(offset, []byte(string(ch)))
	e.treeSitterHighlighter.InsertTextEdit(e.code.Value(), offset, len(string(ch)))

	// undo operation type is switched between insert and addchar types
	var op Action

	// change op if this is a new line character
	if ch == '\n' { op = Insert } else { op = AddChar }

	// record the operation on the undo stack. Note that we're creating a new EditOperation
	// and adding all the Operations to it
	e.Undo = append(e.Undo, EditOperation{{op, []byte(string(ch)), offset, CursorMove{line, pos}}})
}

func (e *Editor) InsertString(line, pos int, linestring string) {
	// modify string to remove tabs and space from beginning
	l := RemoveLeadingTabsSpaces(linestring)
	offset := LineOffset(e.code.Value(), line) + pos + 1
	e.code.Insert(offset, []byte(l))
	e.treeSitterHighlighter.InsertTextEdit(e.code.Value(), offset, len(l))

	// record the operation on the undo stack. Note that we're creating a new EditOperation
	// and adding all the Operations to it
	e.Undo = append(e.Undo, EditOperation{{Insert, []byte(l), offset, CursorMove{line, pos}}})
}

func (e *Editor) DeleteCharacter(line, pos int) {
	offset := LineOffset(e.code.Value(), line) + pos
	ch := e.code.At(offset)
	e.code.Remove(offset, offset + 1)
	e.treeSitterHighlighter.RemoveTextEdit(e.code.Value(), offset, len(string(ch)))

	// record the operation on the undo stack. Note that we're creating a new EditOperation
	// and adding all the Operations to it
	e.Undo = append(e.Undo, EditOperation{{Delete, []byte{ch}, offset, CursorMove{line, pos}}})
}

func (e *Editor) ReplaceString(line, from, end int, instext string) {
	offset := LineOffset(e.code.Value(), line)
	begin_idx := offset + from + 1
	end_idx := offset + end + 1

	deltext := make([]byte, end_idx - begin_idx)
	copy(deltext, e.code.Slice(begin_idx, end_idx))

	e.code.Remove(begin_idx, end_idx)
	e.treeSitterHighlighter.RemoveTextEdit(e.code.Value(), begin_idx, len(deltext))

	e.code.Insert(begin_idx, []byte(instext))
	e.treeSitterHighlighter.InsertTextEdit(e.code.Value(), begin_idx, len(instext))

	// record the operation on the undo stack. Note that we're creating a new EditOperation
	// and adding all the Operations to it
	e.Undo = append(e.Undo, EditOperation{
		{Delete, []byte(deltext), begin_idx, CursorMove{line, from}},
		{Insert, []byte(instext), begin_idx, CursorMove{line, from}},
	})
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
}

func (e *Editor) MaybeAddPair(line, pos int, ch rune) (rune, bool) {
	pairMap := map[rune]rune{
		'(': ')', '{': '}', '[': ']',
		'"': '"', '\'': '\'', '`': '`',
	}

	if e.code.Len() == 0 { return rune('\000'), false }

	offset_start := LineOffset(e.code.Value(), line)
	offset_end := LineOffset(e.code.Value(), line + 1)
	line_str := e.code.Slice(offset_start + 1, offset_end)
	current_char := e.code.At(offset_start + 1 + pos - 2)

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
