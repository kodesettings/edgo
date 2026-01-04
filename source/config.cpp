#include "config.h"

config ParseLang(const ptree& pt) {
	config conf;

	for (const auto& item : pt.get_child("langs")) {
		lang l;
		l.name     = item.second.get<std::string>("name");
		l.lsp      = item.second.get<std::string>("lsp");
		l.comment  = item.second.get<std::string>("comment");
		l.tabwidth = item.second.get<int>("tabwidth");
		l.cmd      = item.second.get<std::string>("cmd");
		l.cmdargs  = item.second.get<std::string>("cmdargs");
		conf.langs[item.first] = l;
	}

	conf.theme = pt.get<std::string>("theme");
	size_t count = pt.get_child("langs").size();
	return count == 0 ? default_config : conf;
}

void OverrideDefaultConfig(config *provided_config) {
	const auto m_langs = provided_config->langs;

	// override default config
	for (auto it = m_langs.begin(); it != m_langs.end(); ++it) {
		//set default tab width and comment if not specified
		lang langConf = it->second;
		std::string langName = it->first;
		if (langConf.tabwidth == 0) { langConf.tabwidth = 2; }
		if (langConf.comment == "") { langConf.comment = "//"; }
		provided_config->langs[langName] = langConf;
	}
}

config GetConfig(void) {
	// override default config
	OverrideDefaultConfig(&default_config);
	default_config.theme = "vesper";

	const char *conffilename = getenv("EDGO_CONF");
	if (conffilename == NULL) {
		conffilename = "config.json";
	}

	std::ifstream file(conffilename);
	if (!file.is_open()) {
		return default_config;
	}

	ptree pt;
	config json_config = ParseLang(pt);

	// read json config and override
	OverrideDefaultConfig(&json_config);

	if (json_config.theme != "") { default_config.theme = json_config.theme; }

	return default_config;
}
