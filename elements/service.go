package elements

import "fmt"

// Service struct
type Service struct {
	StructName  string
	Constructor string
	Arguments   []interface{}
	Calls       []string
	Returns     interface{}
}

// NewService create new Service struct
func NewService(serviceMap map[interface{}]interface{}) *Service {
	service := &Service{
		StructName:  getStruct(serviceMap),
		Constructor: getConstructor(serviceMap),
		Arguments:   getArguments(serviceMap),
		Calls:       getCalls(serviceMap),
		Returns:     getReturns(serviceMap),
	}

	if len(service.Constructor) == 0 && len(service.StructName) == 0 {
		panic(fmt.Errorf("struct or construct is required"))
	}

	return service
}

func getReturns(serviceMap map[interface{}]interface{}) interface{} {
	val, ok := serviceMap["returns"]
	if ok {
		return val
	}

	return "service"
}

func getCalls(serviceMap map[interface{}]interface{}) []string {
	val, ok := serviceMap["calls"]
	if ok {
		val, ok := val.([]interface{})
		if !ok {
			panic(fmt.Errorf("calls accepts only array of strings"))
		}
		var arr []string
		for _, elem := range val {
			switch v := elem.(type) {
			case string:
				arr = append(arr, v)
			default:
				panic(fmt.Errorf("calls accepts only array of strings"))
			}
		}

		return arr
	}

	return nil
}

func getArguments(serviceMap map[interface{}]interface{}) []interface{} {
	val, ok := serviceMap["arguments"]
	if ok {
		val, ok := val.([]interface{})
		if !ok {
			panic(fmt.Errorf("arguments accepts only array"))
		}

		return val
	}

	return nil
}

func getStruct(serviceMap map[interface{}]interface{}) string {
	val, ok := serviceMap["struct"]
	if ok {
		val, ok := val.(string)
		if !ok {
			panic(fmt.Errorf("struct accepts only string"))
		}

		return val
	}

	return ""
}

func getConstructor(serviceMap map[interface{}]interface{}) string {
	val, ok := serviceMap["constructor"]
	if ok {
		val, ok := val.(string)
		if !ok {
			panic(fmt.Errorf("constructor accepts only string"))
		}

		return val
	}

	return ""
}
