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
#include "screen_logo.h"

#include <iomanip>
using namespace std;

string OnLangLinesCount(const string &dirpath) {
	int filesprocessedcount, totallinescount;

	// fetch all the files for the provided directory path
	auto files = LinesCountOnDir(dirpath, &filesprocessedcount, &totallinescount);

	langlinescountresult_v results;

	// extract the all the filesystem data for each languages
	LangCount(files, &results);

	stringstream report;
	report << setw(10) << "Language Lines Report" << endl;
	report << setw(10) << "Total Files : " << setw(5) << filesprocessedcount << endl;
	report << setw(10) << "Total Rows  : " << setw(5) << totallinescount << endl;
	report << setw(10) << "Language" << setw(10) << "Files" << setw(10) << "Lines" << setw(10) << "Empty/Code" << endl;
	report << endl;
	int index = 0;
row:
	auto r = results[index];
	report << setw(10) << r.lang << setw(10) << r.filescount << setw(10) << r.linescount << setw(10) << r.emptylinescount << endl;
	if (index < (int)results.size()) { index++; goto row; }
	return report.str();
}
