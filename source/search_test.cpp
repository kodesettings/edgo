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


#include "search.h"
#include <gtest/gtest.h>

TEST(SearchTests, TestSearch) {
	line_v lines {line_t{"Hello, World!"}, line_t{"This is a test"}, line_t{"Another test"}};
	searchresult_v results = Search(lines, "test");
	searchresult_v expected = searchresult_v{{1, 10}, {2, 8}};
	EXPECT_EQ(results[0].line, 1);
	EXPECT_EQ(results[0].position, 10);
	EXPECT_EQ(results[1].line, 2);
	EXPECT_EQ(results[1].position, 8);

	results = Search(lines, "");
	EXPECT_EQ(results.size(), 0);

	results = Search(lines, "foo");
	EXPECT_EQ(results.size(), 0);

	results = Search(lines, "Hello");
	expected = searchresult_v{{0, 0}};
	EXPECT_EQ(results[0].line, 0);
	EXPECT_EQ(results[0].position, 0);

	results = Search(lines, "is");
	expected = searchresult_v{{1, 2}, {1, 5}};
	EXPECT_EQ(results[0].line, 1);
	EXPECT_EQ(results[0].position, 2);
	EXPECT_EQ(results[1].line, 1);
	EXPECT_EQ(results[1].position, 5);

	results = Search(lines, "o");
	expected = searchresult_v{{0, 4}, {0, 8}, {2, 2}};
	EXPECT_EQ(results[0].line, 0);
	EXPECT_EQ(results[0].position, 4);
	EXPECT_EQ(results[1].line, 0);
	EXPECT_EQ(results[1].position, 8);
	EXPECT_EQ(results[2].line, 2);
	EXPECT_EQ(results[2].position, 2);

	results = Search(lines, "World");
	expected = searchresult_v{{0, 7}};
	EXPECT_EQ(results[0].line, 0);
	EXPECT_EQ(results[0].position, 7);
}

TEST(SearchTests, TestSearchOnFile) {
	int lines;
	auto results = SearchOnFile("/usr/include/strings.h", "__THROW", &lines);
	EXPECT_NE(lines, 0);
	EXPECT_NE(results.size(), 0);
}

TEST(SearchTests, TestSearchOnDir) {
	int files, lines;
	auto results = SearchOnDir(".", "text", &files, &lines);
	EXPECT_NE(results.size(), 0);
	EXPECT_NE(files, 0);
	EXPECT_NE(lines, 0);
}

TEST(SearchTests, TestLinesCountOnFile) {
	int lc, emptylc;
	LineCountOnFile("/usr/include/strings.h", &lc, &emptylc);
	EXPECT_NE(lc, 0);
	EXPECT_NE(emptylc, 0);
}

TEST(SearchTests, TestLinesCountOnDir) {
	int filesprocessedcount, totallinescount;
	auto results = LinesCountOnDir(".", &filesprocessedcount, &totallinescount);
	EXPECT_NE(filesprocessedcount, 0);
	EXPECT_NE(totallinescount, 0);
	langlinescountresult_v llct;
	LangCount(results, &llct);
	EXPECT_NE(llct.size(), 0);
}
