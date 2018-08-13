// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package main

//TODO: spell check and streamline text
const (
	versionMessage   = "carpenter %s (%s) based on rosewood lib version %s\nCopyRight Salah Mahmud 2018\n"
	longUsageMessage = ` 
Carpenter parses and renders Rosewood tables as HTML or DOCX files.

Usage:
carpenter [global options] command [command options]

Commands:
init	  Generates a default carpenter.mdson file in the current folder 	
check     Parse one or more Rosewood files and print any errors
run       Parse and render one or more Rosewood files into HTML or DOCX files
do		  Run using specified mdson file
v1tov2 	  Converts most Rosewood v0.1 files to Rosewood v0.2 files 	
help      Shows this screen or help for a command if help [command] was specified
version   Print executable version

Global options:
-debug, -d  	   controls what information is printed (e.g., -d=0)
			0 prints errors only
			1 printers errors and warnings only [default]
			2 prints information on names of files proceesed etc 
			3 prints internal debug information
-max-threads, -mt  specifies the maximum numbers of Rosewood files that can be processed concurrently. Default=25.	

Notes:
- typing carpetner without any command line options will run carpenter using carpenter.mdson in the current directory (equivalent to 'carpenter do carpenter.mdson')
 `
	runUsageMessage = `
run: prints the results of applying one or more templates to one or more CCV xml files.

Usage:
carpenter [global options] run rosewoodfilenames [command options]

rosewoodfilenames      file name or path pattern that must at least match one Rosewood file

if an input file is not specified, Rosewood source will be read from standard input <stdin>.
if an output file is not specified (using -o), the output will be printed to 
standard output if the input was read from <stdin>, otherwise an html file will be 
created with the same base name as the input file.
if a docx file is specified (e.g., -o mydoc.docx), all output will be written to a docx file.

Options:
-convert-old, co    v0.1 files will be converted to current version and proceesed.
-keep-temp, k       keep all temporary HTML files generated during processing. Useful for debugging.
-no-inlined-css     link to the CSS file instead of embeding CSS in the generated html file. 
-output, -o         path to output file
-replace, -r        any existing output files will be overwritten	
-save-converted, sc save a copy of the converted Rosewood v0.1 file.
-style, -s          Style sheet file name in CSS format
-work-dir, w        used as the default dir name for Rosewood files specified with incomplete path
-config, cfg        specifies filename for saving temporary configuration mdson file created while generating docx.
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
carpenter [global options] init rosewoodfilenames [command options]

Options:
-config, cfg        specifies filename for saving configuration mdson file
-convert-old, co    v0.1 files will be converted to current version and proceesed
-keep-temp, k       keep all temporary HTML files generated during processing. Useful for debugging
-no-inlined-css     link to the CSS file instead of embeding CSS in the generated html file
-output, -o         path to output file
-replace, -r        any existing output files will be overwritten	
-save-converted, sc save a copy of the converted Rosewood v0.1 file
-style, -s          Style sheet file name in CSS format
-work-dir, w        used as the default dir name for Rosewood files specified with incomplete path
`
)
