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
#include "langs.h"

// NOTE: initialize with SetLang() first, then apply these functions
void TestQuery(TSQuery *testquery) {
	uint32_t  err_offset;
	TSQueryError err_type;
	const auto queryLang = MatchTestQueryLang(h.lang);
	testquery = ts_query_new(h.language, queryLang.c_str(), queryLang.size(), &err_offset, &err_type);
}

void TestFinder(testfinder_t root, std::string filename, std::string code, testdata_m *tests) {
	TSQueryCursor *query_cursor = ts_query_cursor_new();
	TSQueryMatch match;
	while (ts_query_cursor_next_match(query_cursor, &match)) {
		for (int i = 0; i < match.capture_count; i++) {
			const TSQueryCapture capture = match.captures[i];
			uint32_t length;

			// find the node
			auto nodename = ts_query_capture_name_for_id(root.testquery, i, &length);
			bool istestfound = std::string(nodename).compare("test-name") == 0;
			if (root.testquery) ts_query_delete(root.testquery); // free query

			// get this data if found
			if (istestfound) {
				int line = ts_node_start_point(capture.node).row;
				(*tests)[line] = testdata_t{
					name: code,
					filename: filename,
					line: line,
				};
			}
		}
	}
	ts_query_cursor_delete(query_cursor); // free query cursor
}
