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
#include "selection.h"
#include <gtest/gtest.h>

TEST(SelectionTests, TestNoSelection) {
	auto text = "Hello, world!\nHow are you doing today?\nI hope you're doing well.";

	selection selection;
	selection.ssx = 0;
	selection.ssy = 0;
	selection.sex = 0;
	selection.sey = 0;

	auto got = selection.GetSelectionString(text);
	auto expected = "";
	
	EXPECT_EQ(got, expected);
}

TEST(SelectionTests, TestSingleCharacterSelection) {
	auto text = "Hello, world!\nHow are you doing today?\nI hope you're doing well.";

	selection selection;
	selection.ssx = 0;
	selection.ssy = 0;
	selection.sex = 1;
	selection.sey = 0;

	auto got = selection.GetSelectionString(text);
	auto expected = "H";
	
	EXPECT_EQ(got, expected);
}

TEST(SelectionTests, TestMultipleCharacterSelection) {
	auto text = "Hello, world!\nHow are you doing today?\nI hope you're doing well.";

	selection selection;
	selection.ssx = 0;
	selection.ssy = 0;
	selection.sex = 5;
	selection.sey = 0;

	auto got = selection.GetSelectionString(text);
	auto expected = "Hello";
	
	EXPECT_EQ(got, expected);
}

TEST(SelectionTests, TestMultipleLineSelection1) {
	auto text = "Hello, world!\nHow are you doing today?\nI hope you're doing well.";

	selection selection;
	selection.ssx = 0;
	selection.ssy = 0;
	selection.sex = 5;
	selection.sey = 1;

	auto got = selection.GetSelectionString(text);
	auto expected = "Hello, world!\nHow a";
	
	EXPECT_EQ(got, expected);
}

TEST(SelectionTests, TestMultipleLineSelection2) {
	auto text = "Hello, world!\nHow are you doing today?\nI hope you're doing well.";

	selection selection;
	selection.ssx = 0;
	selection.ssy = 0;
	selection.sex = 11;
	selection.sey = 1;

	auto got = selection.GetSelectionString(text);
	auto expected = "Hello, world!\nHow are you";
	
	EXPECT_EQ(got, expected);
}

TEST(SelectionTests, TestMultipleLineSelection3) {
	auto text = "Hello, world!\nHow are you doing today?\nI hope you're doing well.";

	selection selection;
	selection.ssx = 6;
	selection.ssy = 0;
	selection.sex = 23;
	selection.sey = 1;

	auto got = selection.GetSelectionString(text);
	auto expected = " world!\nHow are you doing today";
	
	EXPECT_EQ(got, expected);
}

TEST(SelectionTests, TestGetSelectedLines) {
	auto text = "Hello, world!\nHow are you doing today?\nI hope you're doing well.";

	selection selection;
	selection.ssx = 0;
	selection.ssy = 0;
	selection.sex = 0;
	selection.sey = 2;

	auto got = selection.GetSelectedLines(text);
	auto expected = std::set<int>{0,1};

	EXPECT_EQ(got, expected);
}
