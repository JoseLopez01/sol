package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		panic("Please provide a command and a version")
	}

	createOptFolder()

	command := os.Args[1]
	version := os.Args[2]

	switch command {
	case "install":
		install(version)
	case "remove":
		remove(version)
	case "use":
		use(version)
	default:
		panic("Unknown command")
	}
}

func install(version string) {
	url := fmt.Sprintf("https://nodejs.org/download/release/v%s/node-v%s-darwin-arm64.tar.gz", version, version)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatal(fmt.Sprintf("Failed to download file: %s", resp.Status))
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	dest := fmt.Sprintf("%s/.sol/versions/node-v%s", homeDir, version)
	if err := extractFile(resp.Body, dest); err != nil {
		log.Fatal(err)
	}

	nodeBin := fmt.Sprintf("%s/bin", dest)
	bin := fmt.Sprintf("%s/.sol/bin", homeDir)
	os.Remove(bin)

	if err := os.Symlink(nodeBin, bin); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Node.js version %s installed successfully\n", version)
}

func remove(version string) {
	fmt.Printf("Removing version %s\n", version)
}

func use(version string) {
	fmt.Printf("Using version %s\n", version)
}

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
