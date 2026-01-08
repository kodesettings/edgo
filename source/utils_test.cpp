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
	line_v output = GetLinesArrayFromData(data, 3);

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
