package main

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
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
	Privileged     bool                  `json:"privileged"`
	Mounts         []FeatureMount        `json:"mounts"`
	EntryPoint     string                `json:"entrypoint"`
}

type FeatureMount struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Type   string `json:"type"`
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

func ExpandDevContainerFeatureVars(value string) string {
	randomId := randomStringWithCharset(16, "abcdefghijklmnopqrstuvwxyz0123456789")
	return strings.ReplaceAll(value, "${devcontainerId}", randomId)
}

func randomStringWithCharset(length int, charset string) string {
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
