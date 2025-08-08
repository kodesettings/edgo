package main

import (
	config "github.com/vipmax/edgo/internal/config"
	highlighter "github.com/vipmax/edgo/internal/highlighter"
	logger "github.com/vipmax/edgo/internal/logger"
	ui "github.com/vipmax/edgo/internal/ui"
	"fmt"
	"runtime"
)

func main() {
	loghandler := logger.FileHandler{}
	loghandler.SetLogger()
	conf := config.GetConfig()
	highlighter.HighlighterGlobal.SetTheme(conf.Theme)
	editor := ui.Editor{}
	editor.Config = conf

	defer func() {
		if r := recover(); r != nil {
			editor.Exit()
			errMsg := fmt.Sprintf("Recovered from panic. Error: %v\n", r)
			stackTrace := make([]byte, 4096)
			stackSize := runtime.Stack(stackTrace, false)
			fmt.Printf("%s\nStack Trace:\n%s\n", errMsg, stackTrace[:stackSize])
		}
	}()

	editor.Start()
}
