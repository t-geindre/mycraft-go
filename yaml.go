package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

func ParseYamlFile[T any](file string) T {
	rawYaml, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	output := new(T)
	err = yaml.Unmarshal(rawYaml, output)
	if err != nil {
		panic(err)
	}

	return *output
}
