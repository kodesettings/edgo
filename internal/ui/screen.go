package ui

import (
	. "github.com/vipmax/edgo/internal/highlighter"
	. "github.com/vipmax/edgo/internal/search"
	. "github.com/vipmax/edgo/internal/utils"
	"fmt"
	. "github.com/gdamore/tcell"
	"strings"
	"time"
	"unicode/utf8"
	"log/slog"
)

func (e *Editor) DrawEverything() {
	e.mu.Lock()  // Lock the mutex to ensure exclusive access to the method
	defer e.mu.Unlock() // Ensure the mutex is unlocked when the method exits
	e.Screen.Clear() // Clearing the terminal screen before opening anything

	var nlines int = 0
	if e.code != nil { nlines := e.code.Count(0, e.code.Len(), []byte{'\n'})
		e.Lines = GetLinesArrayFromData(e.code.Value(), nlines) }
	if len(e.Lines) == 0 { e.DrawLogo(); return }
	if e.Update == false { return }

	slog.Info("get", "datalen", e.code.Len(), "lines", len(e.Lines), "nlines", nlines)

	countTabsTo := CountTabsTo(e.Lines[e.Row].Buf, e.Col)
	tabcor := countTabsTo * (e.langTabWidth - 1)

	if e.Col < e.X { e.X = e.Col }
	if e.Col + e.LINES_WIDTH + tabcor >= e.X + e.COLUMNS  {
		e.X = e.Col - e.COLUMNS + 1 + e.LINES_WIDTH + tabcor
	}

	/*
		1. Getting bytes ranges and colors from tree-sitter only for visible text
		2. iterating over characters and increment bytesCounter
		3. find range that matches bytesCounter
		4. draw each cell (row,col, char, color)
	*/

	start := time.Now()
	coloredByteRanges := e.treeSitterHighlighter.ColorRanges(e.Y, e.Y+e.TERMINAL_HEIGHT, e.code.Value())

	bytesCounter := 0

	if e.Y > 0 { //  if scrolling needs to recalculate bytesCounter offset
		newlineCount := 0
		for _, c := range e.code.Value() {
			bytesCounter += utf8.RuneLen(rune(c))
			if c == '\n' { newlineCount++ }
			if newlineCount == e.Y { break }
		}
	}

	for row := 0; row < e.TERMINAL_HEIGHT; row++ {
		ry := row + e.Y // index to get right row in characters buffer by scrolling offset Y
		if row >= len(e.Lines) || ry >= len(e.Lines) { break }
		e.DrawLineNumber(ry, row)

		if _, found := e.Tests[ry]; found { e.DrawTest(ry, row) }

		tabsOffset := 0

		for col := 0; true; col++ {
			cx := col + e.X // index to get right column in characters buffer by scrolling offset x

			if cx < 0 { break }
			if col >= len(e.Lines[ry].Buf) { break }
			ch := e.Lines[ry].Buf[col]

			isOutside := col-e.X+e.LINES_WIDTH+tabsOffset > e.COLUMNS
			if isOutside || e.X > col { bytesCounter += utf8.RuneLen(ch); continue }

			style := StyleDefault

			for _, i := range coloredByteRanges {
				if i.StartByte <= bytesCounter && bytesCounter < i.EndByte {
					style = StyleDefault.Foreground(Color(i.Color))
					break
				}
			}

			if e.Selection.IsUnderSelection(col, ry) {
				style = style.Background(Color(SelectionColor))
			}

			if ch == '\t' && e.X == 0 { // draw big cursor for tab
				if ry == e.Row && cx == e.Col {
					style = StyleDefault.Background(Color(AccentColor))
				}
				for i := 0; i < e.langTabWidth; i++ {
					x := col - e.X + e.LINES_WIDTH + tabsOffset
					e.Screen.SetContent(x, row, ' ', nil, style)
					if i != e.langTabWidth-1 { tabsOffset++ }
				}
			} else {
				x := col - e.X + e.LINES_WIDTH + tabsOffset
				e.Screen.SetContent(x, row, ch, nil, style)
			}
			bytesCounter += utf8.RuneLen(ch)
		}

		bytesCounter += 1 // for '/n'
	}

	e.DrawDiagnostic()

	ttr := time.Since(start).String()
	var changes = ""
	if e.IsContentChanged { changes = "*" }
	status := fmt.Sprintf(" %s %s %d %d %s%s ", ttr, e.Lang, e.Row+1, e.Col+1, e.Filename, changes)
	e.DrawStatus(status)

	// if tab under cursor, hide cursor because it has already drawn
	if e.Row < len(e.Lines) && e.Col < len(e.Lines[e.Row].Buf) && e.Lines[e.Row].Buf[e.Col] == '\t' {
		e.Screen.HideCursor()
	} else {
		tabs := CountTabsTo(e.Lines[e.Row].Buf, e.Col) * (e.langTabWidth - 1)
		e.Screen.ShowCursor(e.Col-e.X+e.LINES_WIDTH+tabs, e.Row-e.Y) // show cursor
		if e.X != 0 {
			e.Screen.ShowCursor(e.Col-e.X+e.LINES_WIDTH, e.Row-e.Y) // show cursor
		}
	}

	if e.Row-e.Y >= e.ROWS { e.Screen.HideCursor() }

	e.DrawProcessPanel()

	if e.IsContentSearch {
		e.DrawSearch(e.SearchPattern, len(e.SearchPattern))
	}
}

