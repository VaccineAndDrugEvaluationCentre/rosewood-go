package main

const (
	versionMessage = "Carpenter %s (%s)\nCopyRight VDEC 2017\n"
	usageMessage   = `
Carpenter is a tool to parse and render Rosewood tables as HTML
	
Usage: 
	carpenter [options] [input Rosewood file ...] 
	
	if an input file is not specified, code will be read from standard input <stdin>.
	if an output file is not specified (using -o), the output will be printed to 
	standard output if the input was read from <stdin>, otherwise a file will be 
	created with the same name as the input file and .html extension.

Options:
	css         Style sheet file name
	h           Print help screen
	c, check	Only check for syntax errors, do not render the table
	o, output   Output file name
	r, replace  Overwrite output file
	v, verbose  Output debug information

Examples:
	
	`
)
