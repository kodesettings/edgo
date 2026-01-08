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

#include "lsp_client.h"
#include <gtest/gtest.h>
#include <fstream>
#include <unistd.h>

// make sure that lsp client is initialized for the entire suite
class LspTest : public ::testing::Test {
protected:
	static void SetUpTestSuite() {
		lspclient.lang = "go";
		started = StartLspClient("gopls", "");

		// finding current working directory
		char currentDir[PATH_MAX];
		getcwd(currentDir, sizeof(currentDir));
		InitLspClient(currentDir);

		// reading file content
		file = file.append("lsp_client_test.cpp");
		std::vector<char> filedata;
		std::ifstream testfile(file, std::ios_base::in);
		testfile.read(filedata.data(), 2400);
		testfile.close();

		// accessing didOpen event
		content = std::string(filedata.data());
		DidOpen(file, &content);
	}
public:
	static bool started;             // flag whether the server started
	static std::string file;         // direct file path
	static std::string content;      // file content
};

// defining the test variables
bool LspTest::started;
std::string LspTest::file;
std::string LspTest::content;

TEST_F(LspTest, TestLspClientStarted) {
        EXPECT_EQ(LspTest::started, true);
		EXPECT_NE(LspTest::content.size(), 0);
}

TEST_F(LspTest, TestLspClientHover) {
	hoverresponse_t response = Hover(file, 27-1, 12);
	EXPECT_NE(response.jsonrpc, "");

	auto expected = "void DidOpen(const std::string &file, std::string *text)";
	auto got = response.result.contents.value;
	EXPECT_EQ(expected, got);
}

TEST_F(LspTest, TestLspClientCompletion) {
	completionresponse_t response = Completion(file, 27-1, 12);
	EXPECT_NE(response.jsonrpc, "");

	auto expected = "DidOpen";
	EXPECT_EQ(response.result.items.size(), 1);
	EXPECT_EQ(expected, response.result.items[0].label);
}

TEST_F(LspTest, TestLspClientDefinition) {
	definitionresponse_t response = Definition(file, 27-1, 12);
	EXPECT_NE(response.jsonrpc, "");

	EXPECT_EQ(response.result.size(), 1);
	EXPECT_NE(response.result[0].uri, "");
}

TEST_F(LspTest, TestLspClientSignatureHelp) {
	signaturehelpresponse_t response = SignatureHelp(file, 27-1, 12);
	EXPECT_NE(response.jsonrpc, "");

	auto expected = "DidOpen(const std::string &file, std::string *text)";
	EXPECT_EQ(response.result.signatures.size(), 1);
	EXPECT_EQ(expected, response.result.signatures[0].label);
}

TEST_F(LspTest, TestLspClientReferences) {
	referencesresponse_t response = References(file, 27-1, 12);
	EXPECT_NE(response.jsonrpc, "");

	EXPECT_EQ(response.result.size(), 2);
	EXPECT_NE(response.result[0].uri, "");
}
