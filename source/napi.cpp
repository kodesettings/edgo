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
#include "editor.h"

napi_value add_text(napi_env env, napi_callback_info info) {
	size_t argc = 3;
	napi_value args[3];
	napi_get_cb_info(env, info, &argc, args, NULL, NULL);

	int line = 0, pos = 0;
	char buffer[52000];
	size_t length;
	napi_get_value_int32(env, args[0], &line);
	napi_get_value_int32(env, args[1], &pos);
	napi_get_value_string_utf8(env, args[2], buffer, sizeof(buffer), &length);

	switch (length) {
	case 0: break;
	case 1: AddCharacter(buffer[0]); break; // line and pos ignored
	default: InsertString(line, pos, buffer); break;
	}

	napi_value undefined;
	napi_get_undefined(env, &undefined);
	return undefined;
}

napi_value del_character(napi_env env, napi_callback_info info) {
	size_t argc = 1;
	napi_value args[1];
	napi_get_cb_info(env, info, &argc, args, NULL, NULL);

	int line = 0, pos = 0;
	napi_get_value_int32(env, args[0], &line);
	napi_get_value_int32(env, args[1], &pos);

	DeleteCharacter(line, pos);

	napi_value undefined;
	napi_get_undefined(env, &undefined);
	return undefined;
}

napi_value replace_text(napi_env env, napi_callback_info info) {
	size_t argc = 4;
	napi_value args[4];
	napi_get_cb_info(env, info, &argc, args, NULL, NULL);

	int line = 0, from = 0, end = 0;
	char buffer[52000];
	size_t length;
	napi_get_value_int32(env, args[0], &line);
	napi_get_value_int32(env, args[1], &from);
	napi_get_value_int32(env, args[2], &end);
	napi_get_value_string_utf8(env, args[3], buffer, sizeof(buffer), &length);

	if (length > 0) { ReplaceString(line, from, end, buffer); }

	napi_value undefined;
	napi_get_undefined(env, &undefined);
	return undefined;
}

napi_value shift_with_tabs(napi_env env, napi_callback_info info) {
	size_t argc = 3;
	napi_value args[3];
	napi_get_cb_info(env, info, &argc, args, NULL, NULL);

	int line = 0, pos = 0;
	napi_get_value_int32(env, args[0], &line);
	napi_get_value_int32(env, args[1], &pos);

	bool is_array;
	napi_status status = napi_is_array(env, args[2], &is_array);
	if (status != napi_ok || !is_array) {
		napi_throw_type_error(env, NULL, "third argument must be an array");
		return NULL;
	}

	uint32_t length = 0, index = 0;
	status = napi_get_array_length(env, args[2], &length);
	if (status != napi_ok) return NULL;
loop:
	napi_value element;
	status = napi_get_element(env, args[2], index, &element);
	if (status != napi_ok) return NULL;

	int lx = 0;
	status = napi_get_value_int32(env, element, &lx);
	if (status != napi_ok) {
		napi_throw_type_error(env, NULL, "array elements must be numbers");
		return NULL;
	}

	if (index < length) goto loop;

	napi_value undefined;
	napi_get_undefined(env, &undefined);
	return undefined;
}

napi_value copy(napi_env env, napi_callback_info info) {
	size_t argc = 0;
	napi_value args[1];
	napi_get_cb_info(env, info, &argc, args, NULL, NULL);
	OnCopy(); // TODO: select text
	napi_value undefined;
	napi_get_undefined(env, &undefined);
	return undefined;
}

napi_value paste(napi_env env, napi_callback_info info) {
	size_t argc = 0;
	napi_value args[1];
	napi_get_cb_info(env, info, &argc, args, NULL, NULL);
	OnPaste(); // TODO: input text
	napi_value undefined;
	napi_get_undefined(env, &undefined);
	return undefined;
}

napi_value cut(napi_env env, napi_callback_info info) {
	size_t argc = 1;
	napi_value args[1];
	napi_get_cb_info(env, info, &argc, args, NULL, NULL);

	bool isCopySelected;
	napi_get_value_bool(env, args[0], &isCopySelected);
	Cut(isCopySelected);

	napi_value undefined;
	napi_get_undefined(env, &undefined);
	return undefined;
}

napi_value duplicate(napi_env env, napi_callback_info info) {
	size_t argc = 0;
	napi_value args[1];
	napi_get_cb_info(env, info, &argc, args, NULL, NULL);
	Duplicate(); // TODO: input text
	napi_value undefined;
	napi_get_undefined(env, &undefined);
	return undefined;
}

napi_value undo(napi_env env, napi_callback_info info) {
	size_t argc = 0;
	napi_value args[1];
	napi_get_cb_info(env, info, &argc, args, NULL, NULL);
	OnUndo();
	napi_value undefined;
	napi_get_undefined(env, &undefined);
	return undefined;
}

napi_value redo(napi_env env, napi_callback_info info) {
	size_t argc = 0;
	napi_value args[1];
	napi_get_cb_info(env, info, &argc, args, NULL, NULL);
	OnRedo();
	napi_value undefined;
	napi_get_undefined(env, &undefined);
	return undefined;
}

