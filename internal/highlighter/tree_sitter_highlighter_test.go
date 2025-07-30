package highlighter

import (
	. "github.com/vipmax/edgo/internal/utils"
	"fmt"
	"testing"
	"time"
	assert "github.com/stretchr/testify/assert"
)


func TestTreeSitterHighlighterColors(t *testing.T) {
	treeSitterHighlighter := NewTreeSitter()
	treeSitterHighlighter.SetLang("go")

	filecode, _ := ReadFileToString("../../internal/ui/editor.go")
	code := filecode

	treeSitterHighlighter.Parse(&code)

	start := time.Now()
	colors := treeSitterHighlighter.ColorRanges(1, 2, []byte(code))
	fmt.Println("colorized, elapsed", time.Since(start).Nanoseconds())

	for i, colorsLine := range colors {
		fmt.Println(i, "line", colorsLine)
	}
}

func TestColorFromString(t *testing.T) {
	h := NewTreeSitter()
	col := h.ParseColor("#fc9994")
	fmt.Println(col)

	want := 33331604
	assert.Equal(t, want, col, "got %v want %v", want, col)
}
