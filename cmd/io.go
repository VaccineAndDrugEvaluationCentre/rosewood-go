package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/drgo/fileutils"
	rosewood "github.com/drgo/rosewood/lib"
)

//TODO: move to fileUtils
type RwInputDescriptor struct {
	FileName           string
	MinFileSize        int
	AllowedFileTypes   string
	ConvertFromVersion string
}

func DefaultRwInputDescriptor(settings *rosewood.Settings) *RwInputDescriptor {
	return &RwInputDescriptor{
		FileName:           "",
		MinFileSize:        settings.SectionsPerTable * len(settings.SectionSeparator),
		AllowedFileTypes:   "text",
		ConvertFromVersion: "",
	}
}

func (iDesc *RwInputDescriptor) SetFileName(fileName string) *RwInputDescriptor {
	iDesc.FileName = fileName
	return iDesc
}

func (iDesc *RwInputDescriptor) SetConvertFromVersion(convertFromVersion string) *RwInputDescriptor {
	iDesc.ConvertFromVersion = convertFromVersion
	return iDesc
}

//TODO: fix errors in tests due to returning "------"
//TODO: move to Rosewood lib
func getValidInputReader(iDesc *RwInputDescriptor) (*os.File, error) {
	//TODO: stop supporting piped-in ?
	if iDesc.FileName == "" || iDesc.FileName == "<stdin>" {
		return os.Stdin, nil
	}
	in, err := os.Open(iDesc.FileName)
	if err != nil {
		return nil, err
	}
	//this check here rather than in the interpreter because we need access to *File to rewind it
	//where as the interpreter uses io.Reader which does not support seek
	//TODO: change pass acceptable type(s) as a parameter
	if err = fileutils.CheckTextStream(in, iDesc.MinFileSize); err != nil {
		return nil, err
	}
	//TODO: clean up extra-check
	if err != nil {
		return nil, fmt.Errorf(ErrOpenInFile, iDesc.FileName, err)
	}
	return in, nil
}

func CheckFileVersion(r io.ReadSeeker) (version string, err error) {
	first3Bytes := make([]byte, 3)
	n, err := r.Read(first3Bytes)
	defer func() { //defer rewinding the file stream
		_, serr := r.Seek(0, 0)
		if err != nil {
			err = serr
		}
	}()
	switch {
	case err != nil && err != io.EOF:
		return "", err
	case n < 3:
		return "", fmt.Errorf("stream is empty or does not contain sufficient data, size=%d", n)
	default:
		//TODO: replace with GetFileVersion
		switch string(first3Bytes) {
		case "---":
			return "v0.1", nil
		case "+++":
			return "v0.2", nil
		default:
			return "unknown", nil
		}
	}
}

//TODO: replace with temp file to prevent creating output file when run fails
func getOutputFile(fileName string, overWrite bool) (*os.File, error) {
	if fileName == "" || fileName == "<stdout>" {
		return os.Stdout, nil
	}
	out, err := fileutils.CreateFile(fileName, overWrite)
	if err != nil {
		return nil, fmt.Errorf(ErrOpenOutFile, fileName, err)
	}
	return out, nil
}

//TODO: clean temp dir if any on start
func start() {
	tmp := filepath.Join(os.TempDir(), "mydir")
	os.RemoveAll(tmp)
	t, err := ioutil.TempFile(tmp, "prefix")
	if err != nil {
		// ...
	}
	defer os.Remove(t.Name())
	defer t.Close()
}

//without removing the dir itself
//use os.RemoveAll("/tmp/") to remove folder and its contents
func RemoveDirContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}
