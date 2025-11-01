package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
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
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	result := make(Environment, len(files))

	for _, file := range files {
		fname := file.Name()
		if file.IsDir() || strings.Contains(fname, "=") {
			continue
		}

		fpath := filepath.Join(dir, fname)
		openFile, err := os.Open(fpath)
		if err != nil {
			return nil, err
		}
		defer openFile.Close()

		infoFile, err := openFile.Stat()
		if err != nil {
			return nil, err
		}

		if infoFile.Size() == 0 {
			result[fname] = EnvValue{"", true}
			continue
		}
		scanner := bufio.NewScanner(openFile)

		var value string
		if scanner.Scan() {
			value = scanner.Text()
			value = strings.ReplaceAll(value, "\x00", "\n")
			value = strings.TrimRight(value, " \t\r\n")
		}

		result[fname] = EnvValue{value, false}
	}

	return result, nil
}
