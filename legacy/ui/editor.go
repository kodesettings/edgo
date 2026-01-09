package ui

import (
	. "github.com/vipmax/edgo/internal/config"
	. "github.com/vipmax/edgo/internal/highlighter"
	. "github.com/vipmax/edgo/internal/lsp"
	. "github.com/vipmax/edgo/internal/operations"
	. "github.com/vipmax/edgo/internal/process"
	. "github.com/vipmax/edgo/internal/search"
	. "github.com/vipmax/edgo/internal/selection"
	. "github.com/vipmax/edgo/internal/tests"
	. "github.com/vipmax/edgo/internal/utils"
	"fmt"
	"github.com/atotto/clipboard"
	. "github.com/gdamore/tcell"
	"github.com/zyedidia/rope"
	"github.com/gdamore/tcell/encoding"
	"log/slog"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type Editor struct {
	COLUMNS     int // terminal size columns
	ROWS        int // terminal size rows
	LINES_WIDTH int // draw file lines number

	TERMINAL_HEIGHT int
	TERMINAL_WIDHT  int

	Row int // cursor position row
	Col int // cursor position column
	Y   int // row offset for scrolling
	X   int // col offset for scrolling

	Lines  []Line // Lines of text characters
	Screen Screen // Screen for drawing

	Lang         string // current file language
	Config       Config // config, lsp, tabs, comments, etc
	langConf     Lang   // current lang conf
	langTabWidth int    // current lang tabs indentation  '\t' -> "    "

	Selection Selection // selection

	Undo []EditOperation // stack for undo operations
	Redo []EditOperation // stack for redo operations

	Cwd              string // current dir
	InputFile        string // exact user input
	Filename         string // current file name
	AbsoluteFilePath string // current file name and directory
	IsContentChanged bool   // shows * if file is changed
	IsColorize       bool   // colorize text is true by default
	Update           bool   // for Screen updates,  if false it will not draw
	IsOverlay        bool   // true if overlay is active (completion, hover, errors...)
	IsStartupScreen  bool   // true if no file was opened

	IsContentSearch   bool
	SearchPattern     []rune // pattern for search in a buffer
	SearchResults     []SearchResult
	SearchResultIndex int

	// process panel vars
	ProcessPanelHeight            int
	ProcessPanelWidth             int
	ProcessContent                []Line
	ProcessPanelScroll            int
	ProcessPanelHScroll           int
	IsProcessPanelMoving          bool
	IsProcessPanelFocused         bool
	Process                       *Process
	ProcessPanelSpacing           int
	ProcessPanelCursorX           int
	ProcessPanelCursorY           int
	ProcessPanelSelection         Selection
	IsProcessPanelSearch          bool
	ProcessPanelSearchPattern     []rune // pattern for search in a buffer
	ProcessPanelSearchResults     []SearchResult
	ProcessPanelSearchResultIndex int

	lsp2lang map[string]*LspClient
	lspver map[string]int // version number sequencing

	treeSitterHighlighter *TreeSitterHighlighter
	code *rope.Node // storing a rope node element

	Tests      map[int]TestData
	TestFinder TestFinder
	Test       Test

	TreePath *Path

	// drawingWg sync.WaitGroup
	mu sync.Mutex
}

func (e *Editor) Start() {
	slog.Info("starting edgo")

	e.Init()

	cwd, _ := os.Getwd()
	e.Cwd = cwd

	// reading file from cmd args
	if len(os.Args) == 1 {
		// if no args, open current dir
		e.DrawLogo()
		e.IsStartupScreen = true
	} else {
		e.Filename = os.Args[1]
		e.InputFile = e.Filename

		if e.IsFileExists(e.InputFile) {
			// if arg is file, open file
			err := e.OpenFile(e.InputFile)
			if err != nil { slog.Error(err.Error()) }
			e.IsStartupScreen = false
		} else {
			e.IsStartupScreen = true
		}
	}

	// main draw cycle
	for {
		if e.Update && e.Filename != "" {
			e.DrawEverything()
			e.Screen.Show()
		}
		e.HandleEvents()
	}
}

func (e *Editor) Exit() {
	e.Screen.Fini()
}

func (e *Editor) HandleEvents() {
	//e.Update = false
	e.Update = true
	ev := e.Screen.PollEvent()
	switch ev := ev.(type) {
	case *EventResize:
		e.COLUMNS, e.ROWS = e.Screen.Size()
		e.TERMINAL_HEIGHT = e.ROWS
		e.TERMINAL_WIDHT = e.COLUMNS
		e.ROWS -= e.ProcessPanelHeight
		e.DrawEverything()
		e.Screen.Show()

	case *EventMouse:
		mx, my := ev.Position()
		buttons := ev.Buttons()
		modifiers := ev.Modifiers()

		e.HandleMouse(mx, my, buttons, modifiers)

	case *EventKey:
		key := ev.Key()
		modifiers := ev.Modifiers()

		e.HandleKeyboard(key, ev, modifiers)
	}
}

func (e *Editor) OpenFile(fname string) error {
	if strings.HasPrefix(fname, "~/") {
		currentUser, _ := user.Current()
		homeDir := currentUser.HomeDir
		fname = filepath.Join(homeDir, fname[2:])
	}

	absoluteDir, err := filepath.Abs(path.Dir(fname))
	if err != nil { return err }
	//directory := absoluteDir;
	e.Filename = filepath.Base(fname)
	e.AbsoluteFilePath = path.Join(absoluteDir, e.Filename)

	slog.Info("open", "file", e.AbsoluteFilePath)

	newLang := DetectLang(e.AbsoluteFilePath)
	slog.Info("new lang is", "lang", newLang)

	if newLang != "" && newLang != e.Lang {
		e.Lang = newLang

		_, found := e.lsp2lang[newLang]
		if !found {
			lsp := LspClient{Lang: newLang}
			e.lsp2lang[newLang] = &lsp
			go e.InitLsp(e.Lang)
		}
	}

	conf, found := e.Config.Langs[e.Lang]
	if !found { conf = DefaultLangConfig }
	e.langConf = conf
	e.langTabWidth = conf.TabWidth

	e.code = rope.New(e.ReadFile(e.AbsoluteFilePath))
	e.treeSitterHighlighter = NewTreeSitter()
	e.treeSitterHighlighter.SetTheme(e.Config.Theme)
	e.treeSitterHighlighter.SetLang(e.Lang)
	e.treeSitterHighlighter.Parse(e.code.Value())

	e.Undo = []EditOperation{}
	e.Redo = []EditOperation{}

    e.Row = 0; e.Col = 0; e.Y = 0; e.X = 0
	e.Selection = Selection{-1,-1,-1,-1,false }
	e.SearchResults = []SearchResult{}

	// Opening file for LSP
	e.UpdateLsp(true, string(e.code.Value()))

	e.FindTests()

	return nil
}

func (e *Editor) Init() {
	encoding.Register()
	screen, err := NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	e.Screen = screen

	err2 := e.Screen.Init()
	if err2 != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err2)
		os.Exit(1)
	}

	e.Screen.EnableMouse()
	e.Screen.Clear()

	e.COLUMNS, e.ROWS = e.Screen.Size()
	e.TERMINAL_HEIGHT = e.ROWS
	e.TERMINAL_WIDHT = e.COLUMNS
	e.LINES_WIDTH = 6
	e.Update = true
	e.IsColorize = true
	e.lsp2lang = map[string]*LspClient{}
	e.lspver = map[string]int{}

	e.treeSitterHighlighter = NewTreeSitter()
	e.treeSitterHighlighter.SetTheme(e.Config.Theme)

	e.Tests = make(map[int]TestData)
	e.TestFinder = TestFinder{}

	return
}

