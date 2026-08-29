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

#include "editor.h"
#include "query.h"
#include "highlighter.h"

void UpdateLsp(const std::string &text, bool isOpen) {
	if (!lspclient.isReady) return;

	if (!isOpen) {
		e.lspver[e.absoluteFilePath]++;
		auto version = e.lspver[e.absoluteFilePath];
		DidChange(e.absoluteFilePath, text, version);
	} else {
		e.lspver[e.absoluteFilePath] = 1;
		DidOpen(e.absoluteFilePath, text);
	}
}

void FindTests(void) {
	e.tests.clear();
	bool alloc = TestQueryAlloc();
	if (!alloc) return;
	TestFinder(e.absoluteFilePath, &e.tests);
}

void GetColor(char ch, int *fg, int *bg, colorindexer_t indexer) {
	for (auto i : indexer.ranges) {
		if (i.startbyte <= indexer.counter && indexer.counter < i.endbyte) {
			*fg = i.color;
			break;
		}
	}

	if (e.__selection.IsUnderSelection(indexer.col, indexer.row)) {
		*bg = SELECTIONCOLOR;
	}

	if (ch == '\t' && e.x == 0) { // draw big cursor for tab
		if ((int)indexer.row == e.row && (int)indexer.col == e.col) {
			*bg = ACCENTCOLOR;
			*fg = -1;
		}
	}
}

void TabsSpaces(int fg, int bg, std::string *line, std::string *str) {
	// TODO: handling tab spaces when we write into the tab area:
	// we need to count the characters and move our tab forward
	// when it reaches the last character inside tab width
	for (int i = 0; i < e.langTabWidth; i++) {
		if (fg == -1)
			ColorizeTextChunk(0, bg, str);
		line->append(*str);
	}
}

void Colorize(char b, std::string *line, bool colorize, colorindexer_t indexer) {
	std::string str = " ";
	int fg = 0, bg = 0;

	if (colorize) {
		GetColor(b, &fg, &bg, indexer);

		if (fg == 0 && b != '\t') goto nocolor;
		if (fg == 0 && b == '\t') goto nocolortabs;
		if (fg > 0 && b != '\t') goto color;

		// coloring tabs if cursor is there
		TabsSpaces(fg, bg, line, &str);
		return;
color:
		str = std::string(1, b);
		ColorizeTextChunk(fg, bg, &str);
		line->append(str);
	} else {
nocolor:
		line->append(std::string(1, b));
		return;
nocolortabs:
		TabsSpaces(fg, bg, line, &str);
	}
}

void SetTerminalDims(int *rows, int *cols) {
	struct winsize size;
	ioctl(STDOUT_FILENO, TIOCGWINSZ, &size);
	*rows = size.ws_row;
	*cols = size.ws_col;
	LOG(INFO) << "refresh terminal dimensions " << *cols << "x" << *rows;
}

void SetFileAttributes(const std::string &file, std::string *content, bool isOpen) {
	char currentdir[PATH_MAX];
	getcwd(currentdir, sizeof(currentdir));
	e.cwd = currentdir;

	if (isOpen) {
		boost::filesystem::path p(file);
		e.filename = p.filename().string();
		e.absoluteFilePath = p.string();
		*content = ReadFileToString(e.absoluteFilePath);
	} else {
		std::string absFilePath;
		absFilePath.append(currentdir);
		absFilePath.append("/");
		absFilePath.append(file);
		e.filename = file;
		e.absoluteFilePath = absFilePath;;
		if (!SaveToFile(e.absoluteFilePath, e.code_str())) {
			LOG(ERROR) << "unable to save to file " << e.absoluteFilePath;
		}
	}
}

bool HandleFile(const std::string &filepath, bool isOpen) {
	EDGO_LOGGING_INIT;      // init google logging system
	EDGO_LOGGING_NO_STDERR; // make sure it doesn't print to stderr
	if (!e.isLogging)       // only log if enabled via env variable
		EDGO_LOGGING_SUPPRESS;

	std::string content;
	SetFileAttributes(filepath, &content, isOpen);
	std::string log_text = isOpen ? "open file:" : "new file:";
	LOG(INFO) << log_text << e.absoluteFilePath;

	// Set language and content
	e.lang = DetectLang(e.absoluteFilePath);
	e.config = GetConfig();
	e.code = rope<char>(content.c_str());
	e.langConf = e.config.langs[e.lang.name];
	e.langTabWidth = e.langConf.tabwidth;

	switch (e.lang.name.empty()) {
	case true:
		LOG(INFO) << "unknown language";
		return false;
	break;
	case false:
		LOG(INFO) << "new language is " << e.lang.name;
		StartLspClient(e.langConf.cmd, e.langConf.lsp);
		InitLspClient(e.cwd);
		SetTerminalDims(&e.ROWS, &e.COLUMNS);
		SetLang(e.lang.lang);
		Parse(content);
		SetTheme(e.config.theme);
	break;
	}

	// Opening file for LSP
	UpdateLsp(e.code_str(), true);
	FindTests();
	return true;
}

bool SaveFile(void) {
	if (!SaveToFile(e.absoluteFilePath, e.code_str())) {
		LOG(ERROR) << "unable to save to file " << e.absoluteFilePath;
		return false;
	}

	UpdateLsp(e.code_str(), false);
	FindTests();
	return true;
}
