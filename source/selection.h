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

#ifndef _SELECTION_H_
#define _SELECTION_H_

#include <string>
#include <vector>

class selection {
public:
	int ssx;          // selection Start x
	int ssy;          // selection Start y
	int sex;          // selection end x
	int sey;          // selection end y
	bool is_selected; // true if selection is active
private:
	typedef struct {int x = 0; int y = 0; } indices_t;
public:
	void CleanSelection(void);
	bool IsSelectionNonEmpty(void);
	bool IsUnderSelection(int x, int y);
	bool GreaterThan(int x, int y, int x1, int y1);
	bool LessThan(int x, int y, int x1, int y1);
	bool GreaterEqual(int x, int y, int x1, int y1);
	bool Equal(int x, int y, int x1, int y1);
	std::vector<indices_t> GetSelectedIndices(std::string text);
	std::string GetSelectionString(std::string text);
	std::vector<int> GetSelectedLines(std::string text);
};

#endif // _SELECTION_H_
