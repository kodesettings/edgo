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
#ifndef _LSP_CLIENT_H
#define _LSP_CLIENT_H

#include <string>
#include <glog/logging.h>
#include <map>
#include <sstream>
#include <boost/version.hpp>
#if BOOST_VERSION < 108800
#include <boost/process.hpp>
namespace bp = boost::process;
#else
#include <boost/process/v1/child.hpp>
#include <boost/process/v1/io.hpp>
#include <boost/process/v1/pipe.hpp>
#include <boost/process/v1/start_dir.hpp>
#include <boost/process/v1/search_path.hpp>
namespace bp = boost::process::v1;
#endif
#include <thread>
#include <chrono>

#include "logging.h"
#include "lsp_model.h"

typedef struct {
	std::string lang;
	bool isReady; // flag for lsp initialization

	bp::opstream stdin;
	bp::ipstream stdout;

	std::stringstream userMessages;
	std::stringstream diagnosticsChannel;

	int id;
	std::map<std::string, diagnosticparams_t> file2diagnostic;
} lspclient_t;

/* storing lspclient_t in memory */
static lspclient_t lspclient;

bool StartLspClient(const std::string &cmd, std::string args...);

//
// Method for calling serialization and sending data into stdin
//
template<typename T> void send(T obj) {
	std::string message; // TODO: serialize object to json
	LOG(INFO) << "send json" << message;
	char f[message.size()];
	sprintf(f, "Content-Length: %ld\r\n\r\n%s", message.size(), message.c_str());
	lspclient.stdin << std::string(f);
}

//
// Method for deserializing message and setting a timeout for that
//
template<typename T> T WaitForRequest(std::stringstream &chan, long int timeout) {
	// TODO: deserialize message and timeout
	return T{};
}

typedef struct {diagnosticresponse_t dr; std::string message; } diagnostics_t;
diagnostics_t receiveDiagnostics(void);
void receiveLoop(void);
diagnosticparams_t GetDiagnostic(const std::string &filename);
void InitLspClient(const std::string &dir);

// lsp user actions
hoverresponse_t Hover(std::string file, const int line, const int character);
completionresponse_t Completion(std::string file, int line, int character);
definitionresponse_t Definition(std::string file, int line, int character);
signaturehelpresponse_t SignatureHelp(std::string file, int line, int character);
referencesresponse_t References(std::string file, int line, int character);
preparerenameresponse_t PrepareRename(std::string file, int line, int character);
renameresponse_t Rename(std::string file, std::string newname, int line, int character);
codeactionresponse_t CodeAction(std::string file, int spc, int spl, int epc, int epl);

// basic lsp file events
void DidOpen(const std::string &file, std::string *text);
void DidChange(const std::string &file, std::string *text, int version);
void DidClose(const std::string &file);

#endif // _LSP_CLIENT_H