func (e *Editor) CleanProcessPanel() {
	for j := e.ROWS; j < e.TERMINAL_HEIGHT; j++ {
		for i := 0; i < e.COLUMNS; i++ {
			e.Screen.SetContent(i, j, ' ', nil, StyleDefault)
		}
	}
}

func (e *Editor) DrawProcessPanel() {
	if e.langConf.Cmd != "" && (e.Process == nil || e.Process != nil && e.Process.IsStopped()) {
		e.Screen.SetContent(e.COLUMNS-2, 0, '▶', nil, StyleDefault.Foreground(Color(HighlighterGlobal.GetRunButtonStyle())))
	}

	for i := 0; i < e.COLUMNS-7; i++ {
		e.Screen.SetContent(i, e.ROWS, '─', nil, SeparatorStyle)
	}

	e.Screen.SetContent(e.COLUMNS-7, e.ROWS, ' ', nil, StyleDefault)

	e.Screen.SetContent(e.COLUMNS-6, e.ROWS, '▶', nil, StyleDefault.Foreground(Color(HighlighterGlobal.GetRunButtonStyle())))

	e.Screen.SetContent(e.COLUMNS-5, e.ROWS, ' ', nil, StyleDefault)

	if e.Process != nil && e.Process.IsStopped() {
		e.Screen.SetContent(e.COLUMNS-4, e.ROWS, ' ', nil, StyleDefault)
	} else {
		e.Screen.SetContent(e.COLUMNS-4, e.ROWS, '■', nil, StyleDefault.Foreground(Color(AccentColor)))
	}
	e.Screen.SetContent(e.COLUMNS-3, e.ROWS, ' ', nil, StyleDefault)
	e.Screen.SetContent(e.COLUMNS-2, e.ROWS, 'x', nil, StyleDefault)

	screenCols, screenRows := e.Screen.Size()

	if e.ProcessPanelCursorX < e.ProcessPanelHScroll {
		e.ProcessPanelHScroll = e.ProcessPanelCursorX
	}
	if e.ProcessPanelCursorX >= e.ProcessPanelHScroll+screenCols-e.ProcessPanelSpacing {
		e.ProcessPanelHScroll = e.ProcessPanelCursorX - screenCols + e.ProcessPanelSpacing + 1
	}

	for index := 0; index < len(e.ProcessContent); index++ {
		if index+e.ProcessPanelScroll > len(e.ProcessContent)-1 {
			break
		}
		line := e.ProcessContent[index+e.ProcessPanelScroll].Buf
		y := e.ROWS + index + 1
		if y > screenRows { break }

		for col := 0; col <= e.COLUMNS; col++ {
			cx := col + e.ProcessPanelHScroll // index to get right column in characters buffer by scrolling offset x
			if cx >= len(line) { break }
			ch := line[cx]

			style := StyleDefault
			style = style.Foreground(Color(AccentColor3))
			if e.ProcessPanelSelection.IsUnderSelection(cx, index+e.ProcessPanelScroll) {
				style = style.Background(Color(SelectionColor))
			}

			e.Screen.SetContent(col+e.ProcessPanelSpacing, y, ch, nil, style)
		}

		for i := len(line); i < e.COLUMNS; i++ {
			e.Screen.SetContent(i+e.ProcessPanelSpacing, y, ' ', nil, StyleDefault)
		}
	}

	if e.IsProcessPanelFocused {
		if e.ProcessPanelCursorY-e.ProcessPanelScroll+e.ROWS+1 <= e.ROWS {
			e.Screen.HideCursor()
		} else {
			e.Screen.ShowCursor(e.ProcessPanelCursorX+e.ProcessPanelSpacing-e.ProcessPanelHScroll, e.ProcessPanelCursorY-e.ProcessPanelScroll+e.ROWS+1)
		}
	} else {
		//e.Screen.HideCursor()
	}
}

