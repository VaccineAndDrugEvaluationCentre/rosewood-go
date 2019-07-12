package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/drgo/core/ui"
	"github.com/drgo/mdson"
)

//Job holds info about the current run
const configFileBaseName = "htmldocx.json"
const scriptFileBaseName = "script.htds"

// Job holds all info related to current run
//IMMUTABLE once initialized; TODO: hide internal details using getter functions
type Job struct {
	Command               string `mdson:"-"`
	Debug                 int
	RwFileNames           []string
	OverWriteOutputFile   bool
	OutputFileName        string
	OutputFormat          string
	PreserveWorkFiles     bool
	SaveConvertedFile     bool
	StyleSheetName        string
	WorkDirName           string
	RosewoodSettings      *RosewoodSettings
	HTMLWorkDirName       string //dir where HTMLFileNames are stored
	ConfigFileName        string `mdson:"-"` //MDSon file that was used to load the config
	UI                    ui.UI  `mdson:"-"` // provides access to the UI for lower-level routines
	ExecutableVersion     string `mdson:"-"`
	LibVersion            string `mdson:"-"`
	DefaultConfigFileName string `mdson:"-"`
	DefaultScriptFileName string `mdson:"-"`
}

//DefaultJob returns default job
func DefaultJob(settings *RosewoodSettings) *Job {
	return &Job{
		Command:               "run",
		RosewoodSettings:      settings,
		DefaultConfigFileName: configFileBaseName,
		DefaultScriptFileName: scriptFileBaseName,
		UI:                    ui.NewUI(0), //default ui
	}
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
	if index < 0 || len(job.RwFileNames) >= index {
		return ""
	}
	return job.RwFileNames[index]
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

//LoadFromMDSonFile loads job configuration from a file
func (job *Job) LoadFromMDSonFile(FileName string) error {
	configFile, err := os.Open(FileName)
	if err != nil {
		return fmt.Errorf("failed to load configuration file: %v", err)
	}
	defer configFile.Close()
	if err = job.LoadFromMDSon(configFile); err != nil {
		return fmt.Errorf("failed to parse configuration file %s: %v", FileName, err)
	}
	job.UI.Log("configuration loaded from " + FileName)
	return nil
}

//LoadFromMDSon loads job configuration from an io.Reader
func (job *Job) LoadFromMDSon(r io.Reader) error {
	mdson.SetDebug(job.Debug)
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
		return fmt.Errorf("failed to save job configuration: %v", err)
	}
	return nil
}

// GetValidFormat validates and returns a file format extension (without the period)
func (job *Job) GetValidFormat() (string, error) {
	format := strings.Trim(strings.ToLower(filepath.Ext(job.OutputFileName)), ".")
	switch {
	case job.OutputFileName == "": //no outputfile specified, assume one html per each inputfile
		format = "html"
	case format == "": //outputfile specified but without an extension, return an error for simplicity
		return "", fmt.Errorf("must specify an extension for output file : %s", job.OutputFileName)
	case format == "html": //if an html outputfile is specified, currently >1 input file are not allowed
		if len(job.RwFileNames) > 1 {
			return "", fmt.Errorf("merging html files into one html file is not supported")
		}
	case format == "docx": //any number of inputfiles is acceptable
	default:
		return "", fmt.Errorf("unsupported output format: %s", format)
	}
	return format, nil
}

// GetInputDir returns either job.WorkDirName, or the dir where job.ConfigFileName is held
// or current dir in this order
func (job *Job) GetInputDir() (string, error) {
	dir := filepath.Clean(job.WorkDirName)
	if dir != "." {
		return dir, nil
	}
	dir = filepath.Dir(job.ConfigFileName)
	if dir != "." {
		return dir, nil
	}
	return os.Getwd()
}
