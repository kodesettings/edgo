~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
                          ======= the classic edgo =======
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

    Yet another console text editor, but with native lsp support.
    An edgo derivative: red text editor, check it out at https://github.com/red-rs/red

Installation

    Install golang on your OS and then clone the repository.

    git clone https://github.com/vipmax/edgo && cd edgo
    make

Usage

    edgo [filename]

    # with no args it will open startup screen
    edgo

Configuration

    Please check the manual page in your console for all the key bindings and for how to
    alter the color schemas or LSP configuration of the editor.

    Themes are supported via the configuration file and the theme sources can be found
    in internal/highlighter/themes directory for now.

The LSP Support

    Following lsp features are supported: completion, hover, signature help, definition,
    references, rename, method extraction and diagnostic information.

Syntax highlighting

    Following languages are currently supported:
    bash, c/c++, go, html, css, java, javascript, typescript, python and rust

    Install your language server for interaction with the edgo text editor.

Extra features

    Edgo provides a seamless testing experience with the ability to execute tests using
    a simple button click. Whenever a unit test case is detected, a marker will appear
    on the right side of the pane to launch the test inside the editor. This is useful
    for novice users who would like to write test cases and validate them right away.

Licensing

    This project is based on the original work licensed under MIT. See license.txt
    for the original license. Modifications and new code are licensed under GNU GPLv2.
    See source/license.txt for details.
