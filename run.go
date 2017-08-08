package carpenter

import (
	"fmt"
	"io/ioutil"
	"strconv"
)

func (p *CommandParser) Run(cmdList []*RwCommand) error {
	for _, cmd := range cmdList {
		switch cmd.token {
		case kwSet:
			if err := p.runSetCommand(cmd); err != nil {
				return err
			}
			fmt.Printf("%#v\n", p.settings)
		case kwMerge:
			if err := p.runMergeCommand(cmd); err != nil {
				return err
			}
		default:
			return fmt.Errorf("cannot run unknown command")
		}
	}
	return nil
}

func (p *CommandParser) runMergeCommand(cmd *RwCommand) error {
	if p.runMode == Interactive {
		if len(p.tables) < 1 || p.tables[0] == nil {
			return fmt.Errorf("cannot find an active table")
		}
		t := p.tables[0]
		if err := t.Merge(cmd.cellRange); err != nil {
			return fmt.Errorf("merge command failed %s", err)
		}
	}
	return nil
}

func (p *CommandParser) runSetCommand(cmd *RwCommand) error {
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
	loadTable := func(fileName string) (*Table, error) {
		if data, err := ioutil.ReadFile(fileName); err != nil {
			return nil, fmt.Errorf("failed to load table data %s", err)
		} else {
			return ParseTableData(string(data))
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
			fmt.Println("an error occurred ", err)
			return err
		} else {
			p.tables = append(p.tables, table)
			p.settings.TableFileName = s
			if p.debug {
				fmt.Printf("%v", table)
			}
		}
	case "logfilename":
		if s, err = getArgAsString(1, 1); err != nil {
			return err
		}
		p.settings.LogFileName = s //change to method on CommandParser
	default:
		return fmt.Errorf("unknown option %s", cmd.args[0])
	}
	return nil
}
