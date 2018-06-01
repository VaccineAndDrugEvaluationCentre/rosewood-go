# Carpenter
Reference implementation of RoseWood. 

Version 0.4.0 build 4463e45 on Mon Mar 26 02:29:15 CDT 2018

### Install
- Windows: download cmd/carpenter.exe
- MacOS: download cmd/carpenter
- ensure that you are using the most recent version. Type carpenter version and compare output with version above.

### Compile

- this software requires use of the strings.Builder class, and thus the minimum golang version must at least be v1.10 to compile

- ensure all dependencies in the cmd/vendor can be seen from the GOPATH; see fig dep-graph.png

- clone a copy of the rosewood-go repo in VDEC into the GOPATH directory under "drgo", like so:

```
git clone https://github.com/VaccineAndDrugEvaluationCentre/rosewood-go /path/to/GOPATH/src/github.com/drgo/rosewood
```

- clone the repo where you are currently developing

#### MacOS/linux
- cd to cmd/
- type make build to generate an executable file
- type make buildwin to generate a Windows exe

#### Windows
- cd to cmd/
- type make build to generate an executable file

### Issues
--  


### TODO:
- add support for line continuation \\
- test use of multiple input files
- better html generation and html generation tests

### How to use
```
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

check : Parse one or more Rosewood files and print any errors

Usage:
carpenter [global options] check rosewoodfilenames [command options]

rosewoodfilenames      file name or path pattern that must at least match one Rosewood file


Options:
-sep, -S		  section separator used (default is +++)	
```
