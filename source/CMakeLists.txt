cmake_minimum_required(VERSION 3.12)

project(edgo VERSION 0.0.2 LANGUAGES CXX)

add_library(edgo SHARED "config.cpp" "config.h" "selection.cpp" "selection.h")
add_executable(edgotest "config_test.cpp" "selection_test.cpp")

find_package(GTest REQUIRED)
target_link_libraries(edgotest PUBLIC ${GTEST_LIBRARIES} gtest_main pthread edgo)
