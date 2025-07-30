package ui

import (
	. "github.com/vipmax/edgo/internal/operations"
	. "github.com/vipmax/edgo/internal/utils"
)

func (e *Editor) OnDown(isPaging bool) {
	e.Update = false
	var numberOfLines int

	if e.ROWS / 3 < 30 { numberOfLines = 30 } else { numberOfLines = e.ROWS / 3 } // adjustment
	if !isPaging { numberOfLines = 1; } // normal scroll

	if len(e.Content) == 0 { return }
	if e.Row + numberOfLines >= len(e.Content) {
		if !isPaging {
			e.Y = e.Row - e.ROWS + 1
			if e.Y < 0 { e.Y = 0; }
			return
		} else {
			numberOfLines = len(e.Content) - e.Row - 1
		}
	}

	e.Row += numberOfLines
	if e.Col > len(e.Content[e.Row]) { e.Col = len(e.Content[e.Row]) } // fit to e.Content
	if e.Row < e.Y { e.Y = e.Row }
	if e.Row >= e.Y + e.ROWS { e.Y = e.Row - e.ROWS + 1  }

	e.Update = true
	clear(e.HighlightElements)
}

func (e *Editor) OnUp(isPaging bool) {
	e.Update = false
	var numberOfLines int

	if e.ROWS / 3 < 30 { numberOfLines = 30 } else { numberOfLines = e.ROWS / 3 } // adjustment
	if !isPaging { numberOfLines = 1; } // normal scroll

	if len(e.Content) == 0 { return }
	if e.Row == 0 { e.Y = 0; return }
	if e.Row - numberOfLines <= 0 {e.Row = 0 } else { e.Row -= numberOfLines }
	if e.Col > len(e.Content[e.Row]) { e.Col = len(e.Content[e.Row]) } // fit to e.Content
	if e.Row < e.Y { e.Y = e.Row }
	if e.Row > e.Y + e.ROWS { e.Y = e.Row - e.ROWS + 1  }

	e.Update = true
	clear(e.HighlightElements)
}

func (e *Editor) OnLeft() {
	e.Update = false
	if len(e.Content) == 0 { return }

	if e.Col > 0 {
		e.Col--
		e.Update = true
	} else if e.Row > 0 {
		e.Row--
		e.Col = len(e.Content[e.Row]) // fit to e.Content
		if e.Row < e.Y { e.Y = e.Row }
		e.Update = true
	}
	clear(e.HighlightElements)
}

func (e *Editor) OnRight() {
	e.Update = false
	if len(e.Content) == 0 { return }

	if e.Col < len(e.Content[e.Row]) {
		e.Col++
		e.Update = true
	} else if e.Row < len(e.Content)-1 {
		e.Row++
		e.Col = 0
		if e.Row > e.Y+ e.ROWS { e.Y++  }
		e.Update = true
	}
	clear(e.HighlightElements)
}

func (e *Editor) GoTop() {
	e.Row = 0; e.Col = 0; e.X = 0; e.Y = 0;
}

func (e *Editor) GoBottom() {
	if len(e.Content) == 0 {
		return
	} else {
		e.Row = len(e.Content)-1; e.Col = 0;
		e.X = 0;
		if e.Row > e.TERMINAL_HEIGHT { e.FocusCenter()}
		e.OnDown(false)
	}
}

func (e *Editor) OnScrollUp() {
	e.Update = false
	if len(e.Content) == 0 { return }
	if e.Y == 0 { return }
	e.Y--
	e.Update = true
}

func (e *Editor) OnScrollDown() {
	e.Update = false
	if len(e.Content) == 0 { return }
	if e.Y+ e.ROWS >= len(e.Content) { return }
	e.Y++
	e.Update = true
}

func (e *Editor) OnEnter() {
	var ops = EditOperation{{Enter, '\n', e.Row, e.Col}}
	tabs := CountTabs(e.Content[e.Row], e.Col)
	spaces := CountSpaces(e.Content[e.Row], e.Col)

	after := e.Content[e.Row][e.Col:]
	before := e.Content[e.Row][:e.Col]
	e.Content[e.Row] = before

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
	e.Col = countToInsert

	newline := append(begining, after...)
	e.Content = InsertTo(e.Content, e.Row, newline)

	contentToString := ConvertContentToString(e.Content)
	e.treeSitterHighlighter.AddCharEdit(&contentToString, e.Row, max(e.Col,0), '\n')

	e.Undo = append(e.Undo, ops)
	e.Focus(); if e.Row- e.Y == e.ROWS { e.OnScrollDown() }
	e.OnCursorChanged()
	if len(e.Redo) > 0 { e.Redo = []EditOperation{} }
	e.Update = true
	e.UpdateLsp(false, ConvertContentToString(e.Content))
	e.IsContentChanged = true
	e.FindTests()
}

