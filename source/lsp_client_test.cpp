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

#include "lsp_test.h"

bool LspTest::started;
std::string LspTest::filepath;
std::string LspTest::content;

TEST_F(LspTest, TestLspClientStarted) {
	EXPECT_EQ(LspTest::started, true);
	EXPECT_NE(LspTest::content.size(), 0);
}

TEST_F(LspTest, TestLspClientHover) {
	hoverresponse_t response = Hover(filepath, 223-1, 18);
	EXPECT_NE(response.jsonrpc, "");

	auto expected = "function tss_create";
	auto got = response.result.contents.value;
	EXPECT_NE(got.find(expected), std::string::npos);
}

TEST_F(LspTest, TestLspClientCompletion) {
	completionresponse_t response = Completion(filepath, 223-1, 18);
	EXPECT_NE(response.jsonrpc, "");

	auto expected = " tss_create";
	EXPECT_EQ(response.result.items.size(), 1);

	if (!response.result.items.empty()) {
		EXPECT_EQ(expected, response.result.items[0].label);
	}
}

TEST_F(LspTest, TestLspClientDefinition) {
	definitionresponse_t response = Definition(filepath, 223-1, 18);
	EXPECT_NE(response.jsonrpc, "");

	EXPECT_EQ(response.result.size(), 1);

	if (!response.result.empty()) {
		EXPECT_NE(response.result[0].uri, "");
	}
}

TEST_F(LspTest, TestLspClientSignatureHelp) {
	signaturehelpresponse_t response = SignatureHelp(filepath, 223-1, 34);
	EXPECT_NE(response.jsonrpc, "");
	EXPECT_EQ(response.id, 4);

	EXPECT_EQ(response.result.signatures.size(), 0);
}

TEST_F(LspTest, TestLspClientReferences) {
	referencesresponse_t response = References(filepath, 223-1, 34);
	EXPECT_NE(response.jsonrpc, "");
	EXPECT_EQ(response.id, 5);

	EXPECT_EQ(response.result.size(), 0);
}
