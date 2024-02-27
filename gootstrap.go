package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"strings"
	"text/template"
)

//go:embed templates/*
var templates embed.FS

func main() {
	name := ""
	group := ""
	templateName := ""

	flag.StringVar(&group, "group", "", "group of the service, like 'enrichment' or 'synthplatform', etc")
	flag.StringVar(&name, "name", "", "name of the service")
	flag.StringVar(&templateName, "template", "basic", "name of the template")
	flag.Parse()

	if name == "" || group == "" {
		fmt.Println("'name' and 'group' are obligatory")
		flag.PrintDefaults()
		os.Exit(1)
	}

	generate(name, group, templateName)
}

func generate(name, group, templateName string) {
	type serviceInfo struct {
		Name         string
		Group        string
		ConfigPrefix string
	}

	templateDir := path.Join("templates", templateName)

	info := serviceInfo{
		Name:         name,
		Group:        group,
		ConfigPrefix: strings.ToUpper(name),
	}

	applyTemplate := func(templ string) string {
		// Why: some files, like Github CI files, have notation identical
		// To Go templates. So for now we just try to apply the template and if it
		// fails we assume it should be used as is (not the safest option, but no more time to deal
		// with this right now).
		t, err := template.New("templ").Parse(templ)
		if err != nil {
			return templ
		}
		result := bytes.Buffer{}
		assert(t.Execute(&result, info))
		return result.String()
	}

	err := fs.WalkDir(templates, templateDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if path == templateDir {
			return nil
		}
		targetPath := strings.TrimPrefix(path, templateDir+"/")
		targetPath = strings.TrimSuffix(targetPath, ".template")
		targetPath = applyTemplate(targetPath)

		if d.IsDir() {
			log.Printf("creating dir %q", targetPath)
			assert(os.MkdirAll(targetPath, 0755))
			return nil
		}

		log.Printf("generating file %q", targetPath)

		data, err := fs.ReadFile(templates, path)
		assert(err)

		generated := applyTemplate(string(data))

		log.Printf("writing file %q", targetPath)

		assert(os.WriteFile(targetPath, []byte(generated), 0644))

		return nil
	})

	assert(err)
}

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
