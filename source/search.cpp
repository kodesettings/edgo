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
#include "highlighter.h"
#include "utils.h"

searchresult_v Search(const line_v lines, const std::string &pattern) {
	searchresult_v results;

	if (pattern.empty() || lines.empty()) { return results; }

	for (int i = 0; i < (int)lines.size(); i++) {
		int from = 0;
		auto line = lines[i].buf;
		for (;;) {
			int pos = line.substr(from).find(pattern);
			if (pos == (int)std::string::npos) { break; } else {
				pos = from + pos;
				results.push_back(searchresult_t{i, pos});
				from = pos + 1;
			}
		}
	}
	return results;
}

searchresult_v SearchOnFile(const std::string &filename, const std::string &pattern,
                            int *lineindex) {
	std::stringstream file(ReadFileToString(filename));
	searchresult_v results;

	*lineindex = 1;
	std::string line;
	while (std::getline(file, line)) {
		int pos = line.find(pattern);
		if (pos != (int)std::string::npos) {
			results.push_back(searchresult_t{*lineindex, pos});
		}
		(*lineindex)++;
	}

	return results;
}

void LineCountOnFile(const std::string &filename, int *lines, int *emptylines) {
	std::stringstream file(ReadFileToString(filename));
	searchresult_v results;

	*lines = 0;
	*emptylines = 0;
	std::string line;
	while (std::getline(file, line)) {
		if (line.empty()) { (*emptylines)++; };
		(*lines)++;
	}
}

void ParsePattern(const std::string &pattern, string_v *extensions) {
	if (pattern.find(" -f ") != std::string::npos) {
		auto split = SplitString(pattern, " -f ");
		auto fileExtensions = TrimString(split[0]);
		*extensions = SplitString(fileExtensions, ",");
	}
}

filesearchresult_v SearchOnDir(const std::string &dir, const std::string pattern,
                               int *filesprocessedcount, int *totallinescount) {
	string_v files;
	boost::system::error_code ec;

	string_v allowedExtensions;
	ParsePattern(pattern, &allowedExtensions);

	for (fs::recursive_directory_iterator it(dir), end; it != end; ++it) {
		const std::string abspath = it->path().string();
		if (it->is_directory(ec) && IsIgnored(abspath, ignoredirs)) continue;
		if (!it->is_directory(ec) && !IsMatchExt(abspath, ignoreexts)) {
			if (!allowedExtensions.empty()) {
				if (IsMatchExt(abspath, allowedExtensions)) {
					files.push_back(abspath);
				}
			} else {
				files.push_back(abspath);
			}
		}
	}

	if (ec) { return filesearchresult_v{}; }
	filesearchresult_v results;

	*totallinescount = 0;
	for (auto file : files) {
		int lineindex = 0;
		auto fileResults = SearchOnFile(file, pattern, &lineindex);
		*totallinescount += lineindex;
		if (!fileResults.empty()) {
			results.push_back(filesearchresult_t{file, fileResults});
		}
	}

	*filesprocessedcount = results.size();
	return results;
}


linescountresult_v LinesCountOnDir(const std::string &dir, int *filesprocessedcount,
                                   int *totallinescount) {
	string_v files;
	boost::system::error_code ec;

	for (fs::recursive_directory_iterator it(dir), end; it != end; ++it) {
		const std::string abspath = it->path().string();
		if (it->is_directory(ec) && IsIgnored(abspath, ignoredirs)) continue;
		if (!it->is_directory(ec) && !IsMatchExt(abspath, ignoreexts)) {
			files.push_back(abspath);
		}
	}

	if (ec) { return linescountresult_v{}; }
	linescountresult_v results;

	for (auto file : files) {
		int lines;
		int emptylines;
		LineCountOnFile(file, &lines, &emptylines);
		linescountresult_t result{file, lines, emptylines};
		results.push_back(result);
		*totallinescount += lines;
	}

	*filesprocessedcount = results.size();
	return results;
}

void LangCount(linescountresult_v files, langlinescountresult_v *results) {
	std::map<std::string, int> linescountresults, emptylinescountresults, filescountresults;

	for (auto lcresult : files) {
		auto lang = DetectLang(lcresult.file);
		linescountresults[lang] += lcresult.count;
		emptylinescountresults[lang] += lcresult.emptycount;
		filescountresults[lang]++;
	}

	// adding all the results
	for (const auto &kv : linescountresults) {
		results->push_back(langlinescountresult_t{
			kv.first,
			filescountresults[kv.first],
			kv.second,
			emptylinescountresults[kv.first]}
		);
	}

	// sorting the results vector based on values
	std::sort(results->begin(), results->end(), [](const auto &a, const auto &b) {
		return a.linescount > b.linescount; // descending
	});
}
