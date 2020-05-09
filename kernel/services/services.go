package services

import (
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

	container.Set("test.example_listener", func() interface{} {
		service := &test.ExampleListener{}

		return service
	})

	container.Set("test.my_listener", func() interface{} {
		service := &test.MyListener{}

		return service
	})

	event.PrepareDispatcher(map[string]func(d *event.Data){
		"show_info2": func(d *event.Data) {
			container.Get("test.my_listener").(*test.MyListener).Show(d)
		},
		"show_info": func(d *event.Data) {
			container.Get("test.example_listener").(*test.ExampleListener).ShowInfo(d)
		},
	})

}
