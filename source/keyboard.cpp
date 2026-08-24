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
#include "operations.h"

void OnDown(bool isPaging) {
	int numberOfLines;

	if (!isPaging) {
		numberOfLines = 1; // normal scroll
	} else {
		numberOfLines = e.ROWS / 3 < 30 ? 30 : e.ROWS / 3;
	}

	if (e.lines.empty()) {
		return;
	}

	if (e.row + numberOfLines >= (int)e.lines.size()) {
		if (!isPaging) {
			e.y = e.row - e.ROWS + 1;
			if (e.y < 0) e.y = 0; // reset
			return;
		} else {
			numberOfLines = e.lines.size() - e.row - 1;
		}
	}

	e.row += numberOfLines;
	if (e.col > (int)e.lines[e.row].buf.size()) {
		e.col = e.lines[e.row].buf.size();  // fit to e.lines
	}

	if (e.row < e.y) {
		e.y = e.row;
	} else if (e.row >= e.y + e.ROWS) {
		e.y = e.row - e.ROWS + 1;
	}
}

void OnUp(bool isPaging) {
	int numberOfLines;

	if (!isPaging) {
		numberOfLines = 1; // normal scroll
	} else {
		numberOfLines = e.ROWS / 3 < 30 ? 30 : e.ROWS / 3;
	}

	if (e.lines.size() == 0) {
		return;
	}

	if (e.row == 0) {
		e.y = 0;
		return;
	}

	if (e.row - numberOfLines <= 0) {
		e.row = 0;
	} else {
		e.row -= numberOfLines;
	}

	if (e.col > (int)e.lines[e.row].buf.size()) {
		e.col = e.lines[e.row].buf.size(); // fit to e.lines
	}

	if (e.row < e.y) {
		e.y = e.row;
	} else if (e.row > e.y + e.ROWS) {
		e.y = e.row - e.ROWS + 1;
	}
}

void OnLeft(void) {
	if (e.lines.empty()) {
		return;
	}

	if (e.col > 0) {
		e.col--;
	} else if (e.row > 0) {
		e.row--;
		e.col = e.lines[e.row].buf.size(); // fit to e.Lines
		if (e.row < e.y) {
			e.y = e.row;
		}
	}
}

void OnRight(void) {
	if (e.lines.empty()) {
		return;
	}

	if (e.col < (int)e.lines[e.row].buf.size()) {
		e.col++;
	} else if (e.row < (int)e.lines.size() - 1) {
		e.row++;
		e.col = 0;
		if (e.row > e.y + e.ROWS) {
			e.y++;
		}
	}
}

void GoTop(void) {
	e.row = 0;
	e.col = 0;
	e.x = 0;
	e.y = 0;
}

void GoBottom(void) {
	if (e.lines.empty()) {
		return;
	} else {
		e.row = e.lines.size() - 1; e.col = 0;
		e.x = 0;

		if (e.row > e.TERMINAL_HEIGHT) {
			FocusCenter();
		}

		OnDown(false);
	}
}

void OnScrollUp(void) {
	if (e.lines.empty()) {
		return;
	}

	if (e.y == 0) {
		return;
	}

	e.y--;
}

void OnScrollDown(void) {
	if (e.lines.empty()) {
		return;
	}

	if (e.y + e.ROWS >= (int)e.lines.size()) {
		return;
	}

	e.y++;
}

void OnEnter(void) {
	if (e.__selection.IsSelectionNonEmpty()) {
		Cut(false);
	}

	InsertCharacter(e.row, e.col, '\n');
	e.row++;
	e.col = 0;

	if (e.row - e.y == e.ROWS) {
		OnScrollDown();
	}

	if (!e.redo.empty()) {
		e.redo.clear();
	}

	Focus();
	UpdateLsp(e.code_str(), false);
	FindTests();
}

void OnDelete(void) {
	if (e.__selection.IsSelectionNonEmpty()) {
		Cut(false);
		return;
	}

	DeleteCharacter(e.row, e.col);
	if (e.col > 0) {
		e.col--;
	} else if (e.row > 0) {
		e.row--;
	}

	if (!e.redo.empty())
		e.redo.clear();

	Focus();
	UpdateLsp(e.code_str(), false);
	FindTests();
}

void OnTab(void) {
	auto selectedLines = e.__selection.GetSelectedLines(e.code_str());

	if (selectedLines.empty()) {
		InsertCharacter(e.row, e.col, '\t');
		e.col++;
	} else  {
		ShiftWithTabsToRight(e.row, e.col, selectedLines);
	}

	if (!e.redo.empty())
		e.redo.clear();

	Focus();
	UpdateLsp(e.code_str(), false);
	FindTests();
}

void OnBackTab(void) {
	auto selectedLines = e.__selection.GetSelectedLines(e.code_str());

	// deleting tabs from beginning
	if (selectedLines.empty()) {
		if (e.lines[e.row].buf[0] == '\t')  {
			DeleteCharacter(e.row, 0);
			e.col--;
		}
	} else {
		e.__selection.ssx = 0;
		for (auto linenumber : selectedLines) {
			e.row = linenumber;
			if (!e.lines[e.row].buf.empty() && e.lines[e.row].buf[0] == '\t') {
				DeleteCharacter(e.row, 0);
				e.col = e.lines[e.row].buf.size();
			}
		}
	}

	if (!e.redo.empty())
		e.redo.clear();

	Focus();
	UpdateLsp(e.code_str(), false);
	FindTests();
}
