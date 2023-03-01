package template

import (
	"embed"
	"io/fs"
	"path/filepath"
	"text/template"
)

var (
	rootDir = "generate"

	//go:embed generate/*
	files embed.FS

	templates map[string]*template.Template
)

func GetTemplate(path string) (*template.Template, error) {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	if templates[path] == nil {
		parsedTemplate, err := template.ParseFS(files, rootDir+"/"+path)
		if err != nil {
			return nil, err
		}

		templates[path] = parsedTemplate
	}

	return templates[path], nil
}

func getAllFilenames() ([]string, error) {
	return getFilenamesIn(rootDir)
}

func getFilenamesIn(rootDir string) ([]string, error) {
	entries := make([]string, 0)

	if err := fs.WalkDir(files, rootDir, func(fileName string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		relativePath, err := filepath.Rel(rootDir, fileName)
		if err != nil {
			return err
		}

		entries = append(entries, relativePath)

		return nil
	}); err != nil {
		return nil, err
	}

	return entries, nil
}
