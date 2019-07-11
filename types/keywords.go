// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package types

// RwKeyWord identifies Rosewood keywords
type RwKeyWord int

//RwKeyWord identifiers
const (
	KwInvalid RwKeyWord = iota
	catTableCmdBegin
	KwMerge
	KwStyle
	catTableCmdEnd
	KwSet
	KwUse
)

var keywords = map[string]RwKeyWord{
	"merge": KwMerge,
	"style": KwStyle,
	"set":   KwSet,
	"use":   KwUse,
}

//LookupKeyword returns isKeyWord=true and corresponding keyword id if name is keyword;
//false and KwInvalid if not
func LookupKeyword(name string) (RwKeyWord, bool) {
	keyword, isKeyWord := keywords[name]
	return keyword, isKeyWord
}

//IsTableCommand returns true if the token describes a command that manipulates table contents or format
func IsTableCommand(cmd *Command) bool {
	return cmd.token > catTableCmdBegin && cmd.token < catTableCmdEnd
}
