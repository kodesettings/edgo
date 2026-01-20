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
#include "lsp_client.h"
#include "logging.h"

bool StartLspClient(const std::string &cmd, std::string args...) {
    LOG(INFO) << "starting lsp cmd:" << cmd << " args:" << args;

    auto exe = bp::search_path(cmd);
	if (exe.empty()) { LOG(ERROR) << "lsp not found cmd" << cmd; return -1; }

	std::error_code ec;
	lspclient.c = bp::child(cmd, bp::std_in < lspclient.stdin, bp::std_out > lspclient.stdout, ec);

	if (ec) {
		LOG(ERROR) << "error starting lsp err" << ec.message();
		return false;
	} else {
		std::this_thread::sleep_for(std::chrono::milliseconds(500));
		LOG(INFO) << "lsp started success " << cmd << " " << args;
	}

	lspclient.isReady = true;
	std::thread(&receiveLoop, &lspclient.stdout).detach();

	return true;
}

diagnostics_t receiveDiagnostics(bp::ipstream *stdout) {
	std::string line;
	if (!std::getline(lspclient.stdout, line)) {
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
		std::getline(lspclient.stdout, line);

		// allocate buffer, reading body content and trimming
		// the content to the right size to omit closing bytes
		char temp[nbytes_conv];
		stdout->read(temp, nbytes_conv);
		line = std::string(temp, temp + strlen(temp) - 4);
	}

	const auto dr = recvjson<diagnosticresponse_t>(line);
	return diagnostics_t { dr: dr, message: line };
}

void receiveLoop(bp::ipstream *stdout) {
repeat:
	lspclient.mtx.lock(); // locking this iteration
	diagnostics_t diagnostics = receiveDiagnostics(stdout);
	LOG(INFO) << "recv: " << diagnostics.message;

	if (diagnostics.message.find("publishDiagnostics") != std::string::npos) {
		lspclient.file2diagnostic[diagnostics.dr.params.uri] = diagnostics.dr.params;
		lspclient.diagnosticsChannel.push(diagnostics.message);
	} else if (diagnostics.message.find("result") != std::string::npos) {
		lspclient.userMessages.push(diagnostics.message);
	}

	lspclient.mtx.unlock(); // unlocking iteration
	goto repeat;
}

diagnosticparams_t GetDiagnostic(const std::string &filename) {
	return lspclient.file2diagnostic[filename];
}

void InitLspClient(const std::string &dir) {
	lspclient.id = 0;

	auto initializeRequest = initializerequest_t {
		id: lspclient.id, jsonrpc: "2.0", method: "initialize",
		params: initializeparams_t {
			rootUri: "file://" + dir, rootPath: dir,
			workspaceFolders: workspacefolder_v{{name: "edgo", uri: "file://" + dir}},
			capabilities: capabilities,
			clientInfo: clientinfo_t{name: "edgo", version: "1.0.0"},
		},
	};

	send<initializerequest_t>(initializeRequest);
	auto response = WaitForRequest<initializeresponse_t>(lspclient.userMessages, 10000);

	if (response.result.serverinfo.name.empty() || response.result.serverinfo.version.empty()) {
		LOG(INFO) << "cant get initialize response from lsp server";
		lspclient.isReady = false;
		return;
	}

	auto initializedRequest = initializedrequest_t {
		jsonrpc: "2.0", method: "initialized",
		params: initializedparams_t{},
	};

	send<initializedrequest_t>(initializedRequest);

	LOG(INFO) << "lsp initialized";
	lspclient.isReady = true;
}

void DidOpen(const std::string &file, std::string *text) {
	auto didOpenRequest = didopenrequest_t {
		jsonrpc: "2.0", method: "textDocument/didOpen",
		params: didopentextdocumentparams_t {
			textDocument: textdocument_t {
				languageId: lspclient.lang,
				text: *text,
				uri: "file://" + file,
				version: 1,
			},
		},
	};

	send<didopenrequest_t>(didOpenRequest);
}

void DidChange(const std::string &file, std::string *text, int version) {
	auto didChangeRequest = didchangerequest_t {
		jsonrpc: "2.0", method: "textDocument/didChange",
		params: didchangetextdocumentparams_t {
			contentChanges: textdocumentcontentchangeevent_v {{
//				range: range_t {
//					start: position_t{line: spl, character: spc},
//					end: position_t{line: epl, character: epc},
//				},
				text: *text,
			}},
			textDocument: versionedtextdocumentidentifier_t {
				uri: "file://" + file,
				version: version,
			},
		},
	};

	send<didchangerequest_t>(didChangeRequest);
}

void DidClose(const std::string &file) {
	auto didCloseRequest = didcloserequest_t {
		jsonrpc: "2.0", method: "textDocument/didClose",
		params: didclosetextdocumentparams_t {
			textDocument: textdocumentidentifier_t {
				uri: "file://" + file,
			},
		},
	};

	send<didcloserequest_t>(didCloseRequest);
}
