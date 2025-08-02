package ui

import (
	. "github.com/gdamore/tcell"
	"os"
)

func (e *Editor) HandleKeyboard(key Key, ev *EventKey, modifiers ModMask) {
	if key == KeyCtrlF && !e.IsProcessPanelFocused { e.OnSearch() }
	if ev.Rune() == 'Y' && modifiers&ModAlt != 0 { e.OnLangLinesCount() } // alt + shift + y

	if e.Filename == "" && key != KeyCtrlQ { return }

	if e.IsProcessPanelFocused {
		e.OnProcessKeyHandle(key, ev.Rune())
		return
	}

	intrune := int(ev.Rune())
	if ev.Rune() == '/' && modifiers&ModAlt != 0 || intrune == '÷' {
		// '÷' is option + '/' on Mac
		e.OnCommentLine()
		return
	}
	if intrune == '¨' {
		// '¨' is option + u on Mac
		e.OnRedo()
		return
	}

	if key == KeyUp && modifiers == 3 { e.OnSwapLinesUp(); return } // control + shift + up
	if key == KeyDown && modifiers == 3 { e.OnSwapLinesDown(); return } // control + shift + down
	if key == KeyBacktab { e.OnBackTab(); return }
	if key == KeyTab { e.OnTab(); return }
	if key == KeyCtrlH { e.OnHover(); return }
	if key == KeyCtrlR { e.OnReferences(); return }
	if key == KeyCtrlW { e.OnCodeAction(); return }
	if key == KeyCtrlP { e.OnSignatureHelp(); return }
	if key == KeyCtrlG { e.OnDefinition(); return }
	if key == KeyCtrlE { e.OnErrors(); return }
	if key == KeyCtrlC { e.OnCopy(); return }
	if key == KeyCtrlV { e.OnPaste(); return }
	if key == KeyEscape { e.Selection.CleanSelection(); return }
	if key == KeyCtrlA { e.OnSelectAll(); return }
	if key == KeyCtrlX { e.Cut(true) }
	if key == KeyCtrlD { e.Duplicate() }
	if key == KeyCtrlB { e.Breakpoint() }
	if key == KeyCtrlK { e.GoBottom(); return }
	if key == KeyCtrlJ { e.GoTop(); return }
	if key == KeyCtrlL { e.GoToLine(); return }

	if modifiers&ModShift != 0 && (key == KeyRight || key == KeyLeft || key == KeyUp || key == KeyDown) {

		if e.Selection.Ssx < 0 { e.Selection.Ssx, e.Selection.Ssy = e.Col, e.Row }
		if key == KeyRight { e.OnRight() }
		if key == KeyLeft { e.OnLeft() }
		if key == KeyUp { e.OnUp(false) }
		if key == KeyDown { e.OnDown(false) }
		if e.Selection.Ssx >= 0 {
			e.Selection.Sex, e.Selection.Sey = e.Col, e.Row
			e.Selection.IsSelected = true
		}
		return
	}

	if key == KeyRune && modifiers&ModAlt != 0 && len(e.Lines) > 0 {
		e.HandleSmartMove(ev.Rune())
		return
	}

	if modifiers&ModAlt != 0 && (int(ev.Key()) == 259 || int(ev.Key()) == 260) && len(e.Lines) > 0 {
		e.HandleSmartMoveAlac(int(ev.Key()))
		return
	}

	if key == KeyDown && modifiers&ModAlt != 0 {
		e.OnSelectLessAtCursor()
		//e.HandleSmartMoveDown();
		return
	}
	if key == KeyUp && modifiers&ModAlt != 0 {
		e.OnSelectMoreAtCursor()
		//e.HandleSmartMoveUp();
		return
	}

	if key == KeyRune {
		e.AddChar(ev.Rune())
		if ev.Rune() == '.' {
			e.DrawEverything()
			e.Screen.Show()
			e.OnCompletion()
		}
		//if ev.Rune() == '(' { e.DrawEverything(); e.Screen.Show(); e.OnSignatureHelp(); e.Screen.Clear() }
	}

	if /*key == tcell.KeyEscape ||*/ key == KeyCtrlQ { e.Screen.Fini(); os.Exit(1) }
	if key == KeyCtrlS { e.WriteFile(true) }
	if key == KeyEnter { e.OnEnter(); return }
	if key == KeyBackspace || key == KeyBackspace2 { e.OnDelete() }
	if key == KeyDown { e.OnDown(false); e.Selection.CleanSelection() }
	if key == KeyUp { e.OnUp(false); e.Selection.CleanSelection() }
	if key == KeyPgDn { e.OnDown(true); e.Selection.CleanSelection() }
	if key == KeyPgUp { e.OnUp(true); e.Selection.CleanSelection() }
	if key == KeyLeft { e.OnLeft(); e.Selection.CleanSelection() }
	if key == KeyRight { e.OnRight(); e.Selection.CleanSelection() }
	if key == KeyCtrlT { e.OnFilesTree(true) }
	if key == KeyF18 { e.OnRename() }
	if key == KeyF22 { e.OnProcessRun(true) }
	if key == KeyF23 { e.OnDebug() }
	if key == KeyCtrlZ { e.OnUndo() }
	if key == KeyCtrlY { e.OnRedo() }
	if key == KeyCtrlO { e.OnCursorBack() }
	if key == KeyCtrlRightSq { e.OnCursorBackUndo() }
	if key == KeyCtrlSpace { e.OnCompletion() }

}
