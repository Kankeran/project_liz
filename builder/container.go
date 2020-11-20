package builder

import (
    "os"

    "Liz/domain"
    "Liz/elements"
    "Liz/generators"
    "Liz/parsers"
)

// Container builder struct
type Container struct {
    yamlReader        *domain.YamlReader
    serviceParser     *parsers.Service
    listenerGenerator *generators.Listener
    serviceGenerator  *generators.Service
    fileWriter        *domain.FileWriter
    codeFormatter     *domain.CodeFormatter
}

// NewContainerBuilder initialize Container struct
func NewContainerBuilder(
    yamlReader *domain.YamlReader,
    serviceParser *parsers.Service,
    listenerGenerator *generators.Listener,
    serviceGenerator *generators.Service,
    fileWriter *domain.FileWriter,
    codeFormatter *domain.CodeFormatter,
) *Container {
    return &Container{
        yamlReader:        yamlReader,
        serviceParser:     serviceParser,
        listenerGenerator: listenerGenerator,
        serviceGenerator:  serviceGenerator,
        fileWriter:        fileWriter,
        codeFormatter:     codeFormatter,
    }
}

// Build builds service and listeners code
func (c *Container) Build(path string) {
    var err error
    makeDir(path)
    err = os.Chdir(path)
    if err != nil {
        panic(err)
    }

    var servicesMap, listenersMap interface{}
    servicesMap, err = c.yamlReader.ParseFile("./config/services.yaml")
    if err != nil {
        panic(err)
    }

    listenersMap, err = c.yamlReader.ParseFile("./config/listeners.yaml")
    if err != nil {
        panic(err)
    }

    var code = "package services\n // Build building container container\n func Build() {\n\n"

    if servicesMap != nil {
        c.serviceParser.SetOriginalServicesMap(servicesMap.(map[interface{}]interface{}))
        servicesMap = c.serviceParser.Parse(servicesMap.(map[interface{}]interface{}))
        for serviceName, serviceMap := range servicesMap.(map[interface{}]interface{}) {
            code += "container.Set(\"" + serviceName.(string) + "\", " + c.serviceGenerator.Generate(elements.NewService(serviceMap.(map[interface{}]interface{}))) + ")\n\n"
        }
    }

    if listenersMap != nil {
        listenersMap = c.serviceParser.Parse(listenersMap.(map[interface{}]interface{}))
        for eventName, listenerMap := range listenersMap.(map[interface{}]interface{}) {
            var listeners = make([]*elements.Listener, len(listenerMap.([]interface{})))
            for key, listenerData := range listenerMap.([]interface{}) {
                listeners[key] = elements.NewListener(listenerData.(map[interface{}]interface{}))
            }
            code += "event.Add(\"" + eventName.(string) + "\", " + c.listenerGenerator.Generate(listeners) + ")\n\n"
        }
    }

    code += "}"

    var output []byte
    // println(code)
    output, err = c.codeFormatter.Format(code)
    if err != nil {
        panic(err)
    }

    err = c.fileWriter.Write(output)
    if err != nil {
        panic(err)
    }
}

func makeDir(path string) {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        err = os.MkdirAll(path, os.ModePerm)
        if err != nil {
            panic(err)
        }
    }
}
