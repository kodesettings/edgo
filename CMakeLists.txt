cmake_minimum_required(VERSION 3.12)

project(edgo VERSION 0.0.2 LANGUAGES CXX)

add_library(edgo SHARED "source/config.cpp" "source/config.h" "source/selection.cpp" "source/selection.h")
add_executable(edgotest "source/config_test.cpp" "source/selection_test.cpp")

find_package(GTest REQUIRED)
target_link_libraries(edgotest PUBLIC ${GTEST_LIBRARIES} gtest_main pthread edgo)
