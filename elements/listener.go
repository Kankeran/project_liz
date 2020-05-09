package elements

import (
	"fmt"
)

type Listener struct {
	ServiceGetter, Method string
}

// NewListener create new Service struct
func NewListener(listenerMap map[interface{}]interface{}) *Listener {
	return &Listener{
		ServiceGetter: getServiceGetter(listenerMap),
		Method: getMethod(listenerMap),
	}
}

func getServiceGetter(listenerMap map[interface{}]interface{}) string {
	val, ok := listenerMap["service"]
	if !ok {
		panic(fmt.Errorf("service name is required"))
	}
	typedVal, ok := val.(string)
	if !ok {
		panic(fmt.Errorf("service name must be type string"))
	}

	return typedVal
}

func getMethod(listenerMap map[interface{}]interface{}) string {
	val, ok := listenerMap["method"]
	if !ok {
		panic(fmt.Errorf("method name is required"))
	}
	typedVal, ok := val.(string)
	if !ok {
		panic(fmt.Errorf("method name must be type string"))
	}

	return typedVal
}
