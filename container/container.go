package container

type containerStruct struct {
	services         map[string]interface{}
	servicesCreators map[string]func() interface{}
}

var containerInstance = &containerStruct{
	services:         make(map[string]interface{}),
	servicesCreators: make(map[string]func() interface{}),
}

func Get(serviceName string) interface{} {
	service, ok := containerInstance.services[serviceName]
	if !ok {
		containerInstance.services[serviceName] = containerInstance.servicesCreators[serviceName]()
		service = containerInstance.services[serviceName]
	}

	return service
}

func Has(serviceName string) bool {
	_, ok := containerInstance.servicesCreators[serviceName]

	return ok
}

func Set(serviceName string, serviceCreator func() interface{}) {
	containerInstance.servicesCreators[serviceName] = serviceCreator
}
