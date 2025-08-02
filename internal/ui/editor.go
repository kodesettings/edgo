package ui

import (
	. "github.com/vipmax/edgo/internal/config"
	"github.com/vipmax/edgo/internal/dap"
	. "github.com/vipmax/edgo/internal/highlighter"
	. "github.com/vipmax/edgo/internal/io"
	. "github.com/vipmax/edgo/internal/logger"
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
	"github.com/gdamore/tcell/encoding"
	"github.com/rjeczalik/notify"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
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

	FilesPanelWidth     int        // current width for files panel
	Files               []FileInfo // current dir files
	IsFileSelection     bool       // true if in files selection menu
	FileScrollingOffset int        // for vertical scrolling  in selection menu
	FileSelectedIndex   int        // selected file index
	IsFilesSearch       bool       // true if in files search mode
	IsFilesPanelMoving  bool       // true if in files panel moving mode
	Tree                FileInfo   // files Tree
	FilesSearchPattern  []rune

	IsContentSearch   bool
	SearchPattern     []rune // pattern for search in a buffer
	SearchResults     []SearchResult
	SearchResultIndex int

	//filesInfo []FileInfo
	CursorHistory     []CursorMove
	CursorHistoryUndo []CursorMove

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

	Dap       dap.DapClient
	DebugInfo DebugInfo

	treeSitterHighlighter *TreeSitterHighlighter
	code string // storing UTF-8 representation of buffer

	FileWatcher *FileWatcher
	DirWatcher  *DirWatcher

	Tests      map[int]TestData
	TestFinder TestFinder
	Test       Test

	TreePath *Path

	HighlightElements map[int][]NodeRange

	// drawingWg sync.WaitGroup
	mu sync.Mutex
}

