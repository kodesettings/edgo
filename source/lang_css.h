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

#ifndef _LANG_CSS_H_
#define _LANG_CSS_H_

#define QUERY_CSS "\
(comment) @comment \
\
(tag_name) @tag \
(nesting_selector) @tag \
(universal_selector) @tag \
\
\"~\" @operator \
\">\" @operator \
\"+\" @operator \
\"-\" @operator \
\"*\" @operator \
\"/\" @operator \
\"=\" @operator \
\"^=\" @operator \
\"|=\" @operator \
\"~=\" @operator \
\"$=\" @operator \
\"*=\" @operator \
\
\"and\" @operator \
\"or\" @operator \
\"not\" @operator \
\"only\" @operator \
\
(attribute_selector (plain_value) @string) \
\
((property_name) @variable \
 (#match? @variable \"^--\")) \
((plain_value) @variable \
 (#match? @variable \"^--\")) \
\
(class_name) @property \
(id_name) @property \
(namespace_name) @property \
(property_name) @property \
(feature_name) @property \
\
(pseudo_element_selector (tag_name) @attribute) \
(pseudo_class_selector (class_name) @attribute) \
(attribute_name) @attribute \
\
(function_name) @function \
\
\"@media\" @keyword \
\"@import\" @keyword \
\"@charset\" @keyword \
\"@namespace\" @keyword \
\"@supports\" @keyword \
\"@keyframes\" @keyword \
(at_keyword) @keyword \
(to) @keyword \
(from) @keyword \
(important) @keyword \
\
(string_value) @string \
(color_value) @string.special \
\
(integer_value) @number \
(float_value) @number \
(unit) @type \
\
[ \
  \"#\" \
  \",\" \
  \".\" \
  \":\" \
  \"::\" \
  \";\" \
] @punctuation.delimiter \
\
[ \
  \"{\" \
  \")\" \
  \"(\" \
  \"}\" \
] @punctuation.bracket \
"

#endif // _LANG_CSS_H_