func (e *Editor) FindCursorXPosition(mx int) int {
	count := 0
	realCount := 0 // searching x position
	for _, ch := range e.Lines[e.Row].Buf {
		if count >= mx+e.X { break }
		if ch == '\t' && e.X == 0 {
			count += e.langTabWidth; realCount++
		} else { count++; realCount++ }
	}
	return realCount
}

func (e *Editor) InitLsp(lang string) {
	// Getting the lsp command with args for a language:
	conf, ok := e.Config.Langs[strings.ToLower(lang)]
	if !ok || len(conf.Lsp) == 0 { return }  // lang is not supported.

	lsp := e.lsp2lang[lang]
	split := strings.Split(conf.Lsp, " ")
	started := lsp.Start(split[0], split[1:]...)
	if !started { return }

	//var diagnosticUpdateChan = make(chan string)
	//go lsp.ReceiveLoop(diagnosticUpdateChan, lang)

	currentDir, _ := os.Getwd()

	lsp.Init(currentDir)
	e.UpdateLsp(true, string(e.code.Value()))

	go func() {
		// diagnostic updates
		for range lsp.DiagnosticsChannel {
			if e.IsOverlay { continue }
			e.DrawEverything()
			e.Screen.Show()
		}
	}()
}

func (e *Editor) OnErrors() {
	e.Screen.HideCursor()
	defer e.Screen.ShowCursor(e.Col, e.Row)

	if e.Lang == "" { return }
	lsp := e.lsp2lang[e.Lang]
	if !lsp.IsReady { return }

	maybeDiagnostics, found := lsp.GetDiagnostic("file://" + e.AbsoluteFilePath)

	if !found || len(maybeDiagnostics.Diagnostics) == 0 { return }

	e.IsOverlay = true
	defer e.OverlayFalse()

	var end = false

	// loop until escape or enter pressed
	for !end {

		var options = []string{}
		for i, diagnistic := range maybeDiagnostics.Diagnostics {
			text := fmt.Sprintf("%d/%d [%d:%d] %s ", i+1, len(maybeDiagnostics.Diagnostics),
				int(diagnistic.Range.Start.Line)+1, int(diagnistic.Range.Start.Character+1),
				diagnistic.Message,
			)
			options = append(options, text)
		}

		width := Max(50, MaxString(options)) // width depends on max option len or 30 at min
		height := MinMany(10, len(options))  // depends on min option len or 5 at min or how many rows to the end of e.Screen
		atx := 0 + e.LINES_WIDTH
		aty := 0 // Define the window  position and dimensions
		style := StyleDefault.Foreground(ColorWhite)

		var selectionEnd = false
		var selected = 0
		var selectedOffset = 0

		for !selectionEnd {
			if selected < selectedOffset { selectedOffset = selected } // calculate offsets for scrolling completion
			if selected >= selectedOffset+height { selectedOffset = selected - height + 1 }

			shifty := e.DrawErrors(atx, width, aty, height, options, selectedOffset, selected, style)

			e.Screen.Show()

			switch ev := e.Screen.PollEvent().(type) { // poll and handle event
			case *EventResize:
				e.COLUMNS, e.ROWS = e.Screen.Size()
				//ROWS -= 1
				e.Screen.Sync()
				e.Screen.Clear()
				e.DrawEverything()
				e.Screen.Show()

			case *EventKey:
				key := ev.Key()
				if key == KeyEscape || key == KeyEnter ||
					key == KeyBackspace || key == KeyBackspace2 ||
					key == KeyCtrlE {
					e.Screen.Clear()
					selectionEnd = true
					end = true
				}
				if key == KeyDown { selected = Min(len(options)-1, selected+1) }
				if key == KeyUp { selected = Max(0, selected-1) }
				if key == KeyCtrlC {
					diagnostic := maybeDiagnostics.Diagnostics[selected]
					clipboard.WriteAll(diagnostic.Message)
				}
				if key == KeyRight {
					diagnostic := maybeDiagnostics.Diagnostics[selected]
					e.Row = int(diagnostic.Range.Start.Line)
					e.Col = int(diagnostic.Range.Start.Character)
					e.Focus()
					// add space for errors panel
					if e.Row-e.Y < shifty+height { e.Y -= shifty + height + 1 }
					if e.Y < 0 { e.Y = 0 }
					e.DrawEverything()
					e.Screen.Show()
				}
				if key == KeyEnter {
					selectionEnd = true
					end = true
					diagnostic := maybeDiagnostics.Diagnostics[selected]
					e.Row = int(diagnostic.Range.Start.Line)
					e.Col = int(diagnostic.Range.Start.Character)

					e.Selection.Ssx = e.Col
					e.Selection.Ssy = e.Row
					e.Selection.Sey = int(diagnostic.Range.End.Line)
					e.Selection.Sex = int(diagnostic.Range.End.Character)
					e.Row = e.Selection.Sey
					e.Col = e.Selection.Sex
					e.Selection.IsSelected = true
					e.Focus()
					// add space for errors panel
					if e.Row-e.Y < shifty+height { e.Y -= shifty + height + 1 }
					if e.Y < 0 { e.Y = 0 }
					e.DrawEverything()
					e.Screen.Show()
				}
			}
		}
	}
}

