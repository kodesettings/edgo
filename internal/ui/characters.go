package ui

import (
	. "github.com/vipmax/edgo/internal/operations"
	. "github.com/vipmax/edgo/internal/utils"
)

func (e *Editor) AddChar(ch rune) {
	if len(e.Selection.GetSelectionString(e.Content)) != 0 { e.Cut(false) }

	e.Focus()
	e.InsertCharacter(e.Row, e.Col, ch)
	e.Col++

	e.MaybeAddPair(ch)
	e.OnCursorChanged()

	if len(e.Redo) > 0 { e.Redo = []EditOperation{} }

	e.Update = true
	e.UpdateLsp(false, ConvertContentToString(e.Content))
	e.IsContentChanged = true
	e.FindTests()
}

func (e *Editor) InsertCharacter(line, pos int, ch rune) {
	e.Content[line] = InsertTo(e.Content[line], pos, ch)
	e.UpdateLsp(false, ConvertContentToString(e.Content))
	e.Undo = append(e.Undo, EditOperation{{Insert, ch, e.Row, e.Col}})

	code := ConvertContentToString(e.Content)
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
		e.Content[line] = InsertTo(e.Content[line], pos, ch)
		ops = append(ops, Operation{Insert, ch, line, pos})
		pos++
	}
	e.Col = pos
	e.Undo = append(e.Undo, ops)
}

func (e *Editor) InsertLines(line, pos int, lines []string) {
	var ops = EditOperation{}

	lines[0] = string(e.Content[e.Row][:e.Col]) + RemoveLeadingTabsSpaces(lines[0])

	for _, linestr := range lines {
		e.Col = 0
		if e.Row >= len(e.Content)  { e.Content = append(e.Content, []rune{}) } // if last Line adding empty Line before

		//l := RemoveLeadingTabsSpaces(linestr)
		l := linestr
		//nl := strings.Repeat("\t", tabs) + l
		nl := l
		e.Content = InsertTo(e.Content, e.Row, []rune(nl))

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
	ch := e.Content[line][pos]
	e.Undo = append(e.Undo, EditOperation{
		{MoveCursor, ch, line, pos+1},
		{Delete, ch, line, pos},
	})

	e.Content[line] = Remove(e.Content[line], pos)
	e.UpdateLsp(false, ConvertContentToString(e.Content))

	code := ConvertContentToString(e.Content)
	e.treeSitterHighlighter.RemoveCharEdit(&code, line, pos, ch)
}

func (e *Editor) MaybeAddPair(ch rune) {
	pairMap := map[rune]rune{
		'(': ')', '{': '}', '[': ']',
		'"': '"', '\'': '\'', '`': '`',
	}

	if closeChar, found := pairMap[ch]; found {
		noMoreChars := e.Col >= len(e.Content[e.Row])
		isSpaceNext := e.Col < len(e.Content[e.Row]) && e.Content[e.Row][e.Col] == ' '
		isStringAndClosedBracketNext := closeChar == '"' && e.Col < len(e.Content[e.Row]) && e.Content[e.Row][e.Col] == ')'

		if noMoreChars || isSpaceNext || isStringAndClosedBracketNext {
			e.InsertCharacter(e.Row, e.Col, closeChar)
		}
	}
}
