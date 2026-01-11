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

static std::map<std::string, std::string> languages = {
	{"javascript", QUERY_JAVASCRIPT},
	{"typescript", QUERY_JAVASCRIPT},
	{"python",     QUERY_PYTHON},
	{"rust",       QUERY_RUST},
	{"go",         QUERY_GO},
	{"c",          QUERY_C},
	{"c++",        QUERY_CPP},
	{"cpp",        QUERY_CPP},
	{"css",        QUERY_CSS},
	{"html",       QUERY_HTML},
	{"java",       QUERY_JAVA},
	{"bash",       QUERY_BASH},
};

inline std::string MatchQueryLang(std::string lang) {
	return languages[lang];
}

#endif // _LANGS_H_
