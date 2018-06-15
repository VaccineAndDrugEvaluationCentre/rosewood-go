// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package main

//TODO: update with new options
//TODO: spell check and streamline text
const (
	versionMessage   = "carpenter %s (%s) based on Rosewood version %s\nCopyRight Salah Mahmud 2017\n"
	longUsageMessage = ` 
Carpenter parses and renders Rosewood tables as HTML

Usage:
carpenter [global options] command [command options]

Commands:
init	  Generates a default carpenter.json file in the current folder 	
check     Parse one or more Rosewood files and print any errors
run       Generate a printout using specified template and CCV xml files
help      Shows a list of commands or help for one command
v1tov2 	  Converts most Rosewood v0.1 files to Rosewood v0.2 files. 	
version   Print executable version

Global options:
-debug, -d  controls what information is printed: -d=0 errors only, 1 warning only [default]
	2 prints information on names of files proceesed etc 3 prints internal debug information
-max-threads, -mt specifies the maximum numbers of files that can be processed concurrently. Default=25.	
 `
	runUsageMessage = `
run: prints the results of applying one or more templates to one or more CCV xml files.

Usage:
carpenter [global options] run rosewoodfilenames [command options]

rosewoodfilenames      file name or path pattern that must at least match one Rosewood file

//TODO: revise
if an input file is not specified, code will be read from standard input <stdin> (currently disabled).
if an output file is not specified (using -o), the output will be printed to 
standard output if the input was read from <stdin>, otherwise a file will be 
created with the same name as the input file and .html extension.

Options:
-style, -s        	Style sheet file name in css format
-output, -o       	path to output file or path pattern (in quotes) with * placeholder
-replace, -r      	any existing output files will be overwritten	
-convert-old, co  	v0.1 files will be converted to current version and proceesed.
-save-converted, sc save a copy of the converted Rosewood file.
-keep-temp, k	    keep all temporary files generated during processing. Useful for debugging.
-no-inlined-css   	link to css file insread of including css in the generated html file. 
`
	checkUsageMessage = `
check : Parse one or more Rosewood files and print any errors

Usage:
carpenter [global options] check rosewoodfilenames [command options]

rosewoodfilenames      file name or path pattern that must at least match one Rosewood file

Options:
 None

Notes:
 -check identifies v0.1 files but does not parse or validate them.
`

	v1tov2UsageMessage = `
v1tov2 : attempts to convert Rosewood v0.1 files to Rosewood v0.2 files. Many not work for all v0.1 files.

Usage:
carpenter [global options] v1tov2 v1rosewoodfilenames [command options]

v1rosewoodfilenames      file name or path pattern that must at least match one v0.1 Rosewood file

Options:
-replace, -r      if specified, any existing output files will be overwritten
`

	initUsageMessage = `
init: creates a new Carpenter configuration file

Usage:
carpenter [global options] init configfilename [command options]

configfilename      file name of configuration file to be created. Default=carpenter.json

Options:
-replace, -r      if specified, any existing output files will be overwritten
`
)
