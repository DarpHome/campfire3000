package main

import (
	"encoding/json"
	"os"
	"strings"
)

var (
	Locales map[string]map[string]interface{}
)

func InitializeI18n() {
	Locales = map[string]map[string]interface{}{}
	entries, err := os.ReadDir("locales")
	if err != nil {
		Logger.Fatal("ReadDir", err)
	}
	for _, entry := range entries {
		locale, err := os.ReadFile("./locales/" + entry.Name())
		if err != nil {
			Logger.Fatal("ReadFile", err)
		}
		var m map[string]interface{}
		if err = json.Unmarshal(locale, &m); err != nil {
			Logger.Fatal("Unmarshal", err)
		}
		name := strings.ReplaceAll(entry.Name(), ".json", "")
		Locales[name] = m
	}
}

func GetI18nValue(locale, key string, defaultValue interface{}) interface{} {
	keys, ok := Locales[locale]
	if ok {
		if val, ok := keys[key]; ok {
			return val
		}
		return defaultValue
	}
	keys, ok = Locales["en-US"]
	if !ok {
		Logger.Error("No en-US in locales")
		return defaultValue
	}
	if val, ok := keys[key]; ok {
		return val
	}
	return defaultValue
}

func GetStringI18nValue(locale, key string) string {
	return GetI18nValue(locale, key, key).(string)
}
