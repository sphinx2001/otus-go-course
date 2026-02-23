package main

import "testing"

func TestRunCmd(t *testing.T) {
	tests := []struct {
		name         string
		cmd          []string
		env          Environment
		wantExitCode int
	}{
		{
			name:         "successful command",
			cmd:          []string{"echo", "hello"},
			env:          Environment{},
			wantExitCode: 0,
		},
		{
			name:         "command with custom env",
			cmd:          []string{"/bin/sh", "-c", "echo $TEST_VAR"},
			env:          Environment{"TEST_VAR": {Value: "test_value", NeedRemove: false}},
			wantExitCode: 0,
		},
		{
			name:         "failing command",
			cmd:          []string{"false"},
			env:          Environment{},
			wantExitCode: 1,
		},
		{
			name:         "command not found",
			cmd:          []string{"nonexistent_command_12345"},
			env:          Environment{},
			wantExitCode: -1, // exec.Error возвращает -1
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exitCode := RunCmd(tt.cmd, tt.env)
			if exitCode != tt.wantExitCode {
				t.Errorf("RunCmd() exitCode = %d, want %d", exitCode, tt.wantExitCode)
			}
		})
	}
}

func TestRunCmdWithEnvironment(t *testing.T) {
	env := Environment{
		"GO_ENVDIR_TEST": {Value: "test_value_123", NeedRemove: false},
	}

	cmd := []string{"/bin/sh", "-c", "echo $GO_ENVDIR_TEST"}
	exitCode := RunCmd(cmd, env)

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}
}
