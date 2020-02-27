package main

import (
	"Liz/di"
	"Liz/elements"
	"Liz/generated"
	"Liz/generators"
	"Liz/parsers"
	"fmt"
	"go/token"
	"log"
	"os"
	"path/filepath"

	"github.com/sqs/goreturns/returns"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Build() {
	di.NewContainer(generated.Services)
}

func main() {
	Build()
	servicesMap, err := parsers.ReadYamlFile("./config/services.yaml")
	check(err)

	servicesMap, err = parseReferences(servicesMap.(map[interface{}]interface{}), "./config/services.yaml")
	check(err)

	var generator = di.Container.Get("service_generator").(*generators.Service)
	var code = "package generated\n\nvar Services = map[string]func() interface{} {\n"

	for serviceName, serviceMap := range servicesMap.(map[interface{}]interface{}) {
		code += "\"" + serviceName.(string) + "\": " + generator.Generate(elements.NewService(serviceMap.(map[interface{}]interface{})))
	}
	code += "}"

	var output []byte
	output, err = formatCode(code)
	check(err)

	check(writeToFile(output))

	fill()
}

func formatCode(data string) (output []byte, err error) {
	output, err = returns.Process("./", "", []byte(data), nil)
	if err != nil {
		return nil, err
	}

	output, err = imports.Process("", output, nil)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func parseReferences(source map[interface{}]interface{}, filePath string) (interface{}, error) {
	return di.Container.Get("reference_parser").(*parsers.Reference).Parse(source, filePath)
}

func writeToFile(data []byte) error {
	var (
		err  error
		file *os.File
	)

	if _, err = os.Stat("generated"); os.IsNotExist(err) {
		err = os.MkdirAll("generated", os.ModePerm)
		if err != nil {
			return err
		}
	}

	file, err = os.Create("generated/services.go")
	if err != nil {
		return err
	}

	file.Write(data)

	return nil
}

func absPath(filename string) (string, error) {
	eval, err := filepath.EvalSymlinks(filename)
	if err != nil {
		return "", err
	}
	return filepath.Abs(eval)
}

func fill() {
	path, err := absPath("./")
	if err != nil {
		log.Fatal(err)
	}

	cfg := &packages.Config{
		Mode:  packages.LoadAllSyntax,
		Tests: true,
		Dir:   filepath.Dir(path),
		Fset:  token.NewFileSet(),
		Env:   os.Environ(),
	}

	pkgs, err := packages.Load(cfg)
	if err != nil {
		log.Fatal(err)
	}

	for _, pkg := range pkgs {
		fmt.Printf("%v\n", pkg)
	}
}
