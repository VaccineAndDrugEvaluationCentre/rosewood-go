package rosewood

import (
	"github.com/drgo/rosewood/lib/parser"
	"github.com/drgo/rosewood/lib/types"
)

//this file exposes useful internal types

//Job is an alias for rosewood.Job
type Job = types.Job

//DefaultJob returns a pointer to an initialized job object
func DefaultJob(settings *Settings) *Job {
	return types.DefaultJob(settings)
}

//Settings is an alias for setter settings
type Settings = types.RosewoodSettings

//Debug* aliases for types.Debug*
const (
	DebugSilent  = types.DebugSilent
	DebugWarning = types.DebugWarning
	DebugUpdates = types.DebugUpdates
	DebugAll     = types.DebugAll
)

//DefaultSettings returns a pointer to an initialized settings object
func DefaultSettings() *Settings {
	return types.DefaultRosewoodSettings()
}

//FIXME: replace with package errors
type errorManager struct {
}

var errorsManager = errorManager{}

//Errors provides access to an internal struct that exposes several utility functions for handling
//rosewood errors
func Errors() errorManager {
	return errorsManager
}

//IsParsingError returns true if the error came from parsing rosewood files
func (errors errorManager) IsParsingError(err error) bool {
	_, ok := err.(parser.EmError)
	return ok
}
