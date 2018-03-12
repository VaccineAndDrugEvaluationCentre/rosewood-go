# Carpenter
Reference implementation of RoseWood. 

Version 0.3.0 build 08ff577-dirty on Sat Oct  7 21:03:35 CDT 2017

### Install
- Windows: download cmd/carpenter.exe
- MacOS: download cmd/carpenter

### Compile
- MacOS
-- cd to cmd/
-- type make (to run without building) or make build to generate an executable file

### Issues
--  


### Features
- Support for comment lines starting with // 
- Support for parsing merge and style commands
- Support for 2 set commands:
    - set rangeseparator "char" e.g., set rangeseparator '-' to allow merge row 1-1, default ":"
    - set mandatorycol "true" to make col required, default false.

### TODO:
- add support for line continuation \\
- test use of multiple input files
- better html generation tests

### How to use
	carpenter [options] [input Rosewood file ...] 
	
	if an input file is not specified, code will be read from standard input <stdin>.
	if an output file is not specified (using -o), the output will be printed to 
	standard output if the input was read from <stdin>, otherwise a file will be 
	created with the same name as the input file and .html extension.

Options:
	css         Style sheet file name
	h           Print help screen
	c, check	Only check for syntax errors
	o, output   Output file name
	r, replace  Overwrite output file
	v, verbose  Output debug information

Examples:

	`
)
