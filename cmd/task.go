package main

import rosewood "github.com/drgo/rosewood/lib"

type task struct {
	settings            *rosewood.Settings
	OverWriteOutputFile bool
	inputFileName       string
	outputFileName      string
}

func (task *task) getResult(err error) result {
	return result{
		task:  task,
		error: err,
	}
}

type result struct {
	task  *task
	error error
}

func newResult(task *task, err error) *result {
	return &result{
		task:  task,
		error: err,
	}
}

func (result *result) setError(err error) *result {
	result.error = err
	return result
}

func getTask(inputFile, outputFile string, job *rosewood.Job) *task {
	return &task{
		settings:            job.RosewoodSettings,
		OverWriteOutputFile: job.OverWriteOutputFile,
		inputFileName:       inputFile,
		outputFileName:      outputFile,
	}
}

// func GetDefaultTask(job *rosewood.Job) *task {
// 	return &task{
// 		settings: job.RosewoodSettings,
// 	}
// }
