package parsers

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type YamlFileReader struct {
	sourceFiles map[string]interface{}
}

func NewYamlFileReader(mapping map[string]interface{}) *YamlFileReader {
	return &YamlFileReader{mapping}
}

func (y *YamlFileReader) Read(fileName string) (interface{}, error) {
	unmarshaledData, ok := y.sourceFiles[fileName]
	if ok {
		return unmarshaledData, nil
	}
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal([]byte(data), &unmarshaledData)
	y.sourceFiles[fileName] = unmarshaledData

	return unmarshaledData, err
}
