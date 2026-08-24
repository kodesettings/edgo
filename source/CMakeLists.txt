cmake_minimum_required(VERSION 3.12)

project(edgo VERSION 0.0.2 LANGUAGES CXX)

add_library(edgo SHARED)

# Add your source file to the edgo target. All files
# must be listed here individually to keep track
# of the project structure.
target_sources(edgo PRIVATE config.cpp)
target_sources(edgo PRIVATE config.h)
target_sources(edgo PRIVATE selection.cpp)
target_sources(edgo PRIVATE selection.h)
target_sources(edgo PRIVATE operations.h)
target_sources(edgo PRIVATE lsp_user_actions.cpp)
target_sources(edgo PRIVATE lsp_model.h)
target_sources(edgo PRIVATE lsp_capabilities.h)
target_sources(edgo PRIVATE lsp_client.cpp)
target_sources(edgo PRIVATE lsp_client.h)
target_sources(edgo PRIVATE utils.cpp)
target_sources(edgo PRIVATE utils.h)
target_sources(edgo PRIVATE editor.h)
target_sources(edgo PRIVATE screen.cpp)
target_sources(edgo PRIVATE screen.h)
target_sources(edgo PRIVATE characters.cpp)
target_sources(edgo PRIVATE cursor.cpp)
target_sources(edgo PRIVATE highlighter.cpp)
target_sources(edgo PRIVATE highlighter.h)
target_sources(edgo PRIVATE clipboard.cpp)
target_sources(edgo PRIVATE sdl2_init.h)
target_sources(edgo PRIVATE features_v01.cpp)

# embedded scm files
target_sources(edgo PRIVATE lang_bash.h)
target_sources(edgo PRIVATE lang_c.h)
target_sources(edgo PRIVATE lang_cpp.h)
target_sources(edgo PRIVATE lang_css.h)
target_sources(edgo PRIVATE lang_go.h)
target_sources(edgo PRIVATE lang_html.h)
target_sources(edgo PRIVATE lang_java.h)
target_sources(edgo PRIVATE lang_javascript.h)
target_sources(edgo PRIVATE lang_python.h)
target_sources(edgo PRIVATE lang_rust.h)

target_sources(edgo PRIVATE langs.h)
target_sources(edgo PRIVATE themes.h)
target_sources(edgo PRIVATE test_languages.h)
target_sources(edgo PRIVATE search.cpp)
target_sources(edgo PRIVATE search.h)
target_sources(edgo PRIVATE query.cpp)
target_sources(edgo PRIVATE query.h)
target_sources(edgo PRIVATE io.cpp)
target_sources(edgo PRIVATE keyboard.cpp)

# interface registers
target_sources(edgo PRIVATE intfreg.cpp)
target_sources(edgo PRIVATE lsp.cpp)

find_package(GTest REQUIRED)

add_executable(edgotest)

# Add your source file to the edgotest target. This must come after
# find_package(GTest) which is required by the test suite.
target_sources(edgotest PRIVATE config_test.cpp)
target_sources(edgotest PRIVATE selection_test.cpp)
target_sources(edgotest PRIVATE lsp_client_test.cpp)
target_sources(edgotest PRIVATE lsp_test.h)
target_sources(edgotest PRIVATE utils_test.cpp)
target_sources(edgotest PRIVATE characters_test.cpp)
target_sources(edgotest PRIVATE clipboard_test.cpp)
target_sources(edgotest PRIVATE features_v01_test.cpp)
target_sources(edgotest PRIVATE highlighter_test.cpp)
target_sources(edgotest PRIVATE search_test.cpp)
target_sources(edgotest PRIVATE query_test.cpp)

set(TS_LANGUAGE_LIBS)

# tree-sitter language bindings are going to be
# listed in this section. With each additional
# language support we have to add the dependency here
# for the grammas.
list(APPEND TS_LANGUAGE_LIBS tree-sitter-java)
list(APPEND TS_LANGUAGE_LIBS tree-sitter-bash)
list(APPEND TS_LANGUAGE_LIBS tree-sitter-cpp)
list(APPEND TS_LANGUAGE_LIBS tree-sitter-html)
list(APPEND TS_LANGUAGE_LIBS tree-sitter-css)
list(APPEND TS_LANGUAGE_LIBS tree-sitter-go)
list(APPEND TS_LANGUAGE_LIBS tree-sitter-c)
list(APPEND TS_LANGUAGE_LIBS tree-sitter-javascript)
list(APPEND TS_LANGUAGE_LIBS tree-sitter-rust)
list(APPEND TS_LANGUAGE_LIBS tree-sitter-python)

# target and test builds
target_link_libraries(edgo PUBLIC glog)
target_link_libraries(edgo PUBLIC pthread)
target_link_libraries(edgo PUBLIC boost_filesystem)
target_link_libraries(edgo PUBLIC SDL2 tree-sitter)
target_link_libraries(edgo PUBLIC  ${TS_LANGUAGE_LIBS})
target_link_libraries(edgotest PUBLIC ${GTEST_LIBRARIES})
target_link_libraries(edgotest PUBLIC gtest_main)
target_link_libraries(edgotest PUBLIC pthread)
target_link_libraries(edgotest PUBLIC edgo)
