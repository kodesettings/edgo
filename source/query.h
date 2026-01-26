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

#ifndef _QUERY_H_
#define _QUERY_H_

#include "highlighter.h"

typedef struct {
	std::string name; // NOTE: unused for now
	std::string filename;
	int line;
} testdata_t;

typedef std::map<int, testdata_t> testdata_m;

// NOTE: initialize with SetLang() first, then apply these functions
inline void TestQueryAlloc(void) {
	uint32_t err_offset;
	TSQueryError err_type;
	const auto queryLang = MatchTestQueryLang(h.lang);
	h.testquery = ts_query_new(h.language, queryLang.c_str(), queryLang.size(), &err_offset, &err_type);
}

void TestFinder(const std::string &filename, testdata_m *tests);

#endif // _QUERY_H_
