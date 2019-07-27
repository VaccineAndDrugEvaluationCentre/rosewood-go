// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package parser

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"text/scanner"

	"github.com/drgo/core"
	"github.com/drgo/core/errors"
	"github.com/drgo/rosewood/table"
	"github.com/drgo/rosewood/types"
)

//Position is an alias of scanner.Position
type Position = scanner.Position

var (
	unknownPos = Position{"", -1, -1, -1}
)

// const (
// 	//SectionsPerTable the number of section per table
// 	SectionsPerTable = 4
// 	//RwMinFileSize the size of the smallest possible rosewood file
// 	RwMinFileSize = SectionsPerTable * 3 // at least 3 separator chars per section
// )

//File holds information on currently parsed Rosewood file
type File struct {
	FileName string
	sections []*types.Section //holds raw lines
	parser   *CommandParser
	job      *types.Job
	settings *types.RosewoodSettings
	tables   []*table.Table //holds parsed tables and commands
}

//NewFile returns a Rosewood File
func NewFile(fileName string, job *types.Job) *File {
	return &File{FileName: fileName,
		job:      job,
		parser:   NewCommandParser(job),
		settings: job.RosewoodSettings}
}

//Parse parses an io.ReadSeeker streaming a Rosewood file and returns any found tables
func (f *File) Parse(r io.ReadSeeker) error {
	//TODO: add a test file that starts with empty space or other stuff
	var (
		s       *types.Section
		lineNum int
	)
	if core.IsNil(r) {
		panic("nil io.ReadSeeker passed to file.Parse()")
	}
	f.job.UI.Log("*** file parsing started")
	scanner := bufio.NewScanner(r)
	//check file version
	if !scanner.Scan() {
		if scanner.Err() == nil {
			return NewError(ErrSyntaxError, unknownPos, "file is empty")
		}
		return NewError(ErrSyntaxError, unknownPos, scanner.Err().Error())
	}
	lineNum++ //we found a line
	f.job.UI.Log("first line is" + scanner.Text())
	switch GetFileVersion(strings.TrimSpace(scanner.Text())) {
	case "unknown":
		return NewError(ErrSyntaxError, unknownPos, "file does not start by a valid section separator")
	// case "v0.1":
	// 	f.job.UI.Log("found possible version 0.1 file")
	// 	if !f.settings.ConvertOldVersions {
	// 		return NewError(ErrSyntaxError, unknownPos, "possibly version 0.1 file")
	// 	}
	// 	buf, err := f.convertFromVersionZero1(r)
	// 	if err != nil {
	// 		return NewError(ErrSyntaxError, unknownPos, err.Error())
	// 	}
	// 	scanner = bufio.NewScanner(bufio.NewReader(buf))    //rest the scanner to using the modified buffer
	// 	scanner.Scan()                                      //skip the first section separator
	// 	s = types.NewSection(types.SectionUnknown, lineNum) //found the first section
	case "v0.2":
		s = types.NewSection(types.SectionUnknown, lineNum) //found the first section
	}
	//process the rest of the file
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if f.isSectionSeparatorLine(line) { //start of a new section
			if s != nil { //there is an active section, append it to the sections array
				f.sections = append(f.sections, s)
			}
			s = types.NewSection(types.SectionUnknown, lineNum+1) //create a new section
		} else {
			s.Lines = append(s.Lines, line)
		}
	}
	//check for any scanning errors
	if err := scanner.Err(); err != nil {
		return NewError(ErrSyntaxError, unknownPos, err.Error())
	}
	return f.createTables()
}

func (f *File) isSectionSeparatorLine(line string) bool {
	return strings.HasPrefix(strings.TrimSpace(line), f.settings.SectionSeparator)
}

// func (f *File) convertFromVersionZero1(r io.ReadSeeker) (*bytes.Buffer, error) {
// 	var buf bytes.Buffer //buffer to hold converted code
// 	buf.Grow(100 * 1024)
// 	r.Seek(0, 0) //rewind the stream
// 	if err := ConvertToCurrentVersion(f.settings, RWSyntaxVdotzero1, r, &buf); err != nil {
// 		return nil, NewError(ErrSyntaxError, unknownPos, err.Error())
// 	}
// 	if f.settings.SaveConvertedFile {
// 		v2FileName := files.ConstructFileName(f.FileName, "rw", "", "-autogen")
// 		out, err := files.CreateFile(v2FileName, false /* do not overwrite files*/)
// 		if err != nil {
// 			return nil, fmt.Errorf("unable to create converted v0.2 file [%s]: %s", v2FileName, err)
// 		}
// 		if _, err := out.Write(buf.Bytes()); err != nil {
// 			return nil, err
// 		}
// 		if err := out.Close(); err != nil {
// 			return nil, err
// 		}
// 	}
// 	return &buf, nil
// }

//SectionCount returns the number of sections found in the file
func (f *File) SectionCount() int {
	return len(f.sections)
}

func (f *File) createTables() error {
	if f.SectionCount() == 0 || f.SectionCount()%f.settings.SectionsPerTable != 0 {
		return fmt.Errorf("incorrect number of sections: %d", f.SectionCount())
	}
	var t *table.Table
	var err error
	for i, s := range f.sections {
		ii := i + 1 //i is zero-based, section numbers should be one-based
		s.Kind = types.SectionDescriptor(i%f.settings.SectionsPerTable + 1)
		f.job.UI.Log("**** processing " + s.DebugString())

		switch s.Kind {
		case types.SectionCaption:
			t = table.NewTable(f.job.UI)
			t.Caption = s
		case types.SectionBody:
			if t.Contents, err = table.NewTableContents(s.String()); err != nil {
				return NewError(ErrSyntaxError, unknownPos, fmt.Sprintf("error parsing table in section #%d: %s ", ii, err))
			}
		case types.SectionFootNotes:
			t.Footnotes = s
		case types.SectionControl:
			if t.CmdList, err = f.parser.ParseCommandLines(s); err != nil {
				return err
			}
			f.tables = append(f.tables, t)
		default:
			panic(fmt.Sprintf("invalid switch case [%v] in File.CreateTables()", s.Kind))
		}
	}
	return nil
}

//TableCount returns the number of prased tables in the file
func (f *File) TableCount() int {
	return len(f.tables)
}

//Tables returns an array of pointers to parsed Tables
func (f *File) Tables() []*table.Table {
	return f.tables
}

//Errors returns a list of parsing errors
func (f *File) Errors() *errors.ErrorList {
	return f.parser.Errors()
}

//Err returns a list of parsing errors
func (f *File) Err() error {
	return f.parser.errors.Err()
}
