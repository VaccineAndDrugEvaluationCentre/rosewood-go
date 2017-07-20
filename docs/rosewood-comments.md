Design
cmd/
    contains the shell command
testing/
    contains testing files
script.go
    defines the high-level parser that splits a rosewood file into its sections
table.go
    defines the table section parser
command.go
    defines the format section parser        
render.go
    defines the html renderer






decisions

- comments (in the format section): needed for both documenting code but also for debugging and testing (commenting out commands): both // and /**/
- row and col values must be >0