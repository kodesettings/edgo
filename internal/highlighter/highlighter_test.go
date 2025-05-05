package highlighter

import (
	. "github.com/vipmax/edgo/internal/utils"
	"fmt"
	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/styles"
	"github.com/go-enry/go-enry/v2"
	"testing"
	"time"
	assert "github.com/stretchr/testify/assert"
)

func TestColorize(t *testing.T) {
	file := "highlighter_test.go"
	filecontent, _ := ReadFileToString(file)

	h := Highlighter{}
	characterColors := h.Colorize(filecontent, file)

	for color := range characterColors {
		assert.NotEqual(t, color, "", "color is empty")
	}
}

func TestColorizeJs(t *testing.T) {
	file := "highlighter_test.js"
	const code = `
function hello() { 
	console.log('hello') 
}
`

	h := Highlighter{}

	start := time.Now()
	characterColors := h.Colorize(code, file)
	fmt.Println("colorized, elapsed", time.Since(start))

	for _, color := range characterColors {
		assert.NotEqual(t, color, "", "color is empty")
	}
}

func TestGetStyle(t *testing.T) {
	style := styles.Get("github")
	if style == nil { style = styles.Fallback }

	assert.NotEqual(t, style, "", "style is empty")
	assert.NotEqual(t, style.Get(chroma.Comment).Background, "", "style chroma comment backgroudn is empty")
}

func TestGetColor(t *testing.T) {
	col := ColorFromString("#fc9994")
	assert.Equal(t, col, 33331604, "color does not match")
}

func TestLangDetect(t *testing.T) {
	file := "highlighter_test.go"
	lang := DetectLang(file)
	assert.Equal(t, lang, "go", "output must be go")
}

func TestLangDetect2(t *testing.T) {
	file := "highlighter_test.go"
	lang, safe := enry.GetLanguageByExtension(file)
	assert.Equal(t, lang, "Go", "lang must be Go")
	assert.Equal(t, safe, true, "save must be true")
}
