package utils

import (
	. "github.com/vipmax/edgo/internal/highlighter"
	"fmt"
	"testing"
	assert "github.com/stretchr/testify/assert"
)
		
func TestFormat(t *testing.T) {
	leftText := "Left"; rightText := "Right"; maxWidth := 30
	formattedText := FormatText(leftText, rightText, maxWidth)
	fmt.Println(formattedText)
}

func TestDetectGoLang(t *testing.T) {
	language := DetectLang("highlighter_test.go")
	fmt.Println(language)
	assert.Equal(t, "go", language, "language must be go, got %s", language)
}

func TestDetectPythonLang(t *testing.T) {
	language := DetectLang("test.py")
	fmt.Println(language)
	assert.Equal(t, "python", language, "language must be python, got %s", language)
}
