package ui

import (
	. "github.com/vipmax/edgo/internal/operations"
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
	if e.Selection.IsSelectionNonEmpty() {
		e.Cut(false)
		// TODO: remove extra line
	}

	e.InsertCharacter(e.Row, e.Col, '\n')
	e.Focus();

	if e.Row - e.Y == e.ROWS { e.OnScrollDown() }
	if len(e.Redo) > 0 { e.Redo = []EditOperation{} }

	e.Update = true
	e.UpdateLsp(false, string(e.code.Value()))
	e.IsContentChanged = true
	e.FindTests()
}

func (e *Editor) OnDelete() {

	if e.Selection.IsSelectionNonEmpty() {
		e.Cut(false)
		return
	}

	e.DeleteCharacter(e.Row, e.Col)
	if e.Col > 0 { e.Col-- } else if e.Row > 0 { e.Row-- }

	e.Focus()
	if len(e.Redo) > 0 { e.Redo = []EditOperation{} }
	e.Update = true
	e.UpdateLsp(false, string(e.code.Value()))
	e.IsContentChanged = true
	e.FindTests()
}

func (e *Editor) OnTab() {
	e.Focus()

	selectedLines := e.Selection.GetSelectedLines(e.code.Value())

	if len(selectedLines) == 0 {
		ch := '\t'
		e.InsertCharacter(e.Row, e.Col, ch)
		e.Col++
		e.UpdateLsp(false, string(e.code.Value()))
	} else  {
		e.ShiftWithTabsToRight(e.Row, e.Col, selectedLines)
		e.UpdateLsp(false, string(e.code.Value()))
	}

	if len(e.Redo) > 0 { e.Redo = []EditOperation{} }
	e.Update = true
	e.IsContentChanged = true
	e.FindTests()
}

func (e *Editor) OnBackTab() {
	e.Focus()

	selectedLines := e.Selection.GetSelectedLines(e.code.Value())

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
	e.UpdateLsp(false, string(e.code.Value()))
	e.IsContentChanged = true
	e.FindTests()
}
