package util

import (
	"os"
	"strings"
)

func GetEnv(key, defaultValue string) string {
	val := os.Getenv(key)
	if val != "" {
		return val
	}

	return defaultValue
}

func MergeMap(global, override map[string]interface{}) map[string]interface{} {
	if len(global) == 0 {
		return override
	}

	if len(override) == 0 {
		return global
	}

	for key, val := range override {
		global[key] = val
	}

	return global
}

func ProcessMap(mapToProcess map[string]interface{}) map[string]interface{} {
	for key, val := range mapToProcess {
		if _, ok := val.(string); ok {
			processedVal := strings.Split(val.(string), "|")
			if len(processedVal) > 1 {
				delete(mapToProcess, key)
				mapToProcess[processedVal[0]] = processedVal[1]
			}
		}
	}

	return mapToProcess
}
