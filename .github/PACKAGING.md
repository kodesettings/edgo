> NOTE: project is currently under development and is unstable !!!

# Packaging Guide

This project uses [CMake](https://cmake.org/) for building and packaging.

## Prerequisites

Install:

- CMake 3.x or newer
- A supported C++ compiler

## Directory structure

Two files should be packaged for development packages:

- edgo.h
- lsp.h

All the other header files should be ignored during packaging.

## Additional files

You should include license information into the package. Manual page is located in man/ folder as that has to be added as well.
