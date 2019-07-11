

# parser
`import "."`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
Package parser implements a Rosewood file and command parsers




## <a name="pkg-index">Index</a>
* [Constants](#pkg-constants)
* [Variables](#pkg-variables)
* [func CheckFileVersion(r io.ReadSeeker) (version string, err error)](#CheckFileVersion)
* [func ConvertToCurrentVersion(settings *types.RosewoodSettings, oldSyntax RWSyntax, in io.Reader, out io.Writer) error](#ConvertToCurrentVersion)
* [func ConvertVersion(settings *types.RosewoodSettings, newSyntax RWSyntax, oldSyntax RWSyntax, in io.Reader) ([]string, error)](#ConvertVersion)
* [func GetFileVersion(header string) string](#GetFileVersion)
* [type CommandParser](#CommandParser)
  * [func NewCommandParser(settings *types.RosewoodSettings) *CommandParser](#NewCommandParser)
  * [func (p *CommandParser) ErrorText(index int) string](#CommandParser.ErrorText)
  * [func (p *CommandParser) Errors() *errors.ErrorList](#CommandParser.Errors)
  * [func (p *CommandParser) ParseCommandLines(s *types.Section) ([]*types.Command, error)](#CommandParser.ParseCommandLines)
  * [func (p *CommandParser) Pos() Position](#CommandParser.Pos)
* [type EmError](#EmError)
  * [func NewError(etype int, pos Position, msg string) *EmError](#NewError)
  * [func (e EmError) Error() string](#EmError.Error)
* [type File](#File)
  * [func NewFile(fileName string, settings *types.RosewoodSettings) *File](#NewFile)
  * [func (f *File) Err() error](#File.Err)
  * [func (f *File) Errors() *errors.ErrorList](#File.Errors)
  * [func (f *File) Parse(r io.ReadSeeker) error](#File.Parse)
  * [func (f *File) SectionCount() int](#File.SectionCount)
  * [func (f *File) TableCount() int](#File.TableCount)
  * [func (f *File) Tables() []*types.Table](#File.Tables)
* [type Position](#Position)
* [type RWSyntax](#RWSyntax)


#### <a name="pkg-files">Package files</a>
[cmd_parser.go](/src/target/cmd_parser.go) [convert.go](/src/target/convert.go) [errors.go](/src/target/errors.go) [file_parser.go](/src/target/file_parser.go) [io.go](/src/target/io.go) [parse_table_cmd.go](/src/target/parse_table_cmd.go) [run_set.go](/src/target/run_set.go) 


## <a name="pkg-constants">Constants</a>
``` go
const (
    ErrGeneric int = iota
    ErrSyntaxError
    ErrEmpty
    ErrUnknown
)
```

## <a name="pkg-variables">Variables</a>
``` go
var RWSyntaxVdotzero1 = RWSyntax{
    Version:          "v0.1",
    SectionSeparator: "---",
    ColumnSeparator:  "|",
}
```
RWSyntaxVdotzero1 syntax for Rosewood v0.1

``` go
var RWSyntaxVdotzero2 = RWSyntax{
    Version:          "v0.2",
    SectionSeparator: "+++",
    ColumnSeparator:  "|",
}
```
RWSyntaxVdotzero2 syntax for Rosewood v0.2



## <a name="CheckFileVersion">func</a> [CheckFileVersion](/src/target/io.go?s=293:359#L21)
``` go
func CheckFileVersion(r io.ReadSeeker) (version string, err error)
```


## <a name="ConvertToCurrentVersion">func</a> [ConvertToCurrentVersion](/src/target/convert.go?s=698:815#L36)
``` go
func ConvertToCurrentVersion(settings *types.RosewoodSettings, oldSyntax RWSyntax, in io.Reader, out io.Writer) error
```
ConvertToCurrentVersion utility to convert two current Rosewood version



## <a name="ConvertVersion">func</a> [ConvertVersion](/src/target/convert.go?s=1159:1284#L50)
``` go
func ConvertVersion(settings *types.RosewoodSettings, newSyntax RWSyntax, oldSyntax RWSyntax, in io.Reader) ([]string, error)
```
ConvertVersion utility to convert between two Rosewood versions



## <a name="GetFileVersion">func</a> [GetFileVersion](/src/target/io.go?s=120:161#L10)
``` go
func GetFileVersion(header string) string
```
GetFileVersion returns Rosewood file version based on header info




## <a name="CommandParser">type</a> [CommandParser](/src/target/cmd_parser.go?s=357:621#L18)
``` go
type CommandParser struct {
    // contains filtered or unexported fields
}

```
CommandParser specialized parser for format commands







### <a name="NewCommandParser">func</a> [NewCommandParser](/src/target/cmd_parser.go?s=682:752#L29)
``` go
func NewCommandParser(settings *types.RosewoodSettings) *CommandParser
```
NewCommandParser initializes and returns a CommandParser





### <a name="CommandParser.ErrorText">func</a> (\*CommandParser) [ErrorText](/src/target/cmd_parser.go?s=1280:1331#L48)
``` go
func (p *CommandParser) ErrorText(index int) string
```
ErrorText returns a \n separated list of errors if index =-1, otherwise returns the indexth error;




### <a name="CommandParser.Errors">func</a> (\*CommandParser) [Errors](/src/target/cmd_parser.go?s=1106:1156#L43)
``` go
func (p *CommandParser) Errors() *errors.ErrorList
```
Errors returns a list of parsing errors




### <a name="CommandParser.ParseCommandLines">func</a> (\*CommandParser) [ParseCommandLines](/src/target/cmd_parser.go?s=1683:1768#L62)
``` go
func (p *CommandParser) ParseCommandLines(s *types.Section) ([]*types.Command, error)
```
ParseCommandLines parses a list of strings into list of commands




### <a name="CommandParser.Pos">func</a> (\*CommandParser) [Pos](/src/target/cmd_parser.go?s=1511:1549#L56)
``` go
func (p *CommandParser) Pos() Position
```
Pos returns the current position in the source




## <a name="EmError">type</a> [EmError](/src/target/errors.go?s=237:296#L17)
``` go
type EmError struct {
    Type int
    Position
    Message string
}

```
A EmError is a generic error returned for parsing errors.







### <a name="NewError">func</a> [NewError](/src/target/errors.go?s=929:988#L43)
``` go
func NewError(etype int, pos Position, msg string) *EmError
```
NewError returns a pointer to a new EmError





### <a name="EmError.Error">func</a> (EmError) [Error](/src/target/errors.go?s=340:371#L24)
``` go
func (e EmError) Error() string
```
EmError implements the error interface




## <a name="File">type</a> [File](/src/target/file_parser.go?s=710:911#L34)
``` go
type File struct {
    FileName string
    // contains filtered or unexported fields
}

```
File holds information on currently parsed Rosewood file







### <a name="NewFile">func</a> [NewFile](/src/target/file_parser.go?s=947:1016#L43)
``` go
func NewFile(fileName string, settings *types.RosewoodSettings) *File
```
NewFile returns a Rosewood File





### <a name="File.Err">func</a> (\*File) [Err](/src/target/file_parser.go?s=5891:5917#L195)
``` go
func (f *File) Err() error
```
Err returns a list of parsing errors




### <a name="File.Errors">func</a> (\*File) [Errors](/src/target/file_parser.go?s=5779:5820#L190)
``` go
func (f *File) Errors() *errors.ErrorList
```
Errors returns a list of parsing errors




### <a name="File.Parse">func</a> (\*File) [Parse](/src/target/file_parser.go?s=1205:1248#L50)
``` go
func (f *File) Parse(r io.ReadSeeker) error
```
Parse parses an io.ReadSeeker streaming a Rosewood file and returns any found tables




### <a name="File.SectionCount">func</a> (\*File) [SectionCount](/src/target/file_parser.go?s=4350:4383#L142)
``` go
func (f *File) SectionCount() int
```
SectionCount returns the number of sections found in the file




### <a name="File.TableCount">func</a> (\*File) [TableCount](/src/target/file_parser.go?s=5562:5593#L180)
``` go
func (f *File) TableCount() int
```
TableCount returns the number of prased tables in the file




### <a name="File.Tables">func</a> (\*File) [Tables](/src/target/file_parser.go?s=5676:5714#L185)
``` go
func (f *File) Tables() []*types.Table
```
Tables returns an array of pointers to parsed Tables




## <a name="Position">type</a> [Position](/src/target/file_parser.go?s=321:353#L20)
``` go
type Position = scanner.Position
```
Position is an alias of scanner.Position










## <a name="RWSyntax">type</a> [RWSyntax](/src/target/convert.go?s=205:304#L15)
``` go
type RWSyntax struct {
    Version          string
    SectionSeparator string
    ColumnSeparator  string
}

```
RWSyntax holds info on syntax differences in each Rosewood version














- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
