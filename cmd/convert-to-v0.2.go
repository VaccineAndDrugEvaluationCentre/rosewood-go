package main

import (
	"fmt"
	"io"

	"github.com/drgo/fileutils"

	"github.com/drgo/errors"
	rosewood "github.com/drgo/rosewood/lib"
)

func V1toV2(settings *rosewood.Settings, inFileNames []string) error {
	const minFileSize = 0
	var (
		in  io.ReadCloser
		out io.WriteCloser
		err error
	)
	if settings.Debug > 0 {
		fmt.Printf("Processing %d files\n", len(inFileNames))
	}
	errs := errors.NewErrorList() //gather all errors
	//TODO: remove inputDescriptor?
	iDesc := DefaultRwInputDescriptor(settings).SetConvertFromVersion("v0.1")
	for _, f := range inFileNames {
		in, err = getValidInputReader(iDesc.SetFileName(f))
		if err != nil {
			errs.Add(err)
			continue
		}
		defer in.Close()
		outputFileName := fileutils.ConstructFileName(f, "rw", "", "autogen-v2")
		if out, err = getOutputFile(outputFileName, settings.OverWriteOutputFile); err != nil {
			errs.Add(err)
			continue
		}
		defer out.Close()
		err := rosewood.ConvertToCurrentVersion(settings, in, out)
		in.Close()
		if err != nil {
			errs.Add(fmt.Errorf("error running file %s:\n%s", f, errors.ErrorsToError(err)))
		}
	}
	return errs
}

func annotateError(fileName string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("----------\nerror running file [%s]:\n%s", fileName, err)
}