func (e *Editor) OnDelete() {

	if e.Selection.IsSelectionNonEmpty() {
		e.Cut(false)
		return
	}

	if e.Col > 0 {
		e.Col--
		e.DeleteCharacter(e.Row, e.Col)
		e.OnCursorChanged()
		e.UpdateLsp(false, ConvertContentToString(e.Content))
	} else if e.Row > 0 { // delete line
		e.Undo = append(e.Undo, EditOperation{{DeleteLine, ' ', e.Row -1, len(e.Content[e.Row-1])}})
		left := e.Content[e.Row][e.Col:]
		e.Content = Remove(e.Content, e.Row)

		e.Row--
		e.Col = len(e.Content[e.Row])
		e.Content[e.Row] = append(e.Content[e.Row], left...)

		code := ConvertContentToString(e.Content)
		e.treeSitterHighlighter.RemoveCharEdit(&code, e.Row, e.Col, '\n')
		e.OnCursorChanged()
		e.UpdateLsp(false, ConvertContentToString(e.Content))
	}

	e.Focus()
	if len(e.Redo) > 0 { e.Redo = []EditOperation{} }
	e.Update = true
	e.IsContentChanged = true
	e.FindTests()
}

func (e *Editor) OnTab() {
	e.Focus()

	selectedLines := e.Selection.GetSelectedLines(e.Content)

	if len(selectedLines) == 0 {
		ch := '\t'
		e.InsertCharacter(e.Row, e.Col, ch)
		e.Col++
		e.OnCursorChanged()
		e.UpdateLsp(false, ConvertContentToString(e.Content))
	} else  {
		var ops = EditOperation{}
		e.Selection.Ssx = 0
		for _, linenumber := range selectedLines {
			e.Row = linenumber
			e.Content[e.Row] = InsertTo(e.Content[e.Row], 0, '\t')
			ops = append(ops, Operation{Insert, '\t', e.Row, 0})
			e.Col = len(e.Content[e.Row])
		}
		e.Selection.Sex = e.Col
		e.Undo = append(e.Undo, ops)
		e.UpdateColors()
		e.UpdateLsp(false, ConvertContentToString(e.Content))
	}

	if len(e.Redo) > 0 { e.Redo = []EditOperation{} }
	e.Update = true
	e.IsContentChanged = true
	e.FindTests()
}

func (e *Editor) OnBackTab() {
	e.Focus()

	selectedLines := e.Selection.GetSelectedLines(e.Content)

	// deleting tabs from beginning
	if len(selectedLines) == 0 {
		if e.Content[e.Row][0] == '\t'  {
			e.DeleteCharacter(e.Row,0)
			e.Col--
		}
	} else {
		e.Selection.Ssx = 0
		for _, linenumber := range selectedLines {
			e.Row = linenumber
			if len(e.Content[e.Row]) > 0 && e.Content[e.Row][0] == '\t'  {
				e.DeleteCharacter(e.Row,0)
				e.Col = len(e.Content[e.Row])
			}
		}
		e.UpdateColors()
	}

	if len(e.Redo) > 0 { e.Redo = []EditOperation{} }
	e.Update = true
	e.UpdateLsp(false, ConvertContentToString(e.Content))
	e.IsContentChanged = true
	e.FindTests()
}

func (e *Editor) OnSwapLinesUp() {
	e.Focus()

	if e.Row == 0 { return }
	var ops = EditOperation{}
	ops = append(ops, Operation{MoveCursor, ' ', e.Row, e.Col})

	line1 := e.Content[e.Row]; line2 := e.Content[e.Row-1]

	for i := len(line1)-1; i >= 0; i-- { ops = append(ops, Operation{Delete, line1[i], e.Row, i}) }
	for i := len(line2)-1; i >= 0; i-- { ops = append(ops, Operation{Delete, line2[i], e.Row -1, i}) }
	for i, ch := range line1 { ops = append(ops, Operation{Insert, ch, e.Row -1, i}) }
	for i, ch := range line2 { ops = append(ops, Operation{Insert, ch, e.Row, i}) }

	e.Content[e.Row] = line2; e.Content[e.Row-1] = line1 // swap
	e.Row--

	e.UpdateColors()
	e.Undo = append(e.Undo, ops)
	e.Selection.CleanSelection()
	e.Update = true
	e.IsContentChanged = true
	e.FindTests()
}