func (e *Editor) DrawDiagnostic() {
	if e.Lang == "" { return }
	lsp := e.lsp2lang[e.Lang]
	if !lsp.IsReady { return }

	maybeDiagnostic, found := lsp.GetDiagnostic("file://" + e.AbsoluteFilePath)

	if found {
		style := StyleDefault.Foreground(Color(AccentColor))

		for _, diagnostic := range maybeDiagnostic.Diagnostics {
			dline := int(diagnostic.Range.Start.Line)
			if dline >= len(e.Lines) { continue } // sometimes it out of e.Content
			if dline-e.Y > e.ROWS { continue } // sometimes it out of e.Content

			tabs := CountTabs(e.Lines[dline].Buf, len(e.Lines[dline].Buf))
			var shifty = 0
			errorMessage := "error: " + diagnostic.Message
			errorMessage = PadLeft(errorMessage, e.COLUMNS-len(e.Lines[dline].Buf)-tabs*e.langTabWidth-5-e.LINES_WIDTH)

			// iterate over message characters and draw it
			for i, m := range errorMessage {
				ypos := dline - e.Y
				if ypos < 0 || ypos >= len(e.Lines) { break }

				tabs = CountTabs(e.Lines[dline].Buf, len(e.Lines[dline].Buf))
				xpos := i + e.LINES_WIDTH + len(e.Lines[dline+shifty].Buf) + tabs*e.langTabWidth + 5

				e.Screen.SetContent(xpos, ypos, m, nil, style)
			}
		}
	}
}

func (e *Editor) DrawLineNumber(brw int, row int) {
	var style = StyleDefault.Foreground(247)
	if brw == e.Row { style = StyleDefault }

	lineNumber := CenterNumber(brw+1, e.LINES_WIDTH)
	for index, char := range lineNumber {
		e.Screen.SetContent(index, row, char, nil, style)
	}
}

func (e *Editor) DrawStatus(text string) {
	var style = StyleDefault.Foreground(247)
	e.DrawText(e.ROWS-1, e.COLUMNS-len(text), text, style)
}

func (e *Editor) DrawText(row, col int, text string, style Style) {
	e.Screen.SetContent(col-1, row, ' ', nil, style)
	for _, ch := range []rune(text) {
		if col > e.COLUMNS { break }
		e.Screen.SetContent(col, row, ch, nil, style)
		col++
	}
}

func (e *Editor) DrawTextScreenReport(text string, x, y int) {
	for i, ch := range text {
		e.Screen.SetContent(x+i, y, ch, nil, StyleDefault)
	}
}

