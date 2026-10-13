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

#include <string>
#include <vector>
#include <map>
#include <cereal/cereal.hpp>

#include "lsp_types.h"
#include "utils.h"

struct clientinfo_t {
	string_t name;
	string_t version;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(name);
		ar & CEREAL_NVP(version);
	}
};

struct workspacefolder_t {
	string_t name;
	string_t uri;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(name);
		ar & CEREAL_NVP(uri);
	}
};

struct resolvesupport_t {
	string_v properties;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(properties);
	}
};

struct capabilitiescompletionitem_t {
	bool             resolveProvider;
	bool             snippetSupport;
	bool             insertReplaceSupport;
	bool             labelDetailsSupport;
	resolvesupport_t resolveSupport;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(resolveProvider);
		ar & CEREAL_NVP(snippetSupport);
		ar & CEREAL_NVP(insertReplaceSupport);
		ar & CEREAL_NVP(labelDetailsSupport);
	}
};

struct completion_t {
	capabilitiescompletionitem_t completionItem;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(completionItem);
	}
};

struct signatureinformation2_t {
	string_v documentationFormat;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(documentationFormat);
	}
};

struct signaturehelp_t {
	signatureinformation2_t signatureInformation;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(signatureInformation);
	}
};

struct publishdiagnostics_t {
	bool relatedInformation;
	bool versionSupport;
	bool codeDescriptionSupport;
	bool dataSupport;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(relatedInformation);
		ar & CEREAL_NVP(versionSupport);
		ar & CEREAL_NVP(codeDescriptionSupport);
		ar & CEREAL_NVP(dataSupport);
	}
};

struct hover_t {
	string_v contentFormat;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(contentFormat);
	}
};

struct capabilitiestextdocument_t {
	hover_t              hover;
	publishdiagnostics_t publishDiagnostics;
	signaturehelp_t      signatureHelp;
	completion_t         completion;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(hover);
		ar & CEREAL_NVP(publishDiagnostics);
		ar & CEREAL_NVP(signatureHelp);
		ar & CEREAL_NVP(completion);
	}
};

struct capabilities_t {
	capabilitiestextdocument_t textDocument;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(textDocument);
	}
};

//===========================================================
// the capabilities initalizer is added in this section
// this is for setting the default values
#include "lsp_capabilities.h"

struct initializeparams_t {
	int               processId;
	string_t          rootUri;
	string_t          rootPath;
	workspacefolder_v workspaceFolders;
	string_t          trace;
	string_t          initializationOptions;
	capabilities_t    capabilities;
	clientinfo_t      clientInfo;
	string_t          workDoneToken;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(processId);
		ar & CEREAL_NVP(rootUri);
		ar & CEREAL_NVP(rootPath);
		ar & CEREAL_NVP(workspaceFolders);
		ar & CEREAL_NVP(trace);
		ar & CEREAL_NVP(initializationOptions);
		ar & CEREAL_NVP(capabilities);
		ar & CEREAL_NVP(clientInfo);
		ar & CEREAL_NVP(workDoneToken);
	}
};

struct initializerequest_t {
	int                id;
	string_t           jsonrpc;
	string_t           method;
	initializeparams_t params;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(id);
		ar & CEREAL_NVP(jsonrpc);
		ar & CEREAL_NVP(method);
		ar & CEREAL_NVP(params);
	}
};

struct initializedparams_t {
	template <class Archive>
	void serialize(Archive& ar) {}
};

struct initializedrequest_t {
	string_t            jsonrpc;
	string_t            method;
	initializedparams_t params;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(jsonrpc);
		ar & CEREAL_NVP(method);
		ar & CEREAL_NVP(params);
	}
};

struct serverinfo_t {
	string_t name;
	string_t version;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(name);
		ar & CEREAL_NVP(version);
	}
};

struct initializeparams2_t {
//	capabilities_t capabilities;
	serverinfo_t serverInfo;

	template <class Archive>
	void serialize(Archive& ar) {
//		ar & CEREAL_NVP(capabilities);
		ar & CEREAL_NVP(serverInfo);
	}
};

