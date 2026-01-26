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

void TestFinder(const std::string &filename, testdata_m *tests) {
	static TSQueryCursor *query_cursor = ts_query_cursor_new();
	ts_query_cursor_exec(query_cursor, h.testquery, ts_tree_root_node(h.tree));
	TSQueryMatch match;
	while (ts_query_cursor_next_match(query_cursor, &match)) {
		for (int i = 0; i < match.capture_count; i++) {
			const TSQueryCapture capture = match.captures[i];
			uint32_t length;

			// find the node
			auto nodename = ts_query_capture_name_for_id(h.testquery, capture.index, &length);
			bool istestfound = std::string(nodename).compare("test-name") == 0;

			// get this data if found
			if (istestfound) {
				int line = ts_node_start_point(capture.node).row;
				(*tests)[line] = testdata_t{
					filename: filename,
					line: line,
				};
			}
		}
	}
}
