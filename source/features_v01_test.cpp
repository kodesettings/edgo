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

TEST(FeaturesV01Tests, TestOnCommentLine) {
	e.code = rope<char>("this is a sample text");

	APPLY_HIGHLIGHTER

	e.langConf.comment = "//";
	OnCommentLine();

	auto expected = "//this is a sample text";
	EXPECT_EQ(expected, e.code_str());

	OnUndo();
	auto actual = "this is a sample text";
	EXPECT_EQ(actual, e.code_str());

	OnRedo();
	EXPECT_EQ(expected, e.code_str());
}

TEST(FeaturesV01Tests, TestOnSwapLinesUp) {
	e.code = rope<char>("first line of text\nsecond line of text\nthird line of text\n");

	e.row = 1;
	e.col = 0;

	APPLY_HIGHLIGHTER
	OnSwapLinesUp();

	auto expected = "second line of text\nfirst line of text\nthird line of text\n";
	EXPECT_EQ(expected, e.code_str());

	OnUndo();
	auto actual = "first line of text\nsecond line of text\nthird line of text\n";
	EXPECT_EQ(actual, e.code_str());

	OnRedo();
	EXPECT_EQ(expected, e.code_str());
}

TEST(FeaturesV01Tests, TestOnSwapLinesDown) {
	e.code = rope<char>("first line of text\nsecond line of text\nthird line of text\n");

	e.row = 1;
	e.col = 0;

	APPLY_HIGHLIGHTER
	OnSwapLinesDown();

	auto expected = "first line of text\nthird line of text\nsecond line of text\n";
	EXPECT_EQ(expected, e.code_str());

	OnUndo();
	auto actual = "first line of text\nsecond line of text\nthird line of text\n";
	EXPECT_EQ(actual, e.code_str());

	OnRedo();
	EXPECT_EQ(expected, e.code_str());
}
