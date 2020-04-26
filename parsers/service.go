package parsers

import (
	"regexp"
	"strings"
)

type Service struct {
	servicesMap map[interface{}]interface{}
}

func (s *Service) SetOriginalServicesMap(servicesMap map[interface{}]interface{}) {
	s.servicesMap = servicesMap
}

func (s *Service) Parse(servicesMap interface{}) interface{} {
	switch services := servicesMap.(type) {
	case map[interface{}]interface{}:
		servicesMap = s.parseMap(services)
	case []interface{}:
		servicesMap = s.parseInterface(services)
	}

	return servicesMap
}

func (s *Service) parseInterface(services []interface{}) interface{} {
	for key, value := range services {
		switch v := value.(type) {
		case string:
			services[key] = s.replaceConjunction(s.replaceSelfReference(s.replaceService(v)))
		case map[interface{}]interface{}:
			services[key] = s.parseMap(v)
		case []interface{}:
			services[key] = s.parseInterface(v)
		}
	}

	return services
}

func (s *Service) parseMap(services map[interface{}]interface{}) interface{} {
	for key, value := range services {
		switch v := value.(type) {
		case string:
			services[key] = s.replaceConjunction(s.replaceSelfReference(s.replaceService(v)))
		case map[interface{}]interface{}:
			services[key] = s.parseMap(v)
		case []interface{}:
			services[key] = s.parseInterface(v)
		}
	}

	return services
}

func (s *Service) replaceService(arg string) string {
	serviceMatcher := regexp.MustCompile("@\\(\\S*\\)")
	indexes := serviceMatcher.FindAllStringIndex(arg, -1)
	var replaced string
	var lastId = 0
	for _, index := range indexes {
		replaced += arg[lastId:index[0]]
		replaced += strings.Replace(strings.Replace(arg[index[0]:index[1]], "@(", "container.Get(\"", 1), ")", "\")", 1)
		replaced += ".(*" + s.getStructName(strings.TrimRight(strings.TrimLeft(arg[index[0]:index[1]], "@("), ")")) + ")"
		lastId = index[1]
	}
	if len(replaced) > 0 {
		replaced += arg[lastId:]
		return replaced
	}
	serviceMatcher = regexp.MustCompile("@.+")
	indexes = serviceMatcher.FindAllStringIndex(arg, -1)
	for _, index := range indexes {
		replaced = strings.Replace(arg[index[0]:index[1]], "@", "container.Get(\"", 1) + "\")"
		replaced += ".(*" + s.getStructName(strings.TrimLeft(arg[index[0]:index[1]], "@")) + ")"
	}
	if len(replaced) > 0 {
		return replaced
	}

	return arg
}

func (s *Service) replaceConjunction(arg string) string {
	return strings.ReplaceAll(strings.ReplaceAll(arg, "'", "\""), "~", "+")
}

func (s *Service) replaceSelfReference(arg string) string {
	return strings.ReplaceAll(arg, "$this.", "service.")
}

func (s *Service) getStructName(serviceName string) string {
	serviceData, ok := s.servicesMap[serviceName].(map[interface{}]interface{})
	if !ok {
		panic("service " + serviceName + " doesn't exists")
	}
	if name, ok := serviceData["struct"]; ok {
		return name.(string)
	}

	panic("Cannot get struct name from service " + serviceName + " reason: field struct is empty")
}
