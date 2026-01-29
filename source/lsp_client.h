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
	bool (*changed)(bp::ipstream *stdout);

	int id;
	std::map<std::string, diagnosticparams_t> file2diagnostic;
} lspclient_t;

extern lspclient_t lspclient;

// sleep interval between consecutive lsp process startups,
// this is to prevent parallel processes from being spawned
// if the editor is restarted shoftly after exiting.
#define SLEEP_INTERVAL 500000

//
// Method for calling serialization and sending data into stdin
//
template<typename T> void send(T obj, bool skip) {
	std::ostringstream oss1, oss2;
	cereal::JSONOutputArchive oar(oss1, cereal::JSONOutputArchive::Options::NoIndent());
	obj.serialize(oar);

	std::string json(oss1.str());

	// avoid pretty printing
	// NOTE: the last character is skipped to prevent removal of closing brace
	json.erase(std::remove(json.begin(), json.end(), '\n'), json.end() - 1);
	LOG(INFO) << "send json: " << json;

	oss2 << "Content-Length: " << json.size() << "\r\n\r\n" << json;
	lspclient.stdin << oss2.str();
	lspclient.stdin.flush();
	if (lspclient.changed && !skip) {
		lspclient.changed(&lspclient.stdout);
	}
}

//
// Deserialization method
//
template<typename T> T recvjson(const std::string &json) {
	T obj;
	try {
		std::stringstream is(json);
		cereal::JSONInputArchive iar(is);
		obj.serialize(iar);
	} catch (std::exception &e) {LOG(ERROR) << e.what();}
	return obj;
}

//
// Method for deserializing message from the queue
//
template<typename T> T WaitForRequest(std::queue<std::string> *chan, long int timeout) {
	usleep(timeout); // timeout
	if (chan->empty()) return T{};
	T obj = recvjson<T>(chan->front()); // get the item from queue
	chan->pop(); // removing item from queue
	return obj;
}

inline static void receiveDiagnostics(bp::ipstream *stdout, std::string *message) {
	if (!std::getline(*stdout, *message)) {
		LOG(ERROR) << "io readline error";
	} else if (message->find("Content-Length") == std::string::npos) {
		LOG(ERROR) << "unable to parse header";
	} else {
		// get content length from header
		auto nbytes_str = message->substr(16);
		auto nbytes_conv = std::atoi(nbytes_str.c_str());

		// reading line separator
		std::getline(*stdout, *message);

		// allocate buffer, reading body content
		std::string temp(nbytes_conv, '\0');
		stdout->read(&temp[0], nbytes_conv);
		message->clear();
		message->append(temp);
	}
}

inline static bool receiveDiagnostics(bp::ipstream *stdout) {
	std::string message;
	receiveDiagnostics(stdout, &message);
	LOG(INFO) << "recv: " << message;

	if (message.find("publishDiagnostics") != std::string::npos) {
		auto dr = recvjson<diagnosticresponse_t>(message);
		lspclient.file2diagnostic[dr.params.uri] = dr.params;
		lspclient.diagnosticsChannel.push(message);
		return true;
	} else if (message.find("result") != std::string::npos) {
		lspclient.userMessages.push(message);
		return true;
	}

	return false;
}

// lsp initialization
bool StartLspClient(const std::string &cmd, std::string args...);
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
