package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		exit("Please provide a command (install, remove, use, ls)")
	}

	createOptFolder()

	command := os.Args[1]
	switch command {
	case "install":
		if len(os.Args) < 3 {
			exit("Please provide a version to install")
		}

		version := os.Args[2]
		install(version)
	case "remove":
		if len(os.Args) < 3 {
			exit("Please provide a version to remove")
		}

		version := os.Args[2]
		remove(version)
	case "use":
		if len(os.Args) < 3 {
			exit("Please provide a version to use")
		}

		version := os.Args[2]
		use(version)
	case "ls":
		list()
	default:
		exit("Unknown command")
	}
}

func install(version string) {
	if isInstalled(version) {
		exit(fmt.Sprintf("Version %s is already installed\n", version))
	}

	url := fmt.Sprintf("https://nodejs.org/download/release/v%s/node-v%s-darwin-arm64.tar.gz", version, version)
	resp, err := http.Get(url)
	if err != nil {
		exit(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		exit(fmt.Sprintf("Failed to download file: %s", resp.Status))
	}

	dest := getHomeBasedPath(".sol", "versions", fmt.Sprintf("v%s", version))
	if err := extractFile(resp.Body, dest); err != nil {
		exit(err.Error())
	}

	nodeBin := fmt.Sprintf("%s/bin", dest)
	bin := getHomeBasedPath(".sol", "bin")
	os.Remove(bin)

	if err := os.Symlink(nodeBin, bin); err != nil {
		exit(err.Error())
	}

	fmt.Printf("Node.js version %s installed successfully\n", version)
}

func remove(version string) {
	if !isInstalled(version) {
		exit(fmt.Sprintf("Version %s is not installed\n", version))
	}

	bin := getHomeBasedPath(".sol", "bin")
	symlink, err := os.Readlink(bin)
	if err != nil {
		exit(err.Error())
	}

	dest := getHomeBasedPath(".sol", "versions", fmt.Sprintf("v%s", version))
	nodeBin := fmt.Sprintf("%s/bin", dest)
	if nodeBin == symlink {
		if err := os.Remove(bin); err != nil {
			exit(err.Error())
		}
	}

	if err := os.RemoveAll(dest); err != nil {
		exit(err.Error())
	}
}

func use(version string) {
	if !isInstalled(version) {
		exit(fmt.Sprintf("Version %s is not installed\n", version))
	}

	bin := getHomeBasedPath(".sol", "bin")
	nodeBin := getHomeBasedPath(".sol", "versions", fmt.Sprintf("v%s", version))
	if err := os.Remove(bin); err != nil {
		exit(err.Error())
	}

	if err := os.Symlink(nodeBin, bin); err != nil {
		exit(err.Error())
	}

	fmt.Printf("Node.js version %s is now in use\n", version)
}

func list() {
	versionsDir := getHomeBasedPath(".sol", "versions")
	entries, err := os.ReadDir(versionsDir)
	if os.IsNotExist(err) {
		exit("No versions installed")
		return
	}

	if err != nil {
		log.Fatal(err)
	}

	if len(entries) == 0 {
		exit("No versions installed")
	}

	currentVersion := ""
	bin := getHomeBasedPath(".sol", "bin")
	if target, err := os.Readlink(bin); err == nil {
		currentVersion = target
	}

	fmt.Println("Installed versions:")
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		version := entry.Name()
		if version == currentVersion {
			fmt.Printf("  %s (current)\n", version)
		} else {
			fmt.Printf("  %s\n", version)
		}
	}
}
