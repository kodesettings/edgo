package ui

import (
	. "github.com/vipmax/edgo/internal/search"
	"fmt"
	. "github.com/gdamore/tcell"
	"os"
	"strconv"
	"time"
)

func (e *Editor) OnLangLinesCount() {
	end := false
	start := time.Now()
	results, totalFilesProcessed, totalRowsProcessed := LineCountOnDirParallel(e.Cwd)
	langCount := LangCount(results)

	elapsed := time.Since(start)

	for !end {
		for j := 0; j < e.TERMINAL_HEIGHT; j++ {
			for i := e.FilesPanelWidth; i < e.COLUMNS; i++ {
				e.Screen.SetContent(i, j, ' ', nil, StyleDefault)
			}
		}

		e.DrawTextScreenReport("Lang lines count report, elapsed "+elapsed.String(), e.FilesPanelWidth+1, 0)
		e.DrawTextScreenReport("Total files "+strconv.Itoa(totalFilesProcessed), e.FilesPanelWidth+1, 2)
		e.DrawTextScreenReport("Total rows "+strconv.Itoa(totalRowsProcessed), e.FilesPanelWidth+1, 3)

		maxLangLen := 0
		maxFilesLen := 0
		maxLinesLen := 0
		for _, result := range langCount {
			if len(result.Lang) > maxLangLen {
				maxLangLen = len(result.Lang)
			}
			if len(strconv.Itoa(result.FilesCount)) > maxFilesLen {
				maxFilesLen = len(strconv.Itoa(result.FilesCount))
			}
			if len(strconv.Itoa(result.LinesCount)) > maxLinesLen {
				maxLinesLen = len(strconv.Itoa(result.LinesCount))
			}
		}

		l := fmt.Sprintf("%-*s", maxLangLen+5, "Language")
		l2 := fmt.Sprintf("%-*s", maxFilesLen+5, "Files")
		l3 := fmt.Sprintf("%-*s", maxLinesLen+3, "Lines")
		line := fmt.Sprintf("%s %s %s %s", l, l2, l3, "Empty/Code")
		e.DrawTextScreenReport(line, e.FilesPanelWidth+1, 5)

		for i, result := range langCount {
			l := fmt.Sprintf("%-*s", maxLangLen+5, result.Lang)
			l2 := fmt.Sprintf("%-*s", maxFilesLen+5, strconv.Itoa(result.FilesCount))
			l3 := fmt.Sprintf("%-*s", maxLinesLen+3, strconv.Itoa(result.LinesCount))
			line := fmt.Sprintf("%s %s %s %d/%d", l, l2, l3, result.EmptyLinesCount, result.LinesCount-result.EmptyLinesCount)
			e.DrawTextScreenReport(line, e.FilesPanelWidth+1, i+7)
		}
		e.Screen.Show()

		switch ev := e.Screen.PollEvent().(type) { // poll and handle event
		case *EventResize:
			e.COLUMNS, e.ROWS = e.Screen.Size()

		case *EventKey:
			key := ev.Key()

			if key == KeyCtrlQ {
				e.Screen.Fini()
				os.Exit(1)
			}

			if key == KeyESC || key == KeyEnter {
				end = true
				e.Screen.Clear()
				e.DrawEverything()
			}
		}
	}
}
