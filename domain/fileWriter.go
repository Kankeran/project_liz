package domain

import (
	"os"
)

type FileWriter struct {
	path, fileName string
}

func NewFileWriter(path string, fileName string) *FileWriter {
	return &FileWriter{path: path, fileName: fileName}
}

func (fw *FileWriter) Write(data []byte) error {
	var (
		err  error
		file *os.File
	)

	if _, err = os.Stat(fw.path); os.IsNotExist(err) {
		err = os.MkdirAll(fw.path, os.ModePerm)
		if err != nil {
			return err
		}
	}

	file, err = os.Create(fw.path + "/" + fw.fileName)
	if err != nil {
		return err
	}

	_, err = file.Write(data)

	if err != nil {
		return err
	}

	return nil
}
