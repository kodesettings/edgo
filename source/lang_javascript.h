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

#ifndef _LANG_JAVASCRIPT_H
#define _LANG_JAVASCRIPT_H


#define QUERY_JAVASCRIPT "\
; Special identifiers \
;-------------------- \
\
([ \
    (identifier) \
    (shorthand_property_identifier) \
    (shorthand_property_identifier_pattern) \
 ] @constant \
 (#match? @constant \"^[A-Z_][A-Z\\d_]+$\")) \
\
\
((identifier) @constructor \
 (#match? @constructor \"^[A-Z]\")) \
\
((identifier) @variable.builtin \
 (#match? @variable.builtin \"^(arguments|module|console|window|document)$\") \
 (#is-not? local)) \
\
((identifier) @function.builtin \
 (#eq? @function.builtin \"require\") \
 (#is-not? local)) \
\
; Function and method definitions \
;-------------------------------- \
\
(function_expression \
  name: (identifier) @function) \
(function_declaration \
  name: (identifier) @function) \
(method_definition \
  name: (property_identifier) @function.method) \
\
(pair \
  key: (property_identifier) @function.method \
  value: [(function_expression) (arrow_function)]) \
\
(assignment_expression \
  left: (member_expression \
    property: (property_identifier) @function.method) \
  right: [(function_expression) (arrow_function)]) \
\
(variable_declarator \
  name: (identifier) @function \
  value: [(function_expression) (arrow_function)]) \
\
(assignment_expression \
  left: (identifier) @function \
  right: [(function_expression) (arrow_function)]) \
\
; Function and method calls \
;-------------------------- \
\
(call_expression \
  function: (identifier) @function) \
\
(call_expression \
  function: (member_expression \
    property: (property_identifier) @function.method)) \
\
; Variables \
;---------- \
\
(identifier) @variable \
\
; Properties \
;----------- \
\
(property_identifier) @property \
\
; Literals \
;--------- \
\
(this) @variable.builtin \
(super) @variable.builtin \
\
[ \
  (true) \
  (false) \
  (null) \
  (undefined) \
] @constant.builtin \
\
(comment) @comment \
\
[ \
  (string) \
  (template_string) \
] @string \
\
(regex) @string.special \
(number) @number \
\
; Tokens \
;------- \
\
(template_substitution \
  \"${\" @punctuation.special \
  \"}\" @punctuation.special) @embedded \
\
[ \
  \";\" \
  (optional_chain) \
  \".\" \
  \",\" \
] @punctuation.delimiter \
\
[ \
  \"-\" \
  \"--\" \
  \"-=\" \
  \"+\" \
  \"++\" \
  \"+=\" \
  \"*\" \
  \"*=\" \
  \"**\" \
  \"**=\" \
  \"/\" \
  \"/=\" \
  \"%\" \
  \"%=\" \
  \"<\" \
  \"<=\" \
  \"<<\" \
  \"<<=\" \
  \"=\" \
  \"==\" \
  \"===\" \
  \"!\" \
  \"!=\" \
  \"!==\" \
  \"=>\" \
  \">\" \
  \">=\" \
  \">>\" \
  \">>=\" \
  \">>>\" \
  \">>>=\" \
  \"~\" \
  \"^\" \
  \"&\" \
  \"|\" \
  \"^=\" \
  \"&=\" \
  \"|=\" \
  \"&&\" \
  \"||\" \
  \"??\" \
  \"&&=\" \
  \"||=\" \
  \"\?\?=\" \
] @operator \
\
[ \
  \"(\" \
  \")\" \
  \"[\" \
  \"]\" \
  \"{\" \
  \"}\" \
]  @punctuation.bracket \
\
[ \
  \"as\" \
  \"async\" \
  \"await\" \
  \"break\" \
  \"case\" \
  \"catch\" \
  \"class\" \
  \"const\" \
  \"continue\" \
  \"debugger\" \
  \"default\" \
  \"delete\" \
  \"do\" \
  \"else\" \
  \"export\" \
  \"extends\" \
  \"finally\" \
  \"for\" \
  \"from\" \
  \"function\" \
  \"get\" \
  \"if\" \
  \"import\" \
  \"in\" \
  \"instanceof\" \
  \"let\" \
  \"new\" \
  \"of\" \
  \"return\" \
  \"set\" \
  \"static\" \
  \"switch\" \
  \"target\" \
  \"throw\" \
  \"try\" \
  \"typeof\" \
  \"var\" \
  \"void\" \
  \"while\" \
  \"with\" \
  \"yield\" \
] @keyword \
"

// QUERIES: syntax definition for test identification
// -------------------------------------------------------------------------

#define QUERY_JAVASCRIPT_TEST "\
                (expression_statement\
                (call_expression\
                  function: (identifier) @method-name\
                  (#match? @method-name \"^(describe|test|it)\")\
                  arguments: (arguments [\
                        ((string) @test-name)\
                        ((template_string) @test-name)\
                  ]\
                )))\
"

#endif // _LANG_JAVASCRIPT_H
