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

	if len(e.Lines) == 0 { return }
	if e.Row + numberOfLines >= len(e.Lines) {
		if !isPaging {
			e.Y = e.Row - e.ROWS + 1
			if e.Y < 0 { e.Y = 0; }
			return
		} else {
			numberOfLines = len(e.Lines) - e.Row - 1
		}
	}

	e.Row += numberOfLines
	if e.Col > len(e.Lines[e.Row].Buf) { e.Col = len(e.Lines[e.Row].Buf) } // fit to e.Lines
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

	if len(e.Lines) == 0 { return }
	if e.Row == 0 { e.Y = 0; return }
	if e.Row - numberOfLines <= 0 {e.Row = 0 } else { e.Row -= numberOfLines }
	if e.Col > len(e.Lines[e.Row].Buf) { e.Col = len(e.Lines[e.Row].Buf) } // fit to e.Lines
	if e.Row < e.Y { e.Y = e.Row }
	if e.Row > e.Y + e.ROWS { e.Y = e.Row - e.ROWS + 1  }

	e.Update = true
	clear(e.HighlightElements)
}

func (e *Editor) OnLeft() {
	e.Update = false
	if len(e.Lines) == 0 { return }

	if e.Col > 0 {
		e.Col--
		e.Update = true
	} else if e.Row > 0 {
		e.Row--
		e.Col = len(e.Lines[e.Row].Buf) // fit to e.Lines
		if e.Row < e.Y { e.Y = e.Row }
		e.Update = true
	}
	clear(e.HighlightElements)
}

func (e *Editor) OnRight() {
	e.Update = false
	if len(e.Lines) == 0 { return }

	if e.Col < len(e.Lines[e.Row].Buf) {
		e.Col++
		e.Update = true
	} else if e.Row < len(e.Lines)-1 {
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
	if len(e.Lines) == 0 {
		return
	} else {
		e.Row = len(e.Lines)-1; e.Col = 0;
		e.X = 0;
		if e.Row > e.TERMINAL_HEIGHT { e.FocusCenter()}
		e.OnDown(false)
	}
}

func (e *Editor) OnScrollUp() {
	e.Update = false
	if len(e.Lines) == 0 { return }
	if e.Y == 0 { return }
	e.Y--
	e.Update = true
}

func (e *Editor) OnScrollDown() {
	e.Update = false
	if len(e.Lines) == 0 { return }
	if e.Y+ e.ROWS >= len(e.Lines) { return }
	e.Y++
	e.Update = true
}

func (e *Editor) OnEnter() {
	var ops = EditOperation{{Enter, '\n', e.Row, e.Col}}
	tabs := CountTabs(e.Lines[e.Row].Buf, e.Col)
	spaces := CountSpaces(e.Lines[e.Row].Buf, e.Col)

	after := e.Lines[e.Row].Buf[e.Col:]
	before := e.Lines[e.Row].Buf[:e.Col]
	e.Lines[e.Row].Buf = before

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
	e.Lines = InsertTo(e.Lines, e.Row, Line{newline})

	contentToString := ConvertLinesToString(e.Lines)
	e.treeSitterHighlighter.AddCharEdit(&contentToString, e.Row, max(e.Col,0), '\n')

	e.Undo = append(e.Undo, ops)
	e.Focus(); if e.Row- e.Y == e.ROWS { e.OnScrollDown() }
	if len(e.Redo) > 0 { e.Redo = []EditOperation{} }
	e.Update = true
	e.UpdateLsp(false, ConvertLinesToString(e.Lines))
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
		e.UpdateLsp(false, ConvertLinesToString(e.Lines))
	} else if e.Row > 0 { // delete line
		e.Undo = append(e.Undo, EditOperation{{DeleteLine, ' ', e.Row -1, len(e.Lines[e.Row-1].Buf)}})
		left := e.Lines[e.Row].Buf[e.Col:]
		e.Lines = Remove(e.Lines, e.Row)

		e.Row--
		e.Col = len(e.Lines[e.Row].Buf)
		e.Lines[e.Row].Buf = append(e.Lines[e.Row].Buf, left...)

		code := ConvertLinesToString(e.Lines)
		e.treeSitterHighlighter.RemoveCharEdit(&code, e.Row, e.Col, '\n')
		e.UpdateLsp(false, ConvertLinesToString(e.Lines))
	}

	e.Focus()
	if len(e.Redo) > 0 { e.Redo = []EditOperation{} }
	e.Update = true
	e.IsContentChanged = true
	e.FindTests()
}

