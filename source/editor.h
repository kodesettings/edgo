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
#include "search.h"
#include "highlighter.h"
#include "query.h"
// NOTE: including gnu rope specification, it seems although
// this is legacy  this is still the best one
#include <ext/rope>
// using the namespace for simplicity
using __gnu_cxx::rope;

typedef struct editor_t {
	int COLUMNS;                        // terminal size columns
	int ROWS;                           // terminal size rows
	int LINES_WIDTH;                    // draw file lines number

	int TERMINAL_HEIGHT;
	int TERMINAL_WIDHT;

	int row;                            // cursor position row
	int col;                            // cursor position column
	int y;                              // row offset for scrolling
	int x;                              // col offset for scrolling

	line_v lines;                       // Lines of text characters

	std::string lang;                   // current file language
	config_t config;                    // config, lsp, tabs, comments, etc
	lang_t langConf;                    // current lang conf
	int langTabWidth;                   // current lang tabs indentation  '\t' -> "    "

	selection __selection;              // selection

	std::vector<editoperation_t> undo;  // stack for undo operations
	std::vector<editoperation_t> redo;  // stack for redo operations

	std::string cwd;                    // current dir
	std::string inputFile;              // exact user input
	std::string filename;               // current file name
	std::string absoluteFilePath;       // current file name and directory
	bool (*isContentChanged)(void);     // shows * if file is changed
	bool isColorize;                    // colorize text is true by default
	bool isOverlay;                     // true if overlay is active (completion, hover, errors...)
	bool isStartupScreen;               // true if no file was opened

	bool isContentSearch;
	std::string searchPattern;          // pattern for search in a buffer
	searchresult_v searchResults;
	int searchResultIndex;

	// registering callback for content change
	editor_t() { isContentChanged = { /* registering dummy callback */ };}
#ifdef PROCESS_PANEL
	// process panel vars
	int processPanelHeight;
	int processPanelWidth;
	line_v processContent;
	int processPanelScroll;
	int processPanelHScroll;
	bool isProcessPanelMoving;
	bool isProcessPanelFocused;
	process *process;
	int processPanelSpacing;
	int processPanelCursorX;
	int processPanelCursorY;
	selection processPanelSelection;
	bool isProcessPanelSearch;
	std::string processPanelSearchPattern; // pattern for search in a buffer
	searchresult_v processPanelSearchResults;
	int processPanelSearchResultIndex;
#endif
	std::map<std::string, lspclient_t> lsp2lang;
	std::map<std::string, int> lspver; // version number sequencing

	treesitterhighlighter_t h;
	rope<char> code; // storing a rope node element
	std::string code_str(void) { return std::string(code.begin(), code.end()); }

	std::map<int, testdata_t> tests;
} editor_t;

extern editor_t e;

// character manipulation routines
void AddCharacter(char ch);
void InsertCharacter(int line, int pos, char ch);
void InsertString(int line, int pos, std::string linestring);
void DeleteCharacter(int line, int pos);
void ReplaceString(int line, int from, int end, std::string instext);
void ShiftWithTabsToRight(int line, int pos, std::set<int> selectedLines);
bool MaybeAddPair(int line, int pos, char ch, char *ret);

// cursor routines
void Focus(void);
void FocusCenter(void);
#ifdef PROCESS_PANEL
void FocusProcessPanel(void);
#endif
// clipboard routines
void set_update_parameters(bool changed);
void OnCopy(void);
void OnPaste(void);
void Cut(bool isCopySelected);
void Duplicate(void);
void OnUndo(void);
void OnRedo(void);

// features v01 routines
void OnCommentLine(void);
void OnSwapLinesUp(void);
void OnSwapLinesDown(void);

// io routines
void UpdateLsp(bool isOpen, std::string text);
void FindTests(void);

// keyboard routines
void OnDown(bool isPaging);
void OnUp(bool isPaging);
void OnLeft(void);
void OnRight(void);
void GoTop(void);
void GoBottom(void);
void OnScrollUp(void);
void OnScrollDown(void);
void OnEnter(void);
void OnDelete(void);
void OnTab(void);
void OnBackTab(void);

#endif // _EDITOR_H_
