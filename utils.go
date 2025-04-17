package main

import (
	"fmt"
	"os"
)

func createOptFolder() {
	_, err := os.Stat("/opt/sol")
	if err == nil {
		return
	}

	if !os.IsNotExist(err) {
		panic(err)
	}

	if err = os.MkdirAll("/opt/sol", 0o755); err != nil {
		panic(err)
	}
}

func isInstalled(version string) bool {
	dest := getHomeBasedPath(".sol", "versions", fmt.Sprintf("v%s", version))
	if _, err := os.Stat(dest); err == nil {
		return true
	}

	return false
}

func exit(message string) {
	fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}

func getHomeBasedPath(parts ...string) string {
	homeDir := os.Getenv("HOME")
	result := homeDir
	if len(parts) == 0 {
		return result
	}

	for _, part := range parts {
		result = fmt.Sprintf("%s/%s", result, part)
	}

	return result
}
