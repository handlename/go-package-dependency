package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alecthomas/kingpin/v2"
)

func main() {
	var (
		app                = kingpin.New("go-package-depends", "Generate dependency.gen.go files based on DEPENDENCY.md")
		dependencyFilePath = app.Arg("dependency-file", "Path to the DEPENDENCY.md file").Required().String()
	)

	app.HelpFlag.Short('h')
	app.Version(Version)

	kingpin.MustParse(app.Parse(os.Args[1:]))

	parser := NewParser()
	config, err := parser.ParseDependencyFile(*dependencyFilePath)
	if err != nil {
		fmt.Printf("Error parsing dependency file: %v\n", err)
		os.Exit(1)
	}

	baseDir := filepath.Dir(*dependencyFilePath)
	generator := NewGenerator()
	err = generator.GenerateDependencyFiles(baseDir, config)
	if err != nil {
		fmt.Printf("Error generating dependency files: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Generated dependency.gen.go files successfully")
}
