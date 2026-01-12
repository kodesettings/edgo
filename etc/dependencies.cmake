# This file is to list any in-repository dependencies. This is usually not needed
# and hopefully it's use is temporary only as we would like to fetch these from package repos.
include(ExternalProject)

ExternalProject_Add(bash SOURCE_DIR ../../lib/tree-sitter-bash BINARY_DIR ../../lib/tree-sitter-bash/build)
ExternalProject_Add(c SOURCE_DIR ../../lib/tree-sitter-c BINARY_DIR ../../lib/tree-sitter-c/build)
ExternalProject_Add(cpp SOURCE_DIR ../../lib/tree-sitter-cpp BINARY_DIR ../../lib/tree-sitter-cpp/build)
ExternalProject_Add(go SOURCE_DIR ../../lib/tree-sitter-go BINARY_DIR ../../lib/tree-sitter-go/build)
ExternalProject_Add(html SOURCE_DIR ../../lib/tree-sitter-html BINARY_DIR ../../lib/tree-sitter-html/build)
ExternalProject_Add(css SOURCE_DIR ../../lib/tree-sitter-css BINARY_DIR ../../lib/tree-sitter-css/build)
ExternalProject_Add(java SOURCE_DIR ../../lib/tree-sitter-java BINARY_DIR ../../lib/tree-sitter-java/build)
ExternalProject_Add(javascript SOURCE_DIR ../../lib/tree-sitter-javascript BINARY_DIR ../../lib/tree-sitter-javascript/build)
ExternalProject_Add(python SOURCE_DIR ../../lib/tree-sitter-python BINARY_DIR ../../lib/tree-sitter-python/build)
ExternalProject_Add(rust SOURCE_DIR ../../lib/tree-sitter-rust BINARY_DIR ../../lib/tree-sitter-rust/build)
