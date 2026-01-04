#ifndef _OPERATIONS_H_
#define _OPERATIONS_H_

#include <string>
#include <vector>

typedef enum {
	Insert = 0x01,
	AddChar = 0x02,
	Delete = 0x03,
	Start = 0x04,
	End = 0x05
} Action;

typedef struct {
	Action action;
	std::string text;
	int offset;
	CursorMove cursor;
} Operation;

typedef std::vector<Operation> EditOperation;

typedef struct {
	int line;
	int column;
} CursorMove;

#endif // _OPERATIONS_H_
