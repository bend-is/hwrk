package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
	"unicode"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	res := make(Environment, len(dirEntries))

	for _, entry := range dirEntries {
		if entry.IsDir() || !entry.Type().IsRegular() || strings.Contains(entry.Name(), "=") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			return nil, err
		}

		if info.Size() == 0 {
			res[entry.Name()] = EnvValue{NeedRemove: true}

			continue
		}

		value, err := readFirstLine(path.Join(dir, info.Name()))
		if err != nil {
			return nil, err
		}

		res[entry.Name()] = EnvValue{Value: value}
	}

	return res, nil
}

func readFirstLine(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Scan()

	value := strings.ReplaceAll(scanner.Text(), "\x00", "\n")
	value = strings.TrimRightFunc(value, unicode.IsSpace)

	return value, nil
}

func (e Environment) Apply() error {
	for k, v := range e {
		if v.NeedRemove {
			if err := os.Unsetenv(k); err != nil {
				return fmt.Errorf("unset environment %s: %w", k, err)
			}

			continue
		}

		if err := os.Setenv(k, v.Value); err != nil {
			return fmt.Errorf("set environment %s: %w", k, err)
		}
	}

	return nil
}
