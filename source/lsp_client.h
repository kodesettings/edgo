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

#ifndef _LSP_CLIENT_H_
#define _LSP_CLIENT_H_

#include <string>
#include <glog/logging.h>
#include <map>
#include <sstream>
#include <queue>
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
#include <cereal/archives/json.hpp>
#include <cereal/types/string.hpp>
#include <cereal/types/vector.hpp>
#include <cereal/types/map.hpp>
#include <thread>
#include <chrono>

#include "logging.h"
#include "lsp_model.h"

typedef struct {
	std::string lang;
	bool isReady; // flag for lsp initialization

	bp::opstream stdin;
	bp::ipstream stdout;

	std::queue<std::string> userMessages;
	std::queue<std::string> diagnosticsChannel;

	bp::child c; // child process
	std::mutex mtx; // mutex

	int id;
	std::map<std::string, diagnosticparams_t> file2diagnostic;
} lspclient_t;

static lspclient_t lspclient;
#define SLEEP_INTERVAL 500000

bool StartLspClient(const std::string &cmd, std::string args...);

// TODO: this is for temporary use, because the serialization library
// produces weird output that we have to alter
namespace cereal {
inline void epilogue(cereal::JSONInputArchive&, const initializeresponse_t&){}
inline void prologue(cereal::JSONInputArchive&, const initializeresponse_t&){}
}

//
// Method for calling serialization and sending data into stdin
//
template<typename T> void send(T obj) {
	std::ostringstream oss1, oss2;
	cereal::JSONOutputArchive oar(oss1, cereal::JSONOutputArchive::Options::NoIndent());
	oar(obj);

	std::string json(oss1.str());

	// avoid pretty printing
	// TODO: the last brace is also omitted for some reason, why?
	// first line removes line breaks, second line removes beginning key
	// that was added automatically by the serialization library
	json.erase(std::remove(json.begin(), json.end(), '\n'), json.end());
	json.erase(json.begin(), json.begin() + 10);

	LOG(INFO) << "send json: " << json;

	oss2 << "Content-Length: " << json.size() << "\r\n\r\n" << json;
	lspclient.stdin << oss2.str();
	lspclient.stdin.flush();
}

//
// Deserialization method
//
template<typename T> T recvjson(const std::string &json) {
	T obj;
	try {
		std::stringstream is(json);
		cereal::JSONInputArchive iar(is);
		iar(obj);
	} catch (std::exception &e) {LOG(ERROR) << e.what();}
	return obj;
}

//
// Method for deserializing message from the queue
//
template<typename T> T WaitForRequest(std::queue<std::string> *chan, long int timeout) {
	usleep(timeout); // timeout
	if (chan->empty()) return T{};
	lspclient.mtx.lock();
	T obj = recvjson<T>(chan->front()); // get the item from queue
	chan->pop(); // removing item from queue
	lspclient.mtx.unlock();
	return obj;
}

typedef struct {diagnosticresponse_t dr; std::string message; } diagnostics_t;
inline static diagnostics_t receiveDiagnostics(bp::ipstream *stdout) {
	std::string line;
	if (!std::getline(*stdout, line)) {
		LOG(ERROR) << "io readline error";
		return diagnostics_t{};
	}

	if (line.find("Content-Length") == std::string::npos) {
		LOG(ERROR) << "unable to parse header";
		return diagnostics_t{};
	} else {
		// get content length from header
		auto nbytes_str = line.substr(std::string("Content-Length: ").size());
		auto nbytes_conv = std::atoi(nbytes_str.c_str());

		// reading line separator
		std::getline(*stdout, line);

		// allocate buffer, reading body content and trimming
		// the content to the right size to omit closing bytes
		char temp[nbytes_conv];
		stdout->read(temp, nbytes_conv);
		line = std::string(temp, temp + strlen(temp) - 4);
	}

	const auto dr = recvjson<diagnosticresponse_t>(line);
	return diagnostics_t { dr: dr, message: line };
}

inline static void receiveLoop(bp::ipstream *stdout) {
repeat:
	diagnostics_t diagnostics = receiveDiagnostics(stdout);
	LOG(INFO) << "recv: " << diagnostics.message;

	if (diagnostics.message.find("publishDiagnostics") != std::string::npos) {
		lspclient.mtx.lock(); // locking this iteration
		lspclient.file2diagnostic[diagnostics.dr.params.uri] = diagnostics.dr.params;
		lspclient.diagnosticsChannel.push(diagnostics.message);
		lspclient.mtx.unlock(); // unlocking iteration
	} else if (diagnostics.message.find("result") != std::string::npos) {
		lspclient.mtx.lock(); // locking this iteration
		lspclient.userMessages.push(diagnostics.message);
		lspclient.mtx.unlock(); // unlocking iteration
	}

	usleep(SLEEP_INTERVAL); // sleep for 500ms
	goto repeat;
}

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

#endif // _LSP_CLIENT_H_
