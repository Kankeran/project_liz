package generated

import (
	"Liz/generators"
	"Liz/parsers"
)

var Services = map[string]func() interface{}{
	"service_generator": func() interface{} {
		service := &generators.Service{}
		return service
	},
	"reference_parser": func() interface{} {
		service := parsers.NewReference(
			make(map[string]interface{}),
		)
		return service
	},
}