struct initializeresponse_t {
	int                 id;
	string_t            jsonrpc;
	initializeparams2_t result;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(id);
		ar & CEREAL_NVP(jsonrpc);
		ar & CEREAL_NVP(result);
	}
};

struct context_t {
	bool     includeDeclaration;
	string_v only;
	int      triggerKind;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(includeDeclaration);
		ar & CEREAL_NVP(only);
		ar & CEREAL_NVP(triggerKind);
	}
};

struct textdocument_t {
	string_t languageId;
	string_t text;
	string_t uri;
	int      version;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(languageId);
		ar & CEREAL_NVP(text);
		ar & CEREAL_NVP(uri);
		ar & CEREAL_NVP(version);
	}
};

struct position_t {
	int line;
	int character;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(line);
		ar & CEREAL_NVP(character);
	}
};

struct params_t {
	textdocument_t textDocument;
	position_t     position;
	context_t      context;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(textDocument);
		ar & CEREAL_NVP(position);
		ar & CEREAL_NVP(context);
	}
};

struct didopentextdocumentparams_t {
	textdocument_t textDocument;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(textDocument);
	}
};

struct didopenrequest_t {
	string_t                    jsonrpc;
	string_t                    method;
	didopentextdocumentparams_t params;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(jsonrpc);
		ar & CEREAL_NVP(method);
		ar & CEREAL_NVP(params);
	}
};

struct textdocumentidentifier_t {
	string_t uri;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(uri);
	}
};

struct didclosetextdocumentparams_t {
	textdocumentidentifier_t textDocument;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(textDocument);
	}
};

struct didcloserequest_t {
	string_t                     jsonrpc;
	string_t                     method;
	didclosetextdocumentparams_t params;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(jsonrpc);
		ar & CEREAL_NVP(method);
		ar & CEREAL_NVP(params);
	}
};

struct shutdownrequest_t {
	int      id;
	string_t jsonrpc;
	string_t method;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(id);
		ar & CEREAL_NVP(jsonrpc);
		ar & CEREAL_NVP(method);
	}
};

struct exitrequest_t {
	string_t jsonrpc;
	string_t method;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(jsonrpc);
		ar & CEREAL_NVP(method);
	}
};

struct baserequest_t {
	int      id;
	string_t jsonrpc;
	string_t method;
	params_t params;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(id);
		ar & CEREAL_NVP(jsonrpc);
		ar & CEREAL_NVP(method);
		ar & CEREAL_NVP(params);
	}
};

struct versionedtextdocumentidentifier_t {
	string_t uri;
	int      version;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(uri);
		ar & CEREAL_NVP(version);
	}
};

// NOTE: range argument is temporarily disabled
struct textdocumentcontentchangeevent_t {
//	range_t  range;
	string_t text;

	template <class Archive>
	void serialize(Archive& ar) {
//		ar & CEREAL_NVP(range);
		ar & CEREAL_NVP(text);
	}
};

struct didchangetextdocumentparams_t {
	textdocumentcontentchangeevent_v contentChanges;
	versionedtextdocumentidentifier_t textDocument;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(contentChanges);
		ar & CEREAL_NVP(textDocument);
	}
};

struct didchangerequest_t {
	string_t                      jsonrpc;
	string_t                      method;
	didchangetextdocumentparams_t params;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(jsonrpc);
		ar & CEREAL_NVP(method);
		ar & CEREAL_NVP(params);
	}
};

struct range_t {
	position_t start;
	position_t end;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(start);
		ar & CEREAL_NVP(end);
	}
};

struct textedit_t {
	range_t  range;
	range_t  replace;
	range_t  insert;
	string_t newText;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(range);
		ar & CEREAL_NVP(replace);
		ar & CEREAL_NVP(insert);
		ar & CEREAL_NVP(newText);
	}
};

struct completionitem_t {
	string_t   label;
	int        kind;
	string_t   detail;
	string_t   sortText;
	string_t   insertText;
	string_t   filterText;
	int        insertTextFormat;
	textedit_t textEdit;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(label);
		ar & CEREAL_NVP(kind);
		ar & CEREAL_NVP(detail);
		ar & CEREAL_NVP(sortText);
		ar & CEREAL_NVP(insertText);
		ar & CEREAL_NVP(filterText);
		ar & CEREAL_NVP(insertTextFormat);
		ar & CEREAL_NVP(textEdit);
	}
};

