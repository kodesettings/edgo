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
		{7,  {name: "Test1", filename: "example_test.go", line: 7}},
		{13, {name: "Test2", filename: "example_test.go", line: 13}},
	};

	SetLang(GO);
	testfinder_t testfinder;
	TestQuery(testfinder.testquery);
	TestFinder(testfinder, "example_test.go", GO_TEST_SAMPLE, &test);
	EXPECT_EQ(test.size(), expectedtest.size());

	for (const auto &kv : test) {
		EXPECT_EQ(test[kv.first].name, expectedtest[kv.first].name);
		EXPECT_EQ(test[kv.first].filename, expectedtest[kv.first].filename);
		EXPECT_EQ(test[kv.first].line, expectedtest[kv.first].line);
	}
}

TEST(QueryTests, TestJavascriptFindTest) {
	testdata_m test, expectedtest = {
		{6,  {name: "math tests", filename: "test.js", line: 6}},
		{8, {name: "positive", filename: "test.js", line: 8}},
		{12, {name: "negative", filename: "test.js", line: 12}},
		{16, {name: "failed", filename: "test.js", line: 16}},
	};

	SetLang(JAVASCRIPT);
	testfinder_t testfinder;
	TestQuery(testfinder.testquery);
	TestFinder(testfinder, "test.js", JAVASCRIPT_TEST_SAMPLE, &test);
	EXPECT_EQ(test.size(), expectedtest.size());

	for (const auto &kv : test) {
		EXPECT_EQ(test[kv.first].name, expectedtest[kv.first].name);
		EXPECT_EQ(test[kv.first].filename, expectedtest[kv.first].filename);
		EXPECT_EQ(test[kv.first].line, expectedtest[kv.first].line);
	}
}

TEST(QueryTests, TestJavaFindTest) {
	testdata_m test, expectedtest = {
		{12,  {name: "addition", filename: "test.java", line: 12}},
	};

	SetLang(JAVA);
	testfinder_t testfinder;
	TestQuery(testfinder.testquery);
	TestFinder(testfinder, "test.java", JAVA_TEST_SAMPLE, &test);
	EXPECT_EQ(test.size(), expectedtest.size());

	for (const auto &kv : test) {
		EXPECT_EQ(test[kv.first].name, expectedtest[kv.first].name);
		EXPECT_EQ(test[kv.first].filename, expectedtest[kv.first].filename);
		EXPECT_EQ(test[kv.first].line, expectedtest[kv.first].line);
	}
}

TEST(QueryTests, TestPythonFindTest) {
	testdata_m test, expectedtest = {
		{3,  {name: "TestYo", filename: "test_yo.py", line: 3}},
		{7,  {name: "test_pass", filename: "test_yo.py", line: 7}},
		{10,  {name: "test_fail_sometimes", filename: "test_yo.py", line: 10}},
	};

	SetLang(PYTHON);
	testfinder_t testfinder;
	TestQuery(testfinder.testquery);
	TestFinder(testfinder, "test_yo.py", PYTHON_TEST_SAMPLE, &test);
	EXPECT_EQ(test.size(), expectedtest.size());

	for (const auto &kv : test) {
		EXPECT_EQ(test[kv.first].name, expectedtest[kv.first].name);
		EXPECT_EQ(test[kv.first].filename, expectedtest[kv.first].filename);
		EXPECT_EQ(test[kv.first].line, expectedtest[kv.first].line);
	}
}
