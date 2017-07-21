package carpenter

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
)

//Lib version
const Version = "0.1.0"

const (
	SectionCapacity  = 100
	SectionsPerTable = 5
	SectionSeparator = "+++"
	ColumnSeparator  = "|"
)

const (
	TableCaption = iota
	TableHeader
	TableBody
	TableFooter
	TableControl
)

type rwSection struct {
	rawLines []string
}

//todo: change section to byte.buffer or []bytes
func newSection() *rwSection {
	return &rwSection{rawLines: make([]string, 0, SectionCapacity)}
}

func (s *rwSection) String() string {
	return strings.Join(s.rawLines, "\n")
}

type RwScript struct {
	fileName      string
	sections      [SectionsPerTable]*rwSection //holds raw lines
	rwCommandList []RwCommand                  //holds parsed commands
}

func NewRwScript() *RwScript {
	var rs RwScript
	for i := 0; i < SectionsPerTable; i++ {
		rs.sections[i] = newSection()
	}
	rs.rwCommandList = make([]RwCommand, 0, SectionCapacity)
	return &rs
}

func (t *RwScript) String() string {
	var b bytes.Buffer
	for i := 0; i < SectionsPerTable; i++ {
		b.WriteString(SectionSeparator + "\n")
		b.WriteString(t.sections[i].String())
		b.WriteString("\n")
	}
	return b.String()
}

//TODO: proper error handling
func (rs *RwScript) Parse(r io.Reader) error {
	sectionCount := -1
	scanner := bufio.NewScanner(r)
	//scanner.Split(bufio.ScanBytes) //parse byte-by-byte
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(strings.TrimSpace(line), SectionSeparator) {
			sectionCount++
			if sectionCount >= SectionsPerTable {
				break
			}
		} else {
			rs.sections[sectionCount].rawLines = append(rs.sections[sectionCount].rawLines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func (rs *RwScript) ParseFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return rs.Parse(file)
}
