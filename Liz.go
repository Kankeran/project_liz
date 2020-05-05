package main

import (
	"os"

	"Liz/elements"
	"Liz/generators"
	"Liz/kernel/container"
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
	services.Build()

	servicesMap, err := container.Get("yaml_file_reader").(*parsers.YamlFileReader).Read("./config/services.yaml")
	check(err)

	servicesMap, err = parseReferences(servicesMap.(map[interface{}]interface{}), "./config/services.yaml")
	check(err)
	servicesMap = servicesMap.(map[interface{}]interface{})["services"]
	servicesMap = parseServices(servicesMap.(map[interface{}]interface{}))

	var generator = container.Get("service_generator").(*generators.Service)
	var code = `package services

	// Build building container container
	func Build() {

		`

	for serviceName, serviceMap := range servicesMap.(map[interface{}]interface{}) {
		code += "container.Set(\"" + serviceName.(string) + "\", " + generator.Generate(elements.NewService(serviceMap.(map[interface{}]interface{}))) + ")\n\n"
	}
	code += "}"

	var output []byte
	// println(code)
	output, err = formatCode(code)
	check(err)

	check(writeToFile(output))

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

func parseServices(source map[interface{}]interface{}) map[interface{}]interface{} {
	serviceParser := container.Get("service_parser").(*parsers.Service)
	serviceParser.SetOriginalServicesMap(source)

	return serviceParser.Parse(source).(map[interface{}]interface{})
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
