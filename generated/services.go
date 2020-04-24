package generated

import (
	"Liz/container"
	"Liz/generators"
	"Liz/parsers"
)

// Build building container container
func Build() {

	container.Set("yaml_file_reader", func() interface{} {
		service := parsers.NewYamlFileReader(
			make(map[string]interface{}),
		)

		return service
	})

	container.Set("service_generator", func() interface{} {
		service := &generators.Service{}

		return service
	})

	container.Set("reference_parser", func() interface{} {
		service := parsers.NewReference(
			make(map[string]interface{}),
		)

		return service
	})

	container.Set("service_parser", func() interface{} {
		service := &parsers.Service{}

		return service
	})

}
