package domain

import (
	"os"
)

var workingPath = "./"

// SetWorkingPath set start path
func SetWorkingPath(path string) {
	if path[len(path)-1] != '/' {
		path += "/"
	}
	workingPath = path
}

// GetWorkingPath returns start path
func GetWorkingPath() string {
	return workingPath
}

// FileWriter struct used to write data to specified file
type FileWriter struct {
	path, fileName string
}

// NewFileWriter initialize new file writer
func NewFileWriter(path string, fileName string) *FileWriter {
	return &FileWriter{path: path, fileName: fileName}
}

// Write writes data to specified file
func (fw *FileWriter) Write(data []byte) error {
	var (
		err  error
		file *os.File
	)

	if _, err = os.Stat(workingPath + fw.path); os.IsNotExist(err) {
		err = os.MkdirAll(fw.path, os.ModePerm)
		if err != nil {
			return err
		}
	}

	file, err = os.Create(workingPath + fw.path + "/" + fw.fileName)
	if err != nil {
		return err
	}

	_, err = file.Write(data)

	if err != nil {
		return err
	}

	return nil
}
