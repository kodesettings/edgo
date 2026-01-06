#ifndef _LSP_CAPABILITIES_H
#define _LSP_CAPABILITIES_H

//===================================================================
// capabilities are initially sent when connected to the LSP server
// this struct contains the initialization of such request

const auto capabilities = capabilities_t {
	capabilitiestextdocument_t {
		hover_t {
			.contentFormat = string_v{"plaintext", "markdown"},
		},
		publishdiagnostics_t {
			.relatedInformation = false,
			.versionSupport = false,
			.codeDescriptionSupport = true,
			.dataSupport = true
		},
		signaturehelp_t {
			signatureinformation2_t {
				.documentationFormat = string_v{"plaintext", "markdown"},
			}
		},
		completion_t {
			capabilitiescompletionitem_t {
				.resolveProvider = true,
				.snippetSupport = true,
				.insertReplaceSupport = true,
				.labelDetailsSupport = true,
				resolvesupport_t {
					.properties = string_v{"documentation", "detail", "additionalTextEdits"},
				}
			}
		}
	}
};

#endif // _LSP_CAPABILITIES_H
