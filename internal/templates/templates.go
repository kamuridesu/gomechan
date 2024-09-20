package templates

import (
	"fmt"
	"html/template"
	"os"
	"path"
	"strings"
)

type Template struct {
	folder    string
	templates []string
}

func LoadTemplateFolder(folder string) (*Template, error) {
	fsFolder, err := os.ReadDir(folder)
	if err != nil {
		return nil, fmt.Errorf("error reading from folder: %s", err)
	}
	files := []string{}
	for _, file := range fsFolder {
		if strings.HasSuffix(file.Name(), ".tmpl") {
			files = append(files, file.Name())
		}
	}
	return &Template{
		folder:    folder,
		templates: files,
	}, nil
}

func (t *Template) GetTemplate(filename string) (string, error) {
	for _, file := range t.templates {
		if file == filename {
			c, err := os.ReadFile(path.Join(t.folder, file))
			if err != nil {
				return "", err
			}
			return string(c), nil
		}
	}
	return "", fmt.Errorf("template %s not found", filename)
}

func (t *Template) LoadHTML(name string, variables map[string]any) string {
	buff := new(strings.Builder)

	content, err := t.GetTemplate(name)
	if err != nil {
		return ""
	}

	tmpl, err := template.New("template").Parse(content)
	if err != nil {
		return ""
	}

	err = tmpl.Execute(buff, variables)
	if err != nil {
		return ""
	}

	return buff.String()
}
