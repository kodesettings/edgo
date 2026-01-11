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

#ifndef _LANGS_H_
#define _LANGS_H_

#include <string>
#include <map>
#include <vector>

//---------------------------------------------------------------

#include "lang_javascript.h"
#include "lang_python.h"
#include "lang_rust.h"
#include "lang_go.h"
#include "lang_c.h"
#include "lang_cpp.h"
#include "lang_css.h"
#include "lang_html.h"
#include "lang_java.h"
#include "lang_bash.h"

typedef struct {
	std::string query;
} query_t;

// Supported languages for syntax highlighting
enum languages {
	BASH          = 0x034,
	C             = 0x035,
	CPP           = 0x036,
	GO            = 0x037,
	HTML          = 0x038,
	CSS           = 0x039,
	JAVA          = 0x040,
	JAVASCRIPT    = 0x041,
	TYPESCRIPT    = 0x042,
	PYTHON        = 0x043,
	RUST          = 0x044,
};

static std::map<enum languages, std::string> languages = {
	{JAVASCRIPT, QUERY_JAVASCRIPT},
	{TYPESCRIPT, QUERY_JAVASCRIPT},
	{PYTHON,     QUERY_PYTHON},
	{RUST,       QUERY_RUST},
	{GO,         QUERY_GO},
	{C,          QUERY_C},
	{CPP,        QUERY_CPP},
	{CSS,        QUERY_CSS},
	{HTML,       QUERY_HTML},
	{JAVA,       QUERY_JAVA},
	{BASH,       QUERY_BASH},
};

static std::string MatchQueryLang(enum languages &lang) {
	return languages[lang];
}

#endif // _LANGS_H_
