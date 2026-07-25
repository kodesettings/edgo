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

#include "utils.h"
#include <gtest/gtest.h>

TEST(UtilsTests, TestGetLinesArrayFromData) {
	const std::string data = "some data\nmore data\nand more data";
	line_v output = __build_line_vec(data, 3, false);

	EXPECT_EQ(output.size(), 3);
}

TEST(UtilsTests, TestLineOffset) {
	const std::string text = "some data\nmore data\nand more data";
	int offset = LineOffset(text, 3);

	EXPECT_EQ(offset, text.length());
}

TEST(UtilsTests, TestRemoveLeadingTabsSpaces) {
	const std::string text = "\t\t  sometext";
	std::string output = RemoveLeadingTabsSpaces(text);

	EXPECT_EQ(output, "sometext");
}

TEST(UtilsTests, TestFoundCharacterOccurances) {
	std::string s("asd\nasd\n");
	int val = FindCharacterOccurances(s, '\n');

	EXPECT_EQ(val, 2);
}

TEST(UtilsTests, TestHexToRGB) {
	std::string hex("#FFAA33");
	int r, g, b;
	HexToRGB(hex, &r, &g, &b);

	EXPECT_EQ(r, 255);
	EXPECT_EQ(g, 170);
	EXPECT_EQ(b, 51);
}

TEST(UtilsTests, TestColorizeTextChunk) {
	std::string str("some text");
	int r = 250, g = 88, b = 75;
	ColorizeTextChunk(r, g, b, &str);

	EXPECT_EQ(str, "\x1B[38;2;250;88;75msome text\x1B[0m");
}
