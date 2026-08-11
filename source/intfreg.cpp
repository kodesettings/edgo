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

#include "../include/edgo.h"
#include "editor.h"
#include "utils.h"
#ifdef __cplusplus
extern "C" {
#endif
char* add_text(size_t line, size_t pos, const char* buf, size_t length) {
	switch (length) {
	case 0: break;
	case 1: AddCharacter(buf[0]); break; // line and pos ignored
	default: InsertString(line, pos, buf); break;
	}

	return __export_screen(e.code_str(), e.y);
}

char* remove_text(size_t line, size_t pos, size_t length) {
	switch (length) {
	case 0: break;
	case 1: DeleteCharacter(line, pos); break;
	default: Cut(false); break; // don't copy to clipboard
	}
	return __export_screen(e.code_str(), e.y);
}

char* replace_text(size_t line, size_t from, size_t end, const char* buf, size_t length) {
	switch (length) {
	case 0: break;
	default: ReplaceString(line, from, end, buf);
	}
	return __export_screen(e.code_str(), e.y);
}

char* shift_with_tabs(size_t line, size_t pos, int *arr, size_t arr_len) {
	if (!arr_len) {
		printf("third argument must be an array");
		return NULL;
	} else {
		auto lines = std::set<int>(arr, arr + arr_len);
		ShiftWithTabsToRight(line, pos, lines);
	}

	return __export_screen(e.code_str(), e.y);
}

char* display_screen_report(const char* dirpath, size_t length) {
	if (!length) {
		printf("provide dir path");
		return NULL;
	}

	std::string report = OnLangLinesCount(dirpath);
	return __export_screen(report, 0, false);
}

char* clipboard(enum clipboard_ops ops) {
	switch (ops) {
	case COPY: OnCopy(); break;
	case PASTE: OnPaste(); break;
	case CUT: Cut(true); break;
	case UNDO: OnUndo(); break;
	case REDO: OnRedo(); break;
	default: break;
	}
	return __export_screen(e.code_str(), e.y);
}

char* duplicate(void) {
	Duplicate();
	return __export_screen(e.code_str(), e.y);
}

char* commentline(void) {
	OnCommentLine();
	return __export_screen(e.code_str(), e.y);
}

char* swaplinesup(void) {
	OnSwapLinesUp();
	return __export_screen(e.code_str(), e.y);
}

char* swaplinesdn(void) {
	OnSwapLinesDown();
	return __export_screen(e.code_str(), e.y);
}

char* on_keypress(enum nav_keys keys, void* is_paging) {
	bool isp = is_paging == NULL ? false : static_cast<bool>(is_paging);
	switch (keys) {
	default: e.lines = __build_line_vec(e.code_str(), -1, false);
	case DOWN: OnDown(isp); break;
	case UP: OnUp(isp); break;
	case LEFT: OnLeft(); break;
	case RIGHT: OnRight(); break;
	case TOP: GoTop(); break;
	case BOTTOM: GoBottom(); break;
	case ENTER: OnEnter(); break;
	case DEL: OnDelete(); break;
	case TAB: OnTab(); break;
	case BACKTAB: OnBackTab(); break;
	}
	return __export_screen(e.code_str(), e.y);
}

char* on_scroll(enum nav_keys keys) {
	switch (keys) {
	case UP: OnScrollUp(); break;
	case DOWN: OnScrollDown(); break;
	default: break;
	}
	return __export_screen(e.code_str(), e.y);
}

char* open_file(const char* filepath) {
	if (!HandleFile(filepath, true)) {
		printf("cannot open filepath");
		return NULL;
	}

	return __export_screen(e.code_str(), e.y);
}

int new_file(const char* filename) {
	return (int)HandleFile(filename, false) == 1;
}

int save_file(void) {
	return (int)SaveFile() == 1;
}
#ifdef __cplusplus
}
#endif
