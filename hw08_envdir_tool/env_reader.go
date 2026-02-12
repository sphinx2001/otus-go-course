package main

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
	return nil, nil
}

func PrepareEnv(env Environment) []string {
	if env == nil {
		return nil
	}

	result := make([]string, 0, len(env))
	return result
}
