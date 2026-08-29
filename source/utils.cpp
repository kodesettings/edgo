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

void ColorizeTextChunk(int fg, int bg, std::string *str) {
	int r, g, b;
	std::stringstream ss;

	ColorToRGB(fg, &r, &g, &b);
	ss << "\033[38;2;" << r << ";" << g << ";" << b;
	if (bg > 0) {
		ColorToRGB(bg, &r, &g, &b);
		ss << ";48;2;" << r << ";" << g << ";" << b;
	}
	ss << "m" << *str << "\033[0m";
	*str = ss.str();
}
