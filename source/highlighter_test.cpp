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

#include "highlighter.h"
#include "utils.h"
#include "test_languages.h"
#include <gtest/gtest.h>

TEST(HighlighterTests, TestTreeSitterGo) {
	auto sourcecode = GO_SAMPLE;
	SetLang(GO); // Setup new parser
	Parse(sourcecode); // Parse source code
	EXPECT_NE(h.tree, nullptr);
}

TEST(HighlighterTests, TestTreeSitterPython) {
	auto sourcecode = PYTHON_SAMPLE;
	SetLang(PYTHON); // Setup new parser
	Parse(sourcecode); // Parse source code
	EXPECT_NE(h.tree, nullptr);
}

TEST(HighlighterTests, TestTreeSitterJs) {
	auto sourcecode = JAVASCRIPT_SAMPLE;
	SetLang(JAVASCRIPT); // Setup new parser
	Parse(sourcecode); // Parse source code
	EXPECT_NE(h.tree, nullptr);
}

TEST(HighlighterTests, TestTreeSitterJsEdit) {
	auto sourcecode = JAVASCRIPT_SAMPLE;
	SetLang(JAVASCRIPT); // Setup new parser
	Parse(sourcecode); // Parse source code
	EXPECT_NE(h.tree, nullptr);

	TSNode node = ts_tree_root_node(h.tree);
	EXPECT_EQ(ts_node_is_null(node), false);

	// Edit input
	InsertTextEdit("2", 14, 1);
	node = ts_tree_root_node(h.tree);
	EXPECT_EQ(ts_node_is_null(node), false);
}

TEST(HighlighterTests, TestTreeSitterJsEditDelete) {
	auto sourcecode = JAVASCRIPT_SAMPLE;
	SetLang(JAVASCRIPT); // Setup new parser
	Parse(sourcecode); // Parse source code
	EXPECT_NE(h.tree, nullptr);

	TSNode node = ts_tree_root_node(h.tree);
	EXPECT_EQ(ts_node_is_null(node), false);

	// Edit input
	RemoveTextEdit("lo", 14, 2);
	node = ts_tree_root_node(h.tree);
	EXPECT_EQ(ts_node_is_null(node), false);
}

TEST(HighlighterTests, TestTreeSitterJsEditDeleteMultiple) {
	auto sourcecode = JAVASCRIPT_SAMPLE;
	SetLang(JAVASCRIPT); // Setup new parser
	Parse(sourcecode); // Parse source code
	EXPECT_NE(h.tree, nullptr);

	TSNode node = ts_tree_root_node(h.tree);
	EXPECT_EQ(ts_node_is_null(node), false);

	// Edit input
	RemoveTextEdit("lo", 14, 2);
	node = ts_tree_root_node(h.tree);
	EXPECT_EQ(ts_node_is_null(node), false);

	// Edit input
	RemoveTextEdit("el", 12, 2);
	node = ts_tree_root_node(h.tree);
	EXPECT_EQ(ts_node_is_null(node), false);
}

TEST(HighlighterTests, TestTreeSitterJsEditEnter) {
	auto sourcecode = JAVASCRIPT_SAMPLE;
	SetLang(JAVASCRIPT); // Setup new parser
	Parse(sourcecode); // Parse source code
	EXPECT_NE(h.tree, nullptr);

	TSNode node = ts_tree_root_node(h.tree);
	EXPECT_EQ(ts_node_is_null(node), false);

	// Edit input
	InsertTextEdit("\n", 19, 1);
	node = ts_tree_root_node(h.tree);
	EXPECT_EQ(ts_node_is_null(node), false);
}
