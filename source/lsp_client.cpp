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

lspclient_t lspclient;

bool StartLspClient(const std::string &cmd, std::string args...) {
    LOG(INFO) << "starting lsp cmd:" << cmd << " args:" << args;

    auto exe = bp::search_path(cmd);
	if (exe.empty()) { LOG(ERROR) << "lsp not found cmd" << cmd; return -1; }

	std::error_code ec;
	lspclient.c = bp::child(cmd, bp::std_in < lspclient.stdin, bp::std_out > lspclient.stdout, ec);
	lspclient.changed = receiveDiagnostics;

	if (ec) {
		LOG(ERROR) << "error starting lsp err" << ec.message();
		return false;
	} else {
		usleep(SLEEP_INTERVAL); // sleep for 500ms
		LOG(INFO) << "lsp started success " << cmd << " " << args;
	}

	return true;
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

	send<initializerequest_t>(initializeRequest, false);
	auto response = WaitForRequest<initializeresponse_t>(&lspclient.userMessages, 10000);

	if (response.result.serverInfo.name.empty() || response.result.serverInfo.version.empty()) {
		LOG(INFO) << "cant get initialize response from lsp server";
		lspclient.isReady = false;
		return;
	}

	auto initializedRequest = initializedrequest_t {
		jsonrpc: "2.0", method: "initialized",
		params: initializedparams_t{},
	};

	send<initializedrequest_t>(initializedRequest, true);

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

	send<didopenrequest_t>(didOpenRequest, false);
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

	send<didchangerequest_t>(didChangeRequest, false);
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

	send<didcloserequest_t>(didCloseRequest, false);
}
