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

///
/// Convert LSP range between API and internal LSP models
///
#define CONVERT_LSP_RANGE(range, orig)                                 \
	range.start.line      = orig.start.line;                           \
	range.start.character = orig.start.character;                      \
	range.end.line        = orig.end.line;                             \
	range.end.character   = orig.end.character;

///
/// String allocation is fixed to 2048 bytes
/// This value should be enough for most of the cases
///
#define ALLOC_STRING(name, value)                                      \
	static char name[2048];                                            \
	strcpy(name, value.c_str());

///
/// String allocation for source file outputs
///
#define ALLOC_STRING_L(name, value)                                    \
	static char name[24128];                                           \
	strcpy(name, value.c_str());

///
/// LSP Text Edit is a special field in the API that also has LSP range
/// field copied using previous macro called CONVERT_LSP_RANGE
///
#define LSP_TEXT_EDIT(destination, edit)                               \
	struct lsp_text_edit lsp_text_edit;                                \
	struct lsp_range lsp_range;                                        \
	CONVERT_LSP_RANGE(lsp_range, edit.range);                          \
	lsp_text_edit.range = lsp_range;                                   \
	lsp_text_edit.new_text = edit.newText.c_str();                     \
	destination = lsp_text_edit;

