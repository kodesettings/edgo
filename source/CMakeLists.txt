cmake_minimum_required(VERSION 3.12)

project(edgo VERSION 0.0.2 LANGUAGES CXX)

find_package(GTest REQUIRED)

# edgo library files
add_library(edgo SHARED
	"config.cpp"
	"config.h"
	"selection.cpp"
	"selection.h"
	"operations.h"
	"lsp_user_actions.cpp"
	"lsp_model.h"
	"lsp_capabilities.h"
	"lsp_client.cpp"
	"lsp_client.h"
	"utils.cpp"
	"utils.h"
	"editor.h"
	"screen_logo.h"
	"characters.cpp"
	"cursor.cpp"
	"highlighter.cpp"
	"highlighter.h"
	"clipboard.cpp"
	"sdl2_init.h"
	"features_v01.cpp"

# embedded scm files
	"lang_bash.h"
	"lang_c.h"
	"lang_cpp.h"
	"lang_css.h"
	"lang_go.h"
	"lang_html.h"
	"lang_java.h"
	"lang_javascript.h"
	"lang_python.h"
	"lang_rust.h"

	"langs.h"
	"themes.h"
	"test_languages.h"
	"search.cpp"
	"search.h"
	"query.cpp"
	"query.h"
	"io.cpp"
)

# test case files
add_executable(edgotest
	"config_test.cpp"
	"selection_test.cpp"
	"lsp_client_test.cpp"
	"utils_test.cpp"
	"characters_test.cpp"
	"clipboard_test.cpp"
	"features_v01_test.cpp"
	"highlighter_test.cpp"
	"search_test.cpp"
	"query_test.cpp"
)

# tree-sitter language bindings
set(TS_LANGUAGE_LIBS
	tree-sitter-java
	tree-sitter-bash
	tree-sitter-cpp
	tree-sitter-html
	tree-sitter-css
	tree-sitter-go
	tree-sitter-c
	tree-sitter-javascript
	tree-sitter-rust
	tree-sitter-python
)

# target and test builds
target_link_libraries(edgo PUBLIC glog pthread boost_filesystem SDL2 tree-sitter ${TS_LANGUAGE_LIBS})
target_link_libraries(edgotest PUBLIC ${GTEST_LIBRARIES} gtest_main pthread edgo)
