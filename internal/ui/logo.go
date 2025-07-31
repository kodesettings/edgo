package ui

import (
	. "github.com/vipmax/edgo/internal/highlighter"
	. "github.com/gdamore/tcell"
	"strings"
)

// draw logo
func (e *Editor) DrawLogo() {
	// https://patorjk.com/software/taag/#p=testall&h=0&f=3D%20Diagonal&t=edgo
	logo :=
`
 _______ ______   ______  _____ 
 |______ |     \ |  ____ |     |
 |______ |_____/ |_____| |_____|
                                
`
	// split logo by new line
	lines := strings.Split(logo, "\n")

	logoWidth := len(lines[1])
	//fromx := e.COLUMNS / 2 - logoWidth / 2
	fromy := e.ROWS/2 - len(lines)/2

	fromx := (e.COLUMNS+28-2)/2 - logoWidth/2

	for i, line := range lines {
		for j, ch := range line {
			e.Screen.SetContent(j+fromx, i+fromy, ch, nil,
				StyleDefault.Foreground(Color(AccentColor2)))
		}
	}

	e.Screen.Show()
}

