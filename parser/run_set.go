// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package parser

import (
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/drgo/rosewood/table"
	"github.com/drgo/rosewood/types"
)

func (p *CommandParser) runSetCommand(cmd *types.Command) error {
	getArgAsString := func(argIndex int, reqLen int) (string, error) {
		s := cmd.Arg(argIndex)
		if len(s) < reqLen {
			return "", fmt.Errorf("invalid argument in set command")
		}
		return s, nil
	}
	// openFile := func(fileName string) (*os.File, error) {
	// 	return os.Open(fileName)
	// }
	loadTable := func(fileName string) (*table.TableContents, error) {
		data, err := ioutil.ReadFile(fileName)
		if err != nil {
			return nil, fmt.Errorf("failed to load table data %s", err)
		}
		return table.NewTableContents(string(data))
	}

	var s string
	var err error
	switch cmd.Arg(0) { //setting name
	case "rangeseparator":
		if s, err = getArgAsString(1, 1); err != nil {
			return err
		}
		p.job.RosewoodSettings.RangeOperator = rune((s)[0])
	case "mandatorycol":
		if s, err = getArgAsString(1, 4); err != nil {
			return err
		}
		if p.job.RosewoodSettings.MandatoryCol, err = strconv.ParseBool(s); err != nil {
			return err
		}
	case "tablefilename":
		if s, err = getArgAsString(1, 1); err != nil {
			return err
		}
		table, err := loadTable(s)
		if err != nil {
			p.job.UI.Log("an error occurred ", err)
			return err
		}
		p.tables = append(p.tables, table)
		p.job.UI.Logf("%v", table)
	case "logfilename":
		if s, err = getArgAsString(1, 1); err != nil {
			return err
		}
		//		p.job.RosewoodSettings.LogFileName = s //change to method on CommandParser
	default:
		return fmt.Errorf("unknown option %s", cmd.Arg(0))
	}
	return nil
}
