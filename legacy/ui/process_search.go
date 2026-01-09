package ui

import (
	. "github.com/vipmax/edgo/internal/search"
	. "github.com/vipmax/edgo/internal/utils"
	. "github.com/gdamore/tcell"
	"os"
)

func (e *Editor) OnProcessSearch() {
	var end = false
	if e.ProcessPanelSearchPattern == nil { e.ProcessPanelSearchPattern = []rune{} }

	var patternx = len(e.ProcessPanelSearchPattern)
	var isChanged = true

	// loop until escape or enter pressed
	for !end {

		e.DrawProcessPanelSearch(e.ProcessPanelSearchPattern, patternx)
		e.Screen.Show()

		if isChanged && len(e.ProcessPanelSearchPattern) > 0 && len(e.ProcessPanelSearchResults) > 0 {

			var sy, sx = -1, -1
			result := e.ProcessPanelSearchResults[e.ProcessPanelSearchResultIndex]
			sy = result.Line
			sx = result.Position

			if sx != -1 && sy != -1 {
				e.ProcessPanelCursorY = sy
				e.ProcessPanelCursorX = sx
				e.ProcessPanelSelection.Ssx = sx
				e.ProcessPanelSelection.Ssy = sy
				e.ProcessPanelSelection.Sex = sx + len(e.ProcessPanelSearchPattern)
				e.ProcessPanelSelection.Sey = sy
				e.ProcessPanelSelection.IsSelected = true
				e.CleanProcessPanel()
				e.FocusProcessPanel()
				e.DrawProcessPanel()
				e.DrawProcessPanelSearch(e.ProcessPanelSearchPattern, patternx)
				e.Screen.Show()
			} else {
				e.ProcessPanelSelection.CleanSelection()
				e.CleanProcessPanel()
				e.DrawProcessPanel()
				e.DrawProcessPanelSearch(e.ProcessPanelSearchPattern, patternx)
				e.Screen.Show()
			}
			isChanged = false
		}

		switch ev := e.Screen.PollEvent().(type) { // poll and handle event
		case *EventResize:
			e.COLUMNS, e.ROWS = e.Screen.Size()

		case *EventKey:
			isChanged = false
			key := ev.Key()

			if key == KeyCtrlQ { e.Screen.Fini(); os.Exit(1) }

			if key == KeyRune {
				e.ProcessPanelSearchPattern = InsertTo(e.ProcessPanelSearchPattern, patternx, ev.Rune())
				patternx++
				isChanged = true
				e.ProcessPanelSearchResults = Search(e.ProcessContent, string(e.ProcessPanelSearchPattern))
				e.ProcessPanelSearchResultIndex = len(e.ProcessPanelSearchResults) - 1
			}
			if key == KeyBackspace2 && patternx > 0 && len(e.ProcessPanelSearchPattern) > 0 {
				patternx--
				e.ProcessPanelSearchPattern = Remove(e.ProcessPanelSearchPattern, patternx)
				isChanged = true
				e.ProcessPanelSearchResults = Search(e.ProcessContent, string(e.ProcessPanelSearchPattern))
				e.ProcessPanelSearchResultIndex = len(e.ProcessPanelSearchResults) - 1
			}
			if key == KeyLeft && patternx > 0 { patternx-- }
			if key == KeyRight && patternx < len(e.ProcessPanelSearchPattern) { patternx++ }
			if key == KeyDown {
				e.ProcessPanelSearchResultIndex++
				if e.ProcessPanelSearchResultIndex >= len(e.ProcessPanelSearchResults) {
					e.ProcessPanelSearchResultIndex = 0
				}
				isChanged = true
			}
			if key == KeyUp {
				e.ProcessPanelSearchResultIndex--
				if e.ProcessPanelSearchResultIndex < 0 {
					e.ProcessPanelSearchResultIndex = len(e.ProcessPanelSearchResults) - 1
				}
				isChanged = true
			}
			if key == KeyCtrlX {
				e.ProcessPanelSearchPattern = []rune{}
				patternx = 0
			}

			if key == KeyESC || key == KeyCtrlF || key == KeyEnter {
				end = true
			}
		}
	}

	e.IsProcessPanelSearch = false
}