func (e *Editor) Start() {
	Log.Info("starting edgo")

	e.Init()

	cwd, _ := os.Getwd()
	e.Cwd = cwd

	// reading file from cmd args
	if len(os.Args) == 1 {
		// if no args, open current dir
		e.DrawLogo()
		e.OnFilesTree(false)
	} else {
		e.Filename = os.Args[1]
		e.InputFile = e.Filename

		info, err := os.Stat(e.InputFile)
		//if err != nil { log.Fatal(err); return }

		if info != nil && info.IsDir() {
			// if arg is dir, go to dir and open
			err = os.Chdir(e.InputFile)
			if err != nil { log.Fatal(err) }
			e.OnFilesTree(true)
		} else {
			// if arg is file, open file
			err := e.OpenFile(e.InputFile)
			if err != nil { log.Fatal(err) }
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
		if e.Dap.IsStarted {
			e.OnDebugKeyHandle(key, ev, 1)
			return
		}

		//c := strconv.Itoa(int(key))
		//Log.Info("EventKey", c)

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

	Log.Info("open", e.AbsoluteFilePath)

	newLang := DetectLang(e.AbsoluteFilePath)
	Log.Info("new lang is", newLang)

	if newLang != "" && newLang != e.Lang {
		e.Lang = newLang

		_, found := e.lsp2lang[newLang]
		if !found {
			lsp := LspClient{Lang: newLang}
			e.lsp2lang[newLang] = &lsp
			go e.InitLsp(e.Lang)
		}

		if e.Dap.Port > 0 {
			e.Dap = dap.DapClient{Lang: newLang, Conntype: "tcp", Port: e.Dap.Port + 1}
		} else {
			e.Dap = dap.DapClient{Lang: newLang, Conntype: "tcp", Port: 54752}
		}

		e.DebugInfo = DebugInfo{stopline: -1}
	}

	conf, found := e.Config.Langs[e.Lang]
	if !found { conf = DefaultLangConfig }
	e.langConf = conf
	e.langTabWidth = conf.TabWidth

	e.code = e.ReadFile(e.AbsoluteFilePath)
	e.treeSitterHighlighter = NewTreeSitter()
	e.treeSitterHighlighter.SetTheme(e.Config.Theme)
	e.treeSitterHighlighter.SetLang(e.Lang)
	e.treeSitterHighlighter.Parse(&e.code)
	clear(e.HighlightElements)

	e.Undo = []EditOperation{}
	e.Redo = []EditOperation{}

	e.UpdateFilesOpenStats(fname)

    e.Row = 0; e.Col = 0; e.Y = 0; e.X = 0
	e.Selection = Selection{-1,-1,-1,-1,false }
	e.SearchResults = []SearchResult{}

	// Opening file for LSP
	e.UpdateLsp(true, ConvertLinesToString(e.Lines))

	e.FileWatcher.UpdateFile(e.AbsoluteFilePath)
	e.FileWatcher.UpdateStats()

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
	e.FileSelectedIndex = -1
	e.CursorHistory = []CursorMove{}
	e.lsp2lang = map[string]*LspClient{}
	e.lspver = map[string]int{}
	e.DebugInfo = DebugInfo{}

	e.treeSitterHighlighter = NewTreeSitter()
	e.treeSitterHighlighter.SetTheme(e.Config.Theme)

	e.FileWatcher = NewFileWatcher(1000)
	e.FileWatcher.StartWatch(e.OnFileUpdate)

	e.DirWatcher = NewDirWatcher(".")
	e.DirWatcher.StartWatch(e.OnFilesTreeUpdate)

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
	e.UpdateLsp(true, ConvertLinesToString(e.Lines))

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
		atx := 0 + e.LINES_WIDTH + e.FilesPanelWidth
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

func (e *Editor) OnFilesTree(forceOpen bool) {
	e.IsFileSelection = true

	if e.FilesPanelWidth == 0 {
		tree, _ := ReadDirTree(e.Cwd, "", false, 0)
		e.Tree = tree
		if len(tree.Childs) == 0 { return }
		e.FilesPanelWidth = 28
		// root is always opened
		e.Tree.IsDirOpen = true
	}

	if e.Filename != "" { e.DrawEverything() }

	var end = false
	var patternx = len(e.FilesSearchPattern)

	// loop until escape or enter pressed
	for !end {
		//_, screenRows := e.Screen.Size()
		_, screenRows := 0, e.ROWS

		if e.FileSelectedIndex != -1 && e.FileSelectedIndex < e.FileScrollingOffset {
			e.FileScrollingOffset = e.FileSelectedIndex
		}
		if e.FileSelectedIndex >= e.FileScrollingOffset+screenRows {
			e.FileScrollingOffset = e.FileSelectedIndex - screenRows + 1
		}

		treeSize := TreeSize(e.Tree, 0)
		var aty = 0
		var fileindex = 0

		for row := 0; row < screenRows; row++ {
			for col := 0; col < e.FilesPanelWidth-2; col++ { // clean
				e.Screen.SetContent(col, row, ' ', nil, StyleDefault)
			}
			//e.Screen.SetContent(e.FilesPanelWidth-2, row, '▕', nil, SeparatorStyle)
		}

		e.DrawTree(e.Tree, 0, &fileindex, &aty)
		e.DrawTreeSearch(e.FilesSearchPattern, patternx)
		e.Screen.Show()

		switch ev := e.Screen.PollEvent().(type) { // poll and handle event
		case *EventMouse:
			mx, my := ev.Position()
			buttons := ev.Buttons()
			//modifiers := ev.Modifiers()

			if mx > e.FilesPanelWidth-3 { end = true; continue }
			if my >= screenRows { end = true; continue }

			if buttons&WheelDown != 0 && treeSize > screenRows {
				if e.FileScrollingOffset < treeSize-screenRows {
					e.FileScrollingOffset++
				}
			}
			if buttons&WheelUp != 0 && e.FileScrollingOffset > 0 { e.FileScrollingOffset-- }

			if my < treeSize { e.FileSelectedIndex = my + e.FileScrollingOffset }
			if buttons&Button1 == 1 {
				e.FileSelectedIndex = my + e.FileScrollingOffset
				if e.FileSelectedIndex < 0 { continue }
				if e.FileSelectedIndex >= treeSize { continue }
				if e.FileSelectedIndex >= treeSize { continue }
				if !e.IsMouseUnderFile(mx) { continue }
				end = e.SelectAndOpenFile()
				if end && e.IsFilesSearch {
					tree, _ := ReadDirTree(e.Cwd, "", false, 0)
					e.Tree = tree
					e.Tree.IsDirOpen = true
					SetDirOpenFlag(&e.Tree, e.InputFile)
					e.IsFilesSearch = false
					e.FilesSearchPattern = []rune{}
				}
			}

		case *EventKey:
			key := ev.Key()

			if key == KeyCtrlQ { e.Screen.Fini(); os.Exit(1) }
			if key == KeyCtrlN { e.NewFileOrDir() }
			if key == KeyCtrlF { e.IsFilesSearch = !e.IsFilesSearch }
			if key == KeyEscape && !e.IsFilesSearch { end = true; e.FilesPanelWidth = 0 }
			if key == KeyEscape && e.IsFilesSearch {
				end = true
				e.IsFilesSearch = false
				e.CleanFilesSearch()
				e.Screen.Show()
			}
			if key == KeyDown { e.FileSelectedIndex = Min(treeSize-1, e.FileSelectedIndex+1) }
			if key == KeyUp { e.FileSelectedIndex = Max(0, e.FileSelectedIndex-1) }
			if key == KeyLeft && e.IsFilesSearch && patternx > 0 { patternx-- }
			if key == KeyRight && e.IsFilesSearch && patternx < len(e.FilesSearchPattern) { patternx++ }
			if key == KeyCtrlT {
				end = true
				e.IsFilesSearch = false
				e.FilesPanelWidth = 0
			}
			if key == KeyBackspace2 && e.IsFilesSearch && patternx > 0 && len(e.FilesSearchPattern) > 0 {
				patternx--
				e.FilesSearchPattern = Remove(e.FilesSearchPattern, patternx)
				if len(e.FilesSearchPattern) != 0 {
					tree, _ := ReadDirTree(e.Cwd, string(e.FilesSearchPattern), true, 0)
					tree = FilterIfLeafEmpty(tree)
					e.Tree = tree
				} else {
					tree, _ := ReadDirTree(e.Cwd, "", false, 0)
					e.Tree = tree
				}

				e.Tree.IsDirOpen = true
				e.FileScrollingOffset = 0
				_, i := FindFirstFile(e.Tree, 0)
				e.FileSelectedIndex = i
			}
			if key == KeyRune {
				e.IsFilesSearch = true
				e.FilesSearchPattern = InsertTo(e.FilesSearchPattern, patternx, ev.Rune())
				patternx++
				tree, _ := ReadDirTree(e.Cwd, string(e.FilesSearchPattern), true, 0)
				tree = FilterIfLeafEmpty(tree)
				e.Tree = tree
				e.Tree.IsDirOpen = true
				e.FileScrollingOffset = 0
				_, i := FindFirstFile(e.Tree, 0)
				e.FileSelectedIndex = i
			}

			if key == KeyEnter || !e.IsFilesSearch && (key == KeyLeft || key == KeyRight) {
				end = e.SelectAndOpenFile()
				if end && e.IsFilesSearch {
					tree, _ := ReadDirTree(e.Cwd, "", false, 0)
					e.Tree = tree
					e.Tree.IsDirOpen = true
					SetDirOpenFlag(&e.Tree, e.InputFile)
					e.IsFilesSearch = false
					e.FilesSearchPattern = []rune{}
				}
			}
		}
	}

	e.IsFileSelection = false
}

func (e *Editor) SelectAndOpenFile() bool {
	found, selectedFile := GetSelected(e.Tree, e.FileSelectedIndex)
	if found {
		if selectedFile.IsDir {
			selectedFile.IsDirOpen = !selectedFile.IsDirOpen
			return false
		} else {
			e.InputFile = selectedFile.FullName
			e.OpenFile(e.InputFile)
			return true
		}
	}
	return false
}

func (e *Editor) IsMouseUnderFile(mx int) bool {
	found, selectedFile := GetSelected(e.Tree, e.FileSelectedIndex)
	if found {
		if selectedFile.Level+len(selectedFile.Name)+1 >= mx {
			if mx < selectedFile.Level+2 { return false } // 2 is spacing
			return true
		} else {
			return false
		}
	}
	return false
}

func (e *Editor) OverlayFalse() {
	e.IsOverlay = false
}

func (e *Editor) UpdateColors() {
	code := ConvertLinesToString(e.Lines)
	e.treeSitterHighlighter.ReParse(&code) //  todo:: optimize
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

func (e *Editor) OnFilesTreeUpdate(event notify.EventInfo) {
	fullname := event.Path()
	name := filepath.Base(fullname)
	dir := filepath.Dir(fullname)
	parentNode := FindByFullName(&e.Tree, dir)
	if parentNode == nil { return }

	switch event.Event() {
	case notify.Create:
		// find and add node to parent
		fileInfo, err := os.Stat(fullname)
		if err != nil { return }
		isDir := fileInfo.IsDir()

		f := FileInfo{
			Name: name, FullName: fullname,
			IsDir: isDir, IsDirOpen: false,
			Childs: []FileInfo{}, Level: parentNode.Level + 1,
		}

		parentNode.Childs = append(parentNode.Childs, f)
		SortTree(*parentNode)

	case notify.Remove:
		// find and remove node from parent
		for i, child := range parentNode.Childs {
			if child.FullName == fullname {
				parentNode.Childs = Remove(parentNode.Childs, i)
				break
			}
		}
	default:
		// update all
		//tree, _ := ReadDirTree(e.Cwd, "", false, 0)
		//tree.IsDirOpen = true
		//e.Tree = tree
	}

	e.DrawEverything()
	e.Screen.Show()
}

func (e *Editor) NewFileOrDir() {
	inputName := make([]rune, 0)
	cursorx := 0
	end := false
	pref := " create "
	directory := ""

	// add current
	found, selectedFile := GetSelected(e.Tree, e.FileSelectedIndex)
	if found {
		rel, err := filepath.Rel(e.Cwd, selectedFile.FullName)
		if err != nil { return }

		if selectedFile.IsDir {
			directory = rel + string(filepath.Separator)
		} else {
			directory = filepath.Dir(rel)
		}
	}

	for !end {

		e.Screen.ShowCursor(len(pref)+len(directory)+cursorx, e.ROWS-1)

		for col := 0; col <= len(pref)+len(inputName) || col < e.FilesPanelWidth; col++ { // clean
			e.Screen.SetContent(col, e.ROWS-1, ' ', nil, StyleDefault)
		}

		for i, ch := range pref { // draw prefix
			e.Screen.SetContent(i, e.ROWS-1, ch, nil, StyleDefault)
		}
		for i, ch := range directory { // draw directory
			e.Screen.SetContent(i+len(pref), e.ROWS-1, ch, nil, StyleDefault)
		}

		for i, ch := range inputName { // draw inputName
			e.Screen.SetContent(i+len(pref)+len(directory), e.ROWS-1, ch, nil, StyleDefault)
		}

		e.Screen.Show()

		switch ev := e.Screen.PollEvent().(type) {
		case *EventKey:
			key := ev.Key()

			if key == KeyRune {
				inputName = InsertTo(inputName, cursorx, ev.Rune())
				cursorx++
			}
			if key == KeyBackspace2 && cursorx > 0 && len(inputName) > 0 {
				cursorx--
				inputName = Remove(inputName, cursorx)
			}
			if key == KeyLeft && cursorx > 0 { cursorx-- }
			if key == KeyRight && cursorx < len(inputName) { cursorx++ }
			if key == KeyEscape { end = true }

			if key == KeyEnter {
				name := filepath.Join(directory, string(inputName))
				nameabs, err := filepath.Abs(name)
				if err != nil {
					return
				}

				isCreateDir := strings.HasSuffix(string(inputName), string(filepath.Separator))

				if isCreateDir {
					// dir/  - create dir
					err := os.MkdirAll(nameabs, 0750)
					if err != nil { return }
				} else {  // create file
					file, err := os.Create(nameabs)
					defer file.Close()
					if err != nil { return }
				}

				end = true
			}
		}
	}

	for col := 0; col <= len(pref)+len(inputName); col++ { // clean
		e.Screen.SetContent(col, e.ROWS-1, ' ', nil, StyleDefault)
	}
}

func (e *Editor) OnCursorChanged() {
	clear(e.HighlightElements)

	if e.Selection.IsSelectionNonEmpty() { return }

	start := time.Now()
	nodename, noderange := e.treeSitterHighlighter.GetNodeAt(e.Row, e.Col, e.Row, e.Col)

	if strings.Contains(nodename, "identifier") {
		if len(e.Lines) >= noderange.Ssy { return }
		runes := e.Lines[noderange.Ssy].Buf[noderange.Ssx:noderange.Sex]
		content := string(runes)

		searchResults := Search(e.Lines, content)
		e.HighlightElements = make(map[int][]NodeRange)

		if len(searchResults) > 1 {
			for _, searchResult := range searchResults {
				searchnodename, searchnoderange := e.treeSitterHighlighter.GetNodeAt(
					searchResult.Line, searchResult.Position, searchResult.Line, searchResult.Position+len(content))
				if nodename == searchnodename {
					if len(content) != searchnoderange.Sex-searchnoderange.Ssx { continue }
					e.HighlightElements[searchResult.Line] = append(
							e.HighlightElements[searchResult.Line], searchnoderange)
				}
			}
		}
	}

	elapsed := time.Since(start).String()
	Log.Info(nodename,
		fmt.Sprintf("%d %d %d %d elapsed %s", noderange.Ssx,
			noderange.Ssy, noderange.Sex, noderange.Sey, elapsed))
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
