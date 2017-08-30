package shared

import "os"

func GetEnvOrDefault(env, def string) string {
	value := os.Getenv(env)
	if value == "" {
		value = def
	}
	return value
}
