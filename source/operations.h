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

typedef std::vector<operation_t> EditOperation;

#endif // _OPERATIONS_H_