func (e *Editor) OnTab() {
	e.Focus()

	selectedLines := e.Selection.GetSelectedLines(e.Lines)

	if len(selectedLines) == 0 {
		ch := '\t'
		e.InsertCharacter(e.Row, e.Col, ch)
		e.Col++
		e.UpdateLsp(false, ConvertLinesToString(e.Lines))
	} else  {
		var ops = EditOperation{}
		e.Selection.Ssx = 0
		for _, linenumber := range selectedLines {
			e.Row = linenumber
			e.Lines[e.Row].Buf = InsertTo(e.Lines[e.Row].Buf, 0, '\t')
			ops = append(ops, Operation{Insert, '\t', e.Row, 0})
			e.Col = len(e.Lines[e.Row].Buf)
		}
		e.Selection.Sex = e.Col
		e.Undo = append(e.Undo, ops)
		e.UpdateLsp(false, ConvertLinesToString(e.Lines))
	}

	if len(e.Redo) > 0 { e.Redo = []EditOperation{} }
	e.Update = true
	e.IsContentChanged = true
	e.FindTests()
}

func (e *Editor) OnBackTab() {
	e.Focus()

	selectedLines := e.Selection.GetSelectedLines(e.Lines)

	// deleting tabs from beginning
	if len(selectedLines) == 0 {
		if e.Lines[e.Row].Buf[0] == '\t'  {
			e.DeleteCharacter(e.Row,0)
			e.Col--
		}
	} else {
		e.Selection.Ssx = 0
		for _, linenumber := range selectedLines {
			e.Row = linenumber
			if len(e.Lines[e.Row].Buf) > 0 && e.Lines[e.Row].Buf[0] == '\t'  {
				e.DeleteCharacter(e.Row,0)
				e.Col = len(e.Lines[e.Row].Buf)
			}
		}
	}

	if len(e.Redo) > 0 { e.Redo = []EditOperation{} }
	e.Update = true
	e.UpdateLsp(false, ConvertLinesToString(e.Lines))
	e.IsContentChanged = true
	e.FindTests()
}

func (e *Editor) OnSwapLinesUp() {
	e.Focus()

	if e.Row == 0 { return }
	var ops = EditOperation{}
	ops = append(ops, Operation{MoveCursor, ' ', e.Row, e.Col})

	line1 := e.Lines[e.Row].Buf; line2 := e.Lines[e.Row-1].Buf

	for i := len(line1)-1; i >= 0; i-- { ops = append(ops, Operation{Delete, line1[i], e.Row, i}) }
	for i := len(line2)-1; i >= 0; i-- { ops = append(ops, Operation{Delete, line2[i], e.Row -1, i}) }
	for i, ch := range line1 { ops = append(ops, Operation{Insert, ch, e.Row -1, i}) }
	for i, ch := range line2 { ops = append(ops, Operation{Insert, ch, e.Row, i}) }

	e.Lines[e.Row].Buf = line2; e.Lines[e.Row-1].Buf = line1 // swap
	e.Row--

	e.Undo = append(e.Undo, ops)
	e.Selection.CleanSelection()
	e.Update = true
	e.IsContentChanged = true
	e.FindTests()
}

