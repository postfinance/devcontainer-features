package main

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/roemer/gotaskr/log"
)

func BuildReadmeForFeature(featurePath string) error {
	log.Informationf("Build readme for %s", featurePath)
	// Read the specification file
	jsonData, err := ParseFeatureJson(featurePath)
	if err != nil {
		return err
	}

	// Build the template data object
	data := readmeTemplateData{
		Id:          jsonData.Id,
		Name:        jsonData.Name,
		Description: jsonData.Description,
		Version:     jsonData.Version,
		Customizations: readmeTemplateCustomizations{
			VsCodeExtensions: jsonData.Customizations.VsCode.Extensions,
		},
	}
	for _, key := range jsonData.Options.Order {
		option := jsonData.Options.Map[key]
		newOption := readmeTemplateOption{
			Name:        key,
			Description: option.Description,
			Type:        option.Type,
		}
		if newOption.Type == "boolean" {
			newOption.Proposals = "true, false"
			newOption.Default = strconv.FormatBool(option.Default.(bool))
			newOption.DefaultNoQuotes = strconv.FormatBool(option.Default.(bool))
		} else {
			defaultValue := option.Default.(string)
			defaultValueOrEmpty := defaultValue
			if defaultValueOrEmpty == "" {
				defaultValueOrEmpty = "<empty>"
			}
			newOption.Proposals = strings.Join(option.Proposals, ", ")
			newOption.Default = strconv.Quote(defaultValue)
			newOption.DefaultNoQuotes = defaultValueOrEmpty
		}
		data.Options = append(data.Options, newOption)
	}

	// Read the template
	t1, err := template.New("README.md.tmpl").ParseFiles("./build/templates/README.md.tmpl")
	if err != nil {
		return err
	}

	// Execute the template to a memory buffer
	var buf bytes.Buffer
	t1.Execute(&buf, data)

	// Convert the buffer to string
	content := buf.String()

	// Append notes if any
	notesPath := filepath.Join(featurePath, "NOTES.md")
	if _, err := os.Stat(notesPath); err == nil {
		notesContent, err := os.ReadFile(notesPath)
		if err != nil {
			return err
		}
		content += "\n" + string(notesContent)
	}

	// Cleanup multiple newlines
	multiNewLineRegex := regexp.MustCompile(`(?m)^\n{1,}$`)
	content = multiNewLineRegex.ReplaceAllString(content, "")

	// Write the content to the file
	readmePath := filepath.Join(featurePath, "README.md")
	if err := os.WriteFile(readmePath, []byte(content), os.ModePerm); err != nil {
		return err
	}

	return nil
}

type readmeTemplateData struct {
	Id             string
	Name           string
	Description    string
	Version        string
	Options        []readmeTemplateOption
	Customizations readmeTemplateCustomizations
}

type readmeTemplateOption struct {
	Name            string
	Description     string
	Type            string
	Default         string
	DefaultNoQuotes string
	Proposals       string
}

type readmeTemplateCustomizations struct {
	VsCodeExtensions []string
}
