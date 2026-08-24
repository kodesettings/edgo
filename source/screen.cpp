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

#include "search.h"
#include "utils.h"
#include "screen.h"
#include "highlighter.h"
#include "editor.h"

#define LONG_W std::setw(20)
#define SHORT_W std::setw(10)

std::string OnLangLinesCount(const std::string &dirpath) {
	int fpc, tlc;

	// fetch all the files for the provided directory path
	auto files = LinesCountOnDir(dirpath, &fpc, &tlc);

	langlinescountresult_v results;

	// extract the all the filesystem data for each languages
	LangCount(files, &results);

	std::stringstream report;
	report << "      " << "Language and File Statistics Report" << std::endl;
	report << std::endl;
	report << LONG_W << "Total Files : " << SHORT_W << fpc << std::endl;
	report << LONG_W << "Total Rows  : " << SHORT_W << tlc << std::endl;
	report << std::endl;
	report << LONG_W << "Language" << LONG_W << "Files";
	report << LONG_W << "Lines" << LONG_W << "Empty/Code" << std::endl;
	report << std::endl;

	int index = 0;
row_0:
	auto r = results[index];
	report << LONG_W << r.lang;
	report << LONG_W << r.filescount;
	report << LONG_W << r.linescount;
	report << LONG_W << r.emptylinescount;
	report << std::endl;

	if (index < (int)results.size() - 1) { index++; goto row_0; }
	return report.str();
}

line_v BuildLineVec(const std::string &data, int lineNum, bool colorize) {
	colorindexer_t indexer;
	line_v lines;
	std::string line;

	if (colorize) {
		indexer.ranges = ColorRanges(e.y, e.y + e.TERMINAL_HEIGHT);
		indexer.counter = 0;
	}

	for (auto b : data) {
		if (b == '\n') {
			lines.push_back(line_t{line});
			line.erase();
			indexer.counter++;
		} else if (lineNum == (int)lines.size()) {
			break;
		} else {
			Colorize(b, &line, colorize, &indexer);
		}
	}

	if (!line.empty())
		lines.push_back(line_t{line}); // add last line

	return lines;
}

struct screen ExportScreen(const std::string &s, int offset, bool colorize) {
	line_v lines = __build_line_vec(s, -1, colorize);
	std::stringstream ss;

	for (int index = offset;
		index < (int)lines.size() && index < offset + e.ROWS;
		index++) {
		ss << lines[index].buf << "\n";
	}

	// screen content
	memset(output, 0, sizeof(output));
	strcpy(output, ss.str().c_str());

	// screen fields
	screen.filename = e.filename.data();
	screen.language = e.lang.name.data();
	screen.content = output;
	screen.cursor_x = e.col;
	screen.cursor_y = e.row;
	screen.changed = e.isContentChanged;
	return screen;
}
