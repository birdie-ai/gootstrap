package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"
	"text/template"
)

//go:embed templates
var templates embed.FS

func main() {
	name := ""
	group := ""

	flag.StringVar(&group, "group", "", "group of the service, like 'enrichment' or 'synthplatform', etc")
	flag.StringVar(&name, "name", "", "name of the service")
	flag.Parse()

	if name == "" || group == "" {
		fmt.Println("'name' and 'group' are obligatory")
		flag.PrintDefaults()
		os.Exit(1)
	}

	generate(name, group)
}

func generate(name, group string) {
	templateDir := "templates/simple"

	type serviceInfo struct {
		Name         string
		Group        string
		ConfigPrefix string
	}

	info := serviceInfo{
		Name:         name,
		Group:        group,
		ConfigPrefix: strings.ToUpper(name),
	}

	applyTemplate := func(templ string) string {
		t := template.Must(template.New("templ").Parse(templ))
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
