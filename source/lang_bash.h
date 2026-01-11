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

#ifndef _LANG_BASH_H_
#define _LANG_BASH_H_

#define QUERY_BASH "\
[ \
 (string) \
 (raw_string) \
 ;(ansi_c_string) \
 (heredoc_body) \
] @string \
\
(variable_assignment (word) @string) \
;(command argument: \"$\" @string) \
\
\
[ \
 \"if\" \
 \"then\" \
 \"else\" \
 \"elif\" \
 \"fi\" \
 \"case\" \
 \"in\" \
 \"esac\" \
] @conditional @keyword \
\
\
[ \
 \"for\" \
 \"do\" \
 \"done\" \
 ;\"select\" \
 ;\"until\" \
 \"while\" \
] @repeat @keyword \
\
[ \
  \"declare\" \
  \"typeset\" \
  \"export\" \
  \"readonly\" \
  \"local\" \
  \"unset\" \
  \"unsetenv\" \
] @keyword \
\
\"function\" @keyword.function \
(special_variable_name) @constant \
\
; trap -l \
((word) @constant.builtin \
 (#match? @constant.builtin \"^SIG(HUP|INT|QUIT|ILL|TRAP|ABRT|BUS|FPE|KILL|USR[12]|SEGV|PIPE|ALRM|TERM|STKFLT|CHLD|CONT|STOP|TSTP|TT(IN|OU)|URG|XCPU|XFSZ|VTALRM|PROF|WINCH|IO|PWR|SYS|RTMIN([+]([1-9]|1[0-5]))?|RTMAX(-([1-9]|1[0-4]))?)$\")) \
\
((word) @boolean \
 (#any-of? @boolean \"true\" \"false\")) \
 \
(comment) @comment @spell \
\
(test_operator) @operator \
\
(function_definition \
  name: (word) @function) \
\
(command_name (word) @function.call) \
\
(variable_name) @variable @identifier \
\
(case_item \
  value: (word) @parameter) \
 \
\
((program . (comment) @preproc) \
 (#lua-match? @preproc \"^#!/\")) \
"

#endif // _LANG_BASH_H_
