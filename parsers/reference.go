package parsers

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

// Reference is tool to include all references to one source
type Reference struct {
	referencesFiles map[string]interface{}
}

// NewReference create new reference parser
func NewReference(maping map[string]interface{}) *Reference {
	return &Reference{maping}
}

// Parse parses all references to one source
func (r *Reference) Parse(source map[interface{}]interface{}, filePath string) (interface{}, error) {
	data, err := r.prepareYaml(source, filePath)
	if err != nil {
		return nil, err
	}
	source = data.(map[interface{}]interface{})
	for serviceName, serviceData := range source {
		source[serviceName], err = r.prepareYaml(serviceData, filePath)
		if err != nil {
			return nil, err
		}
		serviceData, ok := source[serviceName].(map[interface{}]interface{})
		if ok {
			source[serviceName], err = r.Parse(serviceData, filePath)
			if err != nil {
				return nil, err
			}
		}
	}

	return source, nil
}

func (r *Reference) prepareYaml(source interface{}, filePath string) (interface{}, error) {
	var err error
	switch typedValue := source.(type) {
	case map[interface{}]interface{}:
		if ref, ok := typedValue["$ref"]; ok {
			source, err = r.readReferences(ref, source, filePath)
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

func (r *Reference) readReferences(reference interface{}, destination interface{}, filePath string) (interface{}, error) {
	var err error
	switch typedValue := reference.(type) {
	case []interface{}:
		for _, ref := range typedValue {
			if refPath, ok := ref.(string); ok {
				destination, err = r.readReference(refPath, destination, filePath)
				if err != nil {
					return nil, err
				}
				continue
			}
			return nil, fmt.Errorf("Bad type reference '%v', must be string with format: 'file_path#path_to_element'", reference)
		}
	case string:
		destination, err = r.readReference(typedValue, destination, filePath)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("Bad type reference '%v', must be string or array of strings with format: 'file_path#path_to_element'", reference)
	}

	return destination, nil
}

func (r *Reference) readReference(reference string, currentFile interface{}, filePath string) (interface{}, error) {
	externalFilePath, elementName := r.prepareReferencePath(reference)
	var refData = currentFile
	if len(externalFilePath) == 0 {
		externalFilePath = filePath
	}

	refData, err := r.getExternalFileData(externalFilePath)
	if err != nil {
		return nil, err
	}

	if len(elementName) > 0 {
		mapedRefData, ok := refData.(map[interface{}]interface{})
		if !ok {
			return nil, fmt.Errorf("Referenced source is not an object")
		}
		refData, err = r.getData(mapedRefData, elementName)
		if err != nil {
			return nil, err
		}
	}

	switch currentFile.(type) {
	case map[interface{}]interface{}:
		switch refData.(type) {
		case map[interface{}]interface{}:
			r.mergeMaps(refData.(map[interface{}]interface{}), currentFile.(map[interface{}]interface{}))
		default:
			if len(currentFile.(map[interface{}]interface{})) > 1 {
				return nil, fmt.Errorf("Cannot merge object with non object element")
			}
			currentFile = refData
		}
	default:
		currentFile = refData
	}

	return currentFile, nil
}

func (r *Reference) mergeMaps(source map[interface{}]interface{}, destination map[interface{}]interface{}) {
	for key, value := range source {
		destination[key] = value
	}
}

func (r *Reference) getExternalFileData(filePath string) (interface{}, error) {
	refData, ok := r.referencesFiles[filePath]
	if !ok {
		var err error
		refData, err = ReadYamlFile(filePath)
		if err != nil {
			return nil, err
		}
		r.referencesFiles[filePath] = refData
		if _, ok := refData.(map[interface{}]interface{}); ok {
			refData, err = r.Parse(r.referencesFiles[filePath].(map[interface{}]interface{}), filePath)
			if err != nil {
				return nil, err
			}
			r.referencesFiles[filePath] = refData
		}
	}

	return refData, nil
}

func (r *Reference) prepareReferencePath(refPath string) (filePath string, elementName string) {
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

func (r *Reference) getData(source map[interface{}]interface{}, path string) (interface{}, error) {
	var serchedElement interface{} = source
	for _, elementName := range strings.Split(path, "/") {

		source, ok := serchedElement.(map[interface{}]interface{})
		if !ok {
			return nil, fmt.Errorf("element '%s' not found", path)
		}

		element, ok := source[elementName]
		if !ok {
			return nil, fmt.Errorf("element '%s' not found", path)
		}
		serchedElement = element
	}

	return serchedElement, nil
}

// ReadYamlFile open and unmarshal yaml file
func ReadYamlFile(fileName string) (interface{}, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var unmarshaledData interface{}
	err = yaml.Unmarshal([]byte(data), &unmarshaledData)

	return unmarshaledData, err
}
