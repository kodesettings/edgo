#ifndef _CONFIG_H_
#define _CONFIG_H_

#include <fstream>
#include <iostream>
#include <string>
#include <map>

#include <boost/property_tree/ptree.hpp>
#include <boost/property_tree/json_parser.hpp>

typedef struct {
	std::string name;
	std::string lsp;
	std::string comment;
	int tabwidth;
	std::string cmd;
	std::string cmdargs;
} lang;

typedef struct {
	std::map<std::string, lang> langs;
	std::string theme;
} config;

static config default_config { .langs = {
	{ "go",         { lsp: "gopls", tabwidth: 4, cmd: "go", cmdargs: "run" }},
	{ "python",     { lsp: "pyright-langserver --stdio", comment: "#", tabwidth: 4, cmd: "python3" }},
	{ "typescript", { lsp: "typescript-language-server --stdio", cmd: "tsx" }},
	{ "javascript", { lsp: "typescript-language-server --stdio", cmd: "tsx" }},
	{ "html",       { lsp: "vscode-html-language-server --stdio" }},
	{ "vue",        { lsp: "vscode-html-language-server --stdio" }},
	{ "rust",       { lsp: "rust-analyzer", tabwidth: 4}},
	{ "c",          { lsp: "clangd" }},
	{ "c++",        { lsp: "clangd" }},
	{ "d",          { lsp: "serve-d", cmd: "dmd", cmdargs: "-run"}},
	{ "java",       { lsp: "jdtls", tabwidth: 4, cmd: "java" }},
	{ "swift",      { lsp: "xcrun sourcekit-lsp", cmd: "swift" }},
	{ "haskell",    { lsp: "haskell-language-server-wrapper --lsp", comment: "--" }},
	{ "zig",        { lsp: "zls", tabwidth: 4, cmd: "zig", cmdargs: "run" }},
	{ "lua",        { lsp: "lua-language-server", cmd: "lua" }},
	{ "yaml",       { comment: "#", tabwidth: 4 }},
	{ "ocaml",      { lsp: "ocamllsp", }},
	{ "nim",        { lsp: "nimlangserver", }},
	{ "bash",       { lsp: "bash-language-server start", comment: "#", tabwidth: 2, cmd: "bash" }},
	{ "shell",      { lsp:"bash-language-server start", comment: "#", tabwidth: 2, cmd: "bash" }},
}
};

lang default_lang_config { name: "",.lsp = "",.comment = "//",.tabwidth = 2 };

using boost::property_tree::ptree;

config ParseLang(const ptree& pt);
void OverrideDefaultConfig(config *provided_config);
config GetConfig(void);

#endif // _CONFIG_H_
