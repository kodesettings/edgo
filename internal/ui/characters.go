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
	e.OnCursorChanged()

	if len(e.Redo) > 0 { e.Redo = []EditOperation{} }

	e.Update = true
	e.UpdateLsp(false, ConvertLinesToString(e.Lines))
	e.IsContentChanged = true
	e.FindTests()
}

func (e *Editor) InsertCharacter(line, pos int, ch rune) {
	e.Lines[line].Buf = InsertTo(e.Lines[line].Buf, pos, ch)
	e.UpdateLsp(false, ConvertLinesToString(e.Lines))
	e.Undo = append(e.Undo, EditOperation{{Insert, ch, e.Row, e.Col}})

	code := ConvertLinesToString(e.Lines)
	e.treeSitterHighlighter.AddCharEdit(&code, line, pos, ch)
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

		//l := RemoveLeadingTabsSpaces(linestr)
		l := linestr
		//nl := strings.Repeat("\t", tabs) + l
		nl := l
		e.Lines = InsertTo(e.Lines, e.Row, Line{[]rune(nl)})

		ops = append(ops, Operation{Enter, '\n', e.Row, e.Col})
		for _, ch := range nl {
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
	e.UpdateLsp(false, ConvertLinesToString(e.Lines))

	code := ConvertLinesToString(e.Lines)
	e.treeSitterHighlighter.RemoveCharEdit(&code, line, pos, ch)
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
