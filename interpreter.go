package rosewood

import (
	"bufio"
	"fmt"
	"io"
	"runtime"
	"strings"
)

//VERSION of this library
const VERSION = "0.3.0"

var (
	OSEOL string

	trace tracer
)

func init() {
	OSEOL = "\n"
	if runtime.GOOS == "windows" {
		OSEOL = "\r\n"
	}
	trace = newTrace(off, nil) //default trace is off
}

const (
	sectionCapacity  = 100
	sectionsPerTable = 4
	sectionSeparator = "+++"
	columnSeparator  = "|"
)

type Interpreter struct {
	fileName string
	sections []*section //holds raw lines
	settings *Settings
	tables   []*table
	parser   *CommandParser
}

func NewInterpreter(settings *Settings) *Interpreter {
	ri := &Interpreter{}
	//if no custom settings use default ones
	if settings == nil {
		ri.settings = DefaultSettings()
		if ri.settings == nil {
			panic("Interpreter failed to load settings")
		}
	} else {
		ri.settings = settings
	}
	ri.parser = NewCommandParser(ri.settings)
	return ri
}

// func (ri *Interpreter) String() string {
// 	var b bytes.Buffer
// 	for i := 0; i < sectionsPerTable; i++ {
// 		b.WriteString(sectionSeparator + "\n")
// 		b.WriteString(ri.sections[i].String())
// 		b.WriteString("\n")
// 	}
// 	return b.String()
// }

func (ri *Interpreter) sectionCount() int {
	return len(ri.sections)
}

// func (ri *Interpreter) OK() bool {
// 	return true
// }

func (ri *Interpreter) report(message string, status ReportStatus) {
	if ri.settings.Report != nil {
		ri.settings.Report(message, status)
	}
}

func (ri *Interpreter) createTables() error {
	if ri.sectionCount() == 0 || ri.sectionCount()%sectionsPerTable != 0 {
		return fmt.Errorf("incorrect number of sections %d", ri.sectionCount())
	}
	var t *table
	var err error
	for i, s := range ri.sections {
		ii := i + 1 //i is zero-based, section numbers should be one-based
		s.kind = sectionDescriptor(i%sectionsPerTable + 1)
		switch s.kind {
		case sectionCaption:
			t = newTable()
			t.caption = s
		case sectionBody:
			if t.contents, err = NewTableContents(s.String()); err != nil {
				return fmt.Errorf("error parsing table in section # %d: %s ", ii, err)
			}
		case sectionFootNotes:
			t.footnotes = s
		case sectionControl:
			if t.cmdList, err = ri.parser.ParseCommandLines(s.lines); err != nil {
				return fmt.Errorf("error parsing commands in section # %d: %s ", ii, err)
			}
			ri.tables = append(ri.tables, t)
		default:
			panic(fmt.Sprintf("invalid switch case [%v] in Interpreter.CreateTables()", s.kind))
		}
	}
	if len(ri.tables) == 0 {
		return fmt.Errorf("unknown error in Interpreter.CreateTables()")
	}
	return nil
}

//Parse takes an io.Reader containing RoseWood script and an optional script identifier and return an error
func (ri *Interpreter) Parse(r io.Reader, scriptIdentifer string) error {
	err := ri.parse(r, scriptIdentifer)
	if err != nil {
		ri.report(err.Error(), Info)
	}
	return err
}

func (ri *Interpreter) parse(r io.Reader, scriptIdentifer string) error {
	// helper functions
	isSectionSeparatorLine := func(line string) bool {
		return strings.HasPrefix(strings.TrimSpace(line), sectionSeparator)
	}

	var s *section
	lineNum := 0
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if isSectionSeparatorLine(line) { //start a new section
			if s != nil { //there is an active section, append it to the sections array
				ri.sections = append(ri.sections, s)
			}
			s = newSection(sectionUnknown, lineNum+1) //section if any starts on the next line
		} else {
			//TODO: remove this if
			if s == nil { //if text found before a SectionSeparator (at the start of a script)-> a caption section
				s = newSection(sectionCaption, lineNum)
			}
			s.lines = append(s.lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to parse file %s", err)
	}
	if err := ri.createTables(); err != nil {
		return fmt.Errorf("error parsing tables(s): %s", err)
	}

	return nil
}

func (ri *Interpreter) runTableCommands(table *table) error {
	for _, cmd := range table.cmdList {
		//trace.Printf("inside runTable Commands: %d", i)
		switch cmd.token {
		case kwSet:
			//do nothing parser would have handled that
		case kwMerge:
			// if err := table.Merge(cmd.cellRange); err != nil {
			// 	return fmt.Errorf("merge command %s failed %s", cmd, err)
			// }
		case kwStyle:
		default:
			return fmt.Errorf("cannot run unknown command %s ", cmd)
		}
	}
	return nil
}

func (ri *Interpreter) runTables() error {
	for _, t := range ri.tables {
		//trace.Printf("inside runTables: %d", i)
		// if err := ri.runTableCommands(t); err != nil {
		// 	return fmt.Errorf("failed to run commands for table %s", err)
		// }
		if err := t.run(); err != nil {
			return fmt.Errorf("failed to run commands for table %s", err)
		}

	}
	return nil
}

func (ri *Interpreter) renderTables(w io.Writer, hr *HtmlRenderer) error {
	hr.SetWriter(w)
	hr.SetSettings(ri.settings)
	hr.SetTables(ri.tables)
	hr.StartFile()
	for _, t := range ri.tables {
		hr.StartTable(t)
		for _, row := range t.grid.rows {
			hr.StartRow(row)
			for _, cell := range row.cells {
				hr.OutputCell(cell)
			}
			hr.EndRow(row)
		}
		hr.EndTable(t)
	}
	hr.EndFile()
	return nil
}

//Run takes an io.Reader streaming the contents of one or more Rosewood scripts
//and an io.Writer to output the formatted text.
func (ri *Interpreter) Run(src io.Reader, out io.Writer) error {
	var err error
	//TODO: hook up parsing of error messages
	if err = ri.Parse(src, ""); err != nil {
		return err
	}
	if err = ri.runTables(); err != nil {
		return err
	}
	if err = ri.renderTables(out, NewHtmlRenderer()); err != nil {
		return err
	}
	return nil
}
