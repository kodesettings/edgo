package ui

import (
	. "github.com/vipmax/edgo/internal/search"
	. "github.com/vipmax/edgo/internal/utils"
	"fmt"
	. "github.com/gdamore/tcell"
	"os"
	"path/filepath"
	"time"
)

func (e *Editor) OnSearch() {
	clear(e.HighlightElements)

	e.IsContentSearch = true

	var end = false
	if e.SearchPattern == nil { e.SearchPattern = []rune{} }
	if e.Selection.IsSelectionNonEmpty() {
		e.SearchPattern = []rune(e.Selection.GetSelectionString(e.Content))
		e.SearchResults = Search(e.Content, string(e.SearchPattern))
		e.SearchResultIndex = 0
	}

	var patternx = len(e.SearchPattern)
	var isChanged = true

	// loop until escape or enter pressed
	for !end {

		e.DrawSearch(e.SearchPattern, patternx)
		e.Screen.Show()

		if isChanged && len(e.SearchPattern) > 0 && len(e.SearchResults) > 0 {

			var sy, sx = -1, -1
			e.X = 0

			result := e.SearchResults[e.SearchResultIndex]
			sy = result.Line
			sx = result.Position

			if sx != -1 && sy != -1 {
				e.Row = sy
				e.Col = sx
				e.Focus()
				e.Selection.Ssx = sx
				e.Selection.Ssy = sy
				e.Selection.Sex = sx + len(e.SearchPattern)
				e.Selection.Sey = sy
				e.Selection.IsSelected = true
				e.FocusCenter()
				e.DrawEverything()
				e.DrawSearch(e.SearchPattern, patternx)
				e.Screen.Show()
			} else {
				e.Selection.CleanSelection()
				e.DrawEverything()
				e.DrawSearch(e.SearchPattern, patternx)
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

			if key == KeyRune {
				e.SearchPattern = InsertTo(e.SearchPattern, patternx, ev.Rune())
				patternx++
				isChanged = true
				e.SearchResults = Search(e.Content, string(e.SearchPattern))
				e.SearchResultIndex = 0
			}
			if key == KeyBackspace2 && patternx > 0 && len(e.SearchPattern) > 0 {
				patternx--
				e.SearchPattern = Remove(e.SearchPattern, patternx)
				isChanged = true
				e.SearchResults = Search(e.Content, string(e.SearchPattern))
				e.SearchResultIndex = 0
			}
			if key == KeyLeft && patternx > 0 { patternx-- }
			if key == KeyRight && patternx < len(e.SearchPattern) { patternx++ }
			if key == KeyDown {
				e.SearchResultIndex++
				if e.SearchResultIndex >= len(e.SearchResults) { e.SearchResultIndex = 0 }
				isChanged = true
			}
			if key == KeyUp {
				e.SearchResultIndex--
				if e.SearchResultIndex < 0 { e.SearchResultIndex = len(e.SearchResults) - 1 }
				isChanged = true
			}
			if key == KeyCtrlX {
				e.SearchPattern = []rune{}
				patternx = 0
			}
			if key == KeyCtrlG {
				end = e.OnGlobalSearch()
				e.FocusCenter()
				e.DrawEverything()
				e.DrawSearch(e.SearchPattern, patternx)
				if end { e.CleanContentSearch() }
				e.Screen.Show()
			}
			if key == KeyESC || key == KeyCtrlF {
				end = true
				e.CleanContentSearch()
				e.Screen.Show()
			}
			if key == KeyEnter {
				if len(e.Content) == 0 { // global search if no content and enter
					end = e.OnGlobalSearch()
					e.FocusCenter()
					e.DrawEverything()
					e.DrawSearch(e.SearchPattern, patternx)
					if end { e.CleanContentSearch() }
					e.Screen.Show()
				} else {
					end = true
					e.FocusCenter()
					e.CleanContentSearch()
					e.Screen.Show()
				}
			}
		}
	}

	e.IsContentSearch = false
}

func (e *Editor) OnGlobalSearch() bool {
	clear(e.HighlightElements)

	dir, _ := os.Getwd()

	start := time.Now()
	searchResults, filesProcessedCount, totalRowsProcessed := SearchOnDirParallel(dir, string(e.SearchPattern))
	elapsed := time.Since(start).String()

	e.IsOverlay = true
	defer e.OverlayFalse()

	var end = false
	var isChanged = true

	// loop until escape or enter pressed
	cwd, _ := os.Getwd()

	initialLang := e.treeSitterHighlighter.GetLangStr()

	for !end {
		var resultsCount = 0
		for _, searchResult := range searchResults { resultsCount += len(searchResult.Results) }

		var options = []string{}
		for _, searchResult := range searchResults {
			for _, result := range searchResult.Results {
				relativePath, _ := filepath.Rel(cwd, searchResult.File)

				text := fmt.Sprintf("%d/%d [%d:%d] %s ", len(options)+1, resultsCount, result.Line, result.Position, relativePath)
				options = append(options, text)
			}

		}

		height := MinMany(5, len(options)+1) // depends on min option len or 5 at min or how many rows to the end of e.Screen
		atx := 0 + e.FilesPanelWidth
		aty := 0 // Define the window  position and dimensions
		style := StyleDefault

		var selectionEnd = false
		var selected = 0
		var selectedOffset = 0

		for !selectionEnd {
			if selected < selectedOffset { selectedOffset = selected } // calculate offsets for scrolling completion
			if selected >= selectedOffset+height { selectedOffset = selected - height + 1 }

			if isChanged && resultsCount > 0 {
				isChanged = false
				e.DrawCodePreview(atx, aty, height, options, selectedOffset, selected, style, searchResults,
					fmt.Sprintf("global search: '%s', %d rows found, processed %d rows, %d files, elapsed %s",
						string(e.SearchPattern), resultsCount, totalRowsProcessed, filesProcessedCount, elapsed),
				)

				e.Screen.HideCursor()
				e.Screen.Show()
			}

			switch ev := e.Screen.PollEvent().(type) { // poll and handle event
			case *EventResize:
				e.COLUMNS, e.ROWS = e.Screen.Size()
				e.Screen.Sync()
				e.Screen.Clear()
				e.DrawEverything()
				e.Screen.Show()
				isChanged = true

			case *EventKey:
				key := ev.Key()
				if key == KeyEscape || key == KeyBackspace || key == KeyBackspace2 {
					e.Screen.Clear()
					selectionEnd = true
					end = true
					if e.treeSitterHighlighter.GetLangStr() != initialLang {
						e.treeSitterHighlighter.SetLang(initialLang)
						e.UpdateColors()
					}

					return true
				}

				if key == KeyDown && selected < len(options)-1 { selected++; isChanged = true }
				if key == KeyUp && selected > 0 { selected--; isChanged = true }

				if key == KeyEnter {
					end = true
					file, searchResult, found := e.findSearchGlobalOption(searchResults, selected)
					if found {
						if e.AbsoluteFilePath != file { e.OpenFile(file) }
						searchPattern, _ := ParsePattern(string(e.SearchPattern))
						e.Selection.CleanSelection()
						e.Row = searchResult.Line - 1
						e.Col = searchResult.Position + len(searchPattern)
						e.Selection.Ssy = e.Row
						e.Selection.Sey = e.Row
						e.Selection.Ssx = searchResult.Position
						e.Selection.Sex = searchResult.Position + len(searchPattern)
						e.Selection.IsSelected = true
						e.Focus()

						return true
					}

				}
			}
		}
	}

	if e.treeSitterHighlighter.GetLangStr() != initialLang {
		e.treeSitterHighlighter.SetLang(initialLang)
	}

	return false
}

func (e *Editor) findSearchGlobalOption(searchResults []FileSearchResult, selected int) (string, SearchResult, bool) {
	var i = 0
	for _, searchResult := range searchResults {
		for _, result := range searchResult.Results {
			if i == selected {
				return searchResult.File, result, true
			}
			i++
		}
	}
	return "", SearchResult{}, false
}
