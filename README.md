# Package rosewood a set of packages for parsing and rendering Rosewood files.
[![Documentation](https://godoc.org/github.com/yangwenmai/gpt?status.svg)](http://godoc.org/github.com/yangwenmai/gpt)

## Compile
- for dependencies, see link/to/graph TODO

## Design overview
### Interpreter 
- highest-level interface permitting parsing streams of Rosewood tables and rendering the output as html (and potentially other formats).
- see link/to/carpenter for an example of using Interpreter. 
- Renderer is a Go interface for rendering parsed Rosewood tables in any format. See html_render.go for an implementation of this interface for rendering html output.

### Parser
- package responsible for parsing Rosewood files.
- parser.File is the main interface to this package, see link/to/interpreter for an example of using it to parse a Rosewood file.

### Types
- package holding most of the logic pertaining to constructing and rendering Rosewood tables.

### Settings
- packing holding configuration information.

