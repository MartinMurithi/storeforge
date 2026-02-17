package env

import "os"

// getEnv returns the value of the environment variable `key`
// or `fallback` if the variable is not set.
func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}