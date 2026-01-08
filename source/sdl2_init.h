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

#ifndef _H_CLIPBOARD_H
#define _H_CLIPBOARD_H

#include <SDL2/SDL.h>
#include <glog/logging.h>

#define __init_clipboard__ { \
	if (SDL_Init(SDL_INIT_VIDEO) != 0) { \
		LOG(ERROR) << "SDL_Init failed: " << SDL_GetError(); \
		return 1;\
	} \
}

#define __set_clipboard_text(text) SDL_SetClipboardText(text)
#define __get_clipboard_text() SDL_GetClipboardText()
#define __free_clipboard(text) SDL_free(text)

#endif // _H_CLIPBOARD_H
