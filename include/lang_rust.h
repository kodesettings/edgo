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

#ifndef _LANG_RUST_H_
#define _LANG_RUST_H_

#define QUERY_RUST "\
; ------- \
; Tree-Sitter doesn't allow overrides in regards to captures, \
; though it is possible to affect the child node of a captured \
; node. Thus, the approach here is to flip the order so that \
; overrides are unnecessary. \
; ------- \
\
; ------- \
; Types \
; ------- \
\
; --- \
; Primitives \
; --- \
\
(escape_sequence) @constant.character.escape \
(primitive_type) @type.builtin \
(boolean_literal) @constant.builtin.boolean \
(integer_literal) @constant.numeric.integer \
(float_literal) @constant.numeric.float \
(char_literal) @constant.character \
[ \
  (string_literal) \
  (raw_string_literal) \
] @string \
[ \
  (line_comment) \
  (block_comment) \
] @comment \
\
; --- \
; Extraneous \
; --- \
\
(self) @variable.builtin \
(enum_variant (identifier) @type.enum.variant) \
\
(field_initializer \
  (field_identifier) @variable.other.member) \
(shorthand_field_initializer \
  (identifier) @variable.other.member) \
(shorthand_field_identifier) @variable.other.member \
\
(lifetime \
  \"'\" @label \
  (identifier) @label) \
;(loop_label \
;  \"'\" @label \
;  (identifier) @label) \
\
\
; --- \
; Variables \
; --- \
\
(let_declaration \
  pattern: [ \
    ((identifier) @variable) \
    ((tuple_pattern \
      (identifier) @variable)) \
  ]) \
\
; It needs to be anonymous to not conflict with \"call_expression\" further below. \
( \
 value: (field_expression \
  value: (identifier)? @variable \
  field: (field_identifier) @variable.other.member)) \
\
(parameter \
	pattern: (identifier) @variable.parameter) \
(closure_parameters \
	(identifier) @variable.parameter) \
\
; ------- \
; Keywords \
; ------- \
\
(for_expression \
  \"for\" @keyword.control.repeat) \
((identifier) @keyword.control \
  (#match? @keyword.control \"^yield$\")) \
\
\"in\" @keyword.control \
\
[ \
  \"match\" \
  \"if\" \
  \"else\" \
] @keyword.control.conditional \
\
[ \
  \"while\" \
  \"loop\" \
] @keyword.control.repeat \
\
[ \
  \"break\" \
  \"continue\" \
  \"return\" \
  \"await\" \
] @keyword.control.return \
\
\"use\" @keyword.control.import \
(mod_item \"mod\" @keyword.control.import !body) \
(use_as_clause \"as\" @keyword.control.import) \
\
(type_cast_expression \"as\" @keyword.operator) \
\
[ \
  (crate) \
  (super) \
  \"as\" \
  \"pub\" \
  \"mod\" \
  \"extern\" \
\
  \"impl\" \
  \"where\" \
  \"trait\" \
  \"for\" \
\
  \"default\" \
  \"async\" \
] @keyword \
\
[ \
  \"struct\" \
  \"enum\" \
  \"union\" \
  \"type\" \
] @keyword.storage.type \
\
\"let\" @keyword.storage \
\"fn\" @keyword.function \
\"unsafe\" @keyword.special \
\"macro_rules!\" @function.macro \
\
(mutable_specifier) @keyword.storage.modifier.mut \
\
(reference_type \"&\" @keyword.storage.modifier.ref) \
(self_parameter \"&\" @keyword.storage.modifier.ref) \
\
[ \
  \"static\" \
  \"const\" \
  \"ref\" \
  \"move\" \
  \"dyn\" \
] @keyword.storage.modifier \
\
; TODO: variable.mut to highlight mutable identifiers via locals.scm \
\
; ------- \
; Guess Other Types \
; ------- \
\
((identifier) @constant \
 (#match? @constant \"^[A-Z][A-Z\\d_]*$\")) \
\
; --- \
; PascalCase identifiers in call_expressions (e.g. \"Ok()\") \
; are assumed to be enum constructors. \
; --- \
\
(call_expression \
  function: [ \
    ((identifier) @type.enum.variant \
      (#match? @type.enum.variant \"^[A-Z]\")) \
    (scoped_identifier \
      name: ((identifier) @type.enum.variant \
        (#match? @type.enum.variant \"^[A-Z]\"))) \
  ]) \
\
; --- \
; Assume that types in match arms are enums and not \
; tuple structs. Same for \"if let\" expressions. \
; --- \
\
(match_pattern \
    (scoped_identifier \
      name: (identifier) @constructor)) \
(tuple_struct_pattern \
    type: [ \
      ((identifier) @constructor) \
      (scoped_identifier \
        name: (identifier) @constructor) \
      ]) \
(struct_pattern \
  type: [ \
    ((type_identifier) @constructor) \
    (scoped_type_identifier \
      name: (type_identifier) @constructor) \
    ]) \
\
; --- \
; Other PascalCase identifiers are assumed to be structs. \
; --- \
\
((identifier) @type \
  (#match? @type \"^[A-Z]\")) \
\
; ------- \
; Functions \
; ------- \
\
(call_expression \
  function: [ \
    ((identifier) @function) \
    (scoped_identifier \
      name: (identifier) @function) \
    (field_expression \
      field: (field_identifier) @function) \
  ]) \
(generic_function \
  function: [ \
    ((identifier) @function) \
    (scoped_identifier \
      name: (identifier) @function) \
    (field_expression \
      field: (field_identifier) @function.method) \
  ]) \
\
(function_item \
  name: (identifier) @function) \
\
(function_signature_item \
  name: (identifier) @function) \
\
; --- \
; Macros \
; --- \
\
(attribute \
  (identifier) @special \
  arguments: (token_tree (identifier) @type) \
  (#eq? @special \"derive\") \
) \
\
(attribute \
  (identifier) @function.macro) \
(attribute \
  [ \
    (identifier) @function.macro \
    (scoped_identifier \
      name: (identifier) @function.macro) \
  ] \
  (token_tree (identifier) @function.macro)?) \
\
(inner_attribute_item) @attribute \
\
(macro_definition \
  name: (identifier) @function.macro) \
(macro_invocation \
  macro: [ \
    ((identifier) @function.macro) \
    (scoped_identifier \
      name: (identifier) @function.macro) \
  ] \
  \"!\" @function.macro) \
\
(metavariable) @variable.parameter \
(fragment_specifier) @type \
\
; ------- \
; Paths \
; ------- \
\
(use_declaration \
  argument: (identifier) @namespace) \
(use_wildcard \
  (identifier) @namespace) \
(extern_crate_declaration \
  name: (identifier) @namespace) \
(mod_item \
  name: (identifier) @namespace) \
(scoped_use_list \
  path: (identifier)? @namespace) \
(use_list \
  (identifier) @namespace) \
(use_as_clause \
  path: (identifier)? @namespace \
  alias: (identifier) @namespace) \
\
; --- \
; Remaining Paths \
; --- \
\
(scoped_identifier \
  path: (identifier)? @namespace \
  name: (identifier) @namespace) \
(scoped_type_identifier \
  path: (identifier) @namespace) \
\
; ------- \
; Remaining Identifiers \
; ------- \
\
\"?\" @special \
\
(type_identifier) @type \
(identifier) @variable \
(field_identifier) @variable.other.member \
"

#endif // _LANG_RUST_H_
