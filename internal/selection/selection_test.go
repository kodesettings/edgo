package selection

import (
	"testing"
	assert "github.com/stretchr/testify/assert"
)

func TestNoSelection(t *testing.T) {
	var Content = [][]rune{
		[]rune("Hello, world!"),
		[]rune("How are you doing today?"),
		[]rune("I hope you're doing well."),
	}

	var selection = Selection{}
	selection.Ssx = 0
	selection.Ssy = 0
	selection.Sex = 0
	selection.Sey = 0
	got := selection.GetSelectionString(Content)
	expected := ""

	assert.Equal(t, expected, got, "GetSelectionString() = %v, expected %v", got, expected)
}

func TestSingleCharacterSelection(t *testing.T) {
	var Content = [][]rune{
		[]rune("Hello, world!"),
		[]rune("How are you doing today?"),
		[]rune("I hope you're doing well."),
	}

	var selection = Selection{}
	selection.Ssx = 0
	selection.Ssy = 0
	selection.Sex = 1
	selection.Sey = 0
	got := selection.GetSelectionString(Content)
	expected := "H"

	assert.Equal(t, expected, got, "GetSelectionString() = %v, expected %v", got, expected)
}

func TestMultipleCharacterSelection(t *testing.T) {
	var Content = [][]rune{
		[]rune("Hello, world!"),
		[]rune("How are you doing today?"),
		[]rune("I hope you're doing well."),
	}

	var selection = Selection{}
	selection.Ssx = 0
	selection.Ssy = 0
	selection.Sex = 5
	selection.Sey = 0
	got := selection.GetSelectionString(Content)
	expected := "Hello"

	assert.Equal(t, expected, got, "GetSelectionString() = %v, expected %v", got, expected)
}

func TestMultipleLineSelection1(t *testing.T) {
	var Content = [][]rune{
		[]rune("Hello, world!"),
		[]rune("How are you doing today?"),
		[]rune("I hope you're doing well."),
	}

	var selection = Selection{}
	selection.Ssx = 0
	selection.Ssy = 0
	selection.Sex = 5
	selection.Sey = 1
	got := selection.GetSelectionString(Content)
	expected := "Hello, world!\nHow a"

	assert.Equal(t, expected, got, "GetSelectionString() = %v, expected %v", got, expected)
}

func TestMultipleLineSelection2(t *testing.T) {
	var Content = [][]rune{
		[]rune("Hello, world!"),
		[]rune("How are you doing today?"),
		[]rune("I hope you're doing well."),
	}

	var selection = Selection{}
	selection.Ssx = 0
	selection.Ssy = 0
	selection.Sex = 11
	selection.Sey = 1
	got := selection.GetSelectionString(Content)
	expected := "Hello, world!\nHow are you"

	assert.Equal(t, expected, got, "GetSelectionString() = %v, expected %v", got, expected)
}

func TestMultipleLineSelection3(t *testing.T) {
	var Content = [][]rune{
		[]rune("Hello, world!"),
		[]rune("How are you doing today?"),
		[]rune("I hope you're doing well."),
	}

	var selection = Selection{}
	selection.Ssx = 6
	selection.Ssy = 0
	selection.Sex = 23
	selection.Sey = 1
	got := selection.GetSelectionString(Content)
	expected := " world!\nHow are you doing today"

	assert.Equal(t, expected, got, "GetSelectionString() = %v, expected %v", got, expected)
}
