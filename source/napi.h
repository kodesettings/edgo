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

#ifndef _NAPI_H_
#define _NAPI_H_

#include <node/node_api.h>

// Adding text into the buffer
napi_value add_text(napi_env env, napi_callback_info info);

// Deleting text from the buffer
napi_value del_character(napi_env env, napi_callback_info info);

// Insert text operation
napi_value replace_text(napi_env env, napi_callback_info info);

// Shift tabs with right for selected lines
napi_value shift_with_tabs(napi_env env, napi_callback_info info);

// Clipboard operations
napi_value copy(napi_env env, napi_callback_info info);
napi_value paste(napi_env env, napi_callback_info info);
napi_value cut(napi_env env, napi_callback_info info);
napi_value duplicate(napi_env env, napi_callback_info info);

// Undo/redo stack
napi_value undo(napi_env env, napi_callback_info info);
napi_value redo(napi_env env, napi_callback_info info);

// Comment line and swap lines feature
napi_value commentline(napi_env env, napi_callback_info info);
napi_value swaplinesup(napi_env env, napi_callback_info info);
napi_value swaplinesdn(napi_env env, napi_callback_info info);

// Keyboard stuff
napi_value key_down(napi_env env, napi_callback_info info);
napi_value key_up(napi_env env, napi_callback_info info);
napi_value key_left(napi_env env, napi_callback_info info);
napi_value key_right(napi_env env, napi_callback_info info);
napi_value go_top(napi_env env, napi_callback_info info);
napi_value go_bottom(napi_env env, napi_callback_info info);
napi_value scroll_up(napi_env env, napi_callback_info info);
napi_value scroll_dn(napi_env env, napi_callback_info info);
napi_value key_enter(napi_env env, napi_callback_info info);
napi_value key_delete(napi_env env, napi_callback_info info);
napi_value key_tab(napi_env env, napi_callback_info info);
napi_value key_backtab(napi_env env, napi_callback_info info);

#endif // _NAPI_H_
