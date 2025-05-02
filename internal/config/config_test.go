package config

import (
	"fmt"
	"testing"
	assert "github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	// Read conf
	conf := GetConfig()

	// Print config
	for _, lang := range conf.Langs {
		fmt.Printf("Name: %s, Lsp: %s, Comment: %s, Tab Width: %d\n",
			lang.Name, lang.Lsp, lang.Comment, lang.TabWidth,
		)
	}

	golang, ok := conf.Langs["go"]

	assert.Equal(t, ok, true, "Go lang go not found")
	assert.Equal(t, golang.Lsp, "gopls", "Go lang lsp should be gopls")
}