func (e *Editor) OnSwapLinesDown() {
	e.Focus()

	if e.Row+1 == len(e.Lines) { return }

	var ops = EditOperation{}
	ops = append(ops, Operation{MoveCursor, ' ', e.Row, e.Col})

	line1 := e.Lines[e.Row].Buf; line2 := e.Lines[e.Row+1].Buf

	for i := len(line1)-1; i >= 0; i-- { ops = append(ops, Operation{Delete, line1[i], e.Row, i}) }
	for i := len(line2)-1; i >= 0; i-- { ops = append(ops, Operation{Delete, line2[i], e.Row +1, i}) }
	for i, ch := range line1 { ops = append(ops, Operation{Insert, ch, e.Row +1, i}) }
	for i, ch := range line2 { ops = append(ops, Operation{Insert, ch, e.Row, i}) }

	e.Lines[e.Row].Buf = line2; e.Lines[e.Row+1].Buf = line1 // swap
	e.Row++

	e.Undo = append(e.Undo, ops)
	e.Selection.CleanSelection()
	e.Update = true
	e.IsContentChanged = true
	e.FindTests()
}

func (e *Editor) HandleSmartMove(char rune) {
	e.Focus()
	if char == 'f' || char == 'F' {
		nw := FindNextWord(e.Lines[e.Row].Buf, e.Col+ 1)
		e.Col = nw
		e.Col = Min(e.Col, len(e.Lines[e.Row].Buf))
	}
	if char == 'b' || char == 'B' {
		nw := FindPrevWord(e.Lines[e.Row].Buf, e.Col-1)
		e.Col = nw
	}
}

func (e *Editor) HandleSmartMoveAlac(char int) {
	e.Focus()
	if char == 259 {
		nw := FindNextWord(e.Lines[e.Row].Buf, e.Col+ 1)
		e.Col = nw
		e.Col = Min(e.Col, len(e.Lines[e.Row].Buf))
	}
	if char == 260 {
		nw := FindPrevWord(e.Lines[e.Row].Buf, e.Col-1)
		e.Col = nw
	}
}

func (e *Editor) HandleSmartMoveDown() {

	var ops = EditOperation{{Enter, '\n', e.Row, e.Col}}

	// moving down, insert new Line, add same amount of tabs
	tabs := CountTabs(e.Lines[e.Row].Buf, e.Col)
	spaces := CountSpaces(e.Lines[e.Row].Buf, e.Col)

	countToInsert := tabs
	characterToInsert := '\t'
	if tabs == 0 && spaces != 0 { characterToInsert = ' '; countToInsert = spaces }

	e.Row++; e.Col = 0
	e.Lines = InsertTo(e.Lines, e.Row, Line{[]rune{}})
	for i := 0; i < countToInsert; i++ {
		e.Lines[e.Row].Buf = InsertTo(e.Lines[e.Row].Buf, e.Col, characterToInsert)
		ops = append(ops, Operation{Insert, characterToInsert, e.Row, e.Col})
		e.Col++
	}

	e.Focus(); e.OnScrollDown()
	e.Undo = append(e.Undo, ops)
	e.Update = true
	e.IsContentChanged = true
}

func (e *Editor) HandleSmartMoveUp() {
	e.Focus()
	// add new Line and shift all lines, add same amount of tabs/spaces
	tabs := CountTabs(e.Lines[e.Row].Buf, e.Col)
	spaces := CountSpaces(e.Lines[e.Row].Buf, e.Col)

	countToInsert := tabs
	characterToInsert := '\t'
	if tabs == 0 && spaces != 0 { characterToInsert = ' '; countToInsert = spaces }

	var ops = EditOperation{{Enter, '\n', e.Row, e.Col}}
	e.Lines = InsertTo(e.Lines, e.Row, Line{[]rune{}})

	e.Col = 0
	for i := 0; i < countToInsert; i++ {
		e.Lines[e.Row].Buf = InsertTo(e.Lines[e.Row].Buf, e.Col, characterToInsert)
		ops = append(ops, Operation{Insert, characterToInsert, e.Row, e.Col})
		e.Col++
	}

	e.Undo = append(e.Undo, ops)
	e.Update = true
	e.IsContentChanged = true
}
