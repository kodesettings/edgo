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

	// TODO: serialize request to string
	auto response = send<hoverresponse_t>("");

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

	// TODO: serialize request to string
	auto response = send<completionresponse_t>("");

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

	// TODO: serialize request to string
	auto response = send<definitionresponse_t>("");

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

	// TODO: serialize request to string
	auto response = send<signaturehelpresponse_t>("");

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

	// TODO: serialize request to string
	auto response = send<referencesresponse_t>("");

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

	// TODO: serialize request to string
	auto response = send<preparerenameresponse_t>("");

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

	// TODO: serialize request to string
	auto response = send<renameresponse_t>("");

	return response;
}

codeactionresponse_t CodeAction(std::string file, int spc, int spl, int epc, int epl) {
	lspclient.id++;

	auto request = codeactionrequest_t {
		id: lspclient.id, jsonrpc: "2.0", method: "textDocument/codeAction",
		params: codeactionparams_t {
			textDocument: textdocument_t { uri:  "file://" + file },
			context: context_t { only: std::string{{"refactor"}}, triggerKind: 1 },
			range: requestrange_t {
				start: position_t { line: spl, character: spc},
				end: position_t { line: epl, character: epc},
			},
		},
	};

	// TODO: serialize request to string
	auto response = send<codeactionresponse_t>("");

	return response;
}
