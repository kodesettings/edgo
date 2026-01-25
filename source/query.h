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
	std::string name;
	std::string filename;
	int line;
} testdata_t;

typedef struct {
	TSQuery *testquery;
	std::string lang;
} testfinder_t;

void TestQuery(TSQuery *testquery);
typedef std::map<int, testdata_t> testdata_m;
void TestFinder(testfinder_t root, std::string filename, std::string code, testdata_m *tests);

#endif // _QUERY_H_
