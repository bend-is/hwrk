package main

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

var ErrUnsupportedFormat = errors.New("unsupported config file format")

type Config struct {
	Logger LoggerConf `toml:"Logger"`
}

type LoggerConf struct {
	Level  string `toml:"level"`
	Format string `toml:"format"`
}

func ParseConfigFile(filePath string) (Config, error) {
	switch filepath.Ext(filePath) {
	case ".toml":
		var c Config
		_, err := toml.DecodeFile(filePath, &c)
		if err != nil {
			return Config{}, fmt.Errorf("failed to parse config file: %w", err)
		}

		return c, nil
	default:
		return Config{}, ErrUnsupportedFormat
	}
}
