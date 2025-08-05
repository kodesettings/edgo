package highlighter

import (
	"context"
	. "github.com/vipmax/edgo/internal/highlighter/langs"
	"github.com/gdamore/tcell"
	sitter "github.com/tree-sitter/go-tree-sitter"
	bash "github.com/tree-sitter/tree-sitter-bash/bindings/go"
	c "github.com/tree-sitter/tree-sitter-c/bindings/go"
	cpp "github.com/tree-sitter/tree-sitter-cpp/bindings/go"
	golang "github.com/tree-sitter/tree-sitter-go/bindings/go"
	html "github.com/tree-sitter/tree-sitter-html/bindings/go"
	css "github.com/tree-sitter/tree-sitter-css/bindings/go"
	java "github.com/tree-sitter/tree-sitter-java/bindings/go"
	javascript "github.com/tree-sitter/tree-sitter-javascript/bindings/go"
	python "github.com/tree-sitter/tree-sitter-python/bindings/go"
	rust "github.com/tree-sitter/tree-sitter-rust/bindings/go"
	. "gopkg.in/yaml.v2"
	"log"
	"os"
	"strings"
	"unicode/utf8"
)

type TreeSitterHighlighter struct {
	parser         *sitter.Parser
	tree           *sitter.Tree
	lines          []string
	lang           string
	language       *sitter.Language
	query          *sitter.Query
	colorsMap      map[string]string
	themePath      string
}

func NewTreeSitter() *TreeSitterHighlighter {
	parser := sitter.NewParser()

	return &TreeSitterHighlighter{
		parser: parser,
		tree:   nil,
	}
}


var defaultColors =
`
identifier: "#A5FCB6"
field_identifier: "#A5FCB6"
property_identifier: "#A5FCB6"
property: "#A5FCB6"
string: "#b1fce5"
keyword: "#a0a0a0"
constant: "#f6c99f"
number: "#f6c99f"
integer: "#f6c99f"
float: "#f6c99f"
variable: "#ffffff"
#variable.builtin: "#f992e6"
function: "#f6c99f"
function.call: "#f6c99f"
method: "#f6c99f"
comment: "#585858"
namespace: "#f6c99f"
type: "#f6c99f"
tag.attribute: "#c6a5fc"
tag: "#c6a5fc"
error: "#A5FCB6"

accent_color: "#f992e6"
accent_color2: "#A5FCB6"
`

func (h *TreeSitterHighlighter) ParseColor(colour string) int {
	return int(tcell.GetColor(colour))
}

func (h *TreeSitterHighlighter) SetTheme(themePath string) {
	yamlFile, err := os.ReadFile(themePath)
	if err != nil {
		//log.Println("Error reading theme YAML file: %v", err)
		yamlFile = []byte(defaultColors)
	}

	h.themePath = themePath

	err = Unmarshal(yamlFile, &h.colorsMap)
	if err != nil {
		log.Fatalf("Error unmarshaling YAML: %v", err)
	}

	if value, ok := h.colorsMap["accent_color"]; ok {
		AccentColor = h.ParseColor(value)
	}
	if value, ok := h.colorsMap["accent_color2"]; ok {
		AccentColor2 = h.ParseColor(value)
	}
	if value, ok := h.colorsMap["accent_color3"]; ok {
		AccentColor3 = h.ParseColor(value)
	}
}

func GetSitterLang(lang string) *sitter.Language {
	switch lang {
	case "bash": return sitter.NewLanguage(bash.Language())
	case "c": return sitter.NewLanguage(c.Language())
	case "c++": return sitter.NewLanguage(cpp.Language())
	case "cpp": return sitter.NewLanguage(cpp.Language())
	case "go": return sitter.NewLanguage(golang.Language())
	case "html": return sitter.NewLanguage(html.Language())
	case "css": return sitter.NewLanguage(css.Language())
	case "java":return sitter.NewLanguage(java.Language())
	case "javascript": return sitter.NewLanguage(javascript.Language())
	case "typescript": return sitter.NewLanguage(javascript.Language())
	case "python": return sitter.NewLanguage(python.Language())
	case "rust": return sitter.NewLanguage(rust.Language())
	default: return sitter.NewLanguage(css.Language()) // default css works on all file types
	}
}

