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

#include "config.h"
#include <gtest/gtest.h>

TEST(ConfigTests, TestReadConfig) {
	// Read conf
	auto conf = GetConfig();

	// Print config
	for (auto lang : conf.langs) {
		printf("name: %s, lsp: %s, comment: %s, tab width: %d\n",
			lang.first.c_str(), lang.second.lsp.c_str(),
			lang.second.comment.c_str(), lang.second.tabwidth
		);
	}

	auto golang = conf.langs["go"];

	EXPECT_EQ(golang.lsp, "gopls");
	EXPECT_EQ(golang.comment, "//");
	EXPECT_EQ(golang.tabwidth, 4);
}
