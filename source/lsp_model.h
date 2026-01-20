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

#ifndef _LSP_MODEL_H_
#define _LSP_MODEL_H_

#include <string>
#include <vector>
#include <map>
#include <cereal/cereal.hpp>
#include <cereal/types/optional.hpp>

typedef std::string string_t;
typedef struct {
	std::optional<string_t> name;
	std::optional<string_t> version;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(name);
		ar & CEREAL_NVP(version);
	}
} clientinfo_t;

typedef struct {
	string_t name;
	string_t uri;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(name);
		ar & CEREAL_NVP(uri);
	}
} workspacefolder_t;

typedef std::vector<std::string> string_v;
typedef struct {
	string_v properties;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(properties);
	}
} resolvesupport_t;

typedef struct {
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
} capabilitiescompletionitem_t;

typedef struct {
	capabilitiescompletionitem_t completionItem;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(completionItem);
	}
} completion_t;

typedef struct {
	string_v documentationFormat;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(documentationFormat);
	}
} signatureinformation2_t;

typedef struct {
	signatureinformation2_t signatureInformation;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(signatureInformation);
	}
} signaturehelp_t;

typedef struct {
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
} publishdiagnostics_t;

typedef struct {
	string_v contentFormat;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(contentFormat);
	}
} hover_t;

typedef struct {
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
} capabilitiestextdocument_t;

typedef struct {
	capabilitiestextdocument_t textDocument;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(textDocument);
	}
} capabilities_t;

//===========================================================
// the capabilities initalizer is added in this section
// this is for setting the default values
#include "lsp_capabilities.h"

typedef std::vector<workspacefolder_t> workspacefolder_v;
typedef struct {
	std::optional<int>               processId;
	std::optional<string_t>          rootUri;
	std::optional<string_t>          rootPath;
	std::optional<workspacefolder_v> workspaceFolders;
	std::optional<string_t>          trace;
	std::optional<string_t>          initializationOptions;
	std::optional<capabilities_t>    capabilities;
	std::optional<clientinfo_t>      clientInfo;
	std::optional<string_t>          workDoneToken;

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
} initializeparams_t;

typedef struct {
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
} initializerequest_t;

typedef struct {
	template <class Archive>
	void serialize(Archive& ar) {}
} initializedparams_t;

typedef struct {
	string_t                           jsonrpc;
	string_t                           method;
	std::optional<initializedparams_t> params;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(jsonrpc);
		ar & CEREAL_NVP(method);
		ar & CEREAL_NVP(params);
	}
} initializedrequest_t;

typedef struct {
	string_t name;
	string_t version;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(name);
		ar & CEREAL_NVP(version);
	}
} serverinfo_t;

typedef struct {
	capabilities_t capabilities;
	serverinfo_t serverinfo;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(capabilities);
		ar & CEREAL_NVP(serverinfo);
	}
} initializeparams2_t;

typedef struct {
	int                 id;
	string_t            jsonrpc;
	initializeparams2_t result;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(id);
		ar & CEREAL_NVP(jsonrpc);
		ar & CEREAL_NVP(result);
	}
} initializeresponse_t;

typedef struct {
	std::optional<bool>     includeDeclaration;
	std::optional<string_v> only;
	std::optional<int>      triggerKind;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(includeDeclaration);
		ar & CEREAL_NVP(only);
		ar & CEREAL_NVP(triggerKind);
	}
} context_t;

typedef struct {
	std::optional<string_t> languageId;
	std::optional<string_t> text;
	std::optional<string_t> uri;
	std::optional<int>      version;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(languageId);
		ar & CEREAL_NVP(text);
		ar & CEREAL_NVP(uri);
		ar & CEREAL_NVP(version);
	}
} textdocument_t;

typedef struct {
	int line;
	int character;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(line);
		ar & CEREAL_NVP(character);
	}
} position_t;

typedef struct {
	textdocument_t            textDocument;
	std::optional<position_t> position;
	std::optional<context_t>  context;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(textDocument);
		ar & CEREAL_NVP(position);
		ar & CEREAL_NVP(context);
	}
} params_t;

typedef  struct {
	textdocument_t textDocument;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(textDocument);
	}
} didopentextdocumentparams_t;

typedef struct {
	string_t                    jsonrpc;
	string_t                    method;
	didopentextdocumentparams_t params;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(jsonrpc);
		ar & CEREAL_NVP(method);
		ar & CEREAL_NVP(params);
	}
} didopenrequest_t;

typedef struct {
	string_t uri;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(uri);
	}
} textdocumentidentifier_t;

typedef struct {
	textdocumentidentifier_t textDocument;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(textDocument);
	}
} didclosetextdocumentparams_t;

typedef struct {
	string_t                     jsonrpc;
	string_t                     method;
	didclosetextdocumentparams_t params;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(jsonrpc);
		ar & CEREAL_NVP(method);
		ar & CEREAL_NVP(params);
	}
} didcloserequest_t;

