package managers

import (
	"flag"
	"fmt"
	"os"

	"Liz/builder"
	"Liz/kernel/container"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func DispatchCommands() {
	newCmd := flag.NewFlagSet("new", flag.ExitOnError)
	// newType := newCmd.String("type", "core", "type of new project [core, desktop:console]")
	newName := newCmd.String("name", "App", "type new project name")
	newPath := newCmd.String("path", "./", "path to new project")

	buildCmd := flag.NewFlagSet("build", flag.ExitOnError)
	buildPath := buildCmd.String("path", "./", "path to project to build")

	flag.Usage = getUsage(newCmd, buildCmd)
	flag.Parse()

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "new":
		check(newCmd.Parse(os.Args[2:]))
		container.Get("project_starter_builder").(*builder.ProjectStarter).Build(*newName, *newPath)
	case "build":
		check(buildCmd.Parse(os.Args[2:]))
		container.Get("container_builder").(*builder.Container).Build(*buildPath)
	default:
		flag.Usage()
		os.Exit(2)
	}
}

func getUsage(cmds ...*flag.FlagSet) func() {
	return func() {
		_, _ = fmt.Fprintf(os.Stderr, "Usage of %s [command] [flags...]\n", os.Args[0])
		_, _ = fmt.Fprintf(os.Stderr, "\ncommands:\n")
		_, _ = fmt.Fprintf(os.Stderr, "  - new\n")
		_, _ = fmt.Fprintf(os.Stderr, "  - build\n")
		for _, cmd := range cmds {
			_, _ = fmt.Fprintf(
				os.Stderr,
				"\n-----------------------------------------\n\nflags for "+cmd.Name()+"\n",
			)
			cmd.PrintDefaults()
		}
		_, _ = fmt.Fprintf(os.Stderr, "\n-----------------------------------------\n")
	}
}
