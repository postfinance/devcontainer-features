package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
)

// Reads and parses the "devcontainer-feature.json" file.
func ParseFeatureJson(featurePath string) (*FeatureSpec, error) {
	fileContent, err := os.ReadFile(filepath.Join(featurePath, "devcontainer-feature.json"))
	if err != nil {
		return nil, err
	}
	var jsonData *FeatureSpec
	if err := json.Unmarshal(fileContent, &jsonData); err != nil {
		return nil, err
	}
	return jsonData, nil
}

type FeatureSpec struct {
	Id             string                `json:"id"`
	Version        string                `json:"version"`
	Name           string                `json:"name"`
	Description    string                `json:"description"`
	Options        OrderedOptionsMap     `json:"options"`
	Customizations FeatureCustomizations `json:"customizations"`
}

type FeatureOption struct {
	Type        string   `json:"type"`
	Default     any      `json:"default"`
	Description string   `json:"description"`
	Proposals   []string `json:"proposals"`
}

type FeatureCustomizations struct {
	VsCode FeatureCustomizationsVsCode `json:"vscode"`
}

type FeatureCustomizationsVsCode struct {
	Extensions []string `json:"extensions"`
}

type OrderedOptionsMap struct {
	Order []string
	Map   map[string]FeatureOption
}

// Custom unmarshaller that also keeps the order of keys in a slice.
func (om *OrderedOptionsMap) UnmarshalJSON(b []byte) error {
	json.Unmarshal(b, &om.Map)

	index := make(map[string]int)
	for key := range om.Map {
		om.Order = append(om.Order, key)
		esc, _ := json.Marshal(key) //Escape the key
		index[key] = bytes.Index(b, esc)
	}

	sort.Slice(om.Order, func(i, j int) bool { return index[om.Order[i]] < index[om.Order[j]] })
	return nil
}
