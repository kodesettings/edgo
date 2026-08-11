/**
    Copyright (C) 2023 - 2026, edgo authors

    This program is free software; you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation; either version 2 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License along
    with this program; if not, see <https://www.gnu.org/licenses/>.
*/

#ifndef _HIGHLIGHTER_H_
#define _HIGHLIGHTER_H_

#include "langs.h"
#include "themes.h"

// tree-sitter language bindings
#include <tree_sitter/tree-sitter-bash.h>
#include <tree_sitter/tree-sitter-c.h>
#include <tree_sitter/tree-sitter-cpp.h>
#include <tree_sitter/tree-sitter-go.h>
#include <tree_sitter/tree-sitter-html.h>
#include <tree_sitter/tree-sitter-css.h>
#include <tree_sitter/tree-sitter-java.h>
#include <tree_sitter/tree-sitter-javascript.h>
#include <tree_sitter/tree-sitter-python.h>
#include <tree_sitter/tree-sitter-rust.h>

// base tree-sitter API
#include <tree_sitter/api.h>

//--------------------------------------------------------------------------------

typedef struct {
	struct TSParser       *parser;
	struct TSTree         *tree;
	std::string           lines;
	enum languages        lang;
	struct TSLanguage     *language;
	struct TSQuery        *query;
	struct TSQuery        *testquery;
	std::map<std::string, std::string> colorsmap;
	std::string           themepath;
} treesitterhighlighter_t;

extern treesitterhighlighter_t h;

// Hardcoded colors that are fixed in style
#define SELECTIONCOLOR 246 // gray
#define HIGHLIGHTCOLOR 238 // gray
#define OVERLAYCOLOR -1 // transparent
#define ACCENTCOLOR 303 // pink
#define ACCENTCOLOR2 = 30 // aqua
#define ACCENTCOLOR3 -1 // aqua

using colormap = std::map<std::string, std::string>;

//
// This is the theme map that is matching config parameters
// with the provided theme maps in themes.h file.
//
// If you add any new theme into themes.h, they have
// to be added here as well.
//
static std::map<std::string, colormap> themes = {
	{"darcula", THEME_DARCULA},
	{"edgo-new", THEME_EDGO_NEW},
	{"edgo", THEME_EDGO},
	{"material-darker", THEME_MATERIAL_DARKER},
	{"vesper", THEME_VESPER}
};

inline treesitterhighlighter_t NewTreeSitter(void) {
	auto parser = ts_parser_new();

	return treesitterhighlighter_t {
		parser: parser,
		tree: nullptr,
	};
}

typedef struct {
	uint32_t startbyte;
	uint32_t endbyte;
	uint32_t color;
} coloredbyterange_t;

inline uint32_t GetColorCode(std::string colorcode) {
	const char* str = colorcode.substr(1).data();
	return std::strtol(str, NULL, 16); // converting hex to int
}

inline void ColorToRGB(uint32_t color, int *r, int *g, int *b) {
	*r = (color >> 16) & 0xFF;
	*g = (color >> 8) & 0xFF;
	*b = color & 0xFF;
}

langsbyname_t DetectLang(const std::string &filepath);
void SetTheme(const std::string &theme);
void GetSitterLang(enum languages lang, bool *parsed);
void SetLang(enum languages lang);
void InsertTextEdit(const std::string &code, int offset, int length);
void RemoveTextEdit(const std::string &code, int offset, int length);
void Parse(const std::string &code);

#define APPLY_HIGHLIGHTER { SetLang(CSS); Parse(e.code_str()); }

//--------------------------------------------------------------------------------

typedef std::vector<coloredbyterange_t> coloredbyterange_v;
coloredbyterange_v ColorRanges(const int from, const int to);

typedef struct {
	coloredbyterange_v ranges;
	uint32_t counter;
} colorindexer_t;

//--------------------------------------------------------------------------------
// treesitterhighlighter_t getters

inline TSTree *GetTree(void) {
	return h.tree;
}

inline TSLanguage *GetLang(void) {
	return h.language;
}

inline enum languages GetLangStr(void) {
	return h.lang;
}

typedef struct {
	uint32_t ssy;
	uint32_t ssx;
	uint32_t sey;
	uint32_t sex;
} noderange_t;

typedef struct {
	uint32_t         atx;
	uint32_t         aty;
	std::vector<noderange_t> nodes;
	uint32_t         current;
} path_t;

extern path_t p;

//--------------------------------------------------------------------------------
// node manipulators

inline noderange_t CurrentNode(void) {
	return p.nodes[p.current];
}

inline noderange_t next(void) {
	p.current++;
	if (p.current >= p.nodes.size()) { p.current = p.nodes.size() - 1; }
	return p.nodes[p.current];
}

inline noderange_t prev(void) {
	p.current--;
	if (p.current < 0) {
		p.current = 0;
		return noderange_t{p.aty, p.atx, p.aty, p.atx};
	}
	return p.nodes[p.current];
}

path_t GetNodePathAt(int startPointRow, int startPointColumn, int endPointRow, int endPointColumn);
noderange_t GetNodeAt(int startPointRow, int startPointColumn, int endPointRow, int endPointColumn);

#endif // _HIGHLIGHTER_H_
