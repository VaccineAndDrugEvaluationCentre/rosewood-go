package rosewood

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/drgo/rosewood/parser"
	"github.com/drgo/rosewood/types"
	"github.com/drgo/rosewood/utils"
)

//VERSION of this library
const VERSION = "0.3.0"

const (
	sectionCapacity  = 100
	sectionsPerTable = 4
	sectionSeparator = "+++"
	columnSeparator  = "|"
)

type Interpreter struct {
	fileName string
	sections []*types.Section //holds raw lines
	settings *utils.Settings
	tables   []*types.Table
	parser   *parser.CommandParser
}

func NewInterpreter(settings *utils.Settings) *Interpreter {
	ri := &Interpreter{}
	//if no custom settings use default ones
	if settings == nil {
		ri.settings = utils.DefaultSettings()
		if ri.settings == nil {
			panic("Interpreter failed to load settings")
		}
	} else {
		ri.settings = settings
	}
	ri.parser = parser.NewCommandParser(ri.settings)
	return ri
}

//Errors returns a list of rosewood.EmError
func (ri *Interpreter) Errors(err error) (eList []error) {
	//	fmt.Printf("%T\n", err)
	switch cause := err.(type) {
	case *utils.EmError:
		//		fmt.Println(cause.Type)
		switch cause.Type {
		case utils.ErrSyntaxError:
			return ri.parser.Errors()
		default:
			return nil
		}
	default:
		return nil
	}
}

func (ri *Interpreter) sectionCount() int {
	return len(ri.sections)
}

// func (ri *Interpreter) report(message string, status ReportStatus) {
// 	if ri.settings.Report != nil {
// 		ri.settings.Report(message, status)
// 	}
// }

func (ri *Interpreter) createTables() error {
	if ri.sectionCount() == 0 || ri.sectionCount()%sectionsPerTable != 0 {
		return fmt.Errorf("incorrect number of sections %d", ri.sectionCount())
	}
	var t *types.Table
	var err error
	for i, s := range ri.sections {
		ii := i + 1 //i is zero-based, section numbers should be one-based
		kind := types.SectionDescriptor(i%sectionsPerTable + 1)
		switch kind {
		case types.SectionCaption:
			t = types.NewTable()
			t.Caption = s
		case types.SectionBody:
			if t.Contents, err = types.NewTableContents(s.String()); err != nil {
				return fmt.Errorf("error parsing table in section # %d: %s ", ii, err)
			}
		case types.SectionFootNotes:
			t.Footnotes = s
		case types.SectionControl:
			if t.CmdList, err = ri.parser.ParseCommandLines(s); err != nil {
				return err
			}
			ri.tables = append(ri.tables, t)
		default:
			panic(fmt.Sprintf("invalid switch case [%v] in Interpreter.CreateTables()", kind))
		}
	}
	if len(ri.tables) == 0 {
		return utils.NewError(utils.ErrUnknown, ri.parser.Pos(), "unknown error in Interpreter.CreateTables()")
	}
	return nil
}

//Parse takes an io.Reader containing RoseWood script and an optional script identifier and return an error
func (ri *Interpreter) Parse(r io.Reader, scriptIdentifer string) error {
	err := ri.parse(r, scriptIdentifer)
	// if err != nil {
	// 	ri.report(err.Error(), Info)
	// }
	return err
}

func (ri *Interpreter) parse(r io.Reader, scriptIdentifer string) error {
	// helper functions
	isSectionSeparatorLine := func(line string) bool {
		return strings.HasPrefix(strings.TrimSpace(line), sectionSeparator)
	}

	var s *types.Section
	lineNum := 0
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if isSectionSeparatorLine(line) { //start a new section
			if s != nil { //there is an active section, append it to the sections array
				ri.sections = append(ri.sections, s)
			}
			s = types.NewSection(types.SectionUnknown, lineNum+1) //section if any starts on the next line
		} else {
			//TODO: remove this if
			if s == nil { //if text found before a SectionSeparator (at the start of a script)-> a caption section
				s = types.NewSection(types.SectionCaption, lineNum)
			}
			s.Lines = append(s.Lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return utils.NewError(utils.ErrSyntaxError, ri.parser.Pos(), err.Error())
	}
	if err := ri.createTables(); err != nil {
		return utils.NewError(utils.ErrSyntaxError, ri.parser.Pos(), err.Error())
	}
	return nil
}

func (ri *Interpreter) runTables() error {
	for _, t := range ri.tables {
		if err := t.Run(); err != nil {
			return fmt.Errorf("failed to run commands for table %s", err)
		}
	}
	return nil
}

func (ri *Interpreter) renderTables(w io.Writer, hr types.Renderer) error {
	var err error
	hr.SetWriter(w)
	hr.SetSettings(ri.settings)
	hr.SetTables(ri.tables)
	hr.StartFile()
	for _, t := range ri.tables {
		err = t.Render(w, hr)
	}
	hr.EndFile()
	return err
}

//Run takes an io.Reader streaming the contents of one or more Rosewood scripts
//and an io.Writer to output the formatted text.
func (ri *Interpreter) Run(src io.Reader, out io.Writer) error {
	var err error
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

// func (ri *Interpreter) String() string {
// 	var b bytes.Buffer
// 	for i := 0; i < sectionsPerTable; i++ {
// 		b.WriteString(sectionSeparator + "\n")
// 		b.WriteString(ri.sections[i].String())
// 		b.WriteString("\n")
// 	}
// 	return b.String()
// }

// func (ri *Interpreter) OK() bool {
// 	return true
// }

// func (ri *Interpreter) runTableCommands(table *types.Table) error {
// 	for _, cmd := range table.cmdList {
// 		//trace.Printf("inside runTable Commands: %d", i)
// 		switch cmd.token {
// 		case kwSet:
// 			//do nothing parser would have handled that
// 		case kwMerge:
// 			// if err := table.Merge(cmd.cellRange); err != nil {
// 			// 	return fmt.Errorf("merge command %s failed %s", cmd, err)
// 			// }
// 		case kwStyle:
// 		default:
// 			return fmt.Errorf("cannot run unknown command %s ", cmd)
// 		}
// 	}
// 	return nil
// }
