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
#include "utils.h"
#include "sdl2_init.h"
#include "highlighter.h"

// this function sets the required parameters to validate updates
void set_update_parameters(bool changed) {
	e.update = true;
	e.isContentChanged = changed;
//	FindTests();
}

void OnCopy(void) {
	auto selectionString = e.__selection.GetSelectionString(e.code_str());
	__set_clipboard_text(selectionString.c_str());
}

void OnPaste(void) {
	if (e.__selection.IsSelectionNonEmpty()) {
		Cut(false);
	}

	char* text = __get_clipboard_text();

	if (strlen(text) == 0) { return; }
	InsertString(e.row, e.col, std::string(text));
	__free_clipboard(text);
	set_update_parameters(true);
}

void Cut(bool isCopySelected) {
	Focus();

	if (e.code.size() == 0) { e.row = 0; e.col = 0; return; }

	if (isCopySelected) {
		auto selectionString = e.__selection.GetSelectionString(e.code_str());
		if (selectionString.length() == 0) { return; } // skipping function if selection is empty
		__set_clipboard_text(selectionString.c_str());
	}

	auto selectedIndices = e.__selection.GetSelectedIndices(e.code_str());
	int lastElement = selectedIndices.size() - 1;

	auto exd = selectedIndices[lastElement].x;
	auto eyd = selectedIndices[lastElement].y;
	auto sxd = selectedIndices[0].x;
	auto syd = selectedIndices[0].y;

	e.row = syd;
	e.col = sxd;

	int exd_offset = 2;
	if (sxd == 0) { exd_offset++; }

	sxd += LineOffset(e.code_str(), syd) + 1;
	exd += LineOffset(e.code_str(), eyd) + exd_offset;

	auto text = e.code.substr(sxd, exd);

	e.code.erase(sxd, exd - sxd);
	RemoveTextEdit(e.code_str(), sxd, text.length());

	e.undo.push_back({operation_t{DELETE, std::string(text.c_str()), sxd, cursormove_t{e.row, e.col}}});
	e.__selection.CleanSelection();
//	e.UpdateLsp(false, string(e.code.Value()));
	set_update_parameters(true);
}

void Duplicate(void) {
	Focus();

	if (e.code.size() == 0) { return; }
	auto syd = LineOffset(e.code_str(), e.row) + 1;
	auto eyd = LineOffset(e.code_str(), e.row + 1) + 1;
	if (e.row == 0) { syd--; } // this is required for first line only
	auto duplicatedSlice = e.code.substr(syd, eyd);

	e.code.insert(eyd, duplicatedSlice);
	e.row++;

	eyd = LineOffset(e.code_str(), e.row) + 1;

	InsertTextEdit(e.code_str(), syd, duplicatedSlice.size());
	e.undo.push_back({operation_t{INSERT, std::string(duplicatedSlice.c_str()), eyd, cursormove_t{e.row, e.col}}});
//	e.UpdateLsp(false, string(e.code.Value()));
	set_update_parameters(true);
}

void OnUndo(void) {
	if (e.undo.empty()) { /*e.UpdateLsp(true, string(e.code.Value())); */return; }

	auto lastOperation = e.undo[e.undo.size() - 1];
	e.undo.pop_back(); // removing last element
	Focus();

	auto index = lastOperation.size() - 1;
	operation_t o{};
undo:
	if (lastOperation.size() > 0) { o = lastOperation[index]; } else { goto exit; }

	switch (o.action) {
	case INSERT:
	case ADDCHAR:
		e.code.erase(o.offset, o.offset - o.text.size());
		RemoveTextEdit(e.code_str(), o.offset, o.text.size());
		e.row = o.cursor.line; e.col = o.cursor.column;
	break;
	case DELETE:
		e.code.insert(o.offset, o.text.c_str());
		InsertTextEdit(e.code_str(), o.offset, o.text.size());
		e.row = o.cursor.line; e.col = o.cursor.column;
	break;
	default:
	break;
	}

	if (index > 0) { index--; goto undo; }
exit:
	e.redo.push_back(lastOperation);
//	e.UpdateLsp(false, string(e.code.Value()));
	set_update_parameters(true);
}

void OnRedo(void) {
	if (e.redo.size() == 0) { return; }

	auto lastRedoOperation = e.redo[e.redo.size() - 1];
	e.redo.pop_back(); // removing last element

	auto index = 0, offset = 0;
	operation_t o{};
redo:
	if (lastRedoOperation.size() > 0) { o = lastRedoOperation[index]; } else { goto exit; }

	switch (o.action) {
	case ADDCHAR:
		offset++;
	case INSERT:
		e.code.insert(o.offset, o.text.c_str());
		InsertTextEdit(e.code_str(), o.offset, o.text.size());
		e.row = o.cursor.line; e.col = o.cursor.column + offset;
	break;
	case DELETE:
		e.code.erase(o.offset, o.offset - o.text.size());
		RemoveTextEdit(e.code_str(), o.offset, o.text.size());
		e.row = o.cursor.line; e.col = o.cursor.column;
	break;
	default:
	break;
	}

	if (index < (int)lastRedoOperation.size() - 1) { index++; goto redo; }
exit:
	e.undo.push_back(lastRedoOperation);
//	e.UpdateLsp(false, string(e.code.Value()))
	set_update_parameters(true);
}
