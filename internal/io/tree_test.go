package io

import (
	"fmt"
	"os"
	"testing"
	assert "github.com/stretchr/testify/assert"
)

func TestProcessDirectory(t *testing.T) {

	// Get the current directory
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Failed to get current directory:", err)
		return
	}

	// Recursively process the directory
	fileInfo, err := ReadDirTree(dir, "", true,0)
	assert.NotNil(t, err, "failed to process directory: %s", err)

	// Print the result
	//fmt.Printf("%+v\n", fileInfo)

	PrintTree(fileInfo, 0)
}

func TestTreeSize(t *testing.T) {
	dir, _ := os.Getwd()
	tree, _ := ReadDirTree(dir,"", true,0)
	size := TreeSize(tree, 0)
	fmt.Println("size", size)
	assert.Equal(t, true, size > 0, "size is zero")
}

func TestGetSelected(t *testing.T) {
	dir, _ := os.Getwd()
	tree, _ := ReadDirTree(dir, "", true,0)
	found, fi := GetSelected(tree, 13)
	fmt.Println("selected", found, fi)
	assert.NotNil(t, found, "foudn is nil")
	assert.NotNil(t, fi, "fi is nil")
}