func (e *Editor) OnSwapLinesDown() {
	e.Focus()

	if e.Row+1 == len(e.Content) { return }

	var ops = EditOperation{}
	ops = append(ops, Operation{MoveCursor, ' ', e.Row, e.Col})

	line1 := e.Content[e.Row]; line2 := e.Content[e.Row+1]

	for i := len(line1)-1; i >= 0; i-- { ops = append(ops, Operation{Delete, line1[i], e.Row, i}) }
	for i := len(line2)-1; i >= 0; i-- { ops = append(ops, Operation{Delete, line2[i], e.Row +1, i}) }
	for i, ch := range line1 { ops = append(ops, Operation{Insert, ch, e.Row +1, i}) }
	for i, ch := range line2 { ops = append(ops, Operation{Insert, ch, e.Row, i}) }

	e.Content[e.Row] = line2; e.Content[e.Row+1] = line1 // swap
	e.Row++

	e.UpdateColors()
	e.Undo = append(e.Undo, ops)
	e.Selection.CleanSelection()
	e.Update = true
	e.IsContentChanged = true
	e.FindTests()
}

func (e *Editor) HandleSmartMove(char rune) {
	e.Focus()
	if char == 'f' || char == 'F' {
		nw := FindNextWord(e.Content[e.Row], e.Col+ 1)
		e.Col = nw
		e.Col = Min(e.Col, len(e.Content[e.Row]))
	}
	if char == 'b' || char == 'B' {
		nw := FindPrevWord(e.Content[e.Row], e.Col-1)
		e.Col = nw
	}
}

func (e *Editor) HandleSmartMoveAlac(char int) {
	e.Focus()
	if char == 259 {
		nw := FindNextWord(e.Content[e.Row], e.Col+ 1)
		e.Col = nw
		e.Col = Min(e.Col, len(e.Content[e.Row]))
	}
	if char == 260 {
		nw := FindPrevWord(e.Content[e.Row], e.Col-1)
		e.Col = nw
	}
}

func (e *Editor) HandleSmartMoveDown() {

	var ops = EditOperation{{Enter, '\n', e.Row, e.Col}}

	// moving down, insert new Line, add same amount of tabs
	tabs := CountTabs(e.Content[e.Row], e.Col)
	spaces := CountSpaces(e.Content[e.Row], e.Col)

	countToInsert := tabs
	characterToInsert := '\t'
	if tabs == 0 && spaces != 0 { characterToInsert = ' '; countToInsert = spaces }

	e.Row++; e.Col = 0
	e.Content = InsertTo(e.Content, e.Row, []rune{})
	for i := 0; i < countToInsert; i++ {
		e.Content[e.Row] = InsertTo(e.Content[e.Row], e.Col, characterToInsert)
		ops = append(ops, Operation{Insert, characterToInsert, e.Row, e.Col})
		e.Col++
	}

	e.UpdateColors()
	e.Focus(); e.OnScrollDown()
	e.Undo = append(e.Undo, ops)
	e.Update = true
	e.IsContentChanged = true
}

func (e *Editor) HandleSmartMoveUp() {
	e.Focus()
	// add new Line and shift all lines, add same amount of tabs/spaces
	tabs := CountTabs(e.Content[e.Row], e.Col)
	spaces := CountSpaces(e.Content[e.Row], e.Col)

	countToInsert := tabs
	characterToInsert := '\t'
	if tabs == 0 && spaces != 0 { characterToInsert = ' '; countToInsert = spaces }

	var ops = EditOperation{{Enter, '\n', e.Row, e.Col}}
	e.Content = InsertTo(e.Content, e.Row, []rune{})

	e.Col = 0
	for i := 0; i < countToInsert; i++ {
		e.Content[e.Row] = InsertTo(e.Content[e.Row], e.Col, characterToInsert)
		ops = append(ops, Operation{Insert, characterToInsert, e.Row, e.Col})
		e.Col++
	}

	e.UpdateColors()
	e.Undo = append(e.Undo, ops)
	e.Update = true
	e.IsContentChanged = true
}
