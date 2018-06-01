package parser

import (
	"fmt"
	"io"
	"strings"
)

//GetFileVersion returns Rosewood file version based on header info
func GetFileVersion(header string) string {
	switch strings.TrimSpace(header) {
	case "---":
		return "v0.1"
	case "+++":
		return "v0.2"
	default:
		return "unknown"
	}
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
