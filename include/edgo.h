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
#include <stddef.h>
#include "lsp.h"
#ifdef __cplusplus
extern "C" {
#endif
struct screen {
	char *filename; // name of the file currently opened
	char *language; // langauge identifier of the file
	char *content;  // terminal screen content to be displayed
	int  cursor_x;  // cursor position
	int  cursor_y;  // cursor position
	bool changed;   // flag for file modification
};

//
// Text manipulators
//

struct screen add_text(size_t line, size_t pos, const char *buf, size_t length);
struct screen remove_text(size_t line, size_t pos, size_t length);
struct screen replace_text(size_t line, size_t from, size_t end, const char *buf, size_t length);
struct screen shift_with_tabs(size_t line, size_t pos, int *arr, size_t arr_len);

//
// Screen report
//

struct screen display_screen_report(const char *dirpath, size_t length);

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

struct screen clipboard(enum clipboard_ops ops);

//
// Editor features
//

struct screen duplicate(void);
struct screen commentline(void);
struct screen swaplinesup(void);
struct screen swaplinesdn(void);

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

struct screen on_keypress(enum nav_keys keys, void *is_paging);
struct screen on_scroll(enum nav_keys keys);

//
// File handling
//

struct screen open_file(const char *filepath);
int new_file(const char *filename);
int save_file(void);

//
// LSP functions
//

void lsp_client_hover(struct lsp_hover*);
void lsp_client_completion(struct lsp_completion*);
void lsp_client_definition(struct lsp_definition*);
void lsp_client_signature_help(struct lsp_signature_help*);
void lsp_client_references(struct lsp_references*);
void lsp_client_prepare_rename(struct lsp_prepare_rename*);
void lsp_client_rename(const char* newname, struct lsp_rename*);
void lsp_client_code_action(struct lsp_code_action*);

//
// LSP diagnostics
//

void lsp_diagnostics(struct lsp_publish_diagnostics*);
#ifdef __cplusplus
}
#endif
