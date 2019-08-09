package gomailer

import (
	"bytes"
	"html/template"
	"path/filepath"
	"strings"
)

const (
	DefaultLayoutExtension = "html"
	TemplateRoot           = `{{define "root" }} {{ template "main" . }} {{ end }}`
)

type TemplateConfig struct {
	LayoutFiles     []string
	LayoutDirectory string
	LayoutExtension string
}

type Template struct {
	file   string
	config *TemplateConfig
}

// Template Constructor
func NewTemplate(file string, config *TemplateConfig) *Template {
	return &Template{
		file:   file,
		config: config,
	}
}

// Parse template
func (t *Template) Parse(data interface{}) (*string, error) {
	// Parse root template
	var err error
	mainTemplate := template.New("main")
	if mainTemplate, err = mainTemplate.Parse(TemplateRoot); err != nil {
		return nil, err
	}

	// Get all template files
	files, err := t.getFiles()
	if err != nil {
		return nil, err
	}

	// Parse template files
	tmpl, err := mainTemplate.ParseFiles(files...)
	if err != nil {
		return nil, err
	}

	var output bytes.Buffer
	if err := tmpl.Execute(&output, data); err != nil {
		return nil, err
	}
	body := output.String()
	return &body, nil
}

// Get all template files (include layouts)
func (t *Template) getFiles() ([]string, error) {
	layoutFiles, err := t.getLayoutFiles()
	if err != nil {
		return nil, err
	}
	files := append(layoutFiles, t.file)
	return files, nil
}

// Get all template layouts
func (t *Template) getLayoutFiles() ([]string, error) {
	var scannedFiles []string
	var err error
	if len(t.config.LayoutDirectory) > 0 {
		layoutExt := DefaultLayoutExtension
		if len(t.config.LayoutExtension) > 0 {
			layoutExt = t.config.LayoutExtension
		}
		dir := strings.TrimRight(t.config.LayoutDirectory, "/")
		if scannedFiles, err = filepath.Glob(dir + "/*." + layoutExt); err != nil {
			return nil, err
		}
	}
	return append(scannedFiles, t.config.LayoutFiles...), nil
}
