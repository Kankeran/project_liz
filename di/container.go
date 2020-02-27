package di

type container struct {
	services         map[string]interface{}
	servicesCreators map[string]func() interface{}
}

// Container useful for di
var Container *container

// NewContainer creates new di container
func NewContainer(services map[string]func() interface{}) {
	Container = &container{
		services:         make(map[string]interface{}),
		servicesCreators: services,
	}
}

func (c *container) Get(serviceName string) interface{} {
	service, ok := c.services[serviceName]
	if !ok {
		c.services[serviceName] = c.servicesCreators[serviceName]()
		service = c.services[serviceName]
	}

	return service
}
