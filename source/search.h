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

#ifndef _SEARCH_H_
#define _SEARCH_H_

#include <boost/filesystem.hpp>
#include <sstream>
#include <algorithm>
#include "utils.h"

namespace fs = boost::filesystem;

typedef struct  {
	int line;
	int position;
} searchresult_t;
typedef std::vector<searchresult_t> searchresult_v;

searchresult_v Search(const line_v lines, const std::string &pattern);
searchresult_v SearchOnFile(const std::string &filename, const std::string &pattern, int *lineindex);
void LineCountOnFile(const std::string &filename, int *lines, int *emptylines);

typedef struct {
	std::string file;
	searchresult_v results;
} filesearchresult_t;
typedef std::vector<filesearchresult_t> filesearchresult_v;

void ParsePattern(const std::string &pattern, string_v *extensions);
filesearchresult_v SearchOnDir(const std::string &dir, const std::string pattern,
                               int *filesprocessedcount, int *totallinescount);

// ignoring dirs on search dir
static std::vector<std::string> ignoredirs {
	".git", ".idea", "node_modules", "dist", "target", "__pycache__", ".pytest_cache", "build",
	".DS_Store", ".venv", "venv",
};

// files ignored on global search
static std::vector<std::string> ignoreexts {
	".doc", ".docx", ".pdf", ".rtf", ".odt", ".xlsx", ".pptx",
	".jpg", ".png", ".gif", ".bmp", ".svg", ".tiff",
	".mp3", ".wav", ".aac", ".flac", ".ogg",
	".mp4", ".avi", ".mov", ".wmv", ".mkv",
	".zip", ".rar", ".tar.gz", ".7z",
	".exe", ".msi", ".bat", ".sh",
	".ttf", ".otf",
};

typedef struct {
	std::string file;
	int count;
	int emptycount;
} linescountresult_t;
typedef std::vector<linescountresult_t> linescountresult_v;

linescountresult_v LinesCountOnDir(const std::string &dir, int *filesprocessedcount,
                                   int *totallinescount);

typedef struct {
	std::string lang;
	int filescount;
	int linescount;
	int emptylinescount;
} langlinescountresult_t;
typedef std::vector<langlinescountresult_t> langlinescountresult_v;

// counting number of languages in folder
void LangCount(linescountresult_v files, langlinescountresult_v *results);

#endif // _SEARCH_H_
