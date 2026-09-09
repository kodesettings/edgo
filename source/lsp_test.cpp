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

#include "lsp_test.h"
#include "../include/edgo.h"
#include "../include/lsp.h"

TEST_F(LspTest, TestLspApiClientHover) {
	move_cursor(223-1, 18);

	struct lsp_hover lsp_hover;
	lsp_client_hover(&lsp_hover);

	EXPECT_STREQ(lsp_hover.contents.language, "plaintext");
	EXPECT_NE(strlen(lsp_hover.contents.value), 0);

	EXPECT_EQ(lsp_hover.range.start.line, 222);
	EXPECT_EQ(lsp_hover.range.start.character, 11);
	EXPECT_EQ(lsp_hover.range.end.line, 222);
	EXPECT_EQ(lsp_hover.range.end.character, 21);
}

TEST_F(LspTest, TestLspApiClientCompletion) {
	move_cursor(223-1, 18);

	struct lsp_completion lsp_completion;
	lsp_client_completion(&lsp_completion);

	EXPECT_EQ(lsp_completion.is_incomplete, false);
	EXPECT_GT(lsp_completion.count, 0);

	for (size_t i = 0; i < lsp_completion.count; i++) {
		auto item = lsp_completion.items[i];
		EXPECT_STRNE(item.label, "");
		EXPECT_EQ(item.kind, 0);
		EXPECT_STRNE(item.detail, "");
		EXPECT_STRNE(item.documentation, "");
		EXPECT_STRNE(item.insert_text, "");

		EXPECT_EQ(item.text_edit.range.start.line, 222);
		EXPECT_EQ(item.text_edit.range.start.character, 11);
		EXPECT_EQ(item.text_edit.range.end.line, 222);
		EXPECT_EQ(item.text_edit.range.end.character, 21);
	}
}

TEST_F(LspTest, TestLspApiClientDefinition) {
	move_cursor(223-1, 18);

	struct lsp_definition lsp_definition;
	lsp_client_definition(&lsp_definition);

	EXPECT_GT(lsp_definition.count, 0);

	for (size_t i = 0; i < lsp_definition.count; i++) {
		auto location = lsp_definition.locations[i];
		EXPECT_STRNE(location.uri, "");
		EXPECT_EQ(location.range.start.line, 222);
		EXPECT_EQ(location.range.start.character, 11);
		EXPECT_EQ(location.range.end.line, 222);
		EXPECT_EQ(location.range.end.character, 21);
	}
}

TEST_F(LspTest, TestLspApiClientSignatureHelp) {
	move_cursor(223-1, 18);

	struct lsp_signature_help lsp_signature_help;
	lsp_client_signature_help(&lsp_signature_help);

	EXPECT_EQ(lsp_signature_help.signature_count, 0);
	EXPECT_EQ(lsp_signature_help.active_signature, 0);
	EXPECT_EQ(lsp_signature_help.active_parameter, 0);

	for (size_t i = 0; i < lsp_signature_help.signature_count; i++) {
		auto signatures = lsp_signature_help.signatures[i];
		EXPECT_STRNE(signatures.label, "");
		EXPECT_EQ(signatures.parameter_count, 0);

		for (size_t j = 0; j < signatures.parameter_count; j++) {
			auto parameters = signatures.parameters[j];
			EXPECT_STRNE(parameters.label, "");
			EXPECT_STRNE(parameters.documentation, "");
		}
	}
}

TEST_F(LspTest, TestLspApiClientReferences) {
	move_cursor(223-1, 18);

	struct lsp_references lsp_references;
	lsp_client_references(&lsp_references);

	EXPECT_NE(lsp_references.count, 0);

	for (size_t i = 0; i < lsp_references.count; i++) {
		auto locations = lsp_references.locations[i];
		EXPECT_STRNE(locations.uri, "");
		EXPECT_EQ(locations.range.start.line, 222);
		EXPECT_EQ(locations.range.start.character, 11);
		EXPECT_EQ(locations.range.end.line, 222);
		EXPECT_EQ(locations.range.end.character, 21);
	}
}

TEST_F(LspTest, TestLspApiClientPrepareRename) {
	move_cursor(223-1, 18);

	struct lsp_prepare_rename lsp_prepare_rename;
	lsp_client_prepare_rename(&lsp_prepare_rename);

	EXPECT_STRNE(lsp_prepare_rename.placeholder, "");
	EXPECT_EQ(lsp_prepare_rename.range.start.line, 222);
	EXPECT_EQ(lsp_prepare_rename.range.start.character, 11);
	EXPECT_EQ(lsp_prepare_rename.range.end.line, 222);
	EXPECT_EQ(lsp_prepare_rename.range.end.character, 21);
}

TEST_F(LspTest, TestLspApiClientRename) {
	move_cursor(223-1, 18);

	struct lsp_rename lsp_rename;
	lsp_client_rename("new-name", &lsp_rename);

	EXPECT_NE(lsp_rename.document_change_count, 0);

	for (size_t i = 0; i < lsp_rename.document_change_count; i++) {
		auto workspace_edit = lsp_rename.document_changes[i].edit;
		auto text_document = lsp_rename.document_changes[i].text_document;
		EXPECT_STRNE(text_document.uri, "");
		EXPECT_GE(text_document.version, 0);
		EXPECT_GE(workspace_edit.edit_count, 0);

		for (size_t j = 0; j < workspace_edit.edit_count; j++) {
			auto edits = workspace_edit.edits[j];
			EXPECT_STREQ(edits.new_text, "new-name");

			EXPECT_EQ(edits.range.start.line, 222);
			EXPECT_EQ(edits.range.start.character, 11);
			EXPECT_EQ(edits.range.end.line, 222);
			EXPECT_EQ(edits.range.end.character, 21);
		}
	}
}

TEST_F(LspTest, TestLspApiClientCodeAction) {
	move_cursor(223-1, 18);

	struct lsp_code_action lsp_code_action;
	lsp_client_code_action(&lsp_code_action);

	EXPECT_NE(lsp_code_action.code_action_count, 0);

	for (size_t i = 0; i < lsp_code_action.code_action_count; i++) {
		auto item = lsp_code_action.items[i];
		EXPECT_STRNE(item.title, "");
		EXPECT_STRNE(item.kind, "");
		EXPECT_EQ(item.edit.edit_count, 1);
	}
}

TEST_F(LspTest, TestLspApiPublishDiagnostics) {
	struct lsp_publish_diagnostics lsp_publish_diagnostics;
	lsp_diagnostics(&lsp_publish_diagnostics);

	EXPECT_EQ(lsp_publish_diagnostics.count, 0);
	EXPECT_STRNE(lsp_publish_diagnostics.uri, "");

	for (size_t i = 0; i < lsp_publish_diagnostics.count; i++) {
		auto diagnostic = lsp_publish_diagnostics.diagnostics[i];
		EXPECT_GE(diagnostic.severity, -1);
		EXPECT_STRNE(diagnostic.source, "");
		EXPECT_STRNE(diagnostic.message, "");

		EXPECT_EQ(diagnostic.range.start.line, 222);
		EXPECT_EQ(diagnostic.range.start.character, 11);
		EXPECT_EQ(diagnostic.range.end.line, 222);
		EXPECT_EQ(diagnostic.range.end.character, 21);
	}
}
