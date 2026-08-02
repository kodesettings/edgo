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
#include <stdint.h>
#include <stddef.h>
#include <stdbool.h>

//
// Basic types
//

struct lsp_position {
	uint32_t line;
	uint32_t character;
};

struct lsp_range {
	struct lsp_position start;
	struct lsp_position end;
};

struct lsp_location {
	const char *uri;
	struct lsp_range range;
};

//
// Text document
//

struct lsp_text_document_identifier {
	const char *uri;
};

struct lsp_versioned_text_document_identifier {
	const char *uri;
	int version;
};

struct lsp_text_document_position {
	struct lsp_text_document_identifier text_document;
	struct lsp_position position;
};

//
// Text edits
//

struct lsp_text_edit {
	struct lsp_range range;
	const char *new_text;
};

struct lsp_workspace_edit {
	size_t edit_count;
	struct lsp_text_edit *edits;
};

//
// Diagnostics
//

enum lsp_diagnostic_severity {
	LSP_SEVERITY_ERROR,
	LSP_SEVERITY_WARNING,
	LSP_SEVERITY_INFORMATION,
	LSP_SEVERITY_HINT
};

struct lsp_diagnostic {
	struct lsp_range range;
	enum lsp_diagnostic_severity severity;
	const char *source;
	const char *message;
};

struct lsp_publish_diagnostics {
	const char *uri;
	size_t count;
	struct lsp_diagnostic *diagnostics;
};

//
// Completion
//

enum lsp_completion_item_kind {
	LSP_COMPLETION_TEXT,
	LSP_COMPLETION_METHOD,
	LSP_COMPLETION_FUNCTION,
	LSP_COMPLETION_VARIABLE,
	LSP_COMPLETION_CLASS
};

struct lsp_completion_item {
	const char *label;
	enum lsp_completion_item_kind kind;
	const char *detail;
	const char *documentation;
	const char *insert_text;
	struct lsp_text_edit *text_edit;
};

struct lsp_completion_list {
	bool is_incomplete;
	size_t count;
	struct lsp_completion_item *items;
};

//
// Hover
//

struct lsp_marked_string {
	const char *language;
	const char *value;
};

struct lsp_hover {
	size_t count;
	struct lsp_marked_string *contents;
	struct lsp_range *range;
};

//
// Definition
//

struct lsp_definition_params {
	struct lsp_text_document_position position;
};

struct lsp_definition {
	size_t count;
	struct lsp_location *locations;
};

//
// Signature Help
//

struct lsp_parameter_information {
	const char *label;
	const char *documentation;
};

struct lsp_signature_information {
	const char *label;
	const char *documentation;
	size_t parameter_count;
	struct lsp_parameter_information *parameters;
};

struct lsp_signature_help {
	size_t signature_count;
	struct lsp_signature_information *signatures;
	int active_signature;
	int active_parameter;
};

//
// References
//

struct lsp_reference_params {
	struct lsp_text_document_position position;
	bool include_declaration;
};

struct lsp_references {
	size_t count;
	struct lsp_location *locations;
};


//
// Prepare Rename
//

struct lsp_prepare_rename_params {
	struct lsp_text_document_position position;
};

struct lsp_prepare_rename {
	bool valid;
	struct lsp_range range;
	const char *placeholder;
};

//
// Rename
//

struct lsp_rename_params {
	struct lsp_text_document_position position;
	const char *new_name;
};

struct lsp_rename {
	struct lsp_workspace_edit edit;
};

//
// Code Action
//

struct lsp_code_action_context {
	size_t diagnostic_count;
	struct lsp_diagnostic *diagnostics;
	size_t only_count;
	const char **only;
};

struct lsp_code_action_params {
	struct lsp_text_document_identifier text_document;
	struct lsp_range range;
	struct lsp_code_action_context context;
};

struct lsp_code_action {
	const char *title;
	const char *kind;
	bool is_preferred;
	size_t diagnostic_count;
	struct lsp_diagnostic *diagnostics;
	struct lsp_workspace_edit *edit;
	struct lsp_command *command;
};

struct lsp_code_actions {
	size_t count;
	struct lsp_code_action *items;
};
