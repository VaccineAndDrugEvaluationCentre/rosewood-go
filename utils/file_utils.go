package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

//FILE utils

func GetWorkingDir() string {
	dir, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Dir(dir)
}

func Iif(condition bool, t, f interface{}) interface{} {
	if condition {
		return t
	}
	return f
}

func ReplaceFileExt(filename, newExt string) string {
	ext := path.Ext(filename)
	newExt = strings.Trim(newExt, ".")
	return filename[0:len(filename)-len(ext)] + "." + newExt
}

// RemoveIfExists removes a file, returning no error if it does not exist.
func RemoveIfExists(filename string) error {
	err := os.Remove(filename)
	if err != nil && os.IsNotExist(err) {
		err = nil
	}
	return err
}

// CreateFile creates the named file with mode 0666 (before umask), truncating
// it if it already exists and overWrite is true. If successful, methods on the returned
// File can be used for I/O; the associated file descriptor has mode O_RDWR.
// If there is an error, it will be of type *PathError.
func CreateFile(fileName string, overWrite bool) (out *os.File, err error) {
	mode := os.O_RDWR | os.O_CREATE | os.O_TRUNC //read-write, create if none exists or truncate existing one
	if !overWrite {
		mode |= os.O_EXCL //file must not exist
	}
	return os.OpenFile(fileName, mode, 0666)
}

//CheckTextStream returns an error if r does contain text otherwise return nil
func CheckTextStream(r io.Reader, streamMinSize int) error {
	first512Bytes := make([]byte, 512)
	n, err := r.Read(first512Bytes)
	//	fmt.Println(n, " ", streamMinSize)
	switch {
	case err == io.EOF || n < streamMinSize:
		return fmt.Errorf("stream is empty or does not contain sufficient data, size=%d", n)
	case err != nil:
		return err
	case !strings.Contains(http.DetectContentType(first512Bytes), "text"):
		return fmt.Errorf("file does not contain text (possibly a binary file)")
	default:
		return nil
	}
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
