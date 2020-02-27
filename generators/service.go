package generators

import "Liz/elements"

import "fmt"

// Service generator struct
type Service struct{}

// Generate builds code
func (s *Service) Generate(service *elements.Service) string {
	var code = "func() interface{}{\n"
	if len(service.StructName) != 0 {
		code += "service := &" + service.StructName + "{"
		code += s.addArguments(service)
		code += "}\n"
	} else {
		code += "service := " + service.Constructor + "("
		code += s.addArguments(service)
		code += ")\n"
	}

	for _, val := range service.Calls {
		code += "" + val + "\n"
	}
	code += "return " + fmt.Sprint(service.Returns) + "\n},\n"

	return code
}

func (s *Service) addArguments(service *elements.Service) (code string) {
	if len(service.Arguments) > 0 {
		code += "\n"
	}
	for _, val := range service.Arguments {
		code += fmt.Sprint(val) + ",\n"
	}

	return code
}
