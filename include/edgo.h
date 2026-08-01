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

#pragma once
#include <unistd.h>

#ifdef __cplusplus
extern "C" {
#endif
//
// Text manipulators
//
char* add_text(size_t line, size_t pos, const char* buf, size_t length);
char* remove_text(size_t line, size_t pos, size_t length);
char* replace_text(size_t line, size_t from, size_t end, const char* buf, size_t length);
char* shift_with_tabs(size_t line, size_t pos, int* arr, size_t arr_len);

//
// Screen report
//
char* display_screen_report(const char* dirpath, size_t length);

enum clipboard_ops {
	COPY  = 0x034FC,
	PASTE = 0x035FC,
	CUT   = 0x036FD,
	UNDO  = 0x044FD,
	REDO  = 0x044FB
};

//
// Clipboard operations
//
char* clipboard(enum clipboard_ops ops);

//
// Editor features
//
char* duplicate(void);
char* commentline(void);
char* swaplinesup(void);
char* swaplinesdn(void);

enum nav_keys {
	UP      = 0x01DF,
	DOWN    = 0x02DF,
	LEFT    = 0x03DF,
	RIGHT   = 0x04DF,
	TOP     = 0x05DF,
	BOTTOM  = 0x06DF,
	ENTER   = 0x01FC,
	DEL     = 0x02FC,
	TAB     = 0x01FF,
	BACKTAB = 0x02FF
};

//
// Navigation keys
//
char* on_keypress(enum nav_keys keys, void* is_paging);
char* on_scroll(enum nav_keys keys);

//
// File handling
//
char* open_file(const char* filepath);
int new_file(const char* filename);
int save_file(void);
#ifdef __cplusplus
}
#endif