func (e *Editor) DrawErrors(atx int, width int, aty int, height int, options []string,
	selectedOffset int, selected int, style Style) int {

	var shifty = 0
	for row := 0; row < aty+height; row++ {
		if row >= len(options) || row >= height { break }
		var option = options[row+selectedOffset]

		isRowSelected := selected == row+selectedOffset
		if isRowSelected {
			style = style.Background(Color(AccentColor))
		} else {
			style = StyleDefault.Background(Color(OverlayColor))
		}

		shiftx := 0
		runes := []rune(option)
		for j := 0; j < len(runes); j++ {
			ch := runes[j]
			nextWord := FindNextWord(runes, j)
			if shiftx == 0 {
				e.Screen.SetContent(atx, row+aty+shifty, ' ', nil, style)
			}
			if shiftx+atx+(nextWord-j) >= e.COLUMNS {
				for k := shiftx; k <= e.COLUMNS; k++ { // Fill the remaining space
					e.Screen.SetContent(k+atx, row+aty+shifty, ' ', nil, style)
				}
				shifty++
				shiftx = 0
			}
			e.Screen.SetContent(atx+shiftx, row+aty+shifty, ch, nil, style)
			shiftx++
		}

		for col := shiftx; col < e.COLUMNS; col++ { // Fill the remaining space
			e.Screen.SetContent(col+atx, row+aty+shifty, ' ', nil, style)
		}
	}

	for col := 0; col < e.COLUMNS; col++ { // Fill empty Line below
		e.Screen.SetContent(col+atx, height+aty+shifty,
			//' ', nil, StyleDefault.Background(Color(OverlayColor)))
			'─', nil, SeparatorStyle)
	}

	return shifty
}

func (e *Editor) DrawSearch(pattern []rune, patternx int) {
	var prefix = []rune("search: ")

	for i := 0; i < len(prefix); i++ {
		e.Screen.SetContent(i+e.LINES_WIDTH, e.ROWS-1, prefix[i], nil, StyleDefault)
	}

	e.Screen.SetContent(len(prefix)+e.LINES_WIDTH, e.ROWS-1, ' ', nil, StyleDefault)

	for i := 0; i < len(pattern); i++ {
		e.Screen.SetContent(len(prefix)+i+e.LINES_WIDTH, e.ROWS-1, pattern[i], nil, StyleDefault)
	}

	e.Screen.ShowCursor(len(prefix)+patternx+e.LINES_WIDTH, e.ROWS-1)

	for i := len(prefix) + len(pattern) + e.LINES_WIDTH; i < e.COLUMNS; i++ {
		e.Screen.SetContent(i, e.ROWS-1, ' ', nil, StyleDefault)
	}

	if len(e.SearchResults) > 0 {
		status := fmt.Sprintf("  %d/%d", e.SearchResultIndex+1, len(e.SearchResults))

		for i := 0; i < len(status); i++ {
			e.Screen.SetContent(e.LINES_WIDTH+len(prefix)+len(pattern)+i, e.ROWS-1,
				rune(status[i]), nil, StyleDefault)
		}
	}

	e.Screen.ShowCursor(len(prefix)+patternx+e.LINES_WIDTH, e.ROWS-1)
}

func (e *Editor) CleanContentSearch() {
	for i := e.LINES_WIDTH; i < e.COLUMNS; i++ {
		e.Screen.SetContent(i, e.ROWS-1, ' ', nil, StyleDefault)
	}

	if len(e.Lines) == 0 {
		e.Screen.HideCursor()
	}
}

