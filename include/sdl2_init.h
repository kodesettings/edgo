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

#ifndef _SDL2_INIT_H_
#define _SDL2_INIT_H_

#include <SDL2/SDL.h>
#include <glog/logging.h>

static std::string clpbd;
static bool x11enabled = false;

#define EDGO_SDL2_INIT { \
	if (SDL_Init(SDL_INIT_VIDEO) != 0) { \
		LOG(ERROR) << "SDL_Init failed: " << SDL_GetError(); \
		LOG(INFO) << "fallback to internal clipboard"; \
	} else { \
		x11enabled = true; \
	} \
}

#define __set_clipboard_text(text) \
	if (x11enabled) SDL_SetClipboardText(text); else clpbd = text;
#define __get_clipboard_text() \
	x11enabled ? SDL_GetClipboardText() : clpbd.data();
#define __free_clipboard(text) \
	if (x11enabled) SDL_free(text); else clpbd.clear();

#endif // _SDL2_INIT_H_
