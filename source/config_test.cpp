#include "config.h"
#include <gtest/gtest.h>

TEST(ConfigTests, TestReadConfig) {
	// Read conf
	auto conf = GetConfig();

	// Print config
	for (auto lang : conf.langs) {
		printf("name: %s, lsp: %s, comment: %s, tab width: %d\n",
			lang.first.c_str(), lang.second.lsp.c_str(),
			lang.second.comment.c_str(), lang.second.tabwidth
		);
	}

	auto golang = conf.langs["go"];

	EXPECT_EQ(golang.lsp, "gopls");
	EXPECT_EQ(golang.comment, "//");
	EXPECT_EQ(golang.tabwidth, 4);
}
