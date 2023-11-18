package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

const SupportedExt = "conf"

var (
	ErrUnsupportedExt = errors.New("unsupported extension")
	ErrParseConfig    = errors.New("invalid syntax")
)

type Config struct {
	m map[string]string
}

func (c Config) Get(key string) string {
	return c.m[key]
}

func ParseConfig(filename string) (Config, error) {
	cfg := Config{}

	ind := strings.LastIndex(filename, ".")
	if ind == -1 || ind+1 == len(filename) || filename[ind+1:] != SupportedExt {
		return cfg, ErrUnsupportedExt
	}

	file, err := os.Open(filename)
	if err != nil {
		return cfg, fmt.Errorf("parse config %s: %w", filename, err)
	}
	defer file.Close()

	cfg.m = make(map[string]string)

	scan := bufio.NewScanner(file)
	var pair []string
	for scan.Scan() {
		pair = strings.SplitN(scan.Text(), " ", 2)
		if len(pair) != 2 {
			return Config{}, ErrParseConfig
		}
		cfg.m[pair[0]] = pair[1]
	}

	if scan.Err() != nil {
		return Config{}, fmt.Errorf("parse config %s: %w", filename, scan.Err())
	}
	return cfg, nil
}
