package rosewood

import (
	"os"
	"path/filepath"
)

//todo merge with other utils
//FILE utils
func GetWorkingDir() string {
	dir, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Dir(dir)
}

func iif(condition bool, t, f interface{}) interface{} {
	if condition {
		return t
	}
	return f
}

// RemoveIfExists removes a file, returning no error if it does not exist.
func RemoveIfExists(filename string) error {
	err := os.Remove(filename)
	if err != nil && os.IsNotExist(err) {
		err = nil
	}
	return err
}

// func FileCompare(file1, file2 string) (error, bool) {
// 	const chunckSize = 64 * 1024
// 	f1s, err := os.Stat(file1)
// 	if err != nil {
// 		return nil, err
// 	}
// 	f2s, err := os.Stat(file2)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if f1s.Size() != f2s.Size() {
// 		return nil, false
// 	}

// 	f1, err := os.Open(file1)
// 	if err != nil {
// 		return nil, err
// 	}

// 	f2, err := os.Open(file2)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for {
// 		b1 := make([]byte, chunckSize)
// 		_, err1 := f1.Read(b1)

// 		b2 := make([]byte, chunckSize)
// 		_, err2 := f2.Read(b2)

// 		if err1 != nil || err2 != nil {
// 			if err1 == io.EOF && err2 == io.EOF {
// 				return nil, true
// 			} else if err1 == io.EOF && err2 == io.EOF {
// 				return nil, false
// 			} else {
// 				log.Fatal(err1, err2)
// 			}
// 		}

// 		if !bytes.Equal(b1, b2) {
// 			return nil, false
// 		}
// 	}
// }
