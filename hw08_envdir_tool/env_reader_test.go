package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadDir(t *testing.T) {
	tests := []struct {
		name    string
		dir     string
		wantErr bool
		check   func(t *testing.T, env Environment)
	}{
		{
			name:    "valid directory",
			dir:     "./testdata/env",
			wantErr: false,
			check: func(t *testing.T, env Environment) {
				if len(env) == 0 {
					t.Error("Expected non-empty environment")
				}
				// Проверка конкретных переменных
				if v, ok := env["HELLO"]; ok {
					if v.Value != `"hello"` {
						t.Errorf("HELLO: expected '\"hello\"', got '%s'", v.Value)
					}
				}
				if v, ok := env["UNSET"]; ok {
					if !v.NeedRemove {
						t.Error("UNSET should have NeedRemove=true")
					}
				}
			},
		},
		{
			name:    "non-existent directory",
			dir:     "./nonexistent",
			wantErr: true,
			check:   nil,
		},
		{
			name:    "empty directory",
			dir:     t.TempDir(),
			wantErr: false,
			check: func(t *testing.T, env Environment) {
				if len(env) != 0 {
					t.Errorf("Expected empty environment, got %d items", len(env))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env, err := ReadDir(tt.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.check != nil {
				tt.check(t, env)
			}
		})
	}
}

func TestReadValue(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name     string
		content  []byte
		expected string
		wantErr  bool
	}{
		{
			name:     "simple value",
			content:  []byte("hello\n"),
			expected: "hello",
			wantErr:  false,
		},
		{
			name:     "value with trailing spaces",
			content:  []byte("hello   \t\n"),
			expected: "hello",
			wantErr:  false,
		},
		{
			name:     "empty file",
			content:  []byte(""),
			expected: "",
			wantErr:  false,
		},
		{
			name:     "only whitespace",
			content:  []byte("   \t\n"),
			expected: "",
			wantErr:  false,
		},
		{
			name:     "null byte replacement",
			content:  []byte("foo\x00bar\n"),
			expected: "foo\nbar",
			wantErr:  false,
		},
		{
			name:     "multiline - first line only",
			content:  []byte("first\nsecond\nthird\n"),
			expected: "first",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := filepath.Join(tmpDir, tt.name)
			err := os.WriteFile(filePath, tt.content, 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			value, err := ReadValue(filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if value != tt.expected {
				t.Errorf("ReadValue() = '%s', want '%s'", value, tt.expected)
			}
		})
	}
}

func TestPrepareEnv(t *testing.T) {
	tests := []struct {
		name     string
		env      Environment
		checkLen func(t *testing.T, result []string)
	}{
		{
			name: "empty environment",
			env:  Environment{},
			checkLen: func(t *testing.T, result []string) {
				// Должно содержать хотя бы системные переменные
				if len(result) == 0 {
					t.Error("Expected non-zero environment variables")
				}
			},
		},
		{
			name: "with custom variables",
			env: Environment{
				"TEST_VAR": {Value: "test_value", NeedRemove: false},
			},
			checkLen: func(t *testing.T, result []string) {
				found := false
				for _, v := range result {
					if v == "TEST_VAR=test_value" {
						found = true
						break
					}
				}
				if !found {
					t.Error("TEST_VAR not found in prepared environment")
				}
			},
		},
		{
			name: "NeedRemove flag",
			env: Environment{
				"REMOVE_ME": {Value: "", NeedRemove: true},
			},
			checkLen: func(t *testing.T, result []string) {
				for _, v := range result {
					if len(v) > 9 && v[:9] == "REMOVE_ME" {
						t.Error("REMOVE_ME should not be in prepared environment")
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PrepareEnv(tt.env)
			if tt.checkLen != nil {
				tt.checkLen(t, result)
			}
		})
	}
}
