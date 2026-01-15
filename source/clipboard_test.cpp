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
#include "highlighter.h"
#include <gtest/gtest.h>

TEST(ClipboardTests, TestCutActionPosition) {
	e.code = rope<char>("this is a sample text\nand another line of text to cut\none more line of text\n");

	e.__selection.ssx = 0;
	e.__selection.ssy = 1;
	e.__selection.sex = 11;
	e.__selection.sey = 2;

	auto got = e.__selection.GetSelectionString(e.code_str());
	auto expected = "and another line of text to cut\none more li";
	EXPECT_EQ(expected, got);

	auto text = e.code_str();
	APPLY_HIGHLIGHTER
	Cut(true);

	expected = "this is a sample text\ne of text\n";
	EXPECT_EQ(expected, e.code_str());

	OnUndo();
	EXPECT_EQ(text, e.code_str());

	OnRedo();
	EXPECT_EQ(expected, e.code_str());
}

TEST(ClipboardTests, TestCutActionLinesOnly) {
	e.code = rope<char>("this is a sample text\nand another line of text to cut\none more line of text\n");

	e.__selection.ssx = 0;
	e.__selection.ssy = 1;
	e.__selection.sex = 0;
	e.__selection.sey = 3;

	auto got = e.__selection.GetSelectionString(e.code_str());
	auto expected = "and another line of text to cut\none more line of text";
	EXPECT_EQ(expected, got);

	auto text = e.code_str();
	APPLY_HIGHLIGHTER
	Cut(true);

	expected = "this is a sample text\n";
	EXPECT_EQ(expected, e.code_str());

	OnUndo();
	EXPECT_EQ(text, e.code_str());

	OnRedo();
	EXPECT_EQ(expected, e.code_str());
}

TEST(ClipboardTests, TestDuplicateAction) {
	e.row = 0;
	e.col = 0;

	e.code = rope<char>("this is a sample text\nand another line of text to cut\none more line of text");

	auto text = e.code_str();
	APPLY_HIGHLIGHTER
	Duplicate();

	auto expected = "this is a sample text\nthis is a sample text\nand another line of text to cut\none more line of text";
	EXPECT_EQ(expected, e.code_str());

	OnUndo();
	EXPECT_EQ(text, e.code_str());

	OnRedo();
	EXPECT_EQ(expected, e.code_str());
}

TEST(ClipboardTests, TestUndoRedoStack) {
	e.row = 0;
	e.col = 0;

	e.code = rope<char>("");
	APPLY_HIGHLIGHTER

	e.undo.clear();
	e.redo.clear();

	AddCharacter('a');
	AddCharacter('b');
	AddCharacter('c');
	AddCharacter('d');

	auto expected = "abcd";
	EXPECT_EQ(expected, e.code_str());

	OnUndo();
	EXPECT_EQ("abc", e.code_str());
	OnUndo();
	EXPECT_EQ("ab", e.code_str());
	OnUndo();
	OnUndo();
	OnUndo();
	EXPECT_EQ("", e.code_str());

	OnUndo();
	OnRedo();
	EXPECT_EQ("a", e.code_str());

	OnRedo();
	OnRedo();
	OnRedo();
	OnRedo();
	EXPECT_EQ("abcd", e.code_str());
}
