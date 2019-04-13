package main

import (
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	configFile := "config/config.common.yml"
	configuration, err := LoadConfig(configFile)
	if err != nil {
		t.Error("Expected to load config file without error: " + configFile)
	}

	expectedGreetings := []string{
		"Hello World!",
		"Hola Mundo!",
		"Hallo Welt!",
	}
	if !reflect.DeepEqual(expectedGreetings, configuration.Greetings) {
		t.Errorf("Expected configuration to contain %v but was %v", expectedGreetings, configuration.Greetings)
	}
}
