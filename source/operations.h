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

#ifndef _OPERATIONS_H_
#define _OPERATIONS_H_

#include <string>
#include <vector>

typedef enum {
	Insert     = 0x01, // insert text into the undo stack
	AddChar    = 0x02, // add one character into the undo stack
	Delete     = 0x03, // delete text from the undo stack
	Start      = 0x04, // flag for start operation
	End        = 0x05 // flag for end operation
} action_t;

typedef struct {
	int line;
	int column;
} cursormove_t;

typedef struct {
	action_t action;
	std::string text;
	int offset;
	cursormove_t cursor;
} operation_t;

typedef std::vector<operation_t> editoperation_t;

#endif // _OPERATIONS_H_