struct completionresult_t {
	bool isIncomplete;
	completionitem_v items;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(isIncomplete);
		ar & CEREAL_NVP(items);
	}
};

struct completionresponse_t {
	string_t           jsonrpc;
	completionresult_t result;
	int                id;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(jsonrpc);
		ar & CEREAL_NVP(result);
		ar & CEREAL_NVP(id);
	}
};

struct contents_t {
	string_t kind;
	string_t value;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(kind);
		ar & CEREAL_NVP(value);
	}
};

struct hoverresult_t {
	contents_t contents;
	range_t    range;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(contents);
		ar & CEREAL_NVP(range);
	}
};

struct hoverresponse_t {
	string_t      jsonrpc;
	hoverresult_t result;
	int           id;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(jsonrpc);
		ar & CEREAL_NVP(result);
		ar & CEREAL_NVP(id);
	}
};

struct parameterinformation_t {
	string_t label;
	string_t documentation;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(label);
		ar & CEREAL_NVP(documentation);
	}
};

struct signatureinformation_t {
	string_t               label;
	parameterinformation_v parameters;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(label);
		ar & CEREAL_NVP(parameters);
	}
};

struct signaturehelpresult_t {
	signatureinformation_v signatures;
	int                    activeSignature;
	int                    activeParameter;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(signatures);
		ar & CEREAL_NVP(activeSignature);
		ar & CEREAL_NVP(activeParameter);
	}
};

struct signaturehelpresponse_t {
	string_t              jsonrpc;
	signaturehelpresult_t result;
	int                   id;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(jsonrpc);
		ar & CEREAL_NVP(result);
		ar & CEREAL_NVP(id);
	}
};

struct lspsettings_t {
	std::map<string_t, string_t> langs;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(langs);
	}
};

struct codedescription_t {
	string_t href;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(href);
	}
};

struct diagnostic_t {
	range_t           range;
	int               severity;
	string_t          code;
	codedescription_t codeDescription;
	string_t          source;
	string_t          message;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(range);
		ar & CEREAL_NVP(severity);
		ar & CEREAL_NVP(code);
		ar & CEREAL_NVP(codeDescription);
		ar & CEREAL_NVP(source);
		ar & CEREAL_NVP(message);
	}
};

struct diagnosticparams_t {
	string_t      uri;
	int           version;
	diagnostic_v  diagnostics;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(uri);
		ar & CEREAL_NVP(version);
		ar & CEREAL_NVP(diagnostics);
	}
};

struct definitionparams_t {
	textdocument_t textDocument;
	position_t     position;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(textDocument);
		ar & CEREAL_NVP(position);
	}
};

struct diagnosticresponse_t {
	string_t           jsonrpc;
	string_t           method;
	diagnosticparams_t params;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(jsonrpc);
		ar & CEREAL_NVP(method);
		ar & CEREAL_NVP(params);
	}
};

struct definitionrequest_t {
	int                id;
	string_t           jsonrpc;
	string_t           method;
	definitionparams_t params;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(id);
		ar & CEREAL_NVP(jsonrpc);
		ar & CEREAL_NVP(method);
		ar & CEREAL_NVP(params);
	}
};

struct definitionresult_t {
	string_t uri;
	range_t  range;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(uri);
		ar & CEREAL_NVP(range);
	}
};

struct definitionresponse_t {
	string_t           jsonrpc;
	definitionresult_v result;
	int                id;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(jsonrpc);
		ar & CEREAL_NVP(result);
		ar & CEREAL_NVP(id);
	}
};

struct definitionresponse2_t {
	string_t           jsonrpc;
	definitionresult_t result;
	int                id;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(jsonrpc);
		ar & CEREAL_NVP(result);
		ar & CEREAL_NVP(id);
	}
};

struct character_t {
	int character;
	int line;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(character);
		ar & CEREAL_NVP(line);
	}
};

