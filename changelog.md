## 0.5.6
- fixed rendering of combined row and col spans: issue #24
- added option (markdown-render or md) to treat markdown as text: #21.
## 0.5.5
- migrated to go 1.12 and modules
- restructured lib packages
- added support for inlined css to permit multiple css styles per cells
- if a stylesheet not specified, now we default to using carpenter.css at the exe folder.
- fixed several bugs

## 0.5.0
- added support for running jobs from a mdson files
- added support for concurrent processing of input files
- added support for processing Rosewood v0.1 source files directly in run commands. 
- added command v1tov2 to convert from v0.1 to v0.2
- updated to be consistent with new specs.
- several bugs fixes

## 0.4.5
- changed algorithm for merging cells to address an unforeseen corner case (needs testing).
- added support for inlining css to permit conversion to docx files.

## 0.4.0
- refactored so now a renderer can be provided by an external package (e.g., package html).
- cleaned up old code.
- added draft documentation.

## 0.3.6
- restructured packages and created cmd/vendor folder for all dependencies for ease of compilation.

## 0.3.5
- many changes...
- new command line style eg carpenter run filename -options instead of carpenter filename -options run
- changed most tests to fit the new style.
- reading from stdin is disabled.
- streamlined error messages and error reporting

## 0.3.1
- fixed crash due to invalid cell coordinates in a merge or style command
- fixed bug Default section separator did not default to +++