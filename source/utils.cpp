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

#include "utils.h"
#include "editor.h"
#include "screen_logo.h"

// ------------------------------------------------------------------
// io helper methods

std::string ReadFileToString(const std::string &filepath) { 
	std::ifstream file(filepath, std::ios::in);
	std::string str, fcontent;

	while (std::getline(file, str)) {
		fcontent.append(str).append("\n");
	}

	file.close();
	return fcontent;
}

bool SaveToFile(const std::string &filepath, const std::string &content) {
	std::ofstream file(filepath, std::ios::out);

	if (!file.is_open()) {
		return false;
	}

	file << content;
	file.close();
	return true;
}

// ------------------------------------------------------------------
// screen helper methods

line_v __build_line_vec(const std::string &data, int lineNum, bool colorize) {
	line_v lines;
	std::string line;

	colorindexer_t indexer;
	indexer.ranges = ColorRanges(e.y, e.y + e.TERMINAL_HEIGHT);
	indexer.counter = 0;

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

char* __export_screen(const std::string &s, const int offset) {
	line_v lines = __build_line_vec(s, -1, true);
	std::stringstream ss;

	for (int index = offset;
		index < (int)lines.size() && index < offset + e.ROWS;
		index++) {
		ss << lines[index].buf << "\n";
	}

	const std::string out = ss.str();
	static char output[256000];
	strcpy(output, out.c_str());
	return output;
}

// ------------------------------------------------------------------
// editor helper methods

int LineOffset(const std::string &text, int lineNum) {
	int index = 0;
	for (int i = 0; i < (int)text.length(); i++) {
		if ((char)text.at(i) == '\n') {
			index++;
		}

		if (index == lineNum) {
			return i; // start of next line
		}
	}

	return text.length(); // end of text
}

// Function to remove leading tabs and spaces
std::string RemoveLeadingTabsSpaces(const std::string &s) {
	for (int i = 0; i < (int)s.length(); i++) {
		if (s[i] != '\t' && s[i] != ' ') {
			return s.substr(i);
		}
	}

	return "";
}

// Find occurances of a character in the string
int FindCharacterOccurances(const std::string &s, const char ch) {
	int counter = 0;
	size_t pos = 0;

	while ((pos = s.find(ch, pos)) != std::string::npos) {
		counter++;
		pos += 1;
	}

	return counter;
}

bool IsIgnored(const std::string &path, string_v ignorePatterns) {
	for (auto pattern : ignorePatterns) {
		boost::filesystem::path p(path);
		const auto file = p.filename().string();
		bool match = fnmatch(pattern.c_str(), file.c_str(), 0) == 0; // POSIX only
		if (match) { return true; }
	}

	return false;
}

bool IsMatchExt(const std::string &path, string_v ignoreExts) {
	for (auto ignoreExt : ignoreExts) {
		boost::filesystem::path p(path);
		auto ext = p.extension().string();
		if (ext == ignoreExt) { return true; }
	}

	return false;
}

void ColorizeTextChunk(uint32_t color, std::string *str) {
	int r, g, b;
	std::stringstream ss;

	ColorToRGB(color, &r, &g, &b);
	ss << "\033[38;2;" << r << ";" << g << ";" << b << "m" << *str << "\033[0m";
	*str = ss.str();
}
