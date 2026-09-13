# Add your source file to the edgo target. All files
# must be listed here individually to keep track
# of the project structure.
target_sources(edgo PRIVATE config.cpp)
target_sources(edgo PRIVATE selection.cpp)
target_sources(edgo PRIVATE lsp_user_actions.cpp)
target_sources(edgo PRIVATE lsp_client.cpp)
target_sources(edgo PRIVATE utils.cpp)
target_sources(edgo PRIVATE screen.cpp)
target_sources(edgo PRIVATE characters.cpp)
target_sources(edgo PRIVATE cursor.cpp)
target_sources(edgo PRIVATE highlighter.cpp)
target_sources(edgo PRIVATE clipboard.cpp)
target_sources(edgo PRIVATE features_v01.cpp)

target_sources(edgo PRIVATE search.cpp)
target_sources(edgo PRIVATE query.cpp)
target_sources(edgo PRIVATE io.cpp)
target_sources(edgo PRIVATE keyboard.cpp)

# interface registers
target_sources(edgo PRIVATE intfreg.cpp)
target_sources(edgo PRIVATE lsp.cpp)
