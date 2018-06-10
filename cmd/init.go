package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/drgo/rosewood/lib"
	"github.com/drgo/rosewood/lib/setter"
)

func Init(settings *rosewood.Settings, args []string) error {
	settings = rosewood.DefaultSettings()
	var path string
	if len(args) == 0 {
		dir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to create config file in dir %s: %s", dir, err)
		}
		path = filepath.Join(dir, ConfigFileBaseName)
	} else {
		path = args[0]
	}
	return setter.SaveSettings(settings, path, true)
}
