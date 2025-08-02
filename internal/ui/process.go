package ui

import (
	. "github.com/vipmax/edgo/internal/highlighter"
	. "github.com/vipmax/edgo/internal/process"
	. "github.com/vipmax/edgo/internal/utils"
	. "github.com/gdamore/tcell"
	"os"
	"strings"
	"time"
)

func (e *Editor) OnProcessRun(newRun bool) {
	if newRun && (e.Lang == "" || e.langConf.Cmd == "") { return }

	if e.ProcessPanelHeight == 0 {
		e.ProcessPanelHeight = 10
		e.COLUMNS, e.ROWS = e.Screen.Size()
		e.ROWS -= e.ProcessPanelHeight
	}

	var args = []string{e.AbsoluteFilePath}

	if e.langConf.CmdArgs != "" {
		args = append(strings.Split(e.langConf.CmdArgs, " "), e.AbsoluteFilePath)
	}

	cmd := e.langConf.Cmd

	if !newRun && e.Process != nil && e.Process.Cmd != nil {
		// use prev cmd and args
		cmd = e.Process.Cmd.Path
		args = e.Process.Cmd.Args[1:]
	}

	ResetSelectionColor()
	e.Process = NewProcess(cmd, args...)
	e.Process.Cmd.Env = append(os.Environ())

	if e.Lang == "python" {
		// printing immediately
		e.Process.Cmd.Env = append(e.Process.Cmd.Env, "PYTHONUNBUFFERED=1")
	}

	e.ProcessContent = []Line{Line{[]rune{}}}
	e.ProcessPanelScroll = 0
	e.ProcessPanelSpacing = 2

	e.Process.Start()

	go func() {
		for range e.Process.Updates {

			//e.ProcessContent = append(e.ProcessContent, []rune(line))

			//newLines := e.Process.Lines[len(e.ProcessContent):]
			newLines := e.Process.GetLines(len(e.ProcessContent))
			for _, line := range newLines {
				e.ProcessContent = append(e.ProcessContent, Line{[]rune(line)})
				if len(e.ProcessContent) > e.ProcessPanelHeight {
					if e.ProcessPanelScroll >= len(e.ProcessContent)-e.ProcessPanelHeight-1 {
						e.ProcessPanelScroll = len(e.ProcessContent) - e.ProcessPanelHeight + 1 // focusing
						e.ProcessPanelScroll = Max(0, e.ProcessPanelScroll)
					}
				}
			}

			e.DrawProcessPanel()
			e.Screen.Show()

			if e.Process.IsStopped() {
				if len(e.ProcessContent) > e.ProcessPanelHeight { // focusing
					e.ProcessPanelScroll = len(e.ProcessContent) - e.ProcessPanelHeight + 1
				}
				e.DrawProcessPanel()
				e.Screen.Show()
				break
			}
		}
	}()

}

func (e *Editor) OnProcessKeyHandle(key Key, keyrune rune) {
	if key == KeyCtrlF {
		e.OnProcessSearch()
	}

	if key == KeyUp {
		e.ProcessPanelCursorY--
		if e.ProcessPanelCursorY < 0 { e.ProcessPanelCursorY = 0 }
		if e.ProcessPanelCursorX > len(e.ProcessContent[e.ProcessPanelCursorY].Buf) {
			e.ProcessPanelCursorX = len(e.ProcessContent[e.ProcessPanelCursorY].Buf)
		}
		e.FocusProcessPanel()
	}
	if key == KeyDown {
		e.ProcessPanelCursorY++
		if e.ProcessPanelCursorY >= len(e.ProcessContent) { e.ProcessPanelCursorY = len(e.ProcessContent) - 1 }
		if e.ProcessPanelCursorX > len(e.ProcessContent[e.ProcessPanelCursorY].Buf) {
			e.ProcessPanelCursorX = len(e.ProcessContent[e.ProcessPanelCursorY].Buf)
		}
		e.FocusProcessPanel()
	}
	if key == KeyRight {
		e.ProcessPanelCursorX++
		if e.ProcessPanelCursorX > len(e.ProcessContent[e.ProcessPanelCursorY].Buf) {
			e.ProcessPanelCursorX = len(e.ProcessContent[e.ProcessPanelCursorY].Buf)
		}

		//if e.ProcessPanelCursorX >= e.TERMINAL_WIDHT - e.ProcessPanelSpacing {
		//	e.ProcessPanelHScroll = e.ProcessPanelCursorX - e.TERMINAL_WIDHT + e.ProcessPanelSpacing
		//}
	}
	if key == KeyLeft {
		e.ProcessPanelCursorX--
		if e.ProcessPanelCursorX < 0 { e.ProcessPanelCursorX = 0 }
		if e.ProcessPanelCursorX > len(e.ProcessContent[e.ProcessPanelCursorY].Buf) {
			e.ProcessPanelCursorX = len(e.ProcessContent[e.ProcessPanelCursorY].Buf)
		}

		//if e.ProcessPanelCursorX >= e.TERMINAL_WIDHT - e.ProcessPanelSpacing {
		//	e.ProcessPanelHScroll = e.ProcessPanelCursorX - e.TERMINAL_WIDHT + e.ProcessPanelSpacing
		//}
	}
	if key == KeyRune && keyrune == 'l' {
		e.ProcessPanelHScroll++
	}
	if key == KeyRune && keyrune == 'l' {
		e.ProcessPanelHScroll++
	}
	if key == KeyRune && keyrune == 's' {
		e.OnProcessStop()
	}
	if key == KeyRune && keyrune == 'f' {
		if len(e.ProcessContent) > e.ProcessPanelHeight { // focusing
			e.ProcessPanelScroll = len(e.ProcessContent) - e.ProcessPanelHeight + 1
		}
	}
}

func (e *Editor) OnProcessStop() {
	if e.Process != nil && e.Process.Cmd != nil && !e.Process.IsStopped() {
		e.Process.Stop()
		time.Sleep(time.Millisecond * 300) // give a time to show 'kill' message
	}

	if e.Dap.IsStarted {
		e.OnDebugStop()
	}

	if len(e.ProcessContent) > e.ProcessPanelHeight { // focusing
		e.ProcessPanelScroll = len(e.ProcessContent) - e.ProcessPanelHeight + 2
	}
	ResetSelectionColor()
}
