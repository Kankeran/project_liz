package parsers

import (
	"regexp"
	"strings"
)

type Service struct {
}

func (s *Service) Parse(servicesMap map[interface{}]interface{}) map[interface{}]interface{} {
	for _, value := range servicesMap {
		switch v := value.(type) {
		case string:
			value = s.replaceConjunction(s.replaceSelfReference(s.replaceService(v, servicesMap)))
		case map[interface{}]interface{}:
			value = s.Parse(v)
		}
	}

	return servicesMap
}

func (s *Service) replaceService(arg string, servicesMap map[interface{}]interface{}) string {
	serviceMatcher := regexp.MustCompile("@\\(\\S+\\)")
	indexes := serviceMatcher.FindAllStringIndex(arg, -1)
	var replaced string
	var lastId = 0
	for _, index := range indexes {
		replaced += arg[lastId:index[0]]
		replaced += strings.Replace(strings.Replace(arg[index[0]:index[1]], "@(", "container.Get(\"", 1), ")", "\")", 1)
		replaced += s.getStructName(servicesMap, strings.TrimRight(strings.TrimLeft(arg[index[0]:index[1]], "@("), ")"))
		lastId = index[1]
	}
	if len(replaced) > 0 {
		replaced += arg[lastId:]
		return replaced
	}
	serviceMatcher = regexp.MustCompile("@\\S+")
	indexes = serviceMatcher.FindAllStringIndex(arg, -1)
	for _, index := range indexes {
		replaced = strings.Replace(arg[index[0]:index[1]], "@", "container.Get(\"", 1) + "\")"
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

func (s *Service) getStructName(servicesMap map[interface{}]interface{}, serviceName string) string {
	serviceData, ok := servicesMap[serviceName].(map[interface{}]interface{})
	if !ok {
		panic("service " + serviceName + " doesn't exists")
	}
	if name, ok := serviceData["struct"]; ok {
		return name.(string)
	}

	panic("Cannot get struct name from service " + serviceName + " reason: field struct is empty")
}
