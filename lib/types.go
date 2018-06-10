package rosewood

import (
	"github.com/drgo/rosewood/lib/parser"
	"github.com/drgo/rosewood/lib/setter"
)

//Settings is an alias for setter settings
type Settings = setter.Settings

//Debug* aliases for setter.Debug*
const (
	DebugSilent  = setter.DebugSilent
	DebugWarning = setter.DebugWarning
	DebugUpdates = setter.DebugUpdates
	DebugAll     = setter.DebugAll
)

//DefaultSettings returns a pointer to an initialized settings object
func DefaultSettings() *Settings {
	return setter.DefaultSettings()
}

// const (
// 	ErrGeneric     = parser.ErrGeneric
// 	ErrSyntaxError = parser.ErrSyntaxError
// 	ErrEmpty       = parser.ErrEmpty
// 	ErrUnknown     = parser.ErrUnknown
// )

type errorManager struct {
}

var errorsManager = errorManager{}

//Errors() provides access to an internal struct that exposes several utility functions for handling
//rosewood errors
func Errors() errorManager {
	return errorsManager
}

//IsParsingError returns true if the error came from parsing rosewood files
func (errors errorManager) IsParsingError(err error) bool {
	_, ok := err.(parser.EmError)
	return ok
}
