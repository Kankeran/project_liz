package parsers

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// YamlFileReader used to read yaml files and cache them
type YamlFileReader struct {
	sourceFiles map[string]interface{}
}

// NewYamlFileReader creates new yaml reader
func NewYamlFileReader(mapping map[string]interface{}) *YamlFileReader {
	return &YamlFileReader{mapping}
}

// Read reads new yaml file or getting file data from cache
func (y *YamlFileReader) Read(fileName string) (interface{}, error) {
	unmarshalledData, ok := y.sourceFiles[fileName]
	if ok {
		return unmarshalledData, nil
	}
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &unmarshalledData)
	y.sourceFiles[fileName] = unmarshalledData

	return unmarshalledData, err
}
