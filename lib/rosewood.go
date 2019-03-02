// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package rosewood

import (
	"fmt"
	"io"

	"github.com/drgo/rosewood/lib/parser"
	"github.com/drgo/rosewood/lib/types"
)

//Version of this library
const version = "0.5.1"

//LibVersion version of Rosewood lib
func LibVersion() string {
	return version
}

//ConvertToCurrentVersion utility to convert older versions of Rosewood to current version
func ConvertToCurrentVersion(settings *types.RosewoodSettings, in io.Reader, out io.Writer) error {
	switch settings.ConvertFromVersion {
	case "v0.1":
		return parser.ConvertToCurrentVersion(settings, parser.RWSyntaxVdotzero1, in, out)
	}
	return fmt.Errorf("invalid version number: %s", settings.ConvertFromVersion)
}
