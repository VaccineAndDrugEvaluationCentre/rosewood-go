package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	redColor = "\x1b[31m"
)

func crash(err error) {
	fmt.Fprintf(os.Stderr, redColor+"%s\n", err)
	os.Exit(ExitWithError)
}

func crashf(format string, a ...interface{}) {
	err := fmt.Sprintf(format, a...)
	fmt.Fprintf(os.Stderr, redColor+"%s\n", err)
	os.Exit(ExitWithError)
}

func helpMessage(topics []string, versionMessage string) {
	if versionMessage != "" {
		fmt.Fprintf(os.Stderr, versionMessage)
	}
	if len(topics) == 0 {
		fmt.Fprintln(os.Stderr, longUsageMessage)
	}
	for _, topic := range topics {
		switch strings.ToLower(strings.TrimSpace(topic)) {
		case "run":
			fmt.Fprintln(os.Stderr, runUsageMessage)
		default:
			fmt.Fprintln(os.Stderr, longUsageMessage)
		}
	}
}

func usage(topics []string, versionMessage string, exitCode int) {
	helpMessage(topics, versionMessage)
	if exitCode > -1 {
		os.Exit(exitCode)
	}
}

type devNull struct{}

func (devNull) Write(p []byte) (n int, err error) {
	return len(p), nil
}

type Flag struct {
	dest   interface{}
	name   string
	letter string
	value  interface{}
}

func NewCommand(name string, args []Flag) *flag.FlagSet {
	cmd := flag.NewFlagSet(name, flag.ContinueOnError)
	//duplicate args with both name and letter specified, so the command
	//can be invoked by either
	var expanded []Flag
	for _, arg := range args {
		expanded = append(expanded, arg)
		if arg.letter != "" {
			expanded = append(expanded, arg)
			expanded[len(expanded)-1].name = arg.letter
		}
	}
	for _, arg := range expanded {
		switch p := arg.dest.(type) {
		case *string:
			cmd.StringVar(p, arg.name, arg.value.(string), "")
		case *bool:
			cmd.BoolVar(p, arg.name, arg.value.(bool), "")
		default:
			continue
		}
	}
	return cmd
}

// ParseCommandLine parses command line arguments for the appropriate subcommandparses arguments.
// The first command is the default command and can be nil.
func ParseCommandLine(top *flag.FlagSet, subs ...*flag.FlagSet) (*flag.FlagSet, error) {
	return ParseArguments(os.Args[1:], top, subs...)
}

// MustParseCommandLine like ParseCommandLine but exits program on error.
func MustParseCommandLine(top *flag.FlagSet, subs ...*flag.FlagSet) (*flag.FlagSet, error) {
	flg, err := ParseCommandLine(top, subs...)
	if err != nil {
		s := err.Error()
		switch {
		case strings.Contains(s, "flag provided but not defined"):
			s = strings.Replace(s, "provided but not defined", "does not exist", 1)
			return nil, fmt.Errorf(s)
		case strings.Contains(s, "help requested"):
			helpMessage(nil, getVersion())
		default:
		}
	}
	if flg == nil || flg.Name() == "" {
		return nil, fmt.Errorf(ErrWrongCommand)
	}
	return flg, nil
}

//ParseArguments parses arguments (passed as a string array) for the appropriate subcommand
func ParseArguments(args []string, top *flag.FlagSet, subs ...*flag.FlagSet) (*flag.FlagSet, error) {
	if top == nil {
		top = flag.NewFlagSet("", flag.ContinueOnError)
	}
	if err := top.Parse(args); err != nil {
		return nil, err
	}
	args = top.Args()
	if len(args) == 0 || len(subs) == 0 { //nothing left to parse
		return top, nil
	}
	cmdTable := make(map[string]*flag.FlagSet)
	for _, cmd := range subs {
		if cmd != nil {
			cmdTable[cmd.Name()] = cmd
		}
	}
	flagSet, found := cmdTable[args[0]] //retrieve the FlagSet for this subcommand
	if !found {
		return nil, fmt.Errorf("command %v is not found", args[0])
	}
	if len(args) == 1 { //nothing left to parse
		return flagSet, nil
	}
	args = args[1:] //skip over the subcommand name
	//move (positional) arguments to their own array
	posArgs := []string{}
	for len(args[0]) > 1 && args[0][0] != '-' { //loop while the first argument is not a flag
		posArgs = append(posArgs, args[0]) //add it to the positional
		//skip to the next arg if any
		if len(args) == 1 {
			break
		}
		args = args[1:]

	}
	//parse the flags
	if err := flagSet.Parse(args); err != nil {
		return nil, err
	}
	//parse the positional arguments
	if err := flagSet.Parse(posArgs); err != nil {
		return nil, err
	}
	return flagSet, nil
}
