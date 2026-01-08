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

#include "selection.h"
#include "operations.h"
#include <boost/algorithm/string.hpp>

void selection::CleanSelection(void) {
	this->is_selected = false;
	this->ssx = -1, this->ssy = -1, this->sex = -1, this->sey = -1;
}

bool selection::IsSelectionNonEmpty(void) {
	if (this->ssx == -1 || this->ssy == -1  || this->sex == -1 || this->sey == -1) { return false; }
	if (Equal(this->ssx, this->ssy, this->sex, this->sey)) {
		return false;
	}

	return true;
}

bool selection::IsUnderSelection(int x, int y) {
	// Check if there is an active selection
	if (this->ssx == -1 || this->ssy == -1  || this->sex == -1 || this->sey == -1) { return false; }

	int startx = this->ssx, starty = this->ssy;
	int endx = this->sex, endy = this->sey;

	if (GreaterThan(startx, starty, endx, endy)) {
		startx = endx, endx = startx;
		starty = endy, endy = starty;
	}

	return GreaterEqual(x, y, startx, starty) && LessThan(x, y, endx, endy);
}

bool selection::GreaterThan(int x, int y, int x1, int y1) {
	if (y > y1) {
		return true;
	}
	return y == y1 && x > x1;
}

bool selection::LessThan(int x, int y, int x1, int y1) {
	if (y < y1) {
		return true;
	}
	return y == y1 && x < x1;
}

bool selection::GreaterEqual(int x, int y, int x1, int y1) {
	if (y > y1) {
		return true;
	}
	if (y == y1 && x >= x1) {
		return true;
	}
	return false;
}

bool selection::Equal(int x, int y, int x1, int y1) {
	return x == x1 && y == y1;
}

std::vector<selection::indices_t> selection::GetSelectedIndices(std::string text) {
	std::vector<indices_t> selectedIndices;

	// check for empty selection
	if (Equal(this->ssx, this->ssy, this->sex, this->sey)) {
		return selectedIndices;
	}

	// getting selection Start point
	int startx = this->ssx, starty = this->ssy;
	int endx = this->sex, endy = this->sey;

	// swap points if selection is inversed
	if (GreaterThan(startx, starty, endx, endy)) {
		startx = endx, endx = startx;
		starty = endy, endy = starty;
	}

	std::vector<std::string> lines;
	boost::split(lines, text, boost::is_any_of("\n"));
	bool inside = false;

	// iterate over lines, starting from selection Start point until out ouf selection
	for (int j = starty; j < (int)lines.size(); j++) {
		if (lines[j].size() == 0 && this->IsUnderSelection(0, j)) {
			selectedIndices.push_back(indices_t{x: 0, y: j});
			inside = true;
		}
		for (int i = 0; i < (int)lines[j].size(); i++) {
			if (this->IsUnderSelection(i, j)) {
				selectedIndices.push_back(indices_t{x: i, y: j});
				inside = true;
			} else  {
				if (inside == true) { // first time when out ouf selection
					return selectedIndices;
				}
			}
		}
	}
	return selectedIndices;
}

std::string selection::GetSelectionString(std::string text) {
	std::string ret;
	bool in = false;

	// check for empty selection
	if (Equal(this->ssx, this->ssy, this->sex, this->sey)) { return ""; }

	// getting selection Start point
	int startx = this->ssx, starty = this->ssy;
	int endx = this->sex, endy = this->sey;

	if (GreaterThan(startx, starty, endx, endy)) {
		startx = endx, endx = startx; // swap  points if selection inverse
		starty = endy, endy = starty;
	}

	std::vector<std::string> lines;
	boost::split(lines, text, boost::is_any_of("\n"));

	for (int j = starty; j < (int)lines.size(); j++) {
		std::string line = lines[j];
		for (int i = 0; i < (int)line.size(); i++) {
			// if inside selection
			if (GreaterEqual(i, j, startx, starty) && LessThan(i, j, endx, endy)) {
				ret.append(1, line.at(i));
				in = true;
			} else {
				in = false;
				// only one selection area can be, early return
				if (ret.size() > 0) {
					return std::string(ret);
				}
			}
		}
		if (in && LessThan(0, j, endx, endy)) {
			ret.append(1, '\n');
		}
	}

	if (ret.size() > 0 && ret[ret.size()-1] == '\n') { ret = ret.substr(0, ret.size()-1); }
	return std::string(ret);
}

std::set<int> selection::GetSelectedLines(std::string text) {
	std::set<int> lineNumbers;
	bool in = false;

	// check for empty selection
	if (Equal(this->ssx, this->ssy, this->sex, this->sey)) { return lineNumbers; }

	// getting selection Start point
	int startx = this->ssx, starty = this->ssy;
	int endx = this->sex, endy = this->sey;

	if (GreaterThan(startx, starty, endx, endy)) {
		startx = endx, endx = startx; // swap  points if selection inverse
		starty = endy, endy = starty;
	}

	std::vector<std::string> lines;
	boost::split(lines, text, boost::is_any_of("\n"));

	for (int j = starty; j < (int)lines.size(); j++) {
		std::string line = lines[j];
		for (int i = 0; i < (int)line.size(); i++) {
			// if inside selection
			if (GreaterEqual(i, j, startx, starty) && LessThan(i, j, endx, endy)) {
				lineNumbers.insert(j);
				in = true;
			} else {
				in = false;
				// only one selection area can be, early return
				if (lineNumbers.size() > 0) {
					return lineNumbers;
				}
			}
		}
		if (in && LessThan(0, j, endx, endy)) {
			lineNumbers.insert(j);
		}
	}
	return lineNumbers;
}
