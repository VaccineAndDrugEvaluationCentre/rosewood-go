package rosewood

import (
	"fmt"
	"io/ioutil"
	"strconv"
)

func (p *CommandParser) runSetCommand(cmd *Command) error {
	getArgAsString := func(argIndex int, reqLen int) (string, error) {
		s := cmd.args.Arg(argIndex)
		if len(s) < reqLen {
			return "", fmt.Errorf("invalid argument in set command")
		}
		return s, nil
	}
	// openFile := func(fileName string) (*os.File, error) {
	// 	return os.Open(fileName)
	// }
	loadTable := func(fileName string) (*tableContents, error) {
		if data, err := ioutil.ReadFile(fileName); err != nil {
			return nil, fmt.Errorf("failed to load table data %s", err)
		} else {
			return NewTableContents(string(data))
		}
	}

	var s string
	var err error

	switch cmd.args[0] { //setting name
	case "rangeseparator":
		if s, err = getArgAsString(1, 1); err != nil {
			return err
		}
		p.settings.RangeOperator = rune((s)[0])
	case "mandatorycol":
		if s, err = getArgAsString(1, 4); err != nil {
			return err
		}
		if p.settings.MandatoryCol, err = strconv.ParseBool(s); err != nil {
			return err
		}
	case "tablefilename":
		if s, err = getArgAsString(1, 1); err != nil {
			return err
		}
		if table, err := loadTable(s); err != nil {
			trace.Println("an error occurred ", err)
			return err
		} else {
			p.tables = append(p.tables, table)
			trace.Printf("%v", table)
		}
	case "logfilename":
		if s, err = getArgAsString(1, 1); err != nil {
			return err
		}
		//		p.settings.LogFileName = s //change to method on CommandParser
	default:
		return fmt.Errorf("unknown option %s", cmd.args[0])
	}
	return nil
}
