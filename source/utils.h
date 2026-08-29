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
#include <fnmatch.h>
#include <boost/filesystem.hpp>

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

// screen output memory
static char output[44880];

#include "../include/edgo.h"
static struct screen screen;

// ------------------------------------------------------------------
// screen helper methods

#include "screen.h"

#define __export_screen ExportScreen
#define __build_line_vec BuildLineVec

// ------------------------------------------------------------------
// other helper methods

std::string ReadFileToString(const std::string &filepath);
bool SaveToFile(const std::string &filepath, const std::string &content);
int LineOffset(const std::string &text, int lineNum);
std::string RemoveLeadingTabsSpaces(const std::string &s);
int FindCharacterOccurances(const std::string &s, const char ch);
bool IsIgnored(const std::string &path, string_v ignorePatterns);
bool IsMatchExt(const std::string &path, string_v ignoreExts);
void ColorizeTextChunk(int fg, int bg, std::string *str);

#endif // _UTILS_H_
