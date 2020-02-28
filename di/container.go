package di

type container struct {
	services         map[string]interface{}
	servicesCreators map[string]func() interface{}
}

// Container useful for di
var Container = &container{
	services:         make(map[string]interface{}),
	servicesCreators: make(map[string]func() interface{}),
}

func (c *container) Get(serviceName string) interface{} {
	service, ok := c.services[serviceName]
	if !ok {
		c.services[serviceName] = c.servicesCreators[serviceName]()
		service = c.services[serviceName]
	}

	return service
}

func (c *container) Has(serviceName string) bool {
	_, ok := c.servicesCreators[serviceName]

	return ok
}

func (c *container) Set(serviceName string, serviceCreator func() interface{}) {
	c.servicesCreators[serviceName] = serviceCreator
}
