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

#include "../include/edgo.h"
#include "../include/lsp.h"
#include "editor.h"
#define CONVERT_LSP_RANGE(range, orig)                                 \
	range.start.line      = orig.start.line;                           \
	range.start.character = orig.start.character;                      \
	range.end.line        = orig.end.line;                             \
	range.end.character   = orig.end.character;

#define ALLOC_STRING(name, value)                                      \
	static char name[2048];                                            \
	strcpy(name, value.c_str());

#define ALLOC_STRING_L(name, value)                                    \
	static char name[24128];                                           \
	strcpy(name, value.c_str());
#ifdef __cplusplus
extern "C" {
#endif
void lsp_client_hover(struct lsp_hover* lsp_hover) {
	hoverresponse_t hover;
	hover = Hover(e.absoluteFilePath, e.row, e.col);

	struct lsp_range lsp_range;
	CONVERT_LSP_RANGE(lsp_range, hover.result.range);

	ALLOC_STRING(language, hover.result.contents.kind);
	ALLOC_STRING_L(contents, hover.result.contents.value);

	lsp_hover->contents.language = language;
	lsp_hover->contents.value = contents;
	lsp_hover->range = lsp_range;
}

void lsp_client_completion(struct lsp_completion* lsp_completion) {
	completionresponse_t completion;
	completion = Completion(e.absoluteFilePath, e.row, e.col);

	lsp_completion->is_incomplete = completion.result.isIncomplete;
	lsp_completion->count = completion.result.items.size();

	struct lsp_completion_item items[lsp_completion->count];
	lsp_completion->items = &(*items);

	for (size_t i = 0; i < lsp_completion->count; i++) {
		struct lsp_completion_item lsp_completion_item;
		completionitem_t item = completion.result.items[i];;

		ALLOC_STRING(label, item.label);
		ALLOC_STRING(detail, item.detail);
		ALLOC_STRING(documentation, item.sortText);
		ALLOC_STRING(insert_text, item.insertText);

		lsp_completion_item.label = label;
		lsp_completion_item.kind = (lsp_completion_item_kind)item.kind;
		lsp_completion_item.detail = detail;
		lsp_completion_item.documentation = documentation;
		lsp_completion_item.insert_text = insert_text;

		struct lsp_range lsp_range;
		textedit_t text_edit = item.textEdit;
		CONVERT_LSP_RANGE(lsp_range, item.textEdit.range);

		lsp_completion_item.text_edit.range = lsp_range;
		lsp_completion_item.text_edit.new_text = item.textEdit.newText.c_str();
		items[i] = lsp_completion_item;
	}
}

void lsp_client_definition(struct lsp_definition* lsp_definition) {
	definitionresponse_t definition;
	definition = Definition(e.absoluteFilePath, e.row, e.col);

	lsp_definition->count = definition.result.size();

	struct lsp_location locations[lsp_definition->count];
	lsp_definition->locations = &(*locations);

	for (size_t i = 0; i < lsp_definition->count; i++) {
		struct lsp_location lsp_location;
		definitionresult_t result = definition.result[i];

		ALLOC_STRING(uri, result.uri);

		lsp_location.uri = uri;
		CONVERT_LSP_RANGE(lsp_location.range, result.range);
		locations[i] = lsp_location;
	}
}

void lsp_client_signature_help(struct lsp_signature_help* lsp_signature_help) {
	signaturehelpresponse_t signaturehelp;
	signaturehelp = SignatureHelp(e.absoluteFilePath, e.row, e.col);

	lsp_signature_help->signature_count = signaturehelp.result.signatures.size();
	lsp_signature_help->active_signature = signaturehelp.result.activeSignature;
	lsp_signature_help->active_parameter = signaturehelp.result.activeParameter;

	size_t signature_count = lsp_signature_help->signature_count;
	struct lsp_signature_information signatures[signature_count];
	lsp_signature_help->signatures = &(*signatures);

	for (size_t i = 0; i < lsp_signature_help->signature_count; i++) {
		signatureinformation_t signature = signaturehelp.result.signatures[i];
		struct lsp_signature_information lsp_signature_information;

		ALLOC_STRING(label, signature.label);

		lsp_signature_information.label = label;
		lsp_signature_information.parameter_count = signature.parameters.size();

		size_t parameter_count = lsp_signature_information.parameter_count;
		struct lsp_parameter_information parameters[parameter_count];
		lsp_signature_information.parameters = &(*parameters);

		for (size_t j = 0; j < parameter_count; j++) {
			parameterinformation_t parameter = signature.parameters[j];
			struct lsp_parameter_information lsp_parameter_information;

			ALLOC_STRING(label, parameter.label);
			ALLOC_STRING(documentation, parameter.documentation);

			lsp_parameter_information.label = label;
			lsp_parameter_information.documentation = documentation;
			parameters[j] = lsp_parameter_information;
		}

		signatures[i] = lsp_signature_information;
	}
}

void lsp_client_references(struct lsp_references* lsp_references) {
	referencesresponse_t references;
	references = References(e.absoluteFilePath, e.row, e.col);

	lsp_references->count = references.result.size();

	struct lsp_location locations[lsp_references->count];
	lsp_references->locations = &(*locations);

	for (size_t i = 0; i < lsp_references->count; i++) {
		struct lsp_location lsp_location;
		referencesrange_t result = references.result[i];

		ALLOC_STRING(uri, result.uri);

		lsp_location.uri = uri;
		CONVERT_LSP_RANGE(lsp_location.range, result.range);
		locations[i] = lsp_location;
	}
}

void lsp_client_prepare_rename(struct lsp_prepare_rename* lsp_prepare_rename) {
	preparerenameresponse_t preparerename;
	preparerename = PrepareRename(e.absoluteFilePath, e.row, e.col);

	CONVERT_LSP_RANGE(lsp_prepare_rename->range, preparerename.result.range);
	ALLOC_STRING(placeholder, preparerename.result.placeholder);

	lsp_prepare_rename->placeholder = placeholder;
}

#define LSP_TEXT_EDIT(destination, edit)                               \
	struct lsp_text_edit lsp_text_edit;                                \
	struct lsp_range lsp_range;                                        \
	CONVERT_LSP_RANGE(lsp_range, edit.range);                          \
	lsp_text_edit.range = lsp_range;                                   \
	lsp_text_edit.new_text = edit.newText.c_str();                     \
	destination = lsp_text_edit;

void lsp_client_rename(const char* newname, struct lsp_rename* lsp_rename) {
	renameresponse_t rename;
	rename = Rename(e.absoluteFilePath, newname, e.row, e.col);

	lsp_rename->document_change_count = rename.result.documentChanges.size();

	size_t document_change_count = lsp_rename->document_change_count;
	struct lsp_rename_document_change document_changes[document_change_count];
	lsp_rename->document_changes = &(*document_changes);

	for (size_t i = 0; i < lsp_rename->document_change_count; i++) {
		struct lsp_versioned_text_document_identifier text_document;
		documentchange_t documentchange = rename.result.documentChanges[i];

		ALLOC_STRING(uri, documentchange.textDocument.uri);

		text_document.uri = uri;
		text_document.version = documentchange.textDocument.version;
		document_changes[i].text_document = text_document;

		struct lsp_workspace_edit lsp_workspace_edit;
		lsp_workspace_edit.edit_count =	documentchange.edits.size();

		struct lsp_text_edit edits[lsp_workspace_edit.edit_count];
		lsp_workspace_edit.edits = &(*edits);

		for (size_t j = 0; j < lsp_workspace_edit.edit_count; j++) {
			LSP_TEXT_EDIT(edits[j], documentchange.edits[j]);
		}

		document_changes[i].edit = lsp_workspace_edit;
	}
}

void lsp_client_code_action(struct lsp_code_action* lsp_code_action) {
	codeactionresponse_t codeaction;
	codeaction = CodeAction(e.absoluteFilePath, e.__selection.ssx,
		e.__selection.ssy, e.__selection.sex, e.__selection.sey);

	lsp_code_action->code_action_count = codeaction.result.size();

	struct lsp_code_action_item items[lsp_code_action->code_action_count];
	lsp_code_action->items = &(*items);

	for (size_t i = 0; i < lsp_code_action->code_action_count; i++) {
		struct lsp_code_action_item lsp_code_action_item;
		codeactionresult_t result = codeaction.result[i];

		ALLOC_STRING(title, result.title);
		ALLOC_STRING(kind, result.kind);

		lsp_code_action_item.title = title;
		lsp_code_action_item.kind = kind;
		lsp_code_action_item.edit.edit_count = 1;

		LSP_TEXT_EDIT(lsp_code_action_item.edit.edits[0], result.edit);
		items[i] = lsp_code_action_item;
	}
}

void lsp_diagnostics(struct lsp_publish_diagnostics* lsp_publish_diagnostics) {
	// TODO: we only take the first entry here, as we handle only one file
	// per session. Currently multifile support is implemented in the lsp using
	// channels, but not in the API session.
	auto it = lspclient.file2diagnostic.begin();

	diagnosticparams_t diagnosticparams = it->second;
	lsp_publish_diagnostics->count = diagnosticparams.diagnostics.size();

	ALLOC_STRING(uri, diagnosticparams.uri);
	lsp_publish_diagnostics->uri = uri;

	struct lsp_diagnostic diagnostics[lsp_publish_diagnostics->count];
	lsp_publish_diagnostics->diagnostics = &(*diagnostics);

	for (size_t i = 0; i < lsp_publish_diagnostics->count; i++) {
		struct lsp_diagnostic lsp_diagnostic;
		diagnostic_t diagnostic = diagnosticparams.diagnostics[i];
		CONVERT_LSP_RANGE(lsp_diagnostic.range, diagnostic.range);
		lsp_diagnostic.severity = (lsp_diagnostic_severity)diagnostic.severity;

		ALLOC_STRING(source, diagnostic.source);
		ALLOC_STRING(message, diagnostic.message);

		lsp_diagnostic.source = diagnostic.source.c_str();
		lsp_diagnostic.message = diagnostic.message.c_str();
		diagnostics[i] = lsp_diagnostic;
	}
}
#ifdef __cplusplus
}
#endif
