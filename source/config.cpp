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

#include "config.h"

config_t ParseLang(const ptree& pt) {
	config_t conf;
	size_t count = 0;

	for (const auto& item : pt) {
		if (item.first.compare(0, 6, "langs.") != 0) {
			continue;
		}

		lang_t lang;
		lang.name     = item.second.get<std::string>("name", "");
		lang.lsp      = item.second.get<std::string>("lsp", "");
		lang.comment  = item.second.get<std::string>("comment", "//");
		lang.tabwidth = item.second.get<int>("tabwidth", 2);
		lang.cmd      = item.second.get<std::string>("cmd", "");
		lang.cmdargs  = item.second.get<std::string>("cmdargs", "");
		conf.langs[item.first] = lang;
		count++;
	}

	conf.theme = pt.get<std::string>("theme", default_config.theme);
	return count == 0 ? default_config : conf;
}

void OverrideDefaultConfig(config_t *provided_config) {
	const auto m_langs = provided_config->langs;

	// override default config
	for (auto it = m_langs.begin(); it != m_langs.end(); ++it) {
		//set default tab width and comment if not specified
		lang_t langConf = it->second;
		std::string langName = it->first;
		if (langConf.tabwidth == 0) { langConf.tabwidth = 2; }
		if (langConf.comment == "") { langConf.comment = "//"; }
		provided_config->langs[langName] = langConf;
	}
}

config_t GetConfig(void) {
	// override default config
	OverrideDefaultConfig(&default_config);
	default_config.theme = "vesper";

	const char *conffilename = getenv("EDGO_CONF");
	if (conffilename == NULL) {
		conffilename = "/etc/edgo.conf";
	}

	std::ifstream file(conffilename);
	if (!file.is_open()) {
		return default_config;
	}

	ptree pt;
	read_ini(conffilename, pt);
	config_t ini_config = ParseLang(pt);

	// read json config and override
	OverrideDefaultConfig(&ini_config);

	if (ini_config.theme != "") { default_config.theme = ini_config.theme; }

	return default_config;
}
