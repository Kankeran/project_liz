package domain

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

// YamlReader used to read yaml files and cache them
type YamlReader struct {
	sourceFiles map[string]interface{}
}

// NewYamlReader creates new yaml reader
func NewYamlReader(mapping map[string]interface{}) *YamlReader {
	return &YamlReader{mapping}
}

// Read reads new yaml file or getting file data from cache
func (y *YamlReader) Read(fileName string) interface{} {
	unmarshalledData, ok := y.sourceFiles[fileName]
	if ok {
		return unmarshalledData
	}

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(data, &unmarshalledData)
	y.sourceFiles[fileName] = unmarshalledData

	if err != nil {
		panic(err)
	}

	return unmarshalledData
}

func (y *YamlReader) ParseFile(filePath string) (interface{}, error) {
	sourceMap := y.Read(filePath)

	if sourceMap == nil {
		return sourceMap, nil
	}

	return y.Parse(sourceMap.(map[interface{}]interface{}), filePath)
}

// Parse parses all references to one source
func (y *YamlReader) Parse(source map[interface{}]interface{}, filePath string) (interface{}, error) {
	data, err := y.prepareYaml(source, filePath)
	if err != nil {
		return nil, err
	}
	source = data.(map[interface{}]interface{})
	for serviceName, serviceData := range source {
		source[serviceName], err = y.prepareYaml(serviceData, filePath)
		if err != nil {
			return nil, err
		}
		serviceData, ok := source[serviceName].(map[interface{}]interface{})
		if ok {
			source[serviceName], err = y.Parse(serviceData, filePath)
			if err != nil {
				return nil, err
			}
		}
	}

	return source, nil
}

func (y *YamlReader) prepareYaml(source interface{}, filePath string) (interface{}, error) {
	var err error
	switch typedValue := source.(type) {
	case map[interface{}]interface{}:
		if ref, ok := typedValue["$ref"]; ok {
			source, err = y.readReferences(ref, source, filePath)
			if err != nil {
				return nil, err
			}
		}
	}

	switch typedValue := source.(type) {
	case map[interface{}]interface{}:
		delete(typedValue, "$ref")
	}

	return source, nil
}

func (y *YamlReader) readReferences(reference interface{}, destination interface{}, filePath string) (interface{}, error) {
	var err error
	switch typedValue := reference.(type) {
	case []interface{}:
		for _, ref := range typedValue {
			if refPath, ok := ref.(string); ok {
				destination, err = y.readReference(refPath, destination, filePath)
				if err != nil {
					return nil, err
				}
				continue
			}
			return nil, fmt.Errorf("bad type reference '%v', must be string with format: 'file_path#path_to_element'", reference)
		}
	case string:
		destination, err = y.readReference(typedValue, destination, filePath)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("bad type reference '%v', must be string or array of strings with format: 'file_path#path_to_element'", reference)
	}

	return destination, nil
}

func (y *YamlReader) readReference(reference string, currentFile interface{}, filePath string) (interface{}, error) {
	externalFilePath, elementName := y.prepareReferencePath(reference)
	if len(externalFilePath) == 0 {
		externalFilePath = filePath
	}

	refData, err := y.getExternalFileData(externalFilePath)
	if err != nil {
		return nil, err
	}

	if len(elementName) > 0 {
		mapedRefData, ok := refData.(map[interface{}]interface{})
		if !ok {
			return nil, fmt.Errorf("referenced source is not an object")
		}
		refData, err = y.getData(mapedRefData, elementName)
		if err != nil {
			return nil, err
		}
	}

	switch currentFile.(type) {
	case map[interface{}]interface{}:
		switch refData.(type) {
		case map[interface{}]interface{}:
			y.mergeMaps(refData.(map[interface{}]interface{}), currentFile.(map[interface{}]interface{}))
		default:
			if len(currentFile.(map[interface{}]interface{})) > 1 {
				return nil, fmt.Errorf("cannot merge object with non object element")
			}
			currentFile = refData
		}
	default:
		currentFile = refData
	}

	return currentFile, nil
}

func (y *YamlReader) mergeMaps(source map[interface{}]interface{}, destination map[interface{}]interface{}) {
	for key, value := range source {
		destination[key] = value
	}
}

func (y *YamlReader) getExternalFileData(filePath string) (interface{}, error) {
	refData, ok := y.sourceFiles[filePath]
	if !ok {
		var err error
		refData = y.Read(filePath)
		y.sourceFiles[filePath] = refData
		if _, ok := refData.(map[interface{}]interface{}); ok {
			refData, err = y.Parse(y.sourceFiles[filePath].(map[interface{}]interface{}), filePath)
			if err != nil {
				return nil, err
			}
			y.sourceFiles[filePath] = refData
		}
	}

	return refData, nil
}

func (y *YamlReader) prepareReferencePath(refPath string) (filePath string, elementName string) {
	data := strings.Split(refPath, "#")
	filePath = data[0]
	if len(data) > 1 {
		elementName = strings.Trim(data[1], "/")
	}

	if len(filePath) > 0 {
		if filePath[:2] == "./" {
			filePath = "./config/" + filePath[2:]
		} else if !strings.Contains(filePath, "/") {
			filePath = "./config/" + filePath
		} else if filePath[:3] == "../" {
			filePath = "./" + filePath[3:]
		}
	}

	return filePath, elementName
}

func (y *YamlReader) getData(source map[interface{}]interface{}, path string) (interface{}, error) {
	var searchedElement interface{} = source
	for _, elementName := range strings.Split(path, "/") {

		source, ok := searchedElement.(map[interface{}]interface{})
		if !ok {
			return nil, fmt.Errorf("element '%s' not found", path)
		}

		element, ok := source[elementName]
		if !ok {
			return nil, fmt.Errorf("element '%s' not found", path)
		}
		searchedElement = element
	}

	return searchedElement, nil
}
