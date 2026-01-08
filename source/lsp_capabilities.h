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
