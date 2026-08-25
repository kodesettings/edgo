cmake_minimum_required(VERSION 3.12)

project(edgo VERSION 0.0.2 LANGUAGES CXX)

add_library(edgo SHARED)

# list of files for library target
include("target.cmake")

find_package(GTest REQUIRED)

add_executable(edgotest)

# list of files for test target
include("edgotest.cmake")

set(TS_LANGUAGE_LIBS tree-sitter)

# tree-sitter language bindings are going to be
# listed in this section. With each additional
# language support we have to add the dependency here
# for the grammar.
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
target_link_libraries(edgo PUBLIC SDL2)
target_link_libraries(edgo PUBLIC ${TS_LANGUAGE_LIBS})
target_link_libraries(edgotest PUBLIC ${GTEST_LIBRARIES})
target_link_libraries(edgotest PUBLIC gtest_main)
target_link_libraries(edgotest PUBLIC pthread)
target_link_libraries(edgotest PUBLIC edgo)
