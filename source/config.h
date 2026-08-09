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

#ifndef _CONFIG_H_
#define _CONFIG_H_

#include <fstream>
#include <iostream>
#include <string>
#include <map>

#include <boost/property_tree/ptree.hpp>
#include <boost/property_tree/ini_parser.hpp>

typedef struct {
	std::string name;
	std::string lsp;
	std::string comment;
	int tabwidth;
	std::string cmd;
	std::string cmdargs;
} lang_t;

typedef struct {
	std::map<std::string, lang_t> langs;
	std::string theme;
} config_t;

static config_t default_config { .langs = {
	{ "go",         { lsp: "gopls", comment: "//", tabwidth: 4, cmd: "go", cmdargs: "run" }},
	{ "python",     { lsp: "pyright-langserver --stdio", comment: "#", tabwidth: 4, cmd: "python3" }},
	{ "typescript", { lsp: "typescript-language-server --stdio", cmd: "tsx" }},
	{ "javascript", { lsp: "typescript-language-server --stdio", cmd: "tsx" }},
	{ "html",       { lsp: "", cmd: "vscode-html-language-server --stdio" }},
	{ "vue",        { lsp: "", cmd: "vscode-html-language-server --stdio" }},
	{ "rust",       { lsp: "", comment: "//", tabwidth: 4, cmd: "rust-analyzer" }},
	{ "c",          { lsp: "", comment: "//", tabwidth: 4, cmd: "clangd" }},
	{ "cpp",        { lsp: "", comment: "//", tabwidth: 4, cmd: "clangd" }},
	{ "c++",        { lsp: "", comment: "//", tabwidth: 4, cmd: "clangd" }},
	{ "d",          { lsp: "serve-d", comment: "//", tabwidth: 4, cmd: "dmd", cmdargs: "-run"}},
	{ "java",       { lsp: "jdtls", comment: "//", tabwidth: 4, cmd: "java" }},
	{ "swift",      { lsp: "xcrun sourcekit-lsp", comment: "//", tabwidth: 4, cmd: "swift" }},
	{ "haskell",    { lsp: "haskell-language-server-wrapper --lsp", comment: "--", tabwidth: 2 }},
	{ "zig",        { lsp: "zls", comment: "//", tabwidth: 4, cmd: "zig", cmdargs: "run" }},
	{ "lua",        { lsp: "lua-language-server", cmd: "lua" }},
	{ "yaml",       { comment: "#", tabwidth: 4 }},
	{ "ocaml",      { lsp: "", tabwidth: 4, cmd: "ocamllsp" }},
	{ "nim",        { lsp: "", comment: "#", tabwidth: 2, cmd: "nimlangserver" }},
	{ "bash",       { lsp: "bash-language-server start", comment: "#", tabwidth: 2, cmd: "bash" }},
	{ "shell",      { lsp:"bash-language-server start", comment: "#", tabwidth: 2, cmd: "bash" }},
}
};

static lang_t default_lang_config { name: "", lsp: "", comment: "//", tabwidth: 2 };

using boost::property_tree::ptree;

config_t ParseLang(const ptree& pt);
void OverrideDefaultConfig(config_t *provided_config);
config_t GetConfig(void);

#endif // _CONFIG_H_
