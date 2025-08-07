package ui

import (
	. "github.com/vipmax/edgo/internal/operations"
	. "github.com/vipmax/edgo/internal/utils"
)

func (e *Editor) AddChar(ch rune) {
	if len(e.Selection.GetSelectionString(e.Lines)) != 0 { e.Cut(false) }

	e.Focus()
	e.InsertCharacter(e.Row, e.Col, ch)
	e.Col++

	e.MaybeAddPair(ch)

	if len(e.Redo) > 0 { e.Redo = []EditOperation{} }

	e.Update = true
	e.UpdateLsp(false, ConvertLinesToString(e.Lines))
	e.IsContentChanged = true
	e.FindTests()
}

func (e *Editor) InsertCharacter(line, pos int, ch rune) {
	e.Lines[line].Buf = InsertTo(e.Lines[line].Buf, pos, ch)
	e.Undo = append(e.Undo, EditOperation{{Insert, ch, e.Row, e.Col}})

	e.code = ConvertLinesToString(e.Lines)
	e.treeSitterHighlighter.AddCharEdit(&e.code, line, pos, ch)
}

func (e *Editor) InsertString(line, pos int, linestring string) {
	// Convert the string to insert to a slice of runes
	l := RemoveLeadingTabsSpaces(linestring)
	insertRunes := []rune(l)

	// Record the operation on the undo stack. Note that we're creating a new EditOperation
	// and adding all the Operations to it
	var ops = EditOperation{}
	for _, ch := range insertRunes {
		e.Lines[line].Buf = InsertTo(e.Lines[line].Buf, pos, ch)
		ops = append(ops, Operation{Insert, ch, line, pos})
		pos++
	}
	e.Col = pos
	e.Undo = append(e.Undo, ops)
}

func (e *Editor) InsertLines(line, pos int, lines []string) {
	var ops = EditOperation{}

	lines[0] = string(e.Lines[e.Row].Buf[:e.Col]) + RemoveLeadingTabsSpaces(lines[0])

	for _, linestr := range lines {
		e.Col = 0
		if e.Row >= len(e.Lines)  { e.Lines = append(e.Lines, Line{[]rune{}}) } // if last Line adding empty Line before

		e.Lines = InsertTo(e.Lines, e.Row, Line{[]rune(linestr)})

		ops = append(ops, Operation{Enter, '\n', e.Row, e.Col})
		for _, ch := range linestr {
			ops = append(ops, Operation{Insert, ch, e.Row, e.Col})
			e.Col++
		}
		e.Row++
	}

	e.Row--
	e.Undo = append(e.Undo, ops)
}

func (e *Editor) DeleteCharacter(line, pos int) {
	ch := e.Lines[line].Buf[pos]
	e.Undo = append(e.Undo, EditOperation{
		{MoveCursor, ch, line, pos+1},
		{Delete, ch, line, pos},
	})

	e.Lines[line].Buf = Remove(e.Lines[line].Buf, pos)

	e.code = ConvertLinesToString(e.Lines)
	e.treeSitterHighlighter.RemoveCharEdit(&e.code, line, pos, ch)
}

func (e *Editor) DeleteLine(line, pos int) {
	e.Undo = append(e.Undo, EditOperation{{DeleteLine, ' ', line, len(e.Lines[line].Buf)}})
	left := e.Lines[line].Buf[pos:]
	e.Lines = Remove(e.Lines, line)

	e.Col = len(e.Lines[line-1].Buf)
	e.Lines[line-1].Buf = append(e.Lines[line-1].Buf, left...)

	e.code = ConvertLinesToString(e.Lines)
	e.treeSitterHighlighter.RemoveCharEdit(&e.code, line-1, e.Col, '\n')
}

func (e *Editor) InsertEnter(line, pos int) {
	var ops = EditOperation{{Enter, '\n', line, pos}}
	tabs := CountTabs(e.Lines[line].Buf, pos)
	spaces := CountSpaces(e.Lines[line].Buf, pos)

	after := e.Lines[line].Buf[pos:]
	before := e.Lines[line].Buf[:pos]
	e.Lines[line].Buf = before

	e.Row++
	e.Col = 0

	countToInsert := tabs
	characterToInsert := '\t'
	if tabs == 0 && spaces != 0 { characterToInsert = ' '; countToInsert = spaces }

	begining := []rune{}
	for i := 0; i < countToInsert; i++ {
		begining = append(begining, characterToInsert)
		ops = append(ops, Operation{Insert, characterToInsert, e.Row, e.Col + i})
	}

	e.Undo = append(e.Undo, ops)
	e.Col = countToInsert

	newline := append(begining, after...)
	e.Lines = InsertTo(e.Lines, e.Row, Line{newline})

	contentToString := ConvertLinesToString(e.Lines)
	e.treeSitterHighlighter.AddCharEdit(&contentToString, e.Row, max(e.Col,0), '\n')
}

func (e *Editor) ShiftWithTabsToRight(line, pos int, selectedLines []int) {
	var ops = EditOperation{}
	e.Selection.Ssx = 0

	for _, linenumber := range selectedLines {
		line = linenumber
		e.Row = linenumber
		e.Lines[line].Buf = InsertTo(e.Lines[line].Buf, 0, '\t')
		ops = append(ops, Operation{Insert, '\t', line, 0})
		e.Col = len(e.Lines[line].Buf)
	}

	e.Selection.Sex = pos
	e.Undo = append(e.Undo, ops)
}

func (e *Editor) MaybeAddPair(ch rune) {
	pairMap := map[rune]rune{
		'(': ')', '{': '}', '[': ']',
		'"': '"', '\'': '\'', '`': '`',
	}

	if closeChar, found := pairMap[ch]; found {
		noMoreChars := e.Col >= len(e.Lines[e.Row].Buf)
		isSpaceNext := e.Col < len(e.Lines[e.Row].Buf) && e.Lines[e.Row].Buf[e.Col] == ' '
		isStringAndClosedBracketNext := closeChar == '"' && e.Col < len(e.Lines[e.Row].Buf) && e.Lines[e.Row].Buf[e.Col] == ')'

		if noMoreChars || isSpaceNext || isStringAndClosedBracketNext {
			e.InsertCharacter(e.Row, e.Col, closeChar)
		}
	}
}
