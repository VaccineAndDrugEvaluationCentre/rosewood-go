package carpenter

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

//Version of this library
const Version = "0.2.0"

const (
	SectionCapacity  = 100
	SectionsPerTable = 5
	SectionSeparator = "+++"
	ColumnSeparator  = "|"
)

type SectionDescriptor int

const (
	TableUnknown SectionDescriptor = iota
	TableCaption
	TableHeader
	TableBody
	TableFooter
	TableControl
)

type rwSection struct {
	scriptIdentifer string
	offset          int
	kind            SectionDescriptor
	lines           []string
}

//todo: change section to byte.buffer or []bytes
func newSection(scriptIdentifer string, offset int, kind SectionDescriptor) *rwSection {
	return &rwSection{scriptIdentifer: scriptIdentifer, offset: offset, kind: kind}
}

func (s *rwSection) String() string {
	return strings.Join(s.lines, "\n")
}

func (s *rwSection) LineCount() int {
	return len(s.lines)
}

var tablePattern = regexp.MustCompile(`\|\S*|\s\|`) //eg "|text|", "|2131|", "||"

func (s *rwSection) hasTablePattern() bool {
	return s.LineCount() == 0 && tablePattern.MatchString(s.lines[0])
}

type RwInterpreter struct {
	fileName string
	sections []*rwSection //holds raw lines
	//	commandList []RwCommand  //holds parsed commands
}

func NewRwInterpreter() *RwInterpreter {
	return &RwInterpreter{}
}

func (ri *RwInterpreter) String() string {
	var b bytes.Buffer
	for i := 0; i < SectionsPerTable; i++ {
		b.WriteString(SectionSeparator + "\n")
		b.WriteString(ri.sections[i].String())
		b.WriteString("\n")
	}
	return b.String()
}

func (ri *RwInterpreter) SectionCount() int {
	return len(ri.sections)
}

func (ri *RwInterpreter) analyzeSections() error {
	for _, s := range ri.sections {
		if s.kind != TableUnknown { //do not change known kinds
			continue
		}

	}
	return nil
}

//Parse takes an io.Reader containing RoseWood script and an optional script identifier and return an error
func (ri *RwInterpreter) Parse(r io.Reader, scriptIdentifer string) error {
	// helper functions
	isSectionSeparatorLine := func(line string) bool {
		return strings.HasPrefix(strings.TrimSpace(line), SectionSeparator)
	}

	var s *rwSection
	lineNum := 0
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if isSectionSeparatorLine(line) { //start a new section
			if s != nil { //there is an active section, append it to the sections array
				ri.sections = append(ri.sections, s)
			}
			s = newSection(scriptIdentifer, lineNum+1, TableUnknown) //section if any starts on the next line
		} else {
			if s == nil { //if text found before a SectionSeparator (at the start of a script)-> a caption section
				s = newSection(scriptIdentifer, lineNum, TableCaption)
			}
			s.lines = append(s.lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to parse file %s", err)
	}
	return nil
}

//ParseFile takes path to a file containing RoseWood script and parses it possibly returning an error
func (ri *RwInterpreter) ParseFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to parse file %s", err)
	}
	defer file.Close()
	return ri.Parse(file, filename)
}
