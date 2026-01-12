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
#include "logging.h"

treesitterhighlighter_t h;

void SetTheme(const std::string &themepath) {/* setting theme if needed */ }

void GetSitterLang(enum languages lang, bool *parsed) {
	switch (lang) {
	case BASH: *parsed = ts_parser_set_language(h.parser, tree_sitter_bash()); break;
	case C: *parsed = ts_parser_set_language(h.parser, tree_sitter_c()); break;
	case CPP: *parsed = ts_parser_set_language(h.parser, tree_sitter_cpp()); break;
	case GO: *parsed = ts_parser_set_language(h.parser, tree_sitter_go()); break;
	case HTML: *parsed = ts_parser_set_language(h.parser, tree_sitter_html()); break;
	case CSS: *parsed = ts_parser_set_language(h.parser, tree_sitter_css()); break;
	case JAVA: *parsed = ts_parser_set_language(h.parser, tree_sitter_java()); break;
	case JAVASCRIPT: *parsed = ts_parser_set_language(h.parser, tree_sitter_javascript()); break;
	case TYPESCRIPT: *parsed = ts_parser_set_language(h.parser, tree_sitter_javascript()); break;
	case PYTHON: *parsed = ts_parser_set_language(h.parser, tree_sitter_python()); break;
	case RUST: *parsed = ts_parser_set_language(h.parser, tree_sitter_rust()); break;
	default: *parsed = false; break;
	}
}

void SetLang(enum languages lang) {
	if (h.lang == lang) { return; }
	h.lang = lang;

	// allocating parser
	if (h.parser == nullptr) { h.parser = ts_parser_new(); }

	// parsing the language
	bool parsed = false;
	GetSitterLang(lang, &parsed);
	if (!parsed) LOG(ERROR) << "language is not supported";

	uint32_t err_offset;
	TSQueryError err_type;

	const auto queryLang = MatchQueryLang(h.lang);
	h.query = ts_query_new(h.language, queryLang.c_str(), queryLang.size(), &err_offset, &err_type);
	if (err_type != TSQueryError::TSQueryErrorNone) { 
		LOG(FATAL) << "could not parse query language, exiting ...";
		exit(-1);
	}
}

void InsertTextEdit(const std::string &code, int offset, int length) {
	auto input_edit = TSInputEdit{
		start_byte: uint32_t(offset),
		old_end_byte: uint32_t(offset),
		new_end_byte: uint32_t(offset) + uint32_t(length),
		start_point:  TSPoint{row: 0, column: 0},
		old_end_point: TSPoint{row: 0, column: 0},
		new_end_point: TSPoint{row: 0, column: 0},
	};
	if (h.tree != nullptr) { ts_tree_edit(h.tree, &input_edit); }
	h.tree = ts_parser_parse_string(h.parser, h.tree, code.c_str(), length);
	LOG(INFO) << "inserttextedit offset:" << offset << " len:" << length;
}

void RemoveTextEdit(const std::string &code, int offset, int length) {
	auto input_edit = TSInputEdit{
		start_byte: uint32_t(offset),
		old_end_byte: uint32_t(offset) + uint32_t(length),
		new_end_byte: uint32_t(offset),
		start_point:  TSPoint{row: 0, column: 0},
		old_end_point: TSPoint{row: 0, column: 0},
		new_end_point: TSPoint{row: 0, column: 0},
	};
	if (h.tree != nullptr) { ts_tree_edit(h.tree, &input_edit); }
	h.tree = ts_parser_parse_string(h.parser, h.tree, code.c_str(), length);
	LOG(INFO) << "removetextedit offset:" << offset << " len:" << length;
}

// TODO: Convinience method from legacy screen logic,
// why is it needed and where?
void Parse(const std::string &code) {
	h.tree = ts_parser_parse_string(h.parser, h.tree, code.c_str(), code.size());
}

coloredbyterange_v ColorRanges(const int from, const int to) {
	TSPoint start_point {row: uint32_t(from), column: 0};
	TSPoint end_point {row: uint32_t(to), column: 0};

	// setting point range
	TSQueryCursor *query_cursor = ts_query_cursor_new();
	ts_query_cursor_set_point_range(query_cursor, start_point, end_point);

	coloredbyterange_v colors;
	TSQueryMatch match;
	while (ts_query_cursor_next_match(query_cursor, &match)) {
		for (int i = 0; i < match.capture_count; i++) {
			const TSQueryCapture capture = match.captures[i];
			uint32_t length;
			auto name = ts_query_capture_name_for_id(h.query, i, &length);
			auto name_str = std::string(name);
			auto split = name_str.substr(0, name_str.find("."));
			auto color = h.colorsmap[split];

			if (name_str.find("injection")) {
				// We don't colorize embedded content for different languages in the editor.
				// Only languages that can be identified for the entire source are colorized.
				// This usually means markdown or html files or any documentation related content.
				continue;
			}

			colors.push_back(coloredbyterange_t{
				startbyte: ts_node_start_byte(capture.node),
				endbyte:   ts_node_end_byte(capture.node),
				color:     GetColorCode(color),
			});
		}
	}

	ts_query_cursor_delete(query_cursor);
	return colors;
}

path_t GetNodePathAt(int startPointRow, int startPointColumn, int endPointRow, int endPointColumn) {
	TSPoint start_point {row: uint32_t(startPointRow), column: uint32_t(startPointColumn)};
	TSPoint end_point {row: uint32_t(endPointRow), column: uint32_t(endPointColumn)};

	// setting point range
	auto root_node = ts_tree_root_node(h.tree);
	auto node = ts_node_named_descendant_for_point_range(root_node, start_point, end_point);

	// add a single node
	const noderange_t node1{
		ts_node_start_point(node).row,
		ts_node_start_point(node).column,
		ts_node_end_point(node).row,
		ts_node_end_point(node).column
	};

	const path_t path{atx: uint32_t(startPointColumn), aty: uint32_t(startPointRow), nodes: {node1}, current: 0};
	if (node.id == nullptr) return path_t{};

	return path;
}

noderange_t GetNodeAt(int startPointRow, int startPointColumn, int endPointRow, int endPointColumn) {
	TSPoint start_point {row: uint32_t(startPointRow), column: uint32_t(startPointColumn)};
	TSPoint end_point {row: uint32_t(endPointRow), column: uint32_t(endPointColumn)};

	// setting point range
	auto root_node = ts_tree_root_node(h.tree);
	auto node = ts_node_named_descendant_for_point_range(root_node, start_point, end_point);

	// node range as return value in this method
	const noderange_t node_range {
		ts_node_start_point(node).row,
		ts_node_start_point(node).column,
		ts_node_end_point(node).row,
		ts_node_end_point(node).column
	};

	return node_range;
}
