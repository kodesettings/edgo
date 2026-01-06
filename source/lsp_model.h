#ifndef _LSP_MODEL_H
#define _LSP_MODEL_H

#include <string>
#include <vector>
#include <map>

typedef std::string string_t;
typedef struct {
	string_t name;
	string_t version;
} clientinfo_t;

typedef struct {
	string_t name;
	string_t uri;
} workspacefolder_t;

typedef std::vector<std::string> string_v;
typedef struct {
	string_v properties;
} resolvesupport_t;

typedef struct {
	bool             resolveProvider;
	bool             snippetSupport;
	bool             insertReplaceSupport;
	bool             labelDetailsSupport;
	resolvesupport_t resolveSupport;
} capabilitiescompletionitem_t;

typedef struct {
	capabilitiescompletionitem_t completionItem;
} completion_t;

typedef struct {
	string_v documentationFormat;
} signatureinformation2_t;

typedef struct {
	signatureinformation2_t signatureInformation;
} signaturehelp_t;

typedef struct {
	bool relatedInformation;
	bool versionSupport;
	bool codeDescriptionSupport;
	bool dataSupport;
} publishdiagnostics_t;

typedef struct {
	string_v contentFormat;
} hover_t;

typedef struct {
	hover_t              hover;
	publishdiagnostics_t publishDiagnostics;
	signaturehelp_t      signatureHelp;
	completion_t         completion;
} capabilitiestextdocument_t;

typedef struct {
	capabilitiestextdocument_t textDocument;
} capabilities_t;

//===========================================================
// the capabilities initalizer is added in this section
// this is for setting the default values
#include "lsp_capabilities.h"

typedef struct {
	int               processId;
	string_t          rootPath;
	string_t          rootUri;
	workspacefolder_t workspaceFolders;
	clientinfo_t      clientInfo;
	string_t          trace;
	string_t          initializationOptions;
	capabilities_t    capabilities;
	string_t          workDoneToken;
} initializeparams_t;

typedef struct {
	int                id;
	string_t           jsonrpc;
	string_t           method;
	initializeparams_t params;
} initializerequest_t;

typedef struct {
	bool     includeDeclaration;
	string_v only;
	int      triggerKind;
} context_t;

typedef struct {
	string_t languageID;
	string_t text;
	string_t uri;
	int      version;
} textdocument_t;

typedef struct {
	int line;
	int character;
} position_t;

typedef struct {
	textdocument_t textDocument;
	position_t     position;
	context_t      context;
} params_t;

typedef  struct {
	textdocument_t textDocument;
} didopentextdocumentparams_t;

typedef struct {
	string_t                    jsonrpc;
	string_t                    method;
	didopentextdocumentparams_t params;
} didopenrequest_t;

typedef struct {
	string_t uri;
} textdocumentidentifier_t;

typedef struct {
	textdocumentidentifier_t textDocument;
} didclosetextdocumentparams_t;

typedef struct {
	string_t                     jsonrpc;
	string_t                     method;
	didclosetextdocumentparams_t params;
} didcloserequest_t;

typedef struct {
	int      id;
	string_t jsonrpc;
	string_t method;
} shutdownrequest_t;

typedef struct {
	string_t jsonrpc;
	string_t method;
} exitrequest_t;

typedef struct {
	int      id;
	string_t jsonrpc;
	string_t method;
	params_t params;
} baserequest_t;

typedef struct {
	string_t uri;
	int      version;
} versionedtextdocumentidentifier_t;

// NOTE: range argument is temporarily disabled
typedef struct {
//	range_t  range;
	string_t text;
} textdocumentcontentchangeevent_t;

typedef std::vector<textdocumentcontentchangeevent_t> textdocumentcontentchangeevent_v;
typedef struct {
	versionedtextdocumentidentifier_t textDocument;
	textdocumentcontentchangeevent_v contentChanges;
} didchangetextdocumentparams_t;

typedef struct {
	string_t                      jsonrpc;
	string_t                      method;
	didchangetextdocumentparams_t params;
} didchangerequest_t;

typedef struct {
	position_t start;
	position_t end;
} range_t;

typedef struct {
	range_t  range;
	range_t  replace;
	range_t  insert;
	string_t newText;
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
} completionitem_t;

typedef std::vector<completionitem_t> completionitem_v;
typedef struct {
	bool isIncomplete;
	completionitem_v items;
} completionresult_t;

typedef struct {
	string_t           jsonrpc;
	completionresult_t result;
	int                id;
} completionresponse_t;

typedef struct {
	string_t kind;
	string_t value;
} contents_t;

typedef struct {
	contents_t contents;
	range_t    range;
} hoverresult_t;

