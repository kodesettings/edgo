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

#ifndef _UTILS_H_
#define _UTILS_H_

#include <vector>
#include <string>

// Matched characters
static std::vector<char> matched {
	' ', '.', ',', '=', '+', '-', '[', '(', '{', ']', ')', '}', 
	'"', ':', '&', '?','!',';','\t', '/','<','>'
};

// Inner struct to contain a line of text
typedef struct {
	std::string buf;
} line_t;

typedef std::vector<line_t> line_v;
typedef std::vector<std::string> string_v;

std::string ReadFileToString(const std::string &filepath);
line_v GetLinesArrayFromData(const std::string &data, int lineNum);
int LineOffset(const std::string &text, int lineNum);
std::string RemoveLeadingTabsSpaces(const std::string &s);
int FindCharacterOccurances(const std::string &s, const char ch);

#endif // _UTILS_H_
