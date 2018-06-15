package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	rosewood "github.com/drgo/rosewood/lib"
)

//TODO: move to fileUtils; replace with generic file descriptors
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

//TODO: move to Rosewood lib
func getValidInputReader(iDesc *RwInputDescriptor) (*os.File, error) {
	//TODO: stop supporting piped-in ?
	if iDesc.FileName == "" || iDesc.FileName == "<stdin>" {
		return os.Stdin, nil
	}
	return os.Open(iDesc.FileName)
	// if err != nil {
	// 	return nil, err
	// }
	// //this check here rather than in the interpreter because we need access to *File to rewind it
	// //where as the interpreter uses io.Reader which does not support seek
	// //TODO: change pass acceptable type(s) as a parameter; "" means do not check
	// if err = fileutils.CheckTextStream(in, iDesc.MinFileSize); err != nil {
	// 	return nil, err
	// }
	// //TODO: clean up extra-check
	// if err != nil {
	// 	return nil, fmt.Errorf(ErrOpenInFile, iDesc.FileName, err)
	// }
	// return in, nil
}

//TODO: remove overwrite flag; not used
func getOutputWriter(fileName string, overWrite bool) (*os.File, error) {
	if fileName == "<stdout>" {
		return os.Stdout, nil
	}
	out, err := ioutil.TempFile("", "rwtmp") //create in the system default temp folder, a file prefixed with rwtmp
	if err != nil {
		return nil, fmt.Errorf(ErrOpenOutFile, fileName, err)
	}
	return out, nil
}

//DefineOutputBaseDir returns current working directory if PreserveIntermediateFiles is false
// otherwise creates a temp dir in the os default temp dir and return its name
//defer os.RemoveAll(dir)
func GetOutputBaseDir(workDirName string, preserveWorkFiles bool) (baseDir string, err error) {
	if strings.TrimSpace(workDirName) != "" {
		return workDirName, nil
	}
	if preserveWorkFiles { //save in current folder permanently
		baseDir, err = os.Getwd()
		if err != nil {
			err = fmt.Errorf("failed to determine current working directory: %s", err)
		}
		return
	}
	// create temp folder
	baseDir, err = ioutil.TempDir("", "rw-temp101") //create temp dir in the os default temp dir
	if err != nil {
		err = fmt.Errorf("failed to create a temp directory: %s", err)
	}
	return
}

func GetValidFormat(job *Job) (string, error) {
	format := strings.Trim(strings.ToLower(filepath.Ext(job.OutputFile.Name)), ".")
	switch {
	case job.OutputFile.Name == "": //no outputfile specified, assume one html per each inputfile
		format = "html"
	case format == "": //outputfile specified but without an extension, return an error for simplicity
		return "", fmt.Errorf("must specify an extension for output file : %s", job.OutputFile.Name)
	case format == "html": //if an html outputfile is specified, currently >1 input file are not allowed
		if len(job.InputFiles) > 1 {
			return "", fmt.Errorf("merging html files into one html file is not supported")
		}
	case format == "docx": //any number of inputfiles is acceptable
	default:
		return "", fmt.Errorf("unsupported output format: %s", format)
	}
	return format, nil
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
