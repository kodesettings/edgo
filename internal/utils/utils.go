package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"bytes"
)

func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func Min(x, y int) int {
	if x <= y {
		return x
	}
	return y
}

func maxMany(nums ...int) int {
	if len(nums) == 0 {
		panic("max: no arguments provided")
	}
	maxValue := nums[0]
	for _, num := range nums[1:] {
		if num > maxValue {
			maxValue = num
		}
	}
	return maxValue
}

func MinMany(nums ...int) int {
	if len(nums) == 0 {
		panic("min: no arguments provided")
	}
	minValue := nums[0]
	for _, num := range nums[1:] {
		if num < minValue {
			minValue = num
		}
	}
	return minValue
}

func InsertTo[T any](a []T, index int, value T) []T {
	n := len(a)
	if index < 0 {
		index = (index%n + n) % n
	}
	switch {
	case index == n: // nil or empty slice or after last element
		return append(a, value)

	case index < n: // index < len(a)
		a = append(a[:index+1], a[index:]...)
		a[index] = value
		return a

	case index < cap(a): // index > len(a)
		a = a[:index+1]
		var zero T
		for i := n; i < index; i++ {
			a[i] = zero
		}
		a[index] = value
		return a

	default:
		b := make([]T, index+1) // malloc
		if n > 0 {
			copy(b, a)
		}
		b[index] = value
		return b
	}
}

var matched = []rune{
	' ', '.', ',', '=', '+', '-', '[', '(', '{', ']', ')', '}', '"', ':', '&', '?','!',';','\t',
	'/','<','>',
}

func FindNextWord(chars []rune, from int) int {
	// Find the next word index after the specified index
	for i := from; i < len(chars); i++ {
		if Contains(matched, chars[i]) {
			return i
		}
	}

	return len(chars)
}

func FindPrevWord(chars []rune, from int) int {
	// Find the previous word index before the specified index
	for i := from - 1; i >= 0; i-- {
		if Contains(matched, chars[i]) {
			return i + 1
		}
	}

	return 0
}

func Contains[T comparable](slice []T, e T) bool {
	for _, val := range slice {
		if val == e {
			return true
		}
	}
	return false
}

func Remove[T any](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}

func MaxString(arr []string) int {
	maxLength := 0
	for _, str := range arr {
		if len(str) > maxLength {
			maxLength = len(str)
		}
	}
	return maxLength
}

func ReadFileToString(filePath string) (string, error) {
	filecontent, err := os.ReadFile(filePath)
	if err != nil { return "", err }
	return string(filecontent), nil
}

func GetLinesArrayFromData(data []byte, lineNum int) []Line {
	if lineNum == 0 { lineNum = CountNewLines(data) }
	var lines []Line = make([]Line, lineNum)
	var row int = 0
	for _, b := range data {
		if b == '\n' {
			row++
			continue
		} else if row == lineNum {
			break
		} else {
			lines[row].Buf = append(lines[row].Buf, rune(b))
		}
	}
	return lines
}

func CountNewLines(data []byte) int {
	newlineCount := 0
	for _, b := range data {
		if b == '\n' {
			newlineCount++
		}
	}
	return newlineCount
}

func CountTabs(str []byte, stopIndex int) int {
	if stopIndex == 0 { return 0 }

	count := 0
	for i, char := range str {
		if i >= stopIndex { break }
		if char == '\t' { count++ } else { break }
	}
	return count
}

func CountTabsTo(str []byte, stopIndex int) int {
	if stopIndex == 0 { return 0 }

	count := 0
	for i, char := range str {
		if i >= stopIndex { break }
		if char == '\t' { count++ }
	}
	return count
}

func CountSpaces(str []byte, stopIndex int) int {
	if stopIndex == 0 { return 0 }

	count := 0
	for i, char := range str {
		if i >= stopIndex { break }
		if char == ' ' { count++ } else { break }
	}
	return count
}

func FormatText(left, right string, maxWidth int) string {
	left = fmt.Sprintf("%-*s", maxWidth, left)
	right = fmt.Sprintf("%s",  right)
	return fmt.Sprintf("%s %s", left, right)
}

func GetFileSize(filename string) int64 {
	file, err := os.Open(filename) // replace with your file name
	if err != nil { return 0 }
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil { return 0 }

	fileSize := fileInfo.Size() // get the size in bytes
	return fileSize
}

func CenterNumber(brw int, width int) string {
	lineNumber := strconv.Itoa(brw )
	padding := width - len(lineNumber)
	leftPad := fmt.Sprintf("%*s", padding/2, "")
	rightPad := fmt.Sprintf("%*s", padding-(padding/2), "")
	lineNumber = leftPad + lineNumber + rightPad
	return lineNumber
}

type Set map[int]struct{}

// a value to the set.
func (this Set) Add(value int) { this[value] = struct{}{} }
func (s Set) Delete(value int) { delete(s, value) }
func (s Set) Contains(value int) bool {
	_, exists := s[value]
	return exists
}
// returns all keys in the set, sorted.
func (this Set) GetKeys() []int {
	keys := make([]int, 0, len(this))
	for key := range this { keys = append(keys, key) }
	sort.Ints(keys) // Sort the keys
	return keys
}

func PadLeft(str string, length int) string {
	format := fmt.Sprintf("%%%ds", length)
	return fmt.Sprintf(format, str)
}

func GetFirstLines(s []byte, lineNum int) ([]byte, error) {
	scanner := bufio.NewScanner(bytes.NewReader(s))
	count := 0
	var buffer bytes.Buffer
	for scanner.Scan() {
		buffer.WriteString(scanner.Text())
		buffer.WriteString("\n")
		count++
		if count == lineNum {
			break
		}
	}

	if scanner.Err() != nil {
		// handle error.
		return []byte{}, scanner.Err()
	}

	return buffer.Bytes(), nil
}

func LineOffset(text []byte, lineNum int) int {
	index := 0
	for i, ch := range text {
		if ch == '\n' {
			index++
		}

		if index == lineNum {
			return i // start of next line
		}
	}

	return len(text) // end of text
}

func IsIgnored(path string, ignorePatterns []string) bool {
	for _, pattern := range ignorePatterns {
		match, err := filepath.Match(pattern, filepath.Base(path))
		if err != nil { log.Println("Invalid pattern:", pattern); continue }
		if match { return true }
	}
	return false
}

func IsMatchExt(path string, ignoreExts []string) bool {
	for _, ignoreExt := range ignoreExts {
		ext := filepath.Ext(filepath.Base(path))
		if ext == ignoreExt {
			return true
		}
	}
	return false
}

// Inner struct to contain a line of text
type Line struct {
	Buf []rune
}

// Function to remove leading tabs and spaces
func RemoveLeadingTabsSpaces(s string) string {
	for i, c := range s {
		if c != '\t' && c != ' ' {
			return s[i:]
		}
	}
	return ""
}
