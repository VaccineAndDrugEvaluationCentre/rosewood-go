// Package setter loads and saves settings for Rosewood and htmldocx libs and commands
// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.
package setter

import (
	"io"

	"github.com/drgo/mdson"
)

// //Settings a generic interface type
// type Settings interface{}

//FIXME: add support to run options and set debug accordingly

//LoadJob loads settings from an io.Reader
func LoadJob(r io.Reader) (*Job, error) {
	job := DefaultJob(DefaultRosewoodSettings())
	mdson.SetDebug(mdson.DebugAll)
	err := mdson.Unmarshal(r, job)
	if err != nil {
		return nil, err
	}
	return job, nil
}

//FIXME: remove once new tracing is implemented
const (
	//DebugSilent print errors only
	DebugSilent int = iota
	//DebugWarning print warnings and errors
	DebugWarning
	//DebugUpdates print execution updates, warnings and errors
	DebugUpdates
	//DebugAll print internal debug messages, execution updates, warnings and errors
	DebugAll
)
