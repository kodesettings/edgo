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

#ifndef _LANG_JAVA_H_
#define _LANG_JAVA_H_

#define QUERY_JAVA "\
; Methods \
\
(method_declaration \
  name: (identifier) @function.method) \
(method_invocation \
  name: (identifier) @function.method) \
(super) @function.builtin \
\
; Annotations \
\
(annotation \
  name: (identifier) @attribute) \
(marker_annotation \
  name: (identifier) @attribute) \
\
; Types \
\
(interface_declaration \
  name: (identifier) @type) \
(class_declaration \
  name: (identifier) @type) \
(record_declaration \
  name: (identifier) @type) \
(enum_declaration \
  name: (identifier) @type) \
\
((field_access \
  object: (identifier) @type) \
 (#match? @type \"^[A-Z]\")) \
((scoped_identifier \
  scope: (identifier) @type) \
 (#match? @type \"^[A-Z]\")) \
\
(constructor_declaration \
  name: (identifier) @type) \
(compact_constructor_declaration \
  name: (identifier) @type) \
\
(type_identifier) @type \
\
[ \
  (boolean_type) \
  (integral_type) \
  (floating_point_type) \
  (floating_point_type) \
  (void_type) \
] @type.builtin \
\
(type_arguments \
  (wildcard \"?\" @type.builtin)) \
\
; Variables \
\
((identifier) @constant \
 (#match? @constant \"^_*[A-Z][A-Z\\d_]+$\")) \
\
(identifier) @variable \
\
(this) @variable.builtin \
\
; Literals \
\
[ \
  (hex_integer_literal) \
  (decimal_integer_literal) \
  (octal_integer_literal) \
  (binary_integer_literal) \
] @constant.numeric.integer \
\
[ \
  (decimal_floating_point_literal) \
  (hex_floating_point_literal) \
] @constant.numeric.float \
\
(character_literal) @constant.character \
\
[ \
  (string_literal) \
] @string \
\
[ \
  (true) \
  (false) \
  (null_literal) \
] @constant.builtin \
\
(line_comment) @comment \
(block_comment) @comment \
\
; Punctuation \
\
[ \
  \"::\" \
  \".\" \
  \";\" \
  \",\" \
] @punctuation.delimiter \
\
[ \
  \"@\" \
  \"...\" \
] @punctuation.special \
\
[ \
  \"(\" \
  \")\" \
  \"[\" \
  \"]\" \
  \"{\" \
  \"}\" \
] @punctuation.bracket \
\
(type_arguments \
  [ \
    \"<\" \
    \">\" \
  ] @punctuation.bracket) \
\
(type_parameters \
  [ \
    \"<\" \
    \">\" \
  ] @punctuation.bracket) \
\
; Operators \
\
[ \
  \"=\" \
  \">\" \
  \"<\" \
  \"!\" \
  \"~\" \
  \"?\" \
  \":\" \
  \"->\" \
  \"==\" \
  \">=\" \
  \"<=\" \
  \"!=\" \
  \"&&\" \
  \"||\" \
  \"++\" \
  \"--\" \
  \"+\" \
  \"-\" \
  \"*\" \
  \"/\" \
  \"&\" \
  \"|\" \
  \"^\" \
  \"%\" \
  \"<<\" \
  \">>\" \
  \">>>\" \
  \"+=\" \
  \"-=\" \
  \"*=\" \
  \"/=\" \
  \"&=\" \
  \"|=\" \
  \"^=\" \
  \"%=\" \
  \"<<=\" \
  \">>=\" \
  \">>>=\" \
] @operator \
\
; Keywords \
\
[ \
  \"abstract\" \
  \"assert\" \
  \"break\" \
  \"case\" \
  \"catch\" \
  \"class\" \
  \"continue\" \
  \"default\" \
  \"do\" \
  \"else\" \
  \"enum\" \
  \"exports\" \
  \"extends\" \
  \"final\" \
  \"finally\" \
  \"for\" \
  \"if\" \
  \"implements\" \
  \"import\" \
  \"instanceof\" \
  \"interface\" \
  \"module\" \
  \"native\" \
  \"new\" \
  \"non-sealed\" \
  \"open\" \
  \"opens\" \
  \"package\" \
  \"permits\" \
  \"private\" \
  \"protected\" \
  \"provides\" \
  \"public\" \
  \"requires\" \
  \"record\" \
  \"return\" \
  \"sealed\" \
  \"static\" \
  \"strictfp\" \
  \"switch\" \
  \"synchronized\" \
  \"throw\" \
  \"throws\" \
  \"to\" \
  \"transient\" \
  \"transitive\" \
  \"try\" \
  \"uses\" \
  \"volatile\" \
  \"while\" \
  \"with\" \
  \"yield\" \
] @keyword \
"

// QUERIES: syntax definition for test identification
// -------------------------------------------------------------------------

#define QUERY_JAVA_TEST "\
        (method_declaration\
                  (modifiers \
                        (marker_annotation\
                                name: (identifier) @string))\
                  name: (identifier) @test-name\
                )\
"

#endif // _LANG_JAVA_H_
