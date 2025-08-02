package selection

import (
	"testing"
	. "github.com/vipmax/edgo/internal/utils"
	assert "github.com/stretchr/testify/assert"
)

func TestNoSelection(t *testing.T) {
	var lines = []Line{
		Line{[]rune("Hello, world!")},
		Line{[]rune("How are you doing today?")},
		Line{[]rune("I hope you're doing well.")},
	}

	var selection = Selection{}
	selection.Ssx = 0
	selection.Ssy = 0
	selection.Sex = 0
	selection.Sey = 0
	got := selection.GetSelectionString(lines)
	expected := ""

	assert.Equal(t, expected, got, "GetSelectionString() = %v, expected %v", got, expected)
}

func TestSingleCharacterSelection(t *testing.T) {
	var lines = []Line{
		Line{[]rune("Hello, world!")},
		Line{[]rune("How are you doing today?")},
		Line{[]rune("I hope you're doing well.")},
	}

	var selection = Selection{}
	selection.Ssx = 0
	selection.Ssy = 0
	selection.Sex = 1
	selection.Sey = 0
	got := selection.GetSelectionString(lines)
	expected := "H"

	assert.Equal(t, expected, got, "GetSelectionString() = %v, expected %v", got, expected)
}

func TestMultipleCharacterSelection(t *testing.T) {
	var lines = []Line{
		Line{[]rune("Hello, world!")},
		Line{[]rune("How are you doing today?")},
		Line{[]rune("I hope you're doing well.")},
	}

	var selection = Selection{}
	selection.Ssx = 0
	selection.Ssy = 0
	selection.Sex = 5
	selection.Sey = 0
	got := selection.GetSelectionString(lines)
	expected := "Hello"

	assert.Equal(t, expected, got, "GetSelectionString() = %v, expected %v", got, expected)
}

func TestMultipleLineSelection1(t *testing.T) {
	var lines = []Line{
		Line{[]rune("Hello, world!")},
		Line{[]rune("How are you doing today?")},
		Line{[]rune("I hope you're doing well.")},
	}

	var selection = Selection{}
	selection.Ssx = 0
	selection.Ssy = 0
	selection.Sex = 5
	selection.Sey = 1
	got := selection.GetSelectionString(lines)
	expected := "Hello, world!\nHow a"

	assert.Equal(t, expected, got, "GetSelectionString() = %v, expected %v", got, expected)
}

func TestMultipleLineSelection2(t *testing.T) {
	var lines = []Line{
		Line{[]rune("Hello, world!")},
		Line{[]rune("How are you doing today?")},
		Line{[]rune("I hope you're doing well.")},
	}

	var selection = Selection{}
	selection.Ssx = 0
	selection.Ssy = 0
	selection.Sex = 11
	selection.Sey = 1
	got := selection.GetSelectionString(lines)
	expected := "Hello, world!\nHow are you"

	assert.Equal(t, expected, got, "GetSelectionString() = %v, expected %v", got, expected)
}

func TestMultipleLineSelection3(t *testing.T) {
	var lines = []Line{
		Line{[]rune("Hello, world!")},
		Line{[]rune("How are you doing today?")},
		Line{[]rune("I hope you're doing well.")},
	}

	var selection = Selection{}
	selection.Ssx = 6
	selection.Ssy = 0
	selection.Sex = 23
	selection.Sey = 1
	got := selection.GetSelectionString(lines)
	expected := " world!\nHow are you doing today"

	assert.Equal(t, expected, got, "GetSelectionString() = %v, expected %v", got, expected)
}
