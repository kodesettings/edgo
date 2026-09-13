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

#include "edgo.h"
#include "lsp.h"
#include "lsp_test.h"

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
	EXPECT_EQ(lsp_completion.count, 1);

	if (lsp_completion.count == 0)
		return;

	auto item = lsp_completion.items[0];
	EXPECT_STREQ(item.label, " tss_create");
	EXPECT_EQ(item.kind, 3);
	EXPECT_STREQ(item.detail, "int");
	EXPECT_STREQ(item.documentation, "40df55d4tss_create");
	EXPECT_STREQ(item.insert_text,
		"tss_create(${1:tss_t *tss_id}, ${2:tss_dtor_t destructor})");

	EXPECT_EQ(item.text_edit.range.start.line, 222);
	EXPECT_EQ(item.text_edit.range.start.character, 11);
	EXPECT_EQ(item.text_edit.range.end.line, 222);
	EXPECT_EQ(item.text_edit.range.end.character, 18);
}

TEST_F(LspTest, TestLspApiClientDefinition) {
	move_cursor(223-1, 18);

	struct lsp_definition lsp_definition;
	lsp_client_definition(&lsp_definition);

	EXPECT_EQ(lsp_definition.count, 1);

	auto location = lsp_definition.locations[0];
	EXPECT_STREQ(location.uri, "file:///usr/include/threads.h");

	EXPECT_EQ(location.range.start.line, 222);
	EXPECT_EQ(location.range.start.character, 11);
	EXPECT_EQ(location.range.end.line, 222);
	EXPECT_EQ(location.range.end.character, 21);
}

TEST_F(LspTest, TestLspApiClientSignatureHelp) {
	move_cursor(223-1, 18);

	struct lsp_signature_help lsp_signature_help;
	lsp_client_signature_help(&lsp_signature_help);

	EXPECT_EQ(lsp_signature_help.signature_count, 0);
	EXPECT_EQ(lsp_signature_help.active_signature, 0);
	EXPECT_EQ(lsp_signature_help.active_parameter, 0);
}

TEST_F(LspTest, TestLspApiClientReferences) {
	move_cursor(223-1, 18);

	struct lsp_references lsp_references;
	lsp_client_references(&lsp_references);

	EXPECT_EQ(lsp_references.count, 0);
}

TEST_F(LspTest, TestLspApiClientPrepareRename) {
	move_cursor(223-1, 18);

	struct lsp_prepare_rename lsp_prepare_rename;
	lsp_client_prepare_rename(&lsp_prepare_rename);

	EXPECT_STREQ(lsp_prepare_rename.placeholder, "tss_create");
	EXPECT_EQ(lsp_prepare_rename.range.start.line, 222);
	EXPECT_EQ(lsp_prepare_rename.range.start.character, 11);
	EXPECT_EQ(lsp_prepare_rename.range.end.line, 222);
	EXPECT_EQ(lsp_prepare_rename.range.end.character, 21);
}

TEST_F(LspTest, TestLspApiClientRename) {
	move_cursor(223-1, 18);

	struct lsp_rename lsp_rename;
	lsp_client_rename("tss_create2", &lsp_rename);

	EXPECT_EQ(lsp_rename.document_change_count, 1);

	if (lsp_rename.document_change_count == 0)
		return;

	auto changes = lsp_rename.document_changes[0];
	auto workspace_edit = changes.edit;
	auto text_document = changes.text_document;
	EXPECT_STRNE(text_document.uri, "");
	EXPECT_GE(text_document.version, 0);
	EXPECT_GE(workspace_edit.edit_count, 0);

	for (size_t j = 0; j < workspace_edit.edit_count; j++) {
		auto edits = workspace_edit.edits[j];
		EXPECT_STREQ(edits.new_text, "tss_create2");

		EXPECT_EQ(edits.range.start.line, 222);
		EXPECT_EQ(edits.range.start.character, 11);
		EXPECT_EQ(edits.range.end.line, 222);
		EXPECT_EQ(edits.range.end.character, 21);
	}
}

TEST_F(LspTest, TestLspApiClientCodeAction) {
	move_cursor(223-1, 18);

	struct lsp_code_action lsp_code_action;
	lsp_client_code_action(&lsp_code_action);

	EXPECT_EQ(lsp_code_action.code_action_count, 1);

	if (lsp_code_action.code_action_count == 0)
		return;

	auto item = lsp_code_action.items[0];
	EXPECT_STRNE(item.title, "");
	EXPECT_STRNE(item.kind, "");
	EXPECT_EQ(item.edit.edit_count, 1);
}

TEST_F(LspTest, TestLspApiPublishDiagnostics) {
	struct lsp_publish_diagnostics lsp_publish_diagnostics;
	lsp_diagnostics(&lsp_publish_diagnostics);

	EXPECT_EQ(lsp_publish_diagnostics.count, 0);
	EXPECT_STREQ(lsp_publish_diagnostics.uri,
		"file:///usr/include/threads.h");
}
