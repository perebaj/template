package main

import (
	"bytes"
	"embed"
	"flag"
	"io/fs"
	"log"
	"os"
	"strings"
	"text/template"
)

//go:embed templates/core/*
var templatesDir embed.FS

const templateDir = "templates/core"

type Project struct {
	Name     string
	Registry string
}

func main() {

	projectName := flag.String("name", "", "Project name")
	registryName := flag.String("registry", "FAKEREGISTRY", "Registry name")
	output := flag.String("output", ".", "Output directory")
	flag.Parse()

	if *projectName == "" {
		log.Fatal("project name is required")
	}

	var project Project
	project.Name = *projectName
	project.Registry = *registryName

	fs.WalkDir(templatesDir, templateDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if path == templateDir {
			return nil
		}
		target := strings.TrimPrefix(path, templateDir+"/")
		target = strings.TrimSuffix(target, ".template")

		template1 := template.Must(template.New("dir").Parse(target))

		var buf bytes.Buffer
		err = template1.Execute(&buf, project)
		if err != nil {
			log.Fatalf("failed executing template: %s", err)
		}

		outputPath := *output + "/" + buf.String()

		if d.IsDir() {
			log.Println("Creating directory", outputPath)

			err := os.MkdirAll(outputPath, 0755)
			if err != nil {
				log.Fatalf("failed creating directory: %s", err)
			}
		} else {
			log.Println("Creating file", outputPath)
			data, err := os.ReadFile(path)
			if err != nil {
				log.Fatalf("failed reading data from file: %s", err)
			}

			template2 := template.Must(template.New("file").Parse(string(data)))
			var result bytes.Buffer
			err = template2.Execute(&result, project)
			if err != nil {
				log.Fatalf("failed executing template: %s", err)
			}
			err = os.WriteFile(outputPath, result.Bytes(), 0644)
			if err != nil {
				log.Fatalf("failed writing data to file: %s", err)
			}
		}

		return nil
	})
}
