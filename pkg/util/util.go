package util

import (
	"bytes"
	"os"
	"text/template"

	"gopkg.in/yaml.v2"
)

func GetEnv(key, defaultValue string) string {
	val := os.Getenv(key)
	if val != "" {
		return val
	}

	return defaultValue
}

func ReplaceTemplateValues(templateData map[string]interface{},
	values map[string]interface{}) (transformedData map[string]interface{}, err error) {
	yamlData, err := yaml.Marshal(templateData)
	if err != nil {
		return
	}

	tmpl, err := template.New("templateVal").Parse(string(yamlData))
	if err != nil {
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, values)
	if err != nil {
		return
	}

	transformedData = map[string]interface{}{}
	err = yaml.Unmarshal(buf.Bytes(), &transformedData)
	if err != nil {
		return
	}
	return
}
