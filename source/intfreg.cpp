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

#include "napi.h"

napi_value register_napi(napi_env env, napi_value exports) {
	napi_value fn;
	napi_create_function(env, "add_text", NAPI_AUTO_LENGTH, add_text, NULL, &fn);
	napi_create_function(env, "del_character", NAPI_AUTO_LENGTH, del_character, NULL, &fn);
	napi_create_function(env, "shift_with_tabs", NAPI_AUTO_LENGTH, shift_with_tabs, NULL, &fn);
	napi_create_function(env, "cut", NAPI_AUTO_LENGTH, cut, NULL, &fn);
	napi_create_function(env, "copy", NAPI_AUTO_LENGTH, copy, NULL, &fn);
	napi_create_function(env, "paste", NAPI_AUTO_LENGTH, paste, NULL, &fn);
	napi_create_function(env, "duplicate", NAPI_AUTO_LENGTH, duplicate, NULL, &fn);
	napi_create_function(env, "undo", NAPI_AUTO_LENGTH, undo, NULL, &fn);
	napi_create_function(env, "redo", NAPI_AUTO_LENGTH, redo, NULL, &fn);
	napi_create_function(env, "commentline", NAPI_AUTO_LENGTH, commentline, NULL, &fn);
	napi_create_function(env, "swaplinesup", NAPI_AUTO_LENGTH, swaplinesup, NULL, &fn);
	napi_create_function(env, "swaplinesdn", NAPI_AUTO_LENGTH, swaplinesdn, NULL, &fn);

	napi_create_function(env, "key_down", NAPI_AUTO_LENGTH, key_down, NULL, &fn);
	napi_create_function(env, "key_up", NAPI_AUTO_LENGTH, key_up, NULL, &fn);
	napi_create_function(env, "key_left", NAPI_AUTO_LENGTH, key_left, NULL, &fn);
	napi_create_function(env, "key_right", NAPI_AUTO_LENGTH, key_right, NULL, &fn);
	napi_create_function(env, "go_top", NAPI_AUTO_LENGTH, go_top, NULL, &fn);
	napi_create_function(env, "go_bottom", NAPI_AUTO_LENGTH, go_bottom, NULL, &fn);
	napi_create_function(env, "scroll_up", NAPI_AUTO_LENGTH, scroll_up, NULL, &fn);
	napi_create_function(env, "scroll_dn", NAPI_AUTO_LENGTH, scroll_dn, NULL, &fn);
	napi_create_function(env, "key_enter", NAPI_AUTO_LENGTH, key_enter, NULL, &fn);
	napi_create_function(env, "key_delete", NAPI_AUTO_LENGTH, key_delete, NULL, &fn);
	napi_create_function(env, "key_tab", NAPI_AUTO_LENGTH, key_tab, NULL, &fn);
	napi_create_function(env, "key_backtab", NAPI_AUTO_LENGTH, key_backtab, NULL, &fn);

	napi_create_function(env, "screen_report", NAPI_AUTO_LENGTH, screen_report, NULL, &fn);

	napi_set_named_property(env, exports, "add_text", fn);
	napi_set_named_property(env, exports, "del_character", fn);
	napi_set_named_property(env, exports, "shift_with_tabs", fn);
	napi_set_named_property(env, exports, "cut", fn);
	napi_set_named_property(env, exports, "copy", fn);
	napi_set_named_property(env, exports, "paste", fn);
	napi_set_named_property(env, exports, "duplicate", fn);
	napi_set_named_property(env, exports, "undo", fn);
	napi_set_named_property(env, exports, "redo", fn);
	napi_set_named_property(env, exports, "commentline", fn);
	napi_set_named_property(env, exports, "swaplinesup", fn);
	napi_set_named_property(env, exports, "swaplinesdn", fn);

	napi_set_named_property(env, exports, "key_down", fn);
	napi_set_named_property(env, exports, "key_up", fn);
	napi_set_named_property(env, exports, "key_left", fn);
	napi_set_named_property(env, exports, "key_right", fn);
	napi_set_named_property(env, exports, "go_top", fn);
	napi_set_named_property(env, exports, "go_bottom", fn);
	napi_set_named_property(env, exports, "scroll_up", fn);
	napi_set_named_property(env, exports, "scroll_dn", fn);
	napi_set_named_property(env, exports, "key_enter", fn);
	napi_set_named_property(env, exports, "key_delete", fn);
	napi_set_named_property(env, exports, "key_tab", fn);
	napi_set_named_property(env, exports, "key_backtab", fn);

	napi_set_named_property(env, exports, "screen_report", fn);

	return exports;
}

NAPI_MODULE(NODE_GYP_MODULE_NAME, register_napi)
