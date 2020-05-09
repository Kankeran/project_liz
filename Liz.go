package main

import (
	"flag"
	"fmt"
	"os"

	"Liz/elements"
	"Liz/generators"
	"Liz/kernel/container"
	"Liz/kernel/event"
	"Liz/kernel/services"
	"Liz/parsers"

	"github.com/pkg/errors"
	"github.com/sqs/goreturns/returns"
	"golang.org/x/tools/imports"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	newCmd := flag.NewFlagSet("new", flag.ExitOnError)
	newName := newCmd.String("name", "App", "type new project name")

	newCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s [new|build] [flags...]\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nflags for new\n")
		newCmd.PrintDefaults()
	}
	flag.Usage = newCmd.Usage
	flag.Parse()

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(2)
	}

	services.Build()

	switch os.Args[1] {
	case "new":
		check(newCmd.Parse(os.Args[2:]))
		fmt.Println(*newName)
		return
	case "build":
		break
	default:
		flag.Usage()
		os.Exit(2)
	}

	servicesMap, err := container.Get("yaml_file_reader").(*parsers.YamlFileReader).Read("./config/services.yaml")
	check(err)

	servicesMap, err = parseReferences(servicesMap.(map[interface{}]interface{}), "./config/services.yaml")
	check(err)
	listenersMap := servicesMap.(map[interface{}]interface{})["listeners"]
	servicesMap = servicesMap.(map[interface{}]interface{})["services"]

	var code = "package services\n // Build building container container\n func Build() {\n\n"

	if servicesMap != nil {
		servicesMap = parseServices(servicesMap.(map[interface{}]interface{}), servicesMap.(map[interface{}]interface{}))
		var serviceGenerator = container.Get("service_generator").(*generators.Service)
		for serviceName, serviceMap := range servicesMap.(map[interface{}]interface{}) {
			code += "container.Set(\"" + serviceName.(string) + "\", " + serviceGenerator.Generate(elements.NewService(serviceMap.(map[interface{}]interface{}))) + ")\n\n"
		}
	}

	event.DispatchSync("show_info2", nil)

	code += "event.PrepareDispatcher(map[string]func(d *event.Data){"

	if listenersMap != nil {
		code += "\n"
		var listenerGenerator = &generators.Listener{}
		listenersMap = parseServices(listenersMap.(map[interface{}]interface{}), servicesMap.(map[interface{}]interface{}))
		for eventName, listenerMap := range listenersMap.(map[interface{}]interface{}) {
			var listeners = make([]*elements.Listener, len(listenerMap.([]interface{})))
			for key, listenerData := range listenerMap.([]interface{}) {
				listeners[key] = elements.NewListener(listenerData.(map[interface{}]interface{}))
			}
			code += "\"" + eventName.(string) + "\": " + listenerGenerator.Generate(listeners)
		}
	}

	code += "})\n\n}"

	var output []byte
	// println(code)
	output, err = formatCode(code)
	check(err)

	check(writeToFile(output))

	event.DispatchSync("show_info", nil)

	// fill()
}

func formatCode(data string) (output []byte, err error) {
	output, err = returns.Process("./", "", []byte(data), nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	output, err = imports.Process("", output, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return output, nil
}

func parseReferences(source map[interface{}]interface{}, filePath string) (interface{}, error) {
	return container.Get("reference_parser").(*parsers.Reference).Parse(source, filePath)
}

func parseServices(dst, src map[interface{}]interface{}) map[interface{}]interface{} {
	serviceParser := container.Get("service_parser").(*parsers.Service)
	serviceParser.SetOriginalServicesMap(src)

	return serviceParser.Parse(dst).(map[interface{}]interface{})
}

func writeToFile(data []byte) error {
	var (
		err  error
		file *os.File
	)

	if _, err = os.Stat("kernel/services"); os.IsNotExist(err) {
		err = os.MkdirAll("kernel/services", os.ModePerm)
		if err != nil {
			return err
		}
	}

	file, err = os.Create("kernel/services/services.go")
	if err != nil {
		return err
	}

	_, err = file.Write(data)

	if err != nil {
		return err
	}

	return nil
}

// func absPath(filename string) (string, error) {
// 	eval, err := filepath.EvalSymlinks(filename)
// 	if err != nil {
// 		return "", err
// 	}
// 	return filepath.Abs(eval)
// }

// func fill() {
// 	path, err := absPath("./")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	cfg := &packages.Config{
// 		Mode:  packages.LoadAllSyntax,
// 		Tests: true,
// 		Dir:   filepath.Dir(path),
// 		Fset:  token.NewFileSet(),
// 		Env:   os.Environ(),
// 	}

// 	pkgs, err := packages.Load(cfg)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	for _, pkg := range pkgs {
// 		fmt.Printf("%v\n", pkg)
// 	}
// }