typedef struct {
	std::optional<int> id;
	string_t           jsonrpc;
	string_t           method;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(id);
		ar & CEREAL_NVP(jsonrpc);
		ar & CEREAL_NVP(method);
	}
} shutdownrequest_t;

typedef struct {
	string_t jsonrpc;
	string_t method;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(jsonrpc);
		ar & CEREAL_NVP(method);
	}
} exitrequest_t;

typedef struct {
	std::optional<int> id;
	string_t           jsonrpc;
	string_t           method;
	params_t           params;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(id);
		ar & CEREAL_NVP(jsonrpc);
		ar & CEREAL_NVP(method);
		ar & CEREAL_NVP(params);
	}
} baserequest_t;

typedef struct {
	string_t uri;
	int      version;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(uri);
		ar & CEREAL_NVP(version);
	}
} versionedtextdocumentidentifier_t;

// NOTE: range argument is temporarily disabled
typedef struct {
//	range_t  range;
	string_t text;

	template <class Archive>
	void serialize(Archive& ar) {
//		ar & CEREAL_NVP(range);
		ar & CEREAL_NVP(text);
	}
} textdocumentcontentchangeevent_t;

typedef std::vector<textdocumentcontentchangeevent_t> textdocumentcontentchangeevent_v;
typedef struct {
	textdocumentcontentchangeevent_v contentChanges;
	versionedtextdocumentidentifier_t textDocument;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(contentChanges);
		ar & CEREAL_NVP(textDocument);
	}
} didchangetextdocumentparams_t;

typedef struct {
	string_t                      jsonrpc;
	string_t                      method;
	didchangetextdocumentparams_t params;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(jsonrpc);
		ar & CEREAL_NVP(method);
		ar & CEREAL_NVP(params);
	}
} didchangerequest_t;

typedef struct {
	position_t start;
	position_t end;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(start);
		ar & CEREAL_NVP(end);
	}
} range_t;

typedef struct {
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
} textedit_t;

typedef struct {
	string_t   label;
	int        kind;
	string_t   detail;
	bool       preselect;
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
		ar & CEREAL_NVP(preselect);
		ar & CEREAL_NVP(sortText);
		ar & CEREAL_NVP(insertText);
		ar & CEREAL_NVP(filterText);
		ar & CEREAL_NVP(insertTextFormat);
		ar & CEREAL_NVP(textEdit);
	}
} completionitem_t;

typedef std::vector<completionitem_t> completionitem_v;
typedef struct {
	bool isIncomplete;
	completionitem_v items;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(isIncomplete);
		ar & CEREAL_NVP(items);
	}
} completionresult_t;

typedef struct {
	string_t           jsonrpc;
	completionresult_t result;
	int                id;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(jsonrpc);
		ar & CEREAL_NVP(result);
		ar & CEREAL_NVP(id);
	}
} completionresponse_t;

typedef struct {
	string_t kind;
	string_t value;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(kind);
		ar & CEREAL_NVP(value);
	}
} contents_t;

typedef struct {
	contents_t contents;
	range_t    range;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(contents);
		ar & CEREAL_NVP(range);
	}
} hoverresult_t;

typedef struct {
	string_t      jsonrpc;
	hoverresult_t result;
	int           id;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(jsonrpc);
		ar & CEREAL_NVP(result);
		ar & CEREAL_NVP(id);
	}
} hoverresponse_t;

typedef struct {
	string_t label;
	string_t documentation;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(label);
		ar & CEREAL_NVP(documentation);
	}
} parameterinformation_t;

typedef std::vector<parameterinformation_t> parameterinformation_v;
typedef struct {
	string_t               label;
	parameterinformation_v parameters;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(label);
		ar & CEREAL_NVP(parameters);
	}
} signatureinformation_t;

typedef std::vector<signatureinformation_t> signatureinformation_v;
typedef struct {
	signatureinformation_v signatures;
	int                    activeSignature;
	int                    activeParameter;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(signatures);
		ar & CEREAL_NVP(activeSignature);
		ar & CEREAL_NVP(activeParameter);
	}
} signaturehelpresult_t;

typedef struct {
	string_t              jsonrpc;
	signaturehelpresult_t result;
	int                   id;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(jsonrpc);
		ar & CEREAL_NVP(result);
		ar & CEREAL_NVP(id);
	}
} signaturehelpresponse_t;

typedef struct {
	std::map<string_t, string_t> langs;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(langs);
	}
} lspsettings_t;

typedef struct {
	string_t href;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(href);
	}
} codedescription_t;

typedef struct {
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
} diagnostic_t;

typedef std::vector<diagnostic_t> diagnostic_v;
typedef struct {
	string_t      uri;
	int           version;
	diagnostic_v  diagnostics;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(uri);
		ar & CEREAL_NVP(version);
		ar & CEREAL_NVP(diagnostics);
	}
} diagnosticparams_t;

