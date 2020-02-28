package generated

import (
	"Liz/di"
	"Liz/generators"
	"Liz/parsers"
)

// Build building di container
func Build() {

	di.Container.Set("service_generator", func() interface{} {
		service := &generators.Service{}
		return service
	})

	di.Container.Set("reference_parser", func() interface{} {
		service := parsers.NewReference(
			make(map[string]interface{}),
		)
		return service
	})

	di.Container.Set("yaml_file_reader", func() interface{} {
		service := parsers.NewYamlFileReader(
			make(map[string]interface{}),
		)
		return service
	})

}
