# Carpenter
Reference implementation of RoseWood. 

Version 0.2.0

### Features
- Support for both comments both line // and paragraph /*...*/
- Support for parsing merge and style
- Support for 2 set commands:
    - set rangeseparator "char" e.g., set rangeseparator '-' to allow merge row 1-1, default ":"
    - set mandatorycol "true" to make col required, default false.

### How to use
#### Interactive (REPL) mode (for testing RoseWood syntax)
- run Carpenter and type commands at the prompt. Currently only prints back the parsed form of the command.

#### Run mode (not enabled in this version)
- 
### Install
- Windows: download cmd/carpenter.exe
- MacOS: download cmd/carpenter

### Compile
- MacOS
-- cd to cmd/
-- type make (to run without building) or make build

### Issues
-- parsing or executing scripts from files not enabled

### TODO
- add support for printing warning and info (similar to Errors)
- add support for line continuation \\
-  

