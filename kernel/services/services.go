package services

import (
	"Liz/builder"
	"Liz/domain"
	"Liz/generators"
	"Liz/kernel/container"
	"Liz/kernel/event"
	"Liz/parsers"
	"Liz/test"
)

// Build building container container
func Build() {

	container.Set("service_generator", func() interface{} {
		service := &generators.Service{}

		return service
	})

	container.Set("listener_generator", func() interface{} {
		service := &generators.Listener{}

		return service
	})

	container.Set("service_file_writer", func() interface{} {
		service := domain.NewFileWriter(
			"kernel/services",
			"services.go",
		)

		return service
	})

	container.Set("container_builder", func() interface{} {
		service := builder.NewContainerBuilder(
			container.Get("yaml_file_reader").(*parsers.YamlFileReader),
			container.Get("reference_parser").(*parsers.Reference),
			container.Get("service_parser").(*parsers.Service),
			container.Get("listener_generator").(*generators.Listener),
			container.Get("service_generator").(*generators.Service),
			container.Get("service_file_writer").(*domain.FileWriter),
		)

		return service
	})

	container.Set("test.example_listener", func() interface{} {
		service := &test.ExampleListener{}

		return service
	})

	container.Set("test.my_listener", func() interface{} {
		service := &test.MyListener{}

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

	event.PrepareDispatcher(map[string]func(d *event.Data){
		"show_info": func(d *event.Data) {
			container.Get("test.example_listener").(*test.ExampleListener).ShowInfo(d)
		},
		"show_info2": func(d *event.Data) {
			container.Get("test.my_listener").(*test.MyListener).Show(d)
		},
	})

}
