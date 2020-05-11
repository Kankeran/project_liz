package domain

import (
	"io/ioutil"
	"os"
)

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
	var err error

	if _, err = os.Stat(fw.path); os.IsNotExist(err) {
		err = os.MkdirAll(fw.path, 0775)
		if err != nil {
			return err
		}
	}

	err = ioutil.WriteFile(fw.path+"/"+fw.fileName, data, 0775)
	if err != nil {
		return err
	}

	return nil
}