struct changerange_t {
	character_t start;
	character_t end;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(start);
		ar & CEREAL_NVP(end);
	}
};

struct didsaveparams_t {
	textdocument_t textDocument;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(textDocument);
	}
};

struct didsaverequest_t {
	string_t        jsonrpc;
	string_t        method;
	didsaveparams_t params;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(jsonrpc);
		ar & CEREAL_NVP(method);
		ar & CEREAL_NVP(params);
	}
};

struct referencesrange_t {
	string_t uri;
	range_t  range;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(uri);
		ar & CEREAL_NVP(range);
	}
};

struct referencesresponse_t {
	string_t          jsonrpc;
	referencesrange_v result;
	int               id;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(jsonrpc);
		ar & CEREAL_NVP(result);
		ar & CEREAL_NVP(id);
	}
};

struct preparerenamerequest_t {
	int      id;
	string_t jsonrpc;
	string_t method;
	params_t params;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(id);
		ar & CEREAL_NVP(jsonrpc);
		ar & CEREAL_NVP(method);
		ar & CEREAL_NVP(params);
	}
};

struct preparerenameresponse_t {
	int      id;
	string_t jsonrpc;
	struct {
		range_t range;
		string_t placeholder;

		template <class Archive>
		void serialize(Archive& ar) {
			ar & CEREAL_NVP(range);
			ar & CEREAL_NVP(placeholder);
		}
	} result;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(id);
		ar & CEREAL_NVP(jsonrpc);
		ar & CEREAL_NVP(result);
	}
};

struct renameparams_t {
	string_t       newName;
	position_t     position;
	textdocument_t textDocument;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(newName);
		ar & CEREAL_NVP(position);
		ar & CEREAL_NVP(textDocument);
	}
};

struct renamerequest_t {
	int            id;
	string_t       jsonrpc;
	string_t       method;
	renameparams_t params;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(id);
		ar & CEREAL_NVP(jsonrpc);
		ar & CEREAL_NVP(method);
		ar & CEREAL_NVP(params);
	}
};

struct edit_t {
	range_t  range;
	string_t newText;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(range);
		ar & CEREAL_NVP(newText);
	}
};

struct documentchange_t {
	textdocument_t textDocument;
	edit_v         edits;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(textDocument);
		ar & CEREAL_NVP(edits);
	}
};

struct changesresult_t {
	documentchange_v documentChanges;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(documentChanges);
	}
};

struct renameresponse_t {
	string_t        jsonrpc;
	changesresult_t result;
	int             id;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(jsonrpc);
		ar & CEREAL_NVP(result);
		ar & CEREAL_NVP(id);
	}
};

struct requestrange_t {
	position_t start;
	position_t end;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(start);
		ar & CEREAL_NVP(end);
	}
};

struct codeactionparams_t {
	textdocument_t textDocument;
	context_t      context;
	requestrange_t range;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(textDocument);
		ar & CEREAL_NVP(context);
		ar & CEREAL_NVP(range);
	}
};

struct codeactionrequest_t {
	int                id;
	string_t           jsonrpc;
	string_t           method;
	codeactionparams_t params;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(id);
		ar & CEREAL_NVP(jsonrpc);
		ar & CEREAL_NVP(method);
		ar & CEREAL_NVP(params);
	}
};

struct argument_t {
	string_t fix;
	string_t uri;
	range_t range;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(fix);
		ar & CEREAL_NVP(uri);
		ar & CEREAL_NVP(range);
	}
};

struct command_t {
	string_t   title;
	string_t   command;
	argument_v arguments;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(title);
		ar & CEREAL_NVP(command);
		ar & CEREAL_NVP(arguments);
	}
};

struct codeactionresult_t {
	string_t  title;
	string_t  kind;
	edit_t    edit;
	command_t command;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(title);
		ar & CEREAL_NVP(kind);
		ar & CEREAL_NVP(edit);
		ar & CEREAL_NVP(command);
	}
};

struct codeactionresponse_t {
	string_t           jsonrpc;
	codeactionresult_v result;
	int                id;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(jsonrpc);
		ar & CEREAL_NVP(result);
		ar & CEREAL_NVP(id);
	}
};
