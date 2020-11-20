package services

import (
	"Liz/builder"
	"Liz/domain"
	"Liz/generators"
	"Liz/kernel/container"
	"Liz/parsers"
	"Liz/test"
)

// Build building container container
func Build() {

	container.Set("listener_generator", func() interface{} {
		service := &generators.Listener{}

		return service
	})

	container.Set("service_parser", func() interface{} {
		service := &parsers.Service{}

		return service
	})

	container.Set("yaml_reader", func() interface{} {
		service := domain.NewYamlReader(
			make(map[string]interface{}),
		)

		return service
	})

	container.Set("service_file_writer", func() interface{} {
		service := domain.NewFileWriter(
			"kernel/services",
			"services.go",
		)

		return service
	})

	container.Set("dispatcher_file_writer", func() interface{} {
		service := domain.NewFileWriter(
			"kernel/event",
			"dispatcher.go",
		)

		return service
	})

	container.Set("autoload_file_writer", func() interface{} {
		service := domain.NewFileWriter(
			"kernel/autoload",
			"autoload.go",
		)

		return service
	})

	container.Set("project_starter_builder", func() interface{} {
		service := builder.NewProjectStarter(
			container.Get("config_yaml_file_writer").(*domain.FileWriter),
			container.Get("service_file_writer").(*domain.FileWriter),
			container.Get("container_file_writer").(*domain.FileWriter),
			container.Get("dispatcher_file_writer").(*domain.FileWriter),
			container.Get("autoload_file_writer").(*domain.FileWriter),
			container.Get("code_formatter").(*domain.CodeFormatter),
		)

		return service
	})

	container.Set("container_file_writer", func() interface{} {
		service := domain.NewFileWriter(
			"kernel/container",
			"container.go",
		)

		return service
	})

	container.Set("container_builder", func() interface{} {
		service := builder.NewContainerBuilder(
			container.Get("yaml_reader").(*domain.YamlReader),
			container.Get("service_parser").(*parsers.Service),
			container.Get("listener_generator").(*generators.Listener),
			container.Get("service_generator").(*generators.Service),
			container.Get("service_file_writer").(*domain.FileWriter),
			container.Get("code_formatter").(*domain.CodeFormatter),
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

	container.Set("config_yaml_file_writer", func() interface{} {
		service := domain.NewFileWriter(
			"config",
			"services.yaml",
		)

		return service
	})

	container.Set("service_generator", func() interface{} {
		service := &generators.Service{}

		return service
	})

	container.Set("code_formatter", func() interface{} {
		service := &domain.CodeFormatter{}

		return service
	})

}
