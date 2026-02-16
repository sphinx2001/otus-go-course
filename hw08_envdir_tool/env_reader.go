package main

import (
	"bytes"
	"fmt"
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
	// Place your code here
	result := make(Environment)
	fd, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	files, err := fd.Readdir(-1)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if !file.IsDir() {
			path := filepath.Join(dir, file.Name())
			v, err := ReadValue(path)
			if err == nil {
				if v != "" {
					result[file.Name()] = EnvValue{Value: v, NeedRemove: false}
				} else {
					result[file.Name()] = EnvValue{Value: v, NeedRemove: true}
				}
			}
		}
	}
	return result, nil
}

func PrepareEnv(env Environment) []string {
	if env == nil {
		return nil
	}
	result := make([]string, 0, len(env))

	oldEnvs := os.Environ()

	for _, line := range oldEnvs {
		v := strings.Split(line, "=")
		name := v[0]
		value := v[1]
		_, ok := env[name]
		if !ok {
			env[name] = EnvValue{Value: value, NeedRemove: false}
		}
	}

	for key, value := range env {
		if !value.NeedRemove {
			result = append(result, fmt.Sprintf("%s=%s", key, value.Value))
		}
	}
	return result
}

func ReadValue(path string) (string, error) {
	value := ""
	val, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	value = string(val)

	value = strings.Split(value, "\n")[0]
	value = strings.TrimRight(value, " \t\n")

	mask := []byte{byte(0)}
	buf := bytes.ReplaceAll([]byte(value), mask, []byte("\n"))
	value = string(buf)

	CheckValue(path, value)
	return value, nil
}

func CheckValue(path, value string) {
	b := []byte(value)
	if bytes.Contains(b, []byte{byte(0)}) {
		fmt.Println(path, " - contains NUL value")
	}
}
