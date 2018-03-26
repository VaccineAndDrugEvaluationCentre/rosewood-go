// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package main

const (
	versionMessage   = "carpenter %s (%s) based on Rosewood version %s\nCopyRight Salah Mahmud 2017\n"
	longUsageMessage = ` 
Carpenter parses and renders Rosewood tables as HTML

Usage:
carpenter [global options] command [command options]

Command:
check     Parse one or more Rosewood files and print any errors
run       Generate a printout using specified template and CCV xml files
help      Shows a list of commands or help for one command
version   Print executable version

Global options:
--debug, -d     prints information useful for debugging: -d=0 errors only [default], 2 prints everything, 1 is in-between
 `
	runUsageMessage = `
run : prints the results of applying one or more templates to one or more CCV xml files.

Usage:
carpenter [global options] run rosewoodfilenames [command options]

rosewoodfilenames      file name or path pattern that must at least match one Rosewood file

if an input file is not specified, code will be read from standard input <stdin> (currently disabled).
if an output file is not specified (using -o), the output will be printed to 
standard output if the input was read from <stdin>, otherwise a file will be 
created with the same name as the input file and .html extension.

Options:
-sep, -S		  section separator used (default is +++)	
-style, -s        Style sheet file name in css format
-output, -o       path to output file or path pattern (in quotes) with * placeholder
-replace, -r      if specified, any existing output files will be overwritten	
`
	checkUsageMessage = `
check : Parse one or more Rosewood files and print any errors

Usage:
carpenter [global options] check rosewoodfilenames [command options]

rosewoodfilenames      file name or path pattern that must at least match one Rosewood file


Options:
-sep, -S		  section separator used (default is +++)	
`
)
