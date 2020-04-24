package elements

import "fmt"

// Service struct
type Service struct {
	StructName  string
	Constructor string
	Arguments   []interface{}
	Calls       []Call
	Returns     interface{}
	Lifecycle   string
}

type Call struct {
	Method    string
	Arguments []interface{}
}

// NewService create new Service struct
func NewService(serviceMap map[interface{}]interface{}) *Service {
	service := &Service{
		StructName:  getStruct(serviceMap),
		Constructor: getConstructor(serviceMap),
		Arguments:   getArguments(serviceMap),
		Calls:       getCalls(serviceMap),
		Returns:     getReturns(serviceMap),
		Lifecycle:   getLifeCycle(serviceMap),
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

func getCalls(serviceMap map[interface{}]interface{}) []Call {
	val, ok := serviceMap["calls"]
	if ok {
		val, ok := val.([]interface{})
		if !ok {
			panic(fmt.Errorf("calls accepts only array of objects"))
		}
		var arr []Call
		for _, elem := range val {
			switch v := elem.(type) {
			case map[interface{}]interface{}:
				arr = append(arr, Call{v["method"].(string), v["arguments"].([]interface{})})
			default:
				panic(fmt.Errorf("calls accepts only array of objects"))
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

	panic("field 'struct' is require in service")
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

func getLifeCycle(serviceMap map[interface{}]interface{}) string {
	val, ok := serviceMap["lifecycle"]
	if ok {
		val, ok := val.(string)
		if !ok {
			panic(fmt.Errorf("lifecycle accepts only string"))
		}

		return val
	}

	return "perm"
}
