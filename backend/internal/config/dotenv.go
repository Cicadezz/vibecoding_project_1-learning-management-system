package config

import (
	"bufio"
	"os"
	"strings"
)

// LoadDotEnv loads simple KEY=VALUE pairs from the given file path.
// It ignores blank lines and lines starting with '#'.
func LoadDotEnv(path string) error {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		eq := strings.Index(line, "=")
		if eq <= 0 {
			continue
		}

		key := strings.TrimSpace(line[:eq])
		if key == "" {
			continue
		}
		value := strings.TrimSpace(line[eq+1:])
		value = strings.Trim(value, "\"")
		value = strings.Trim(value, "'")

		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, value)
		}
	}

	return scanner.Err()
}
