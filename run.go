package carpenter

import (
	"fmt"
)

func (p *CommandParser) Run(cmdList []*RwCommand) error {
	for _, cmd := range cmdList {
		switch cmd.token {
		case kwSet:
			if err := p.runSetCommand(cmd); err != nil {
				return err
			}
			fmt.Printf("%#v\n", p.settings)
		default:
		}
		return nil
	}
	return nil
}

func (p *CommandParser) runSetCommand(cmd *RwCommand) error {
	switch cmd.args[0] { //setting name
	case "rangeseparator":
		s := cmd.args.UnquoteString(1)
		if s == nil || len(*s) < 1 {
			return fmt.Errorf("invalid argument in set command")
		}
		p.settings.RangeOperator = rune((*s)[0])
	case "mandatorycol":
		s := cmd.args.UnquoteString(1)
		if s == nil || len(*s) < 4 {
			return fmt.Errorf("invalid argument in set command")
		}
		p.settings.MandatoryCol = *s == "true"
	default:
		return fmt.Errorf("Unknown option %s", cmd.args[0])
	}
	return nil
}
