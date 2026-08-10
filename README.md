> NOTICE: WORK IN PROGRESS, PROJECT IS IN DEVELOPMENT STATE

## About this project

**Edgo** is a library for building modern text editors. It provides the core functionality required by code and text editing applications, including efficient text manipulation, syntax highlighting, Language Server Protocol (LSP) integration, and other editor features.

The library is designed to separate editor logic from the terminal user interface. All editing functionality is exposed through a C API, allowing client applications to implement their own rendering, input handling, and user interaction while relying on **Edgo** for the underlying editor behavior. This architecture makes **Edgo** suitable for building terminal editors that require advanced text editing capabilities.

Current features include:

* Fast and efficient text buffer manipulation
* Undo and redo support
* Cursor and selection management
* Syntax highlighting
* Language Server Protocol (LSP) integration
* Extensible editor functionality through a stable C API

## Project history

**Edgo** originated from an earlier standalone text editor project. During development, the editor core was refactored and generalized into a reusable library, making it possible to share the editing engine across multiple applications. This redesign focuses on modularity, portability, and maintainability, enabling developers to build custom editors without reimplementing common editing functionality.

## Purpose

The long-term goal of **Edgo** is to provide a lightweight, embeddable, and extensible foundation for text editing applications, while remaining independent of any particular rendering backend, windowing toolkit, or platform-specific UI framework.