func (h *TreeSitterHighlighter) SetLang(lang string) {
	if h.lang == lang { return }
	h.lang = lang
	h.language = GetSitterLang(lang)
	h.parser.SetLanguage(h.language)

	queryLang := MatchQueryLang(h.lang)
	q, err := sitter.NewQuery(h.language, queryLang)
	if err!= nil { panic(err) }
	h.query = q
}

func (h *TreeSitterHighlighter) matchExpression(expression string, fullexpression string) int {
	if value, ok := h.colorsMap[fullexpression]; ok { return  h.ParseColor(value) }
	if value, ok := h.colorsMap[expression]; ok { return  h.ParseColor(value) }
	return -1
}

/*
	comment for sitter.InputEdit
	The StartByte, OldEndByte, and NewEndByte parameters indicate the range of bytes you're modifying
	The StartPosition, OldEndPosition, and NewEndPosition parameters indicate the range of positions (line, column) affected by the edit.
*/

func (h *TreeSitterHighlighter) AddCharEdit(code *string, row int, col int, ch rune) {
	StartIndex := GetStartIndex(code, row, col)
	runeLen := uint(utf8.RuneLen(ch))

	editInput := sitter.InputEdit{
		StartByte: StartIndex,
		OldEndByte: StartIndex,
		NewEndByte: StartIndex + runeLen,
		StartPosition:  sitter.Point{Row: 0, Column: 0},
		OldEndPosition: sitter.Point{Row: 0, Column: 0},
		NewEndPosition: sitter.Point{Row: 0, Column: 0},
	}
	h.tree.Edit(&editInput)
	h.tree = h.parser.ParseCtx(context.Background(), []byte(*code), h.tree)
}

func (h *TreeSitterHighlighter) RemoveCharEdit(code *string, row int, col int, ch rune) {
	StartIndex := GetStartIndex(code, row, col)
	Row := uint(row); Column := uint(col)
	runeLen := uint(utf8.RuneLen(ch))

	editInput := sitter.InputEdit{
		StartByte: StartIndex,
		OldEndByte: StartIndex + runeLen,
		NewEndByte: StartIndex,
		StartPosition:  sitter.Point{Row: Row, Column: Column},
		OldEndPosition: sitter.Point{Row: Row, Column: Column + runeLen},
		NewEndPosition: sitter.Point{Row: Row, Column: Column},
	}
	h.tree.Edit(&editInput)
	h.tree = h.parser.ParseCtx(context.Background(), []byte(*code), h.tree)
}

func (h *TreeSitterHighlighter) UpdateCharsEdit(code *string, row int, col int, nrow int, ncol int) {
	StartIndex1 := GetStartIndex(code, row, col)
	StartIndex2 := GetStartIndex(code, nrow, ncol)
	Row1 := uint(row); Column1 := uint(col)
	Row2 := uint(nrow); Column2 := uint(ncol)

	// TODO: fixing this
	editInput := sitter.InputEdit{
		StartByte: StartIndex1,
		OldEndByte: StartIndex2,
		NewEndByte: StartIndex1,
		StartPosition:  sitter.Point{Row: Row1, Column: Column1},
		OldEndPosition: sitter.Point{Row: Row2, Column: Column2},
		NewEndPosition: sitter.Point{Row: Row1, Column: Column1},
	}
	h.tree.Edit(&editInput)
	h.tree = h.parser.ParseCtx(context.Background(), []byte(*code), nil)
}

func GetStartIndex(code *string, row int, col int) uint {
	r, c, startIndex := 0, 0, 0
	for _, char := range *code {
		runeLen := utf8.RuneLen(char)

		if r == row && c == col {
			break
		}
		startIndex += runeLen

		if char == '\n' {
			r++
			c = 0
		} else {
			c++
		}
	}
	return uint(startIndex)
}

func Use(vals ...interface{}) { }

