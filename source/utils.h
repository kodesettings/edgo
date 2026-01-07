#ifndef _UTILS_H_
#define _UTILS_H_

#include <vector>
#include <string>

// Matched characters
std::vector<char> matched {
	' ', '.', ',', '=', '+', '-', '[', '(', '{', ']', ')', '}', 
	'"', ':', '&', '?','!',';','\t', '/','<','>'
};

// Inner struct to contain a line of text
typedef struct {
	std::string buf;
} line_t;

typedef std::vector<line_t> line_v;
typedef std::vector<std::string> string_v;

line_v GetLinesArrayFromData(std::string data, int lineNum);
int LineOffset(std::string text, int lineNum);


#endif // _UTILS_H_
