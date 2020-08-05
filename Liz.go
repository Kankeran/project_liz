package main

import (
	"Liz/kernel/managers"

	_ "github.com/joho/godotenv/autoload"

	_ "Liz/kernel/autoload"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// servicesMap, _ := parsers.NewYamlFileReader(make(map[string]interface{})).Read("./config/listeners.yaml")
	// for key, val := range servicesMap.(map[interface{}]interface{}) {
	// 	fmt.Printf("%T", val)
	// 	fmt.Println(key, val)
	// }

	managers.DispatchCommands()
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