typedef struct {
	textdocument_t            textDocument;
	std::optional<position_t> position;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(textDocument);
		ar & CEREAL_NVP(position);
	}
} definitionparams_t;

typedef struct {
	string_t           jsonrpc;
	string_t           method;
	diagnosticparams_t params;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(jsonrpc);
		ar & CEREAL_NVP(method);
		ar & CEREAL_NVP(params);
	}
} diagnosticresponse_t;

typedef struct {
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
} definitionrequest_t;

typedef struct {
	string_t uri;
	range_t  range;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(uri);
		ar & CEREAL_NVP(range);
	}
} definitionresult_t;

typedef std::vector<definitionresult_t> definitionresult_v;
typedef struct {
	string_t           jsonrpc;
	definitionresult_v result;
	int                id;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(jsonrpc);
		ar & CEREAL_NVP(result);
		ar & CEREAL_NVP(id);
	}
} definitionresponse_t;

typedef struct {
	string_t           jsonrpc;
	definitionresult_t result;
	int                id;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(jsonrpc);
		ar & CEREAL_NVP(result);
		ar & CEREAL_NVP(id);
	}
} definitionresponse2_t;

typedef  struct {
	int character;
	int line;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(character);
		ar & CEREAL_NVP(line);
	}
} character_t;

typedef struct {
	character_t start;
	character_t end;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(start);
		ar & CEREAL_NVP(end);
	}
} changerange_t;

typedef struct {
	textdocument_t textDocument;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(textDocument);
	}
} didsaveparams_t;

typedef struct {
	string_t        jsonrpc;
	string_t        method;
	didsaveparams_t params;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(jsonrpc);
		ar & CEREAL_NVP(method);
		ar & CEREAL_NVP(params);
	}
} didsaverequest_t;

typedef struct {
	position_t start;
	position_t end;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(start);
		ar & CEREAL_NVP(end);
	}
} span_t;

typedef struct {
	string_t uri;
	span_t   range;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(uri);
		ar & CEREAL_NVP(range);
	}
} referencesrange_t;

typedef std::vector<referencesrange_t> referencesrange_v;
typedef struct {
	string_t          jsonrpc;
	referencesrange_v result;
	int               id;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(jsonrpc);
		ar & CEREAL_NVP(result);
		ar & CEREAL_NVP(id);
	}
} referencesresponse_t;

typedef struct {
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
} preparerenamerequest_t;

typedef struct {
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
} preparerenameresponse_t;

typedef struct {
	string_t       newName;
	position_t     position;
	textdocument_t textDocument;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(newName);
		ar & CEREAL_NVP(position);
		ar & CEREAL_NVP(textDocument);
	}
} renameparams_t;

typedef struct {
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
} renamerequest_t;

typedef struct {
	range_t  range;
	string_t newText;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(range);
		ar & CEREAL_NVP(newText);
	}
} edit_t;

typedef std::vector<edit_t> edit_v;
typedef struct {
	textdocument_t textDocument;
	edit_v         edits;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(textDocument);
		ar & CEREAL_NVP(edits);
	}
} documentchange_t;

typedef std::vector<documentchange_t> documentchange_v;
typedef struct {
	documentchange_v documentChanges;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(documentChanges);
	}
} changesresult_t;

typedef struct {
	string_t        jsonrpc;
	changesresult_t result;
	int             id;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(jsonrpc);
		ar & CEREAL_NVP(result);
		ar & CEREAL_NVP(id);
	}
} renameresponse_t;

typedef struct {
	position_t start;
	position_t end;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(start);
		ar & CEREAL_NVP(end);
	}
} requestrange_t;

typedef struct {
	textdocument_t textDocument;
	context_t      context;
	requestrange_t range;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(textDocument);
		ar & CEREAL_NVP(context);
		ar & CEREAL_NVP(range);
	}
} codeactionparams_t;

typedef struct {
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
} codeactionrequest_t;

typedef struct {
	string_t fix;
	string_t uri;
	range_t range;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(fix);
		ar & CEREAL_NVP(uri);
		ar & CEREAL_NVP(range);
	}
} argument_t;

typedef std::vector<argument_t> argument_v;
typedef struct {
	string_t   title;
	string_t   command;
	argument_v arguments;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(title);
		ar & CEREAL_NVP(command);
		ar & CEREAL_NVP(arguments);
	}
} command_t;

typedef struct {
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
} codeactionresult_t;

typedef std::vector<codeactionresult_t> codeactionresult_v;
typedef struct {
	string_t           jsonrpc;
	codeactionresult_v result;
	int                id;

	template <class Archive>
	void serialize(Archive& ar) {
		ar & CEREAL_NVP(jsonrpc);
		ar & CEREAL_NVP(result);
		ar & CEREAL_NVP(id);
	}
} codeactionresponse_t;

#endif // _LSP_MODEL_H_
