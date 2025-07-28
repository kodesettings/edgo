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

    # with no args it will open current directory
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
    a simple button click. Debugging is also possible via dap protocol for go (dlv) and
    python (debugpy). Although both the testing capability and debugging capability is
    currently experimental, there is no guarantee that they will be included in the
    future releases of the editor.

License

    Edgo text editor's code, theme assets and documentation is covered by MIT license.

