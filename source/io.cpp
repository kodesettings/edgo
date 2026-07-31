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

void UpdateLsp(bool isOpen, std::string text) {
	if (!lspclient.isReady) return;

	if (isOpen) {
		e.lspver[e.absoluteFilePath]++;
		auto version = e.lspver[e.absoluteFilePath];
		DidChange(e.absoluteFilePath, &text, version);
	} else {
		e.lspver[e.absoluteFilePath] = 1;
		DidOpen(e.absoluteFilePath, &text);
	}
}

void FindTests(void) {
	e.tests.clear();
	bool alloc = TestQueryAlloc();
	if (!alloc) return;
	TestFinder(e.absoluteFilePath, &e.tests);
}

uint32_t GetColor(char ch, int col, int row, uint32_t *bytesCounter) {
	coloredbyterange_v coloredByteRanges;
	coloredByteRanges = ColorRanges(e.y, e.y + e.TERMINAL_HEIGHT);
	uint32_t style;

	for (auto i : coloredByteRanges) {
		if (i.startbyte <= *bytesCounter && *bytesCounter < i.endbyte) {
			style = i.color;
			break;
		}
	}

	if (e.__selection.IsUnderSelection(col, row)) {
		style = SELECTIONCOLOR;
	}

	if (ch == '\t' && e.x == 0) { // draw big cursor for tab
		if (row == e.row && col == e.col) {
			style = ACCENTCOLOR;
		}
	}

	bytesCounter += sizeof(ch);
	return style;
}

bool OpenFile(const std::string &filepath) {
	// finding current working directory
	char currentdir[PATH_MAX];
	getcwd(currentdir, sizeof(currentdir));
	e.cwd = currentdir;

	// reading file content
	std::string content;
	content = ReadFileToString(filepath);

	boost::filesystem::path p(filepath);
	e.filename = p.filename().string();
	e.absoluteFilePath = filepath;

	LOG(INFO) << "open file:" << e.absoluteFilePath;

	e.lang = DetectLang(e.absoluteFilePath);
	switch (e.lang.empty()) {
	case true:
		LOG(INFO) << "unknown language";
		return false;
	break;
	case false:
		LOG(INFO) << "new language is " << e.lang;
		config_t config = GetConfig();
		lang_t lsp = config.langs[e.lang];
		StartLspClient(lsp.cmd, lsp.lsp);
		InitLspClient(e.cwd);

		// TODO: set language:
		// SetLang(JAVASCRIPT);
		Parse(content);
		SetTheme(e.config.theme);
	break;
	}

	// Opening file for LSP
	UpdateLsp(true, e.code_str());
	FindTests();
	return true;
}

void SaveFile(void) {
	if (!SaveToFile(e.absoluteFilePath, e.code_str())) {
		LOG(ERROR) << "unable to save to file " << e.absoluteFilePath;
	}

	UpdateLsp(false, e.code_str());
	FindTests();
}
