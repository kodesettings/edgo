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
target_sources(edgotest PRIVATE test_languages.h)
target_sources(edgotest PRIVATE highlighter_test.cpp)
target_sources(edgotest PRIVATE search_test.cpp)
target_sources(edgotest PRIVATE query_test.cpp)
