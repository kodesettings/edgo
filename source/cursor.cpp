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

#include "editor.h"

void Focus(void) {
	if (e.row > e.y + e.ROWS) {
		e.y = e.row + e.ROWS;
	} else if (e.row < e.y) {
		e.y = e.row;
	}
}

void FocusCenter(void) {
//	e.Screen.Show()
	if (e.row > e.y + e.ROWS) {
		e.y = e.row + e.ROWS;
	} else if (e.row < e.y) {
		e.y = e.row;
	}

	e.y -= e.ROWS / 2;
	if (e.y < 0) { e.y = 0; }

	int centerRow = e.ROWS / 2;
	// Update the cursor row to the center row if necessary
	if (e.row - e.y > centerRow) {
		e.y += e.row - e.y - centerRow;
	}
}
