package ui

import (
	. "github.com/vipmax/edgo/internal/operations"
)

func (e *Editor) Focus() {
	if e.Row > e.Y + e.ROWS { e.Y = e.Row + e.ROWS }
	if e.Row < e.Y { e.Y = e.Row }
}

func (e *Editor) FocusCenter() {
	e.Screen.Show()
	if e.Row > e.Y + e.ROWS {
		e.Y = e.Row + e.ROWS
	}
	if e.Row < e.Y {
		e.Y = e.Row
	}

	e.Y -= e.ROWS/2
	if e.Y < 0 { e.Y = 0 }

	centerRow := e.ROWS / 2
	// Update the cursor row to the center row if necessary
	if e.Row - e.Y > centerRow {
		e.Y += e.Row - e.Y - centerRow
	}
}

func (e *Editor) FocusProcessPanel() {
	e.Screen.Show()
	if e.ProcessPanelCursorY > e.ProcessPanelScroll + e.ROWS {
		e.ProcessPanelScroll = e.ProcessPanelCursorY + e.ROWS
	}
	if e.ProcessPanelCursorY < e.ProcessPanelScroll {
		e.ProcessPanelScroll = e.ProcessPanelCursorY
	}

	e.ProcessPanelScroll -= e.ProcessPanelHeight/2
	if e.ProcessPanelScroll < 0 { e.ProcessPanelScroll = 0 }

	centerRow := e.ProcessPanelHeight / 2
	// Update the cursor row to the center row if necessary
	if e.ProcessPanelCursorY - e.ProcessPanelScroll > centerRow {
		e.ProcessPanelScroll += e.ProcessPanelCursorY - e.ProcessPanelScroll - centerRow
	}
}

func (e *Editor) OnCursorBackUndo() {
	if len(e.CursorHistoryUndo) == 0 { return }

	lastCursor := e.CursorHistoryUndo[len(e.CursorHistoryUndo)-1]
	e.CursorHistoryUndo = e.CursorHistoryUndo[:len(e.CursorHistoryUndo)-1]


	if lastCursor.Filename != e.Filename {
		e.OpenFile(lastCursor.Filename)
	}

	e.Row = lastCursor.Row
	e.Col = lastCursor.Col
	e.Y = lastCursor.Y
	e.X = lastCursor.X
	e.Focus()
	e.OnCursorChanged()
	e.CursorHistory = append(e.CursorHistory, lastCursor)
}
func (e *Editor) OnCursorBack() {
	if len(e.CursorHistory) == 0 { return }

	lastCursor := e.CursorHistory[len(e.CursorHistory)-1]
	e.CursorHistory = e.CursorHistory[:len(e.CursorHistory)-1]
	e.CursorHistoryUndo = append(e.CursorHistoryUndo,
		 CursorMove{e.AbsoluteFilePath, e.Row, e.Col, e.Y, e.X},
	)

	if lastCursor.Filename != e.Filename {
		e.OpenFile(lastCursor.Filename)
	}

	e.Row = lastCursor.Row
	e.Col = lastCursor.Col
	e.Y = lastCursor.Y
	e.X = lastCursor.X
	e.Focus()
	e.OnCursorChanged()
}