typedef struct {
	string_t      jsonrpc;
	hoverresult_t result;
	int           id;
} hoverresponse_t;

typedef struct {
	string_t label;
	string_t documentation;
} parameterinformation_t;

typedef std::vector<parameterinformation_t> parameterinformation_v;
typedef struct {
	string_t               label;
	parameterinformation_v parameters;
} signatureinformation_t;

typedef std::vector<signatureinformation_t> signatureinformation_v;
typedef struct {
	signatureinformation_v signatures;
	int                    activeSignature;
	int                    activeParameter;
} signaturehelpresult_t;

typedef struct {
	string_t              jsonrpc;
	signaturehelpresult_t result;
	int                   id;
} signaturehelpresponse_t;

typedef struct {
	std::map<string_t, string_t> langs;
} lspsettings_t;

typedef struct {
	string_t href;
} codedescription_t;

typedef struct {
	range_t           range;
	int               severity;
	string_t          code;
	codedescription_t codeDescription;
	string_t          source;
	string_t          message;
} diagnostic_t;

typedef std::vector<diagnostic_t> diagnostic_v;
typedef struct {
	string_t      uri;
	int           version;
	diagnostic_v  diagnostics;
} diagnosticparams_t;

typedef struct {
	textdocument_t textDocument;
	position_t     position;
} definitionparams_t;

typedef struct {
	string_t           jsonrpc;
	string_t           method;
	diagnosticparams_t params;
} diagnosticresponse_t;

typedef struct {
	int                id;
	string_t           jsonrpc;
	string_t           method;
	definitionparams_t params;
} definitionrequest_t;

typedef struct {
	string_t uri;
	range_t  range;
} definitionresult_t;

typedef std::vector<definitionresult_t> definitionresult_v;
typedef struct {
	string_t           jsonrpc;
	definitionresult_v result;
	int                id;
} definitionresponse_t;

typedef struct {
	string_t           jsonrpc;
	definitionresult_t result;
	int                id;
} definitionresponse2_t;

typedef  struct {
	int character;
	int line;
} character_t;

typedef struct {
	character_t start;
	character_t end;
} changerange_t;

typedef struct {
	textdocument_t textDocument;
} didsaveparams_t;

typedef struct {
	string_t        jsonrpc;
	string_t        method;
	didsaveparams_t params;
} didsaverequest_t;

typedef struct {
	position_t start;
	position_t end;
} span_t;

typedef struct {
	string_t uri;
	span_t   range;
} referencesrange_t;

typedef std::vector<referencesrange_t> referencesrange_v;
typedef struct {
	string_t          jsonrpc;
	referencesrange_v result;
	int               id;
} referencesresponse_t;

typedef struct {
	int      id;
	string_t jsonrpc;
	string_t method;
	params_t params;
} preparerenamerequest_t;

typedef struct {
	int      id;
	string_t jsonrpc;
	struct result {
		range_t range;
		string_t placeholder;
	} result_t;
} preparerenameresponse_t;

typedef struct {
	string_t       newName;
	position_t     position;
	textdocument_t textDocument;
} renameparams_t;

typedef struct {
	int            id;
	string_t       jsonrpc;
	string_t       method;
	renameparams_t params;
} renamerequest_t;

typedef struct {
	range_t  range;
	string_t newText;
} edit_t;

typedef std::vector<edit_t> edit_v;
typedef struct {
	textdocument_t textDocument;
	edit_v         edits;
} documentchange_t;

typedef std::vector<documentchange_t> documentchange_v;
typedef struct {
	documentchange_v documentChanges;
} changesresult_t;

typedef struct {
	string_t        jsonrpc;
	changesresult_t result;
	int             id;
} renameresponse_t;

typedef struct {
	position_t start;
	position_t end;
} requestrange_t;

typedef struct {
	textdocument_t textDocument;
	context_t      context;
	requestrange_t range;
} codeactionparams_t;

typedef struct {
	int                id;
	string_t           jsonrpc;
	string_t           method;
	codeactionparams_t params;
} codeactionrequest_t;

typedef struct {
	string_t fix;
	string_t uri;
	range_t range;
} argument_t;

typedef std::vector<argument_t> argument_v;
typedef struct {
	string_t   title;
	string_t   command;
	argument_v arguments;
} command_t;

typedef struct {
	string_t  title;
	string_t  kind;
	edit_t    edit;
	command_t command;
} codeactionresult_t;

typedef std::vector<codeactionresult_t> codeactionresult_v;
typedef struct {
	string_t           jsonrpc;
	codeactionresult_v result;
	int                id;
} codeactionresponse_t;

#endif // _LSP_MODEL_H
