package main

import (
	"os"

	"Liz/container"
	"Liz/elements"
	"Liz/generated"
	"Liz/generators"
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
	generated.Build()
	servicesMap, err := parsers.ReadYamlFile("./config/services.yaml")
	check(err)

	servicesMap, err = parseReferences(servicesMap.(map[interface{}]interface{}), "./config/services.yaml")
	check(err)
	servicesMap = parseServices(servicesMap.(map[interface{}]interface{}))

	var generator = container.Get("service_generator").(*generators.Service)
	var code = `package generated

	// Build building container container
	func Build() {

		`

	for serviceName, serviceMap := range servicesMap.(map[interface{}]interface{}) {
		code += "container.Set(\"" + serviceName.(string) + "\", " + generator.Generate(elements.NewService(serviceMap.(map[interface{}]interface{}))) + ")\n\n"
	}
	code += "}"

	var output []byte
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
	return container.Get("service_parser").(*parsers.Service).Parse(source)
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
