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
#include "lsp_model.h"
#include "lsp_client.h"

hoverresponse_t Hover(std::string file, const int line, const int character) {
	lspclient.id++;

	auto request = baserequest_t {
		id: lspclient.id, jsonrpc: "2.0", method:  "textDocument/hover",
		params: params_t {
			textDocument: textdocument_t { uri: "file://" + file },
			position: position_t { line: line, character: character },
		},
	};

	send<baserequest_t>(request);
	auto response = WaitForRequest<hoverresponse_t>(&lspclient.userMessages, 1000);

	return response;
}

completionresponse_t Completion(std::string file, int line, int character) {
	lspclient.id++;

	auto request = baserequest_t {
		id: lspclient.id, jsonrpc: "2.0", method:  "textDocument/completion",
		params: params_t {
			textDocument: textdocument_t { uri:  "file://" + file },
			position: position_t { line: line, character: character },
			context: context_t { triggerKind: 1 },
		},
	};

	send<baserequest_t>(request);
	auto response = WaitForRequest<completionresponse_t>(&lspclient.userMessages, 1000);

	return response;
}

definitionresponse_t Definition(std::string file, int line, int character) {
	lspclient.id++;

	auto request = definitionrequest_t {
		id: lspclient.id, jsonrpc: "2.0", method:  "textDocument/definition",
		params: definitionparams_t {
			textDocument: textdocument_t { uri: "file://" + file },
			position: position_t { line: line, character: character },
		},
	};

	send<definitionrequest_t>(request);
	auto response = WaitForRequest<definitionresponse_t>(&lspclient.userMessages, 1000);

	return response;
}

signaturehelpresponse_t SignatureHelp(std::string file, int line, int character) {
	lspclient.id++;

	auto request = baserequest_t {
		id: lspclient.id, jsonrpc: "2.0", method:  "textDocument/signatureHelp",
		params: params_t {
			textDocument: textdocument_t { uri: "file://" + file },
			position: position_t { line: line, character: character },
		},
	};

	send<baserequest_t>(request);
	auto response = WaitForRequest<signaturehelpresponse_t>(&lspclient.userMessages, 1000);

	return response;
}

referencesresponse_t References(std::string file, int line, int character) {
	lspclient.id++;

	auto request = baserequest_t {
		id: lspclient.id, jsonrpc: "2.0", method:  "textDocument/references",
		params: params_t {
			textDocument: textdocument_t { uri: "file://" + file },
			position: position_t { line: line, character: character },
			context: context_t { includeDeclaration: false },
		},
	};

	send<baserequest_t>(request);
	auto response = WaitForRequest<referencesresponse_t>(&lspclient.userMessages, 3000);

	return response;
}


preparerenameresponse_t PrepareRename(std::string file, int line, int character) {
	lspclient.id++;

	auto request = preparerenamerequest_t {
		id: lspclient.id, jsonrpc: "2.0", method:  "textDocument/prepareRename",
		params: params_t {
			textDocument: textdocument_t { uri:  "file://" + file },
			position: position_t { line: line, character: character },
		},
	};

	send<preparerenamerequest_t>(request);
	auto response = WaitForRequest<preparerenameresponse_t>(&lspclient.userMessages, 10000);

	return response;
}

renameresponse_t Rename(std::string file, std::string newname, int line, int character) {
	lspclient.id++;

	auto request = renamerequest_t {
		id: lspclient.id,  jsonrpc: "2.0", method:  "textDocument/rename",
		params: renameparams_t {
			newName: newname,
			position: position_t { line: line, character: character },
			textDocument: textdocument_t { uri:  "file://" + file },
		},
	};

	send<renamerequest_t>(request);
	auto response = WaitForRequest<renameresponse_t>(&lspclient.userMessages, 10000);

	return response;
}

codeactionresponse_t CodeAction(std::string file, int spc, int spl, int epc, int epl) {
	lspclient.id++;

	auto request = codeactionrequest_t {
		id: lspclient.id, jsonrpc: "2.0", method: "textDocument/codeAction",
		params: codeactionparams_t {
			textDocument: textdocument_t { uri:  "file://" + file },
			context: context_t { only: string_v{"refactor"}, triggerKind: 1 },
			range: requestrange_t {
				start: position_t { line: spl, character: spc},
				end: position_t { line: epl, character: epc},
			},
		},
	};

	send<codeactionrequest_t>(request);
	auto response = WaitForRequest<codeactionresponse_t>(&lspclient.userMessages, 10000);

	return response;
}
