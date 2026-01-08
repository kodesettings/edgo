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
#include <gtest/gtest.h>

TEST(CharacterTests, TestAddCharacter) {
	e.code = rope<char>("this is a sample text");

	e.row = 0;
	e.col = 11;

//	apply_highlighter(e.code.Value(), "", "")
	AddCharacter('g');

	const auto text = "this is a sample text";
	const auto expected = "this is a sgample text";
	EXPECT_EQ(e.code_str(), expected);

	OnUndo();
	EXPECT_EQ(e.code_str(), text);

	OnRedo();
	EXPECT_EQ(e.code_str(), expected);
}

TEST(CharacterTests, TestInsertCharacter) {
	e.code = rope<char>("this is a sample text");

	e.row = 0;
	e.col = 2;

//	apply_highlighter(e.code.Value(), "", "")
	InsertCharacter(0, 2, 'f');

	const auto text = "this is a sample text";
	const auto expected = "thfis is a sample text";
	EXPECT_EQ(e.code_str(), expected);

	OnUndo();
	EXPECT_EQ(e.code_str(), text);

	OnRedo();
	EXPECT_EQ(e.code_str(), expected);
}

TEST(CharacterTests, TestInsertString) {
	e.code = rope<char>("this is a sample text");

	e.row = 0;
	e.col = 1;

//	apply_highlighter(e.code.Value(), "", "")
	InsertString(0, 1, "some");

	const auto text = "this is a sample text";
	const auto expected = "thsomeis is a sample text";
	EXPECT_EQ(e.code_str(), expected);

	OnUndo();
	EXPECT_EQ(e.code_str(), text);

	OnRedo();
	EXPECT_EQ(e.code_str(), expected);
}

TEST(CharacterTests, TestInsertLines) {
	e.code = rope<char>("this is a sample text");

	e.row = 0;
	e.col = 10;

//	apply_highlighter(e.code.Value(), "", "")
	InsertString(0, 10, "new line\nanother line\n");

	const auto text = "this is a sample text";
	const auto expected = "this is a snew line\nanother line\nample text";
	EXPECT_EQ(e.code_str(), expected);

	OnUndo();
	EXPECT_EQ(e.code_str(), text);

	OnRedo();
	EXPECT_EQ(e.code_str(), expected);
}

TEST(CharacterTests, TestDeleteCharacter) {
	e.code = rope<char>("this is a sample text");

	e.row = 0;
	e.col = 2;

//	apply_highlighter(e.code.Value(), "", "")
	DeleteCharacter(0, 2);

	const auto text = "this is a sample text";
	const auto expected = "ths is a sample text";
	EXPECT_EQ(e.code_str(), expected);

	OnUndo();
	EXPECT_EQ(e.code_str(), text);

	OnRedo();
	EXPECT_EQ(e.code_str(), expected);
}

TEST(CharacterTests, TestReplaceString) {
	e.code = rope<char>("this is a sample text");

//	apply_highlighter(e.code.Value(), "", "")
	ReplaceString(0, 3, 11, "asd");

	const auto text = "this is a sample text";
	const auto expected = "thisasdmple text";
	EXPECT_EQ(e.code_str(), expected);

	OnUndo();
	EXPECT_EQ(e.code_str(), text);

	OnRedo();
	EXPECT_EQ(e.code_str(), expected);
}

TEST(CharacterTests, TestDeleteLine) {
	e.code = rope<char>("this is a sample text\n");

//	apply_highlighter(e.code.Value(), "", "")
	DeleteCharacter(0, e.code.size() - 1);

	const auto text = "this is a sample text\n";
	const auto expected = "this is a sample text";
	EXPECT_EQ(e.code_str(), expected);

	OnUndo();
	EXPECT_EQ(e.code_str(), text);

	OnRedo();
	EXPECT_EQ(e.code_str(), expected);
}

TEST(CharacterTests, TestInsertLine) {
	e.code = rope<char>("this is a sample text\n");

//	apply_highlighter(e.code.Value(), "", "")
	InsertCharacter(0, e.code.size() - 1, '\n');

	const auto text = "this is a sample text\n";
	const auto expected = "this is a sample text\n\n";
	EXPECT_EQ(e.code_str(), expected);

	OnUndo();
	EXPECT_EQ(e.code_str(), text);

	OnRedo();
	EXPECT_EQ(e.code_str(), expected);
}

TEST(CharacterTests, TestShiftWithTabsToRight) {
	e.code = rope<char>("this is a sample text\none more line");

	e.__selection.ssx = 0;
	e.__selection.ssy = 0;
	e.__selection.sex = 0;
	e.__selection.sey = 2;

	const auto got = e.__selection.GetSelectionString(e.code_str());
	const auto text = "this is a sample text\none more line";
	EXPECT_EQ(text, got);

	const auto selectedLines = e.__selection.GetSelectedLines(e.code_str());
	ShiftWithTabsToRight(0, 0, selectedLines);

	const auto expected = "\tthis is a sample text\t\none more line";
	EXPECT_EQ(e.code_str(), expected);

	OnUndo();
	EXPECT_EQ(e.code_str(), text);

	OnRedo();
	EXPECT_EQ(e.code_str(), expected);
}

TEST(CharacterTests, TestMaybeAddPair) {
	e.code = rope<char>("test[");

	char val;
	bool found = MaybeAddPair(0, 4, '[', &val);

	EXPECT_EQ(val, ']');
	EXPECT_EQ(found, true);
}
