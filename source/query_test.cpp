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

#include "query.h"
#include "test_languages.h"
#include <gtest/gtest.h>

TEST(QueryTests, TestGoFindTest) {
	testdata_m test, expectedtest = {
		{4,  {name: "", filename: "example_test.go", line: 4}},
		{6, {name: "", filename: "example_test.go", line: 6}},
		{12, {name: "", filename: "example_test.go", line: 12}},
	};

	SetLang(GO);
	Parse(GO_TEST_SAMPLE);

	TestQueryAlloc();
	TestFinder("example_test.go", &test);
	EXPECT_EQ(test.size(), expectedtest.size());

	for (const auto &kv : expectedtest) {
		EXPECT_EQ(test[kv.first].name, expectedtest[kv.first].name);
		EXPECT_EQ(test[kv.first].filename, expectedtest[kv.first].filename);
		EXPECT_EQ(test[kv.first].line, expectedtest[kv.first].line);
	}
}

TEST(QueryTests, TestJavascriptFindTest) {
	testdata_m test, expectedtest = {
		{4,  {name: "", filename: "test.js", line: 4}},
		{5, {name: "", filename: "test.js", line: 5}},
		{9, {name: "", filename: "test.js", line: 9}},
		{13, {name: "", filename: "test.js", line: 13}},
	};

	SetLang(JAVASCRIPT);
	Parse(JAVASCRIPT_TEST_SAMPLE);

	TestQueryAlloc();
	TestFinder("test.js", &test);
	EXPECT_EQ(test.size(), expectedtest.size());

	for (const auto &kv : expectedtest) {
		EXPECT_EQ(test[kv.first].name, expectedtest[kv.first].name);
		EXPECT_EQ(test[kv.first].filename, expectedtest[kv.first].filename);
		EXPECT_EQ(test[kv.first].line, expectedtest[kv.first].line);
	}
}

TEST(QueryTests, TestJavaFindTest) {
	testdata_m test, expectedtest = {
		{9,  {name: "", filename: "test.java", line: 9}},
	};

	SetLang(JAVA);
	Parse(JAVA_TEST_SAMPLE);

	TestQueryAlloc();
	TestFinder("test.java", &test);
	EXPECT_EQ(test.size(), expectedtest.size());

	for (const auto &kv : expectedtest) {
		EXPECT_EQ(test[kv.first].name, expectedtest[kv.first].name);
		EXPECT_EQ(test[kv.first].filename, expectedtest[kv.first].filename);
		EXPECT_EQ(test[kv.first].line, expectedtest[kv.first].line);
	}
}

TEST(QueryTests, TestPythonFindTest) {
	testdata_m test, expectedtest = {
		{3,  {name: "", filename: "test_yo.py", line: 3}},
		{4,  {name: "", filename: "test_yo.py", line: 4}},
		{7,  {name: "", filename: "test_yo.py", line: 7}},
		{10,  {name: "", filename: "test_yo.py", line: 10}},
	};

	SetLang(PYTHON);
	Parse(PYTHON_TEST_SAMPLE);

	TestQueryAlloc();
	TestFinder("test_yo.py", &test);
	EXPECT_EQ(test.size(), expectedtest.size());

	for (const auto &kv : expectedtest) {
		EXPECT_EQ(test[kv.first].name, expectedtest[kv.first].name);
		EXPECT_EQ(test[kv.first].filename, expectedtest[kv.first].filename);
		EXPECT_EQ(test[kv.first].line, expectedtest[kv.first].line);
	}
}
