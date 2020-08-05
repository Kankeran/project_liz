package elements

import (
	"fmt"
)

type flag struct {
	ValType, ValName, ValDescription string
	ValDefault                       interface{}
}

type Command struct {
	Name, ServiceGetter, Method string
	Flags                       []*flag
}

func NewCommand(commandMap map[interface{}]interface{}) *Command {
	return &Command{
		Name:          getName(commandMap),
		ServiceGetter: getServiceGetter(commandMap),
		Method:        getMethod(commandMap),
		Flags:         getFlags(commandMap),
	}
}

func getName(sourceMap map[interface{}]interface{}) string {
	val, ok := sourceMap["name"]
	if !ok {
		panic(fmt.Errorf("name is required"))
	}

	typedVal, ok := val.(string)
	if !ok {
		panic(fmt.Errorf("name must be type string"))
	}

	return typedVal
}

func getFlags(sourceMap map[interface{}]interface{}) []*flag {
	val, ok := sourceMap["flags"]
	if !ok {
		return []*flag{}
	}

	typedVal, ok := val.([]interface{})
	if !ok {
		panic(fmt.Errorf("flags must be an array"))
	}

	var flags []*flag
	for _, val := range typedVal {
		flagData := val.(map[interface{}]interface{})
		flags = append(flags, &flag{
			ValType: getType(flagData),
			ValName: getName(flagData),
			ValDefault: getDefault(flagData),
			ValDescription: getDescription(flagData),
		})
	}

	return flags
}

func getType(sourceMap map[interface{}]interface{}) string {
	val, ok := sourceMap["type"]
	if !ok {
		panic(fmt.Errorf("type is required"))
	}

	typedVal, ok := val.(string)
	if !ok {
		panic(fmt.Errorf("type must be type string"))
	}

	return typedVal
}

func getDefault(sourceMap map[interface{}]interface{}) interface{} {
	val, ok := sourceMap["default"]
	if !ok {
		panic(fmt.Errorf("default is required"))
	}

	return val
}

func getDescription(sourceMap map[interface{}]interface{}) string {
	val, ok := sourceMap["description"]
	if !ok {
		panic(fmt.Errorf("description is required"))
	}

	typedVal, ok := val.(string)
	if !ok {
		panic(fmt.Errorf("description must be type string"))
	}

	return typedVal
}
