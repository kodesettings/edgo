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

#ifndef _LANG_GO_H_
#define _LANG_GO_H_

#define QUERY_GO "\
; Function calls \
 \
(call_expression \
  function: (identifier) @function.builtin \
  (.match? @function.builtin \"^(append|cap|close|complex|copy|delete|imag|len|make|new|panic|print|println|real|recover)$\")) \
\
(call_expression \
  function: (identifier) @function) \
\
(call_expression \
  function: (selector_expression \
    field: (field_identifier) @function.method)) \
\
; Function definitions \
\
(function_declaration \
  name: (identifier) @function) \
\
(method_declaration \
  name: (field_identifier) @function.method) \
\
; Identifiers \
\
(type_identifier) @type \
(field_identifier) @property \
(identifier) @variable \
\
; Operators \
\
[ \
  \"--\" \
  \"-\" \
  \"-=\" \
  \":=\" \
  \"!\" \
  \"!=\" \
  \"...\" \
  \"*\" \
  \"*\" \
  \"*=\" \
  \"/\" \
  \"/=\" \
  \"&\" \
  \"&&\" \
  \"&=\" \
  \"%\" \
  \"%=\" \
  \"^\" \
  \"^=\" \
  \"+\" \
  \"++\" \
  \"+=\" \
  \"<-\" \
  \"<\" \
  \"<<\" \
  \"<<=\" \
  \"<=\" \
  \"=\" \
  \"==\" \
  \">\" \
  \">=\" \
  \">>\" \
  \">>=\" \
  \"|\" \
  \"|=\" \
  \"||\" \
  \"~\" \
] @operator \
\
; Keywords \
\
[ \
  \"break\" \
  \"case\" \
  \"chan\" \
  \"const\" \
  \"continue\" \
  \"default\" \
  \"defer\" \
  \"else\" \
  \"fallthrough\" \
  \"for\" \
  \"func\" \
  \"go\" \
  \"goto\" \
  \"if\" \
  \"import\" \
  \"interface\" \
  \"map\" \
  \"package\" \
  \"range\" \
  \"return\" \
  \"select\" \
  \"struct\" \
  \"switch\" \
  \"type\" \
  \"var\" \
] @keyword \
\
; Literals \
\
[ \
  (interpreted_string_literal) \
  (raw_string_literal) \
  (rune_literal) \
] @string \
\
(escape_sequence) @escape \
\
[ \
  (int_literal) \
  (float_literal) \
  (imaginary_literal) \
] @number \
\
[ \
  (true) \
  (false) \
  (nil) \
  (iota) \
] @constant.builtin \
\
(comment) @comment \
"

// QUERIES: syntax definition for test identification
// -------------------------------------------------------------------------

#define QUERY_GO_TEST "\
         (\
          function_declaration name: (identifier) @test-name\
          (#match? @test-name \"Test*\")\
        )\
"

#endif // _LANG_GO_H_
