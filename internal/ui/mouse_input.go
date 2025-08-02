package ui

import (
	. "github.com/vipmax/edgo/internal/operations"
	. "github.com/vipmax/edgo/internal/utils"
	"github.com/atotto/clipboard"
	. "github.com/gdamore/tcell"
)

func (e *Editor) HandleMouse(mx int, my int, buttons ButtonMask, modifiers ModMask) {
	_, screenRows := e.Screen.Size()

	// upper play button
	if mx == e.COLUMNS-2 && my == 0 && buttons&Button1 == 1 {
		// do not show if process active
		if e.Process == nil || (e.Process != nil && e.Process.IsStopped()) {
			e.OnProcessRun(true)
		}
		return
	}

	// play button on process panel
	if mx == e.COLUMNS-6 && my == e.ROWS && buttons&Button1 == 1 {
		e.OnProcessStop()
		e.OnProcessRun(false)
		return
	}

	// stop button on process panel
	if mx == e.COLUMNS-4 && my == e.ROWS && buttons&Button1 == 1 {
		e.OnProcessStop()
		return
	}

	// close button on process panel
	if mx == e.COLUMNS-2 && my == e.ROWS && buttons&Button1 == 1 {
		e.OnProcessStop()
		e.ROWS = screenRows
		e.ProcessPanelHeight = 0
		e.ProcessContent = []Line{Line{[]rune{}}}
		return
	}

	// detect process panel drag start event
	if !e.IsProcessPanelMoving && buttons&Button1 == 1 &&
		my == e.ROWS && e.ProcessPanelHeight > 0 &&
		len(e.ProcessPanelSelection.GetSelectedLines(e.ProcessContent)) == 0 {

		e.ROWS = my
		e.ProcessPanelHeight = screenRows - e.ROWS
		e.IsProcessPanelMoving = true
		e.Update = true
		return
	}
	// detect process panel dragging event
	if e.IsProcessPanelMoving && buttons&Button1 == 1 && screenRows >= my {
		e.ROWS = my
		e.ProcessPanelHeight = screenRows - e.ROWS
		e.Update = true
		return
	}
	// detect process panel drag stop event
	if e.IsProcessPanelMoving && buttons&Button1 == 0 {
		e.IsProcessPanelMoving = false
		return
	}

	if my >= e.ROWS {
		if mx < 2 || mx > e.COLUMNS { return }
		// in process panel
		e.IsProcessPanelFocused = true
		if buttons&WheelDown != 0 && e.ProcessPanelScroll <= len(e.ProcessContent)-e.ProcessPanelHeight {
			e.ProcessPanelScroll++
		}
		if buttons&WheelUp != 0 && e.ProcessPanelScroll > 0 {
			e.ProcessPanelScroll--
		}

		if buttons&Button1 == 1 {
			if mx < e.ProcessPanelSpacing { return }
			e.ProcessPanelCursorX = mx + e.ProcessPanelHScroll - e.ProcessPanelSpacing
			e.ProcessPanelCursorY = my + e.ProcessPanelScroll - e.ROWS - 1

			if e.ProcessPanelCursorY < 0 { e.ProcessPanelCursorY = 0 }
			// fit cursor
			if e.ProcessPanelCursorY >= len(e.ProcessContent) {
				e.ProcessPanelCursorY = len(e.ProcessContent) - 1
			}
			if e.ProcessPanelCursorY < len(e.ProcessContent) && e.ProcessPanelCursorX > len(e.ProcessContent[e.ProcessPanelCursorY].Buf) {
				e.ProcessPanelCursorX = len(e.ProcessContent[e.ProcessPanelCursorY].Buf)
			}

			if e.ProcessPanelSelection.Ssx < 0 {
				e.ProcessPanelSelection.Ssx, e.ProcessPanelSelection.Ssy =
					e.ProcessPanelCursorX, e.ProcessPanelCursorY
			}
			if e.ProcessPanelSelection.Ssx >= 0 {
				e.ProcessPanelSelection.Sex, e.ProcessPanelSelection.Sey =
					e.ProcessPanelCursorX, e.ProcessPanelCursorY
			}
			return
		}

		if buttons&Button1 == 0 {
			if e.ProcessPanelSelection.IsSelectionNonEmpty() {
				selectionString := e.ProcessPanelSelection.GetSelectionString(e.ProcessContent)
				clipboard.WriteAll(selectionString)
			}

			e.ProcessPanelSelection.CleanSelection()
		}

		return
	}

	e.IsProcessPanelFocused = false

	if !e.IsFilesPanelMoving && buttons&Button1 == 1 &&
		(mx == e.FilesPanelWidth-2 || mx == e.FilesPanelWidth-1) &&
		my < e.ROWS && len(e.Selection.GetSelectedLines(e.Lines)) == 0 {
		e.IsFilesPanelMoving = true
		return
	}

	if e.IsFilesPanelMoving && buttons&Button1 == 1 {
		e.FilesPanelWidth = mx
		return
	}
	if e.IsFilesPanelMoving && buttons&Button1 == 0 {
		e.IsFilesPanelMoving = false
		return
	}
	if mx < e.FilesPanelWidth-3 && buttons&Button1 == 0 && !e.Dap.IsStarted {
		e.OnFilesTree(true)
		return
	}

	if e.Filename == "" { return }

	if buttons&Button1 == 1 && mx == e.COLUMNS-2 { // test button
		line := my + e.Y
		if _, found := e.Tests[line]; found {
			e.RunTest(e.Tests[line])
			return
		}

	}

	mx -= e.LINES_WIDTH + e.FilesPanelWidth

	if mx < 0 { return }
	if my > e.ROWS { return }
	if e.Lines == nil { return }

	// if click with control or option, lookup for definition or references
	if buttons&Button1 == 1 && (modifiers&ModAlt != 0 || modifiers&ModCtrl != 0) {
		e.Row = my + e.Y
		if e.Row > len(e.Lines) - 1 { e.Row = len(e.Lines) - 1 } // fit cursor to e.Lines

		e.Col = e.FindCursorXPosition(mx)

		if len(e.Selection.GetSelectedLines(e.Lines)) > 0 { // if text selected
			e.Selection.Sey = e.Row
			e.Selection.Sex = e.Col
			return
		}
		if modifiers&ModAlt != 0 { e.OnReferences() }
		if modifiers&ModCtrl != 0 { e.OnDefinition() }
		return
	}
	prevRow := e.Row

	if e.Selection.IsSelected && buttons&Button1 == 1 {
		e.Update = true
		e.Row = my + e.Y
		if e.Row > len(e.Lines) - 1 { e.Row = len(e.Lines) - 1 } // fit cursor to e.Lines

		xPosition := e.FindCursorXPosition(mx)

		isTripleClick := e.Selection.IsUnderSelection(xPosition, e.Row) &&
			len(e.Selection.GetSelectedLines(e.Lines)) == 1

		if isTripleClick {
			e.Row = my + e.Y
			e.Col = xPosition
			if e.Row > len(e.Lines) - 1 { e.Row = len(e.Lines) - 1 } // fit cursor to e.Lines
			if e.Col > len(e.Lines[e.Row].Buf) { e.Col = len(e.Lines[e.Row].Buf) }

			e.Selection.Ssx = 0
			e.Selection.Sex = len(e.Lines[e.Row].Buf)
			e.Selection.Ssy = e.Row
			e.Selection.Sey = e.Row

			return
		} else {
			e.Selection.CleanSelection()
		}
	}

	if buttons&WheelDown != 0 {
		e.OnScrollDown()
		return
	}
	if buttons&WheelUp != 0 {
		e.OnScrollUp()
		return
	}
	if buttons&Button1 == 0 && e.Selection.Ssx == -1 {
		e.Update = false
		return
	}

	if buttons&Button1 == 1 {
		e.Update = true

		e.Row = my + e.Y
		if e.Row > len(e.Lines) - 1 { e.Row = len(e.Lines) - 1 } // fit cursor to e.Lines

		xPosition := e.FindCursorXPosition(mx)

		if prevRow == e.Row && e.Col == xPosition && len(e.Selection.GetSelectedLines(e.Lines)) == 0 {
			// double click
			lastChar := len(e.Lines[e.Row].Buf) == e.Col
			if lastChar {
				e.Selection.Ssx, e.Selection.Ssy = e.Col, e.Row
				e.Selection.Sex, e.Selection.Sey = e.Col, e.Row
				return
			}
			prw := FindPrevWord(e.Lines[e.Row].Buf, e.Col)
			nxw := FindNextWord(e.Lines[e.Row].Buf, e.Col)
			e.Selection.Ssx, e.Selection.Ssy = prw, e.Row
			e.Selection.Sex, e.Selection.Sey = nxw, e.Row
			e.Col = nxw
			return
		}

		e.Col = xPosition
		e.OnCursorChanged()
		e.CursorHistory = append(e.CursorHistory,
			CursorMove{e.AbsoluteFilePath, e.Row, e.Col, e.Y, e.X},
		)

		if e.Col < 0 { e.Col = 0 }
		if e.Selection.Ssx < 0 { e.Selection.Ssx, e.Selection.Ssy = e.Col, e.Row }
		if e.Selection.Ssx >= 0 { e.Selection.Sex, e.Selection.Sey = e.Col, e.Row }
	}

	if buttons&Button1 == 0 {
		if e.Selection.Ssx != -1 && e.Selection.Sex != -1 { e.Selection.IsSelected = true }
		e.Update = false
	}
	return
}
