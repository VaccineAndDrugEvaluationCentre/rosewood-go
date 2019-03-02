// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package main

const (
	ExitWithError = 1
	ExitSuccess   = 0
)

const (
	ErrWrongCommand  = "must specify a valid command. For a list of commands, type %s help"
	ErrOpenInFile    = "error opening input file %s: %s"
	ErrOpenOutFile   = "error opening output file %s: %s"
	ErrMissingInFile = "no input file specified"
	ErrRunningFile   = "error running file [%s]:\n %s"
	ErrRunningBatch  = "one or more errors running batch:\n %s"
)