napi_value commentline(napi_env env, napi_callback_info info) {
	size_t argc = 0;
	napi_value args[1];
	napi_get_cb_info(env, info, &argc, args, NULL, NULL);
	OnCommentLine(); // TODO: add input
	napi_value undefined;
	napi_get_undefined(env, &undefined);
	return undefined;
}

napi_value swaplinesup(napi_env env, napi_callback_info info) {
	size_t argc = 0;
	napi_value args[1];
	napi_get_cb_info(env, info, &argc, args, NULL, NULL);
	OnSwapLinesUp(); // TODO: add input
	napi_value undefined;
	napi_get_undefined(env, &undefined);
	return undefined;
}

napi_value swaplinesdn(napi_env env, napi_callback_info info) {
	size_t argc = 0;
	napi_value args[1];
	napi_get_cb_info(env, info, &argc, args, NULL, NULL);
	OnSwapLinesDown(); // TODO: add input
	napi_value undefined;
	napi_get_undefined(env, &undefined);
	return undefined;
}

napi_value key_down(napi_env env, napi_callback_info info) {
	size_t argc = 1;
	napi_value args[1];
	napi_get_cb_info(env, info, &argc, args, NULL, NULL);

	bool isPaging;
	napi_get_value_bool(env, args[0], &isPaging);
	OnDown(isPaging);

	napi_value undefined;
	napi_get_undefined(env, &undefined);
	return undefined;
}

napi_value key_up(napi_env env, napi_callback_info info) {
	size_t argc = 1;
	napi_value args[1];
	napi_get_cb_info(env, info, &argc, args, NULL, NULL);

	bool isPaging;
	napi_get_value_bool(env, args[0], &isPaging);
	OnUp(isPaging);

	napi_value undefined;
	napi_get_undefined(env, &undefined);
	return undefined;
}

napi_value key_left(napi_env env, napi_callback_info info) {
	size_t argc = 0;
	napi_value args[1];
	napi_get_cb_info(env, info, &argc, args, NULL, NULL);
	OnLeft();
	napi_value undefined;
	napi_get_undefined(env, &undefined);
	return undefined;
}

napi_value key_right(napi_env env, napi_callback_info info) {
	size_t argc = 0;
	napi_value args[1];
	napi_get_cb_info(env, info, &argc, args, NULL, NULL);
	OnRight();
	napi_value undefined;
	napi_get_undefined(env, &undefined);
	return undefined;
}

napi_value go_top(napi_env env, napi_callback_info info) {
	size_t argc = 0;
	napi_value args[1];
	napi_get_cb_info(env, info, &argc, args, NULL, NULL);
	GoTop();
	napi_value undefined;
	napi_get_undefined(env, &undefined);
	return undefined;
}

napi_value go_bottom(napi_env env, napi_callback_info info) {
	size_t argc = 0;
	napi_value args[1];
	napi_get_cb_info(env, info, &argc, args, NULL, NULL);
	GoBottom();
	napi_value undefined;
	napi_get_undefined(env, &undefined);
	return undefined;
}

napi_value scroll_up(napi_env env, napi_callback_info info) {
	size_t argc = 0;
	napi_value args[1];
	napi_get_cb_info(env, info, &argc, args, NULL, NULL);
	OnScrollUp();
	napi_value undefined;
	napi_get_undefined(env, &undefined);
	return undefined;
}

napi_value scroll_dn(napi_env env, napi_callback_info info) {
	size_t argc = 0;
	napi_value args[1];
	napi_get_cb_info(env, info, &argc, args, NULL, NULL);
	OnScrollDown();
	napi_value undefined;
	napi_get_undefined(env, &undefined);
	return undefined;
}

napi_value key_enter(napi_env env, napi_callback_info info) {
	size_t argc = 0;
	napi_value args[1];
	napi_get_cb_info(env, info, &argc, args, NULL, NULL);
	OnEnter();
	napi_value undefined;
	napi_get_undefined(env, &undefined);
	return undefined;
}

napi_value key_delete(napi_env env, napi_callback_info info) {
	size_t argc = 0;
	napi_value args[1];
	napi_get_cb_info(env, info, &argc, args, NULL, NULL);
	OnDelete();
	napi_value undefined;
	napi_get_undefined(env, &undefined);
	return undefined;
}

napi_value key_tab(napi_env env, napi_callback_info info) {
	size_t argc = 0;
	napi_value args[1];
	napi_get_cb_info(env, info, &argc, args, NULL, NULL);
	OnTab();
	napi_value undefined;
	napi_get_undefined(env, &undefined);
	return undefined;
}

napi_value key_backtab(napi_env env, napi_callback_info info) {
	size_t argc = 0;
	napi_value args[1];
	napi_get_cb_info(env, info, &argc, args, NULL, NULL);
	OnBackTab();
	napi_value undefined;
	napi_get_undefined(env, &undefined);
	return undefined;
}
