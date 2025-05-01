# the classic edgo

Yet another console text editor, but with native lsp support.

# an edgo derivative

Check out my new code editor, red, at https://github.com/red-rs/red. It's like Edgo but rewritten in Rust!

### Key bindings and features:

Please check the manual page in your console for all the key bindings.

### Installation:

Install golang on your OS and then clone the repository.

```shell
git clone https://github.com/vipmax/edgo && cd edgo
make 
```

### Usage:
```
edgo [filename]

# with no args it will open current directory
edgo 
```

### Configuration

Please check the manual page in your console for usage of config files.

### Themes

Themes are supported via the configuration file and the theme sources can be found
in internal/highlighter/themes directory for now.

### Lsp

Following lsp features are supported:
- completion
- hover
- signature help
- definition
- references
- rename
- method extraction
- diagnostic

### Languages

Following languages are currently supported:

bash, c/c++, c#, go, html, java, javascript, kotlin, lua, python, rust, scala, toml, typescript and yaml files.

Install your language server for interaction with the edgo text editor.

### Support

If you like the project, please support it: https://www.buymeacoffee.com/vipmax/edgo

### Tests

Edgo provides a seamless testing experience with the ability to execute tests using a simple button click.

Edgo supports testing functionality using Tree Sitter for `go`, `python`, `javascript`, `java`.  
WIP for other langs

### Debug

Added Debug first implementation

Debug is working via dap protocol for `go` (dlv)  and `python` (debugpy)
WIP for other langs

Usage:
- `control + b` - set/delete breakpoint  
- `control + F11` to start debug  
In debug mode editing is not allowed  
- `c` - continue to the next breakpoint  
- `q` - quit debug  
