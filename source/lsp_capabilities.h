#ifndef _LSP_CAPABILITIES_H
#define _LSP_CAPABILITIES_H

#include <string>
#include <vector>
#include <map>

//===================================================================
// these are the type definitions for lsp models

typedef std::string string_t;
typedef std::vector<workspaceFolder> workspacefolder_t;
typedef std::vector<std::string> string_v;
typedef std::vector<textdocumentcontentchangeevent_t> textdocumentcontentchangeevent_v;
typedef std::vector<completionitem_t> completionitem_v;
typedef std::vector<parameterinformation_t> parameterinformation_v;
typedef std::vector<signatureinformation_t> signatureinformation_v;
typedef std::vector<diagnostic_t> diagnostic_v;
typedef std::vector<definitionresult_t> definitionresult_v;
typedef std::vector<referencesrange_t> referencesrange_v;
typedef std::vector<documentchange_t> documentchange_v;
typedef std::vector<edit_t> edit_v;
typedef std::vector<codeactionresult_t> codeactionresult_v;
typedef std::vector<argument_t> argument_v;

//===================================================================
// capabilities are initially sent when connected to the LSP server
// this struct contains the initialization of such request

capabilities_t capabilities = {
	.textdocument_t = {
		.hover_t = {
			.contentFormat = string_t{"plaintext", "markdown"},
		},
		.publishdiagnostics_t = {
			.relatedInformation = false,
			.versionSupport = false,
			.codeDescriptionSupport = true,
			.dataSupport = true
		},
		.signaturehelp_t = {
			.signatureinformation_t = {
				.documentationFormat = string_t{"plaintext", "markdown"},
			}
		},
		.completion_t = {
			.completionitem_t = {
				.resolveProvider = true,
				.snippetSupport = true,
				.insertReplaceSupport = true,
				.labelDetailsSupport = true,
				.resolvesupport_t = {
					.properties = string_t{"documentation", "detail", "additionalTextEdits"},
				}
			}
		}
	}
};

#endif // _LSP_CAPABILITIES_H
