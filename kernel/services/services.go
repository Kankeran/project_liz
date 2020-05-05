package services

import (
	"Liz/generators"
	"Liz/kernel/container"
	"Liz/parsers"
)

// Build building container container
func Build() {

	container.Set("service_generator", func() interface{} {
		service := &generators.Service{}

		return service
	})

	container.Set("reference_parser", func() interface{} {
		service := parsers.NewReference(
			make(map[string]interface{}),
			container.Get("yaml_file_reader").(*parsers.YamlFileReader),
		)

		return service
	})

	container.Set("service_parser", func() interface{} {
		service := &parsers.Service{}

		return service
	})

	container.Set("yaml_file_reader", func() interface{} {
		service := parsers.NewYamlFileReader(
			make(map[string]interface{}),
		)

		return service
	})

}
