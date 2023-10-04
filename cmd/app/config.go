package main

import "go.uber.org/config"

func newConfig(filename string) (*config.YAML, error) {
	return config.NewYAML(config.File(filename))
}
