package container

type containerStruct struct {
	services         map[string]interface{}
	servicesCreators map[string]func() interface{}
}

var containerInstance = &containerStruct{
	services:         make(map[string]interface{}),
	servicesCreators: make(map[string]func() interface{}),
}

// Get getting searched service instance
func Get(serviceName string) interface{} {
	service, ok := containerInstance.services[serviceName]
	if !ok {
		containerInstance.services[serviceName] = containerInstance.servicesCreators[serviceName]()
		service = containerInstance.services[serviceName]
	}

	return service
}

// Has check service exists
func Has(serviceName string) bool {
	_, ok := containerInstance.servicesCreators[serviceName]

	return ok
}

// Set sets function to invoking service with specified name
func Set(serviceName string, serviceCreator func() interface{}) {
	containerInstance.servicesCreators[serviceName] = serviceCreator
}
