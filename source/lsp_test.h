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
#ifndef _LSP_TEST_H_
#define _LSP_TEST_H_

#include "lsp_client.h"
#include "utils.h"
#include <gtest/gtest.h>

class LspTest : public ::testing::Test {
protected:
	static void SetUpTestSuite() {
		lspclient.lang = "cpp";
		started = StartLspClient("clangd", "");

		// finding current working directory
		InitLspClient("/usr/include");

		// reading file content
		filepath.append("/usr/include/threads.h");
		content = ReadFileToString(filepath);

		// accessing didOpen event
		DidOpen(filepath, content);
	}

	static void TearDownTestSuite() { lspclient.isReady = false; }
public:
	static bool started;             // flag whether the server started
	static std::string filepath;     // direct file path
	static std::string content;      // file content
};

#endif // _LSP_TEST_H_
