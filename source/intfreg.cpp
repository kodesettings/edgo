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
struct screen add_text(const char *buf, size_t length) {
	switch (length) {
	case 0: break;
	case 1: AddCharacter(buf[0]); break; // line and pos ignored
	default: InsertString(e.row, e.col, buf); break;
	}

	return __export_screen(e.code_str(), e.y);
}

struct screen replace_text(size_t ps, size_t pe, const char *buf, size_t length) {
	switch (length) {
	case 0: break;
	default:
		int start_ch = ps == 0 ? e.__selection.ssx : ps;
		int end_ch = pe == 0 ? e.__selection.sex : pe;
		ReplaceString(e.row, start_ch, end_ch, buf);
	}
	return __export_screen(e.code_str(), e.y);
}

struct screen move_cursor(size_t line, size_t pos) {
	e.row = line; e.col = pos;
	e.__selection.CleanSelection();
	return __export_screen(e.code_str(), e.y);
}

struct screen display_screen_report(const char *dirpath, size_t length) {
	std::string dirpath_s, report;
	switch (length) {
	case 0: dirpath_s = e.cwd; break;
	default: dirpath_s = dirpath == NULL ? e.cwd : dirpath;
	}
	try {
		report = OnLangLinesCount(dirpath_s);
	} catch (...) {}
	return __export_screen(report, 0, false);
}

struct screen clipboard(enum clipboard_ops ops) {
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

struct screen duplicate(void) {
	Duplicate();
	return __export_screen(e.code_str(), e.y);
}

struct screen commentline(void) {
	OnCommentLine();
	return __export_screen(e.code_str(), e.y);
}

struct screen swaplinesup(void) {
	OnSwapLinesUp();
	return __export_screen(e.code_str(), e.y);
}

struct screen swaplinesdn(void) {
	OnSwapLinesDown();
	return __export_screen(e.code_str(), e.y);
}

#define SELECT_START \
{ \
	if (e.__selection.ssx < 0) { \
		e.__selection.ssx = e.col; \
		e.__selection.ssy = e.row; \
	} \
}

#define SELECT_END \
{ \
	if (e.__selection.ssx >= 0) { \
		e.__selection.sex = e.col; \
		e.__selection.sey = e.row; \
		e.__selection.is_selected = true; \
	} \
}

#define CLEAN_SELECT \
	e.__selection.CleanSelection();

struct screen on_keypress(enum nav_keys keys, bool is_paging, bool is_shift) {
	e.lines = __build_line_vec(e.code_str(), -1, false);
	switch (keys) {
	case DOWN:
		if (is_shift) SELECT_START else CLEAN_SELECT
		OnDown(is_paging);
		if (is_shift) SELECT_END
	break;
	case UP:
		if (is_shift) SELECT_START else CLEAN_SELECT
		OnUp(is_paging);
		if (is_shift) SELECT_END
	break;
	case LEFT:
		if (is_shift) SELECT_START else CLEAN_SELECT
		OnLeft();
		if (is_shift) SELECT_END
	break;
	case RIGHT:
		if (is_shift) SELECT_START else CLEAN_SELECT
		OnRight();
		if (is_shift) SELECT_END
	break;
	case TOP: GoTop(); break;
	case BOTTOM: GoBottom(); break;
	case ENTER: OnEnter(); break;
	case DEL: OnDelete(); break;
	case ESC: CLEAN_SELECT; break;
	case TAB: OnTab(); break;
	case BACKTAB: OnBackTab(); break;
	default: break;
	}
	return __export_screen(e.code_str(), e.y);
}

struct screen on_scroll(enum nav_keys keys) {
	switch (keys) {
	case UP: OnScrollUp(); break;
	case DOWN: OnScrollDown(); break;
	default: break;
	}
	return __export_screen(e.code_str(), e.y);
}

struct screen open_file(const char *filepath) {
	if (!HandleFile(filepath, true)) {
		printf("cannot open filepath");
		return screen;
	}

	return __export_screen(e.code_str(), e.y);
}

int new_file(const char *filename) {
	return (int)HandleFile(filename, false) == 1;
}

int save_file(void) {
	return (int)SaveFile() == 1;
}
#ifdef __cplusplus
}
#endif
