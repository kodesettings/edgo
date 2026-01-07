#include "utils.h"

line_v GetLinesArrayFromData(std::string data, int lineNum) {
	line_v lines = line_v(lineNum);
	int row = 0;
	for (auto b : data) {
		if (b == '\n') {
			row++;
			continue;
		} else if (row == lineNum) {
			break;
		} else {
			lines[row].buf.append(&b);
		}
	}
	return lines;
}

int LineOffset(std::string text, int lineNum) {
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
