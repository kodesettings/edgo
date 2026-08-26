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

#pragma once
#include <glog/logging.h>

// Initialization macros for google's logging system.
// By default we have to log into files on disk.
#define EDGO_LOGGING_INIT google::InitGoogleLogging("edgo");
#define EDGO_LOGGING_NO_STDERR FLAGS_logtostderr = 0;
#define EDGO_LOGGING_SUPPRESS FLAGS_minloglevel = 2; // FATAL ONLY
