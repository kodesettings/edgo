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
#include "highlighter.h"

editor_t e;

void AddCharacter(char ch) {
	if (e.__selection.GetSelectionString(e.code_str()).length() != 0) { Cut(false); }

	Focus();
	InsertCharacter(e.row, e.col, ch);
	e.col++;

	char val;
	bool found = MaybeAddPair(e.row, e.col, ch, &val);
	if (found) {
		InsertCharacter(e.row, e.col, val);
	}

	if (e.isContentChanged) {
		UpdateLsp(e.code_str(), false);
		FindTests();
	}
}

void InsertCharacter(int line, int pos, char ch) {
	int offset = LineOffset(e.code_str(), line) + pos;

	// adding one character offset after first line
	// this is required due to the newline characters in the string buffer
	if (line > 0) { offset += 1; }

	e.code.insert(offset, ch);
	InsertTextEdit(e.code_str(), offset, 1);

	// undo operation type is switched between insert and addchar types
	action_t op;

	// change op if this is a new line character
	if (ch == '\n') { op = INSERT; } else { op = ADDCHAR; }

	// record the operation on the undo stack. Note that we're creating a new EditOperation
	// and adding all the Operations to it
	e.undo.push_back({operation_t{op, std::string(1, ch), offset, cursormove_t{line, pos}}});
}

void InsertString(int line, int pos, std::string linestring) {
	// modify string to remove tabs and space from beginning
	std::string l = RemoveLeadingTabsSpaces(linestring);
	int offset = LineOffset(e.code_str(), line) + pos + 1;
	e.code.insert(offset, l.c_str());
	InsertTextEdit(e.code_str(), offset, l.length());

	// record the operation on the undo stack. Note that we're creating a new EditOperation
	// and adding all the Operations to it
	e.undo.push_back({operation_t{INSERT, l, offset, cursormove_t{line, pos}}});
}

void DeleteCharacter(int line, int pos) {
	int offset = LineOffset(e.code_str(), line) + pos;
	char ch = e.code.at(offset);
	e.code.erase(offset, 1);
	RemoveTextEdit(e.code_str(), offset, 1);

	// record the operation on the undo stack. Note that we're creating a new EditOperation
	// and adding all the Operations to it
	e.undo.push_back({operation_t{DELETE, std::string(1, ch), offset, cursormove_t{line, pos}}});
}

void ReplaceString(int line, int from, int end, std::string instext) {
	int offset = LineOffset(e.code_str(), line);
	int begin_idx = offset + from + 1;

	auto deltext = e.code.substr(begin_idx, end - from);
	auto deltext_str = std::string(deltext.begin(), deltext.end());
	auto instext_str = std::string(instext.begin(), instext.end());

	e.code.erase(begin_idx, deltext.size());
	RemoveTextEdit(e.code_str(), begin_idx, deltext.size());

	e.code.insert(begin_idx, instext.c_str());
	InsertTextEdit(e.code_str(), begin_idx, instext.size());

	// record the operation on the undo stack. Note that we're creating a new EditOperation
	// and adding all the Operations to it
	e.undo.push_back(editoperation_t{
		operation_t{DELETE, deltext_str, begin_idx, cursormove_t{line, from}},
		operation_t{INSERT, instext_str, begin_idx, cursormove_t{line, from}}
	});
}

void ShiftWithTabsToRight(int line, int pos, std::set<int> selectedLines) {
	e.__selection.ssx = 0;

	editoperation_t ops{};
	for (int linenumber : selectedLines) {
		int offset = LineOffset(e.code_str(), linenumber);
		e.code.insert(offset, "\t");
		ops.push_back(operation_t{INSERT, "\t", offset, cursormove_t{line, pos}});
	}

	e.__selection.sex = pos;
	e.undo.push_back(ops);
}

bool MaybeAddPair(int line, int pos, char ch, char *ret) {
	std::map<char, char> pairMap = {
		{'(', ')'}, {'{', '}'}, {'[', ']'}, {'"', '"'}, {'\\', '\\'}, {'`', '`'},
	};

	if (e.code.size() == 0) { *ret = '\0'; return false; }

	int offset_start = LineOffset(e.code_str(), line);
	int offset_end = LineOffset(e.code_str(), line + 1);
	auto line_str = e.code.substr(offset_start + 1, offset_end);
	auto current_char = e.code.at(offset_start + 1 + pos - 2);

	if (pairMap.count(ch)) {
		bool noMoreChars = pos >= (int)line_str.size() - 1;
		bool isSpaceNext = pos < (int)line_str.size() - 1 && current_char == ' ';
		bool isStringAndClosedBracketNext = pairMap[ch] == '"' && pos < (int)line_str.size() - 1 && current_char == ')';

		if (noMoreChars || isSpaceNext || isStringAndClosedBracketNext) {
			*ret = pairMap[ch]; // returning the character
			return true; // found character pair
		}
	}

	*ret = '\0'; // returning an empty character
	return false; // could not found
}
