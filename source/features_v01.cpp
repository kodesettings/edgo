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

#include "operations.h"
#include "utils.h"
#include "editor.h"
#include "highlighter.h"

void OnCommentLine(void) {
	Focus();

	bool found = false;
	auto offset = LineOffset(e.code_str(), e.row);
	if (e.row > 0) { offset++; } // due to newline character
	auto ch_1 = e.code.at(offset);
	auto ch_2 = e.code.at(offset + 1);

	if (e.langConf.comment.size() == 1 && ch_1 == e.langConf.comment[0]) {
		e.code.erase(offset, 1);
		RemoveTextEdit(e.code_str(), offset, 1);
		e.undo.push_back({operation_t{DELETE, std::to_string(ch_1), offset, cursormove_t{offset, 0}}});
		found = true;
	} else if (e.langConf.comment.size() == 2 && ch_1 == e.langConf.comment[0] && ch_2 == e.langConf.comment[1]) {
		e.code.erase(offset, 2);
		RemoveTextEdit(e.code_str(), offset, 2);
		e.undo.push_back({operation_t{DELETE, std::to_string(ch_1 + ch_2), offset, cursormove_t{offset, 0}}});
		found = true;
	}

	if (found) { goto exit; }

	e.code.insert(offset, e.langConf.comment.c_str());
	InsertTextEdit(e.code_str(), offset, e.langConf.comment.size());
	e.undo.push_back({operation_t{INSERT, e.langConf.comment, offset, cursormove_t{e.row, 0}}});
exit:
//	UpdateLsp(false, e.code_str());
	set_update_parameters(true);
}

void OnSwapLinesUp(void) {
	Focus();

	if (e.row == 0) { return; }

	auto from = LineOffset(e.code_str(), e.row - 1);
	auto to = LineOffset(e.code_str(), e.row);
	auto offset = LineOffset(e.code_str(), e.row + 1);

	auto line_1 = e.code.substr(from, to - from);
	auto line_2 = e.code.substr(to + 1, offset - to);

	e.code.erase(to, line_2.size()); // remove line_2 from current position
	RemoveTextEdit(e.code_str(), to, line_2.size());

	if (e.row > 1) { offset = from + 1; } else { offset = from; }
	e.code.insert(offset, line_2); // add line_2 to top
	InsertTextEdit(e.code_str(), offset, line_2.size());

	e.undo.push_back(editoperation_t{
		operation_t{DELETE, line_2.c_str(), to + 1, cursormove_t{e.row, -1}},
		operation_t{INSERT, line_2.c_str(), offset, cursormove_t{e.row - 1, -1}}
	});

	e.row--;
//	UpdateLsp(false, e.code_str());
	set_update_parameters(true);
}

void OnSwapLinesDown(void) {
	Focus();

	int nlines = FindCharacterOccurances(e.code_str(), '\n');
	if (e.row == nlines - 1) { return; }

	auto from = LineOffset(e.code_str(), e.row);
	auto to = LineOffset(e.code_str(), e.row + 1);
	auto offset = LineOffset(e.code_str(), e.row + 2);

	auto line_1 = e.code.substr(from, to - from);
	auto line_2 = e.code.substr(to + 1, offset - to);

	// TODO: fix moving from first line

	e.code.erase(from, line_1.size()); // remove line_1 from current position
	RemoveTextEdit(e.code_str(), from, line_1.size());

	offset = LineOffset(e.code_str(), e.row + 1);
	e.code.insert(offset, line_1); // add line_1 to bottom
	InsertTextEdit(e.code_str(), offset, line_1.size());

	e.undo.push_back(editoperation_t{
		operation_t{DELETE, line_1.c_str(), from, cursormove_t{e.row, -1}},
		operation_t{INSERT, line_1.c_str(), offset, cursormove_t{e.row + 1, -1}}
	});

	e.row++;
//	UpdateLsp(false, e.code_str());
	set_update_parameters(true);
}