func (h *TreeSitterHighlighter) Parse(code *string) {
	tree := h.parser.ParseCtx(context.Background(), []byte(*code), h.tree)
	h.tree = tree
}

func (h *TreeSitterHighlighter) ColorRanges(from, to int, codeBytes []byte) []ColoredByteRange {

	queryCursor := sitter.NewQueryCursor()
	matches := queryCursor.Matches(h.query, h.tree.RootNode(), codeBytes)
	queryCursor.SetPointRange(
		sitter.Point{Row: uint(from), Column: 0},
		sitter.Point{Row: uint(to), Column: 0},
	)

	colors := make([]ColoredByteRange, 0)

	for match := matches.Next(); match != nil; match = matches.Next() {
		for _, capture := range match.Captures {
			name := h.query.CaptureNames()[capture.Index]
			split := strings.Split(name, ".")
			color := h.matchExpression(split[0], name)

			contentstr := string(codeBytes[capture.Node.StartByte():capture.Node.EndByte()]); Use(contentstr) // for debug

			if strings.Contains(name, "injection") {
				// We don't colorize embedded content for different languages in the editor.
				// Only languages that can be identified for the entire source are colorized.
				// This usually means excluding markdown or html files or any documentation
				// related content.
				continue
			}

			colors = append(colors, ColoredByteRange{
				StartByte: int(capture.Node.StartByte()),
				EndByte:   int(capture.Node.EndByte()),
				Color:     color,
			})
		}
	}
	return colors
}

type ColoredByteRange struct {
	StartByte int
	EndByte   int
	Color     int
}

func (h *TreeSitterHighlighter) GetTree() *sitter.Tree {
	return h.tree
}

func (h *TreeSitterHighlighter) GetLang() *sitter.Language {
	return h.language
}

func (h *TreeSitterHighlighter) GetLangStr() string {
	return h.lang
}

type NodeRange struct {
	Ssy int
	Ssx int
	Sey int
	Sex int
}

type Path struct {
	Atx int
	Aty int
	Nodes []NodeRange
	Current int
}

func (p *Path) CurrentNode() NodeRange {
	return p.Nodes[p.Current]
}

func (p *Path) Next() NodeRange {
	p.Current += 1
	if p.Current >= len(p.Nodes) { p.Current = len(p.Nodes) - 1 }
	return p.Nodes[p.Current]
}

func (p *Path) Prev() NodeRange {
	p.Current -= 1
	if p.Current < 0 {
		p.Current = 0;
		return NodeRange{p.Aty,p.Atx,p.Aty,p.Atx}
	}
	return p.Nodes[p.Current]
}

func (h *TreeSitterHighlighter) GetNodePathAt(StartPointRow int, StartPointColumn int,
	EndPointRow int, EndPointColumn int) Path {

	rootNode := h.tree.RootNode()
	node := rootNode.NamedDescendantForPointRange(
		sitter.Point{Row: uint(StartPointRow), Column: uint(StartPointColumn)},
		sitter.Point{Row: uint(EndPointRow), Column: uint(EndPointColumn)},
	)

	path := Path{Aty: StartPointRow, Atx: StartPointColumn}

	for node != nil {
		r := NodeRange{int(node.StartPosition().Row),
			int(node.StartPosition().Column),
			int(node.EndPosition().Row),
			int(node.EndPosition().Column),
		}
		path.Nodes = append(path.Nodes, r)
		node = node.Parent()
	}

	return path
}

func (h *TreeSitterHighlighter) GetNodeAt(StartPointRow int, StartPointColumn int,
	EndPointRow int, EndPointColumn int) (string, NodeRange) {
	rootNode := h.tree.RootNode()
	node := rootNode.NamedDescendantForPointRange(
		sitter.Point{Row: uint(StartPointRow), Column: uint(StartPointColumn)},
		sitter.Point{Row: uint(EndPointRow), Column: uint(EndPointColumn)},
	)

	return node.Kind(), NodeRange{int(node.StartPosition().Row),
		int(node.StartPosition().Column),
		int(node.EndPosition().Row),
		int(node.EndPosition().Column),
	}
}