func (e *Editor) DrawCodePreview(atx int, aty int, height int, options []string,
	selectedOffset int, selected int,
	style Style, searchResults []FileSearchResult, status string) {

	// draw options
	for row := aty; row < aty+height; row++ {
		if row >= len(options) || row >= height { break }

		var option = options[row+selectedOffset]

		isRowSelected := selected == row+selectedOffset
		if isRowSelected {
			style = style.Background(Color(AccentColor))
		} else {
			style = StyleDefault.Background(Color(OverlayColor))
		}

		for i, ch := range option {
			e.Screen.SetContent(atx+i, row, ch, nil, style)
		}

		for i := atx + len(option); i < e.COLUMNS; i++ {
			e.Screen.SetContent(i, row, ' ', nil, style)
		}
	}

	for i := atx; i < e.COLUMNS; i++ {
		e.Screen.SetContent(i, height, ' ', nil, StyleDefault)
	}

	file, searchResult, found := e.findSearchGlobalOption(searchResults, selected)
	if found {
		rowsToShow := e.ROWS - height
		previewContent := e.ReadContent(file, searchResult.Line-rowsToShow/2, searchResult.Line+rowsToShow/2)
		lang := DetectLang(file)
		if e.treeSitterHighlighter.GetLangStr() != lang {
			e.treeSitterHighlighter.SetLang(lang)
		}

		// clear
		for j := height + 1; j < e.ROWS; j++ {
			for i := atx; i < e.COLUMNS; i++ {
				e.Screen.SetContent(i, j, ' ', nil, StyleDefault)
			}
		}

		linenumber := searchResult.Line - rowsToShow/2
		if linenumber < 0 { linenumber = 0 }

		// draw preview
		for row := 0; row < len(previewContent); row++ {
			y := row + height + 1
			if y >= e.ROWS { break }
			var shiftTabs = 0

			var lineNumberStyle = StyleDefault.Foreground(ColorDimGray)
			for index, char := range CenterNumber(linenumber+1, e.LINES_WIDTH) {
				e.Screen.SetContent(index, y, char, nil, lineNumberStyle)
			}

			for col := 0; col < len(previewContent[row].Buf); col++ {

				chstyle := StyleDefault

				if linenumber == searchResult.Line-1 {
					chstyle = chstyle.Background(Color(SelectionColor))
				}

				if previewContent[row].Buf[col] == '\n' { continue }
				if previewContent[row].Buf[col] == '\t' {
					for i := 0; i < e.langTabWidth; i++ {
						e.Screen.SetContent(atx+e.LINES_WIDTH+col+shiftTabs, y, ' ', nil, chstyle)
						if i != e.langTabWidth-1 { shiftTabs++ }
					}
				} else {
					e.Screen.SetContent(atx+e.LINES_WIDTH+col+shiftTabs, y, previewContent[row].Buf[col], nil, chstyle)
				}

				if atx+e.LINES_WIDTH+col+shiftTabs >= e.COLUMNS { break }
			}

			for i := atx + len(previewContent[row].Buf) + e.LINES_WIDTH + shiftTabs; i < e.COLUMNS; i++ {
				e.Screen.SetContent(i, y, ' ', nil, StyleDefault)
			}

			linenumber++
		}

		label := append([]rune(status), []rune(strings.Repeat(" ", e.COLUMNS-atx))...)

		for i := 0; i < len(label); i++ {
			e.Screen.SetContent(atx+i, e.ROWS-1, label[i], nil, StyleDefault)
		}
	}
}

func (e *Editor) DrawTest(line int, row int) {
	x := e.COLUMNS - 2
	e.Screen.SetContent(x, row, '▶', nil,
		StyleDefault.Foreground(Color(HighlighterGlobal.GetRunButtonStyle())))
}

func (e *Editor) DrawProcessPanelSearch(pattern []rune, patternx int) {

	var prefix = []rune("  search: ")

	for i := 0; i < len(prefix); i++ {
		e.Screen.SetContent(i, e.TERMINAL_HEIGHT-1, prefix[i], nil, StyleDefault)
	}

	e.Screen.SetContent(len(prefix), e.TERMINAL_HEIGHT-1, ' ', nil, StyleDefault)

	for i := 0; i < len(pattern); i++ {
		e.Screen.SetContent(len(prefix)+i, e.TERMINAL_HEIGHT-1, pattern[i], nil, StyleDefault)
	}

	e.Screen.ShowCursor(len(prefix)+patternx, e.TERMINAL_HEIGHT-1)

	for i := len(prefix) + len(pattern); i < e.COLUMNS; i++ {
		e.Screen.SetContent(i, e.TERMINAL_HEIGHT-1, ' ', nil, StyleDefault)
	}

	if len(e.ProcessPanelSearchResults) > 0 {
		status := fmt.Sprintf("  %d/%d", e.ProcessPanelSearchResultIndex+1, len(e.ProcessPanelSearchResults))

		for i := 0; i < len(status); i++ {
			e.Screen.SetContent(len(prefix)+len(pattern)+i, e.TERMINAL_HEIGHT-1,
				rune(status[i]), nil, StyleDefault)
		}
	}
}
