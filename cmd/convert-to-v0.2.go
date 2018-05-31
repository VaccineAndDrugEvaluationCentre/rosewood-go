package main

import (
	"bufio"
	"fmt"
	"io"

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
	errs := errors.NewErrorList()
	iDesc := DefaultRwInputDescriptor(settings).SetConvertFromVersion("v0.1")
	for _, f := range inFileNames {
		in, err = getValidInputReader(iDesc.SetFileName(f))
		if err != nil {
			errs.Add(err)
			continue
		}
		defer in.Close()
		//TODO: replace ext if not rw
		outputFileName := f + "v2.rw"
		if out, err = getOutputFile(outputFileName, settings.OverWriteOutputFile); err != nil {
			return annotateError(f, err)
		}
		defer out.Close()
		newCode, err := rosewood.ConvertToCurrentVersion(settings, in)
		in.Close()
		if err != nil {
			errs.Add(fmt.Errorf("error running file %s:\n%s", f, errors.ErrorsToError(err)))
		}
		//TODO: move to fileUtils
		//write all modified code to output writer
		w := bufio.NewWriter(out)
		//output header comment
		for _, line := range newCode {
			fmt.Fprintln(w, line)
		}
		if err := w.Flush(); err != nil {
			return err
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
