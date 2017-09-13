package main

const (
	versionMessage = "Carpenter %s (%s)\nCopyRight VDEC 2017\n"
	usageMessage   = `
Carpenter is a tool to parse and render Rosewood tables
	
Usage: 
	carpenter [input Rosewood file] -vho
	
	if an output file is not specified (using -o), the output will be printed to standard output.
	if one or more input files are not specified, input will be read from standard input <stdin>.

Options:
	h           Print help screen
	o, output   Output filename
	css         stylesheet file name
	v, verbose  Output debug information
	`
)
