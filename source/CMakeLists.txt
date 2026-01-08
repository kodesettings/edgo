cmake_minimum_required(VERSION 3.12)

project(edgo VERSION 0.0.2 LANGUAGES CXX)

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
	"highlighter.h"
)

# test case files
add_executable(edgotest
	"config_test.cpp"
	"selection_test.cpp"
	"lsp_client_test.cpp"
	"utils_test.cpp"
	"characters_test.cpp"
)

find_package(GTest REQUIRED)
target_link_libraries(edgo PUBLIC glog pthread boost_filesystem)
target_link_libraries(edgotest PUBLIC ${GTEST_LIBRARIES} gtest_main pthread edgo)
