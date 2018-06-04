# Carpenter
Version 0.5.0

## Introduction
Rosewood is a domain-specific simple language primarily intended to simplify the automatic generation and formatting of statistical tables. The language provides facilities for representing tables in a simple, human-readable format and for manipulating and formatting the structure and style of those tables. Because Rosewood's source files are plain text files, they can be created by hand or automatically generated using a statistical package (e.g., SAS, R, Stata) or other general-purpose scripting or programming languages. 

Rosewood uses a simple markup, inspired by Markdown, to define a table and its contents making it easier to create and maintain than more complicated markup languages such as HTML or XML. The simple non-intrusive markup makes Rosewood tables easily readable even without rendering the source code or using special viewers. This feature was the overriding design goal of Rosewood as analysts need to frequently inspect the tables for correctness while analyzing the data, and forcing them to use a different tool to view the table contents tends to interrupt their work flow. For this reason, the code needed to format (e.g., merge certain cells) and style the table and its content is stored in a different section of the file, so not to reduce the readability of the tabular data. 

Carpenter is a cross-platform (including Linux, macOS, and Windows on x86, amd64, ARM and PowerPC architectures) tool that can be used to parse and render Rosewood source files into html or docx files for viewing or printing.

## Create your first Rosewood file
- Rosewood files are text files, so they can be created using any text editor, e.g., Notepad (PC) or TextEdit (Mac). 
- Create a new text file in your favourite editor and write or copy and paste the following text into it.
```
+++
My first Rosewood table
+++
Item             | Stars |
Butter chicken   |   5   |
Star-anise candy |   4   |
Wilted lettuce   |   0   |
+++
footnote 1: butter chicken wins!
+++
+++
TODO: update all commands
```
- Save the file in a folder (e.g., c:\Desktop), and name it simple.rw.
- Open a Command window or terminal and change to the folder where you saved the file.
- At the command prompt, type 
```
 ./carpenter run simple.rw
```
OR 
```
 carpenter run simple.rw
```
- Carpenter will create a new file named simple.html. Click on this file to view in your browser. The result should look something like this
TODO: insert pic

## Convert version 0.1 files
Carpenter parses and renders Rosewood tables as HTML

Usage:
carpenter [global options] command [command options]

Command:
check     Parse one or more Rosewood files and print any errors
run       Generate a printout using specified template and CCV xml files
help      Shows a list of commands or help for one command
v1tov2 	  Converts most Rosewood v0.1 files to Rosewood v0.2 files. 	
version   Print executable version

Global options:
--debug, -d  controls what information is printed: -d=0 errors only, 1 warning only [default]
	2 prints information on names of files proceesed etc 3 prints internal debug information
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
-replace, -r      any existing output files will be overwritten	
-convert-old, co  v0.1 files will be converted to current version and proceesed.
-save-converted, sc save a copy of the converted Rosewood file.
-no-inlined-css   link to css file insread of including css in the generated html file. 
`
	checkUsageMessage = `
check : Parse one or more Rosewood files and print any errors

Usage:
carpenter [global options] check rosewoodfilenames [command options]

rosewoodfilenames      file name or path pattern that must at least match one Rosewood file

Options:
-sep, -S		  section separator used (default is +++)	
`

	v1tov2UsageMessage = `
v1tov2 : converts most Rosewood v0.1 files to Rosewood v0.2 files. Only complete files (ones with 6 
	section separators) are currently supported.

Usage:
carpenter [global options] v1tov2 v1rosewoodfilenames [command options]

v1rosewoodfilenames      file name or path pattern that must at least match one v0.1 Rosewood file

Options:
-replace, -r      if specified, any existing output files will be overwritten
`
)


