package main

import "go.uber.org/config"

var configFiles = []string{"configs/config.yml"}

func newConfig() (*config.YAML, error) {
	cfgOption := make([]config.YAMLOption, 0, len(configFiles))
	for _, filename := range configFiles {
		cfgOption = append(cfgOption, config.File(filename))
	}
	return config.NewYAML(cfgOption...)
}
