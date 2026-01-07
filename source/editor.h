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

#ifndef _EDITOR_H_
#define _EDITOR_H_

#include "utils.h"
#include "config.h"
#include "selection.h"
#include "operations.h"
#include "lsp_client.h"
// NOTE: including gnu rope specification, it seems although
// this is legacy  this is still the best one
#include <ext/rope>
// using the namespace for simplicity
using __gnu_cxx::rope;

typedef struct {
	int COLUMNS;     // terminal size columns
	int ROWS;        // terminal size rows
	int LINES_WIDTH; // draw file lines number

	int TERMINAL_HEIGHT;
	int TERMINAL_WIDHT;

	int row; // cursor position row
	int col; // cursor position column
	int y;   // row offset for scrolling
	int x;   // col offset for scrolling

	line_v lines; // Lines of text characters

	std::string lang; // current file language
	config_t config; // config, lsp, tabs, comments, etc
	lang_t langConf; // current lang conf
	int langTabWidth; // current lang tabs indentation  '\t' -> "    "

	selection __selection; // selection

	EditOperation undo; // stack for undo operations
	EditOperation redo; // stack for redo operations

	std::string cwd;   // current dir
	std::string inputFile; // exact user input
	std::string filename; // current file name
	std::string absoluteFilePath; // current file name and directory
	bool isContentChanged; // shows * if file is changed
	bool isColorize; // colorize text is true by default
	bool update;  // for Screen updates,  if false it will not draw
	bool isOverlay; // true if overlay is active (completion, hover, errors...)
	bool isStartupScreen; // true if no file was opened

	bool isContentSearch;
	std::string searchPattern; // pattern for search in a buffer
//	searchresult_v searchResults;
	int searchResultIndex;

	// process panel vars
	int processPanelHeight;
	int processPanelWidth;
	line_v processContent;
	int processPanelScroll;
	int processPanelHScroll;
	bool isProcessPanelMoving;
	bool isProcessPanelFocused;
//	process *process;
	int processPanelSpacing;
	int processPanelCursorX;
	int processPanelCursorY;
	selection processPanelSelection;
	bool isProcessPanelSearch;
	std::string processPanelSearchPattern; // pattern for search in a buffer
//	searchresult_v processPanelSearchResults;
	int processPanelSearchResultIndex;

	std::map<std::string, lspclient_t> lsp2lang;
	std::map<std::string, int> lspver; // version number sequencing

//	treesitterhighlighter *treeSitterHighlighter;
	rope<char> code; // storing a rope node element

//	std::map<int, testdata_t> tests;
//	testfinder testFinder;
//	test       test;

//	path *treePath;
} editor_t;

#endif // _EDITOR_H_
