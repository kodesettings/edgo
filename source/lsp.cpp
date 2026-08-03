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
#define LSP_TEXT_EDIT(destination, edit)                               \
	struct lsp_text_edit lsp_text_edit;                                \
	struct lsp_range lsp_range;                                        \
	CONVERT_LSP_RANGE(lsp_range, edit.range);                          \
	lsp_text_edit.range = lsp_range;                                   \
	lsp_text_edit.new_text = edit.newText.c_str();                     \
	destination = lsp_text_edit;
#ifdef __cplusplus
extern "C" {
#endif
void lsp_client_hover(struct lsp_hover* lsp_hover) {
	hoverresponse_t hover;
	hover = Hover(e.absoluteFilePath, e.row, e.col);

	struct lsp_range lsp_range;
	CONVERT_LSP_RANGE(lsp_range, hover.result.range);

	lsp_hover->contents.language = hover.result.contents.kind.c_str();
	lsp_hover->contents.value = hover.result.contents.value.c_str();
	lsp_hover->range = lsp_range;
}

void lsp_client_completion(struct lsp_completion* lsp_completion) {
	completionresponse_t completion;
	completion = Completion(e.absoluteFilePath, e.row, e.col);

	lsp_completion->is_incomplete = completion.result.isIncomplete;
	lsp_completion->count = completion.result.items.size();

	for (int i = 0; i < (int)lsp_completion->count; i++) {
		struct lsp_completion_item lsp_completion_item;
		completionitem_t item = completion.result.items[i];;

		lsp_completion_item.label = item.label.c_str();
		lsp_completion_item.kind = (lsp_completion_item_kind)item.kind;
		lsp_completion_item.detail = item.label.c_str();
		lsp_completion_item.documentation = item.sortText.c_str();
		lsp_completion_item.insert_text = item.insertText.c_str();

		struct lsp_range lsp_range;
		textedit_t text_edit = item.textEdit;
		CONVERT_LSP_RANGE(lsp_range, item.textEdit.range);

		lsp_completion_item.text_edit.range = lsp_range;
		lsp_completion_item.text_edit.new_text = item.textEdit.newText.c_str();
		lsp_completion->items[i] = lsp_completion_item;
	}
}

void lsp_client_definition(struct lsp_definition* lsp_definition) {
	definitionresponse_t definition;
	definition = Definition(e.absoluteFilePath, e.row, e.col);

	lsp_definition->count = definition.result.size();

	for (int i = 0; i < (int)lsp_definition->count; i++) {
		struct lsp_location lsp_location;
		definitionresult_t result = definition.result[i];

		lsp_location.uri = result.uri.c_str();
		CONVERT_LSP_RANGE(lsp_location.range, result.range);
		lsp_definition->locations[i] = lsp_location;
	}
}

void lsp_client_signature_help(struct lsp_signature_help* lsp_signature_help) {
	signaturehelpresponse_t signaturehelp;
	signaturehelp = SignatureHelp(e.absoluteFilePath, e.row, e.col);

	lsp_signature_help->signature_count = signaturehelp.result.signatures.size();
	lsp_signature_help->active_signature = signaturehelp.result.activeSignature;
	lsp_signature_help->active_parameter = signaturehelp.result.activeParameter;

	for (int i = 0; i < (int)lsp_signature_help->signature_count; i++) {
		signatureinformation_t signature = signaturehelp.result.signatures[i];
		struct lsp_signature_information lsp_signature_information;
		lsp_signature_information.label = signature.label.c_str();
		lsp_signature_information.parameter_count = signature.parameters.size();

		for (int j = 0; j < (int)lsp_signature_information.parameter_count; j++) {
			parameterinformation_t parameter = signature.parameters[j];
			struct lsp_parameter_information lsp_parameter_information;
			lsp_parameter_information.label = parameter.label.c_str();
			lsp_parameter_information.documentation = parameter.documentation.c_str();
			lsp_signature_information.parameters[j] = lsp_parameter_information;
		}

		lsp_signature_help->signatures[i] = lsp_signature_information;
	}
}

void lsp_client_references(struct lsp_references* lsp_references) {
	referencesresponse_t references;
	references = References(e.absoluteFilePath, e.row, e.col);

	lsp_references->count = references.result.size();

	for (int i = 0; i < (int)lsp_references->count; i++) {
		struct lsp_location lsp_location;
		referencesrange_t result = references.result[i];

		lsp_location.uri = result.uri.c_str();
		CONVERT_LSP_RANGE(lsp_location.range, result.range);
		lsp_references->locations[i] = lsp_location;
	}
}

void lsp_client_prepare_rename(struct lsp_prepare_rename* lsp_prepare_rename) {
	preparerenameresponse_t preparerename;
	preparerename = PrepareRename(e.absoluteFilePath, e.row, e.col);

	CONVERT_LSP_RANGE(lsp_prepare_rename->range, preparerename.result.range);
	lsp_prepare_rename->placeholder = preparerename.result.placeholder.c_str();
}

void lsp_client_rename(const char* newname, struct lsp_rename* lsp_rename) {
	renameresponse_t rename;
	rename = Rename(e.absoluteFilePath, newname, e.row, e.col);

	lsp_rename->document_change_count = rename.result.documentChanges.size();

	for (int i = 0; i < (int)lsp_rename->document_change_count; i++) {
		struct lsp_versioned_text_document_identifier text_document;
		documentchange_t documentchange = rename.result.documentChanges[i];
		text_document.uri = documentchange.textDocument.uri.c_str();
		text_document.version = documentchange.textDocument.version;
		lsp_rename->document_changes[i].text_document = text_document;

		struct lsp_workspace_edit lsp_workspace_edit;
		lsp_workspace_edit.edit_count =	documentchange.edits.size();

		for (int j = 0; j < (int) lsp_workspace_edit.edit_count; j++) {
			LSP_TEXT_EDIT(lsp_workspace_edit.edits[j], documentchange.edits[j]);
		}

		lsp_rename->document_changes[i].edit = lsp_workspace_edit;
	}
}

void lsp_client_code_action(struct lsp_code_action* lsp_code_action) {
	codeactionresponse_t codeaction;
	codeaction = CodeAction(e.absoluteFilePath, e.__selection.ssx,
		e.__selection.ssy, e.__selection.sex, e.__selection.sey);

	lsp_code_action->code_action_count = codeaction.result.size();

	for (int i = 0; i < (int)lsp_code_action->code_action_count; i++) {
		struct lsp_code_action_item lsp_code_action_item;
		codeactionresult_t result = codeaction.result[i];
		lsp_code_action_item.title = result.title.c_str();
		lsp_code_action_item.kind = result.kind.c_str();
		lsp_code_action_item.edit.edit_count = 1;
		LSP_TEXT_EDIT(lsp_code_action_item.edit.edits[0], result.edit);
		lsp_code_action->items[i] = lsp_code_action_item;
	}
}
#ifdef __cplusplus
}
#endif