func (e *Editor) OverlayFalse() {
	e.IsOverlay = false
}

func (e *Editor) OnFileUpdate() {
	row, col := e.Row, e.Col // save cursor
	x, y := e.X, e.Y         // safe scroll
	e.OpenFile(e.AbsoluteFilePath)

	// if row and col fits to content, restore cursor
	if row < len(e.Lines) {
		e.Row = row
		e.X = x
		if col < len(e.Lines[row].Buf) {
			e.Col = col
			e.Y = y
		}
	}
	e.DrawEverything()
	e.Screen.Show()
}

func (e *Editor) GoToLine() {
	var input = []rune{}
	var patternx = 0

	var end = false
	for !end {

		switch ev := e.Screen.PollEvent().(type) { // poll and handle event

		case *EventKey:
			key := ev.Key()

			if key == KeyRune {
				input = InsertTo(input, patternx, ev.Rune())
				patternx++
			}
			if key == KeyBackspace2 && patternx > 0 && len(input) > 0 {
				patternx--
				input = Remove(input, patternx)
			}
			if key == KeyESC  || key == KeyCtrlL { return }
			if key == KeyEnter {
				end = true
			}
		}
	}

	result, err := strconv.Atoi(string(input))
	if err != nil { return }

	e.Row = result - 1
	e.Col = 0

	if e.Row < 0 { e.Row = 0 } // fit to content
	if e.Row >= len(e.Lines) { e.Row = len(e.Lines) - 1 } // fit to content
	e.Focus()
}
