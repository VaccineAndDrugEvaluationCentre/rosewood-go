// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package fileutils

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

//GetFileList get list of files from a pattern
func GetFileList(pattern string) (matches []string, err error) {
	matches, err = filepath.Glob(pattern)
	return
}

/*
Linux rules
/ 		=> means start at root of filesystem (absolute reference).
./ 		=> current directory
../ 	=> means go up one directory from the current directory then proceed.
../../ 	=> means go up two directories then proceed.
~/ 		=> start from home directory
None of the above: path relative to current directory
*/

//returns filePath if absolute, otherwise it constructs one from base and filepath.
//no error checking here.
func GetUsableFilePath(filePath string, base string) string {
	if strings.TrimSpace(filePath) == "" {
		return filePath
	}
	// if this is a glob, return it with base if one is specified
	if IsFileGlobPattern(filePath) {
		return filepath.Join(base, filePath)
	}
	// if this is a valid filename on its own, use it regardless of the base
	if IsValidFileName(filePath) {
		return filePath
	}
	// if not, perhaps adding the base (if any) would help
	return filepath.Join(base, filePath)
}

func IsValidFileGlob(pattern string) bool {
	// The only possible returned error from Match is ErrBadPattern, when pattern is malformed.
	if _, err := filepath.Match(pattern, "dummyfilename"); err == nil {
		return true
	}
	return false
}

// IsFileGlobPattern determines if "pattern" has 1 or more characters [', ']', '*' or '//' used for globbing in Go.
func IsFileGlobPattern(pattern string) bool {
	globCharFound := false
	for i := 0; i < len(pattern); i++ {
		switch pattern[i] {
		case '\\':
			if runtime.GOOS != "windows" { //in windows, \\ is not used for escaping in patterns
				globCharFound = true
				break
			}
		case '[', ']', '*':
			globCharFound = true
			break
		}
	}
	return globCharFound
}

//may not always work!!
func IsValidFileName(fileName string) bool {
	// Check if file already exists
	if _, err := os.Stat(fileName); err == nil {
		return true
	}

	// Attempt to create it
	var d []byte
	if err := ioutil.WriteFile(fileName, d, 0644); err == nil {
		os.Remove(fileName) // And delete it
		return true
	}

	return false
}

func NewOutWriter(outFileName string) (wc io.WriteCloser, err error) {
	if outFileName == "" {
		return os.Stdout, nil
	}
	wc, err = os.OpenFile(outFileName, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		wc.Close() // don't care about the error here
		return nil, err
	}
	return wc, nil
}
