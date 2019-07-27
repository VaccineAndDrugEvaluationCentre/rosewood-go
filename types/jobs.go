package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/drgo/core"
	"github.com/drgo/core/ui"
	"github.com/drgo/mdson"
)

//ConfigFileBaseName default config file name
const ConfigFileBaseName = "carpenter.mdson"

// Job holds all info related to current run
type Job struct {
	RunOptions *core.RunOptions
	// SaveConvertedFile bool
	RosewoodSettings *RosewoodSettings
	UI               ui.UI `mdson:"-"` // provides access to the UI for lower-level routines
}

//DefaultJob returns default job
func DefaultJob(settings *RosewoodSettings) *Job {
	return &Job{
		RunOptions: &core.RunOptions{
			Command: "run",
		},
		RosewoodSettings: settings,
		UI:               ui.NewUI(0), //default ui
	}
}

func (job Job) String() string {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(job); err != nil {
		panic(fmt.Sprintf("failed to print job configuration: %v", err))
	}
	return buf.String()
}

// SetUI pass a ui.UI pointer
func (job *Job) SetUI(ux ui.UI) {
	if ux == nil {
		panic("nil pointer to ui.UI passed to job.SetUI()")
	}
	job.UI = ux
}

// SetDebugLevel sets the debug level for this job
func (job *Job) SetDebugLevel(debug int) {
	job.UI.SetLevel(debug)
}

//GetNameOfInputFile returns name of input file corresponding to index
func (job *Job) GetNameOfInputFile(index int) string {
	if index < 0 || len(job.RunOptions.InputFileNames) >= index {
		return ""
	}
	return job.RunOptions.InputFileNames[index]
}

//LoadFromMDSonFile loads job configuration from a file
func (job *Job) LoadFromMDSonFile(FileName string) error {
	file, err := os.Open(FileName)
	if err != nil {
		return fmt.Errorf("failed to load configuration file %s: %v", FileName, err)
	}
	defer file.Close()
	if err = job.LoadFromMDSon(file); err != nil {
		return fmt.Errorf("failed to load configuration file %s: %s", FileName, err)
	}
	job.UI.Log("configuration loaded from " + FileName)
	return nil
}

//LoadFromMDSon loads job configuration from an io.Reader
func (job *Job) LoadFromMDSon(r io.Reader) error {
	fmt.Println("loading config. Debug is", job.RunOptions.Debug)
	mdson.SetDebug(job.RunOptions.Debug)
	err := mdson.Unmarshal(r, job)
	if err != nil {
		return err
	}
	return nil
}

//SaveToMDSon saves job configuration to specified MDSon writer
func (job *Job) SaveToMDSon(w io.Writer) error {
	buf, err := mdson.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to save job configuration: %s", err)
	}
	_, err = w.Write(buf)
	if err != nil {
		return fmt.Errorf("failed to save job configuration: %s", err)
	}
	return nil
}

//SaveToMDSonFile save job configuration to specified MDSon file
func (job *Job) SaveToMDSonFile(FileName string, overwrite bool) error {
	err := mdson.MarshalToFile(job, FileName, overwrite)
	if err != nil {
		return fmt.Errorf("failed to save job configuration: %s", err)
	}
	return nil
}

// GetValidFormat validates and returns a file format extension (without the period)
func (job *Job) GetValidFormat() (string, error) {
	format := strings.Trim(strings.ToLower(filepath.Ext(job.RunOptions.OutputFileName)), ".")
	switch {
	case job.RunOptions.OutputFileName == "": //no outputfile specified, assume one html per each inputfile
		format = "html"
	case format == "": //outputfile specified but without an extension, return an error
		return "", fmt.Errorf("must specify an extension for output file : %s", job.RunOptions.OutputFileName)
	case format == "html": //if an html outputfile is specified, currently >1 input file are not allowed
		if len(job.RunOptions.InputFileNames) > 1 {
			return "", fmt.Errorf("merging generated html files into one html file is not supported")
		}
	case format == "docx": //any number of inputfiles is acceptable
	default:
		return "", fmt.Errorf("unsupported output format: %s", format)
	}
	return format, nil
}

// GetWorkDir returns either job.RunOptions.WorkDirName, or the dir where job.RunOptions.ConfigFileName is held
// or current dir in this order
func (job *Job) GetWorkDir() (string, error) {
	dir := filepath.Clean(job.RunOptions.WorkDirName)
	if dir != "." {
		return dir, nil
	}
	dir = filepath.Dir(job.RunOptions.ConfigFileName)
	if dir != "." {
		return dir, nil
	}
	return os.Getwd()
}
