## 0.4.0
- Refactored so now a renderer can be provided by an external package (e.g., package html).

## 0.3.6
- restructured packages and created cmd/vendor folder for all dependencies for ease of compilation.

## 0.3.5
- many changes...
- new command line style eg carpenter run filename -options instead of carpenter filename -options run
- changed most tests to fit the new style.
- reading from stdin is disabled.
- streamlined error messages and error reporting

## 0.3.1
-	Crash due to invalid cell coordinates in a merge or style command
-	Default section separator did not default to +++