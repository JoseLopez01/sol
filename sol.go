package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		panic("Please provide a command")
	}

	createOptFolder()

	command := os.Args[1]

	switch command {
	case "install":
		if len(os.Args) < 3 {
			panic("Please provide a version to install")
		}

		version := os.Args[2]
		install(version)
	case "remove":
		if len(os.Args) < 3 {
			panic("Please provide a version to remove")
		}

		version := os.Args[2]
		remove(version)
	case "use":
		if len(os.Args) < 3 {
			panic("Please provide a version to use")
		}

		version := os.Args[2]
		use(version)
	case "ls":
		list()
	default:
		panic("Unknown command")
	}
}

func install(version string) {
	if isInstalled(version) {
		fmt.Printf("Version %s is already installed\n", version)
		return
	}

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

	dest := fmt.Sprintf("%s/.sol/versions/v%s", homeDir, version)
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
	if !isInstalled(version) {
		fmt.Printf("Version %s is not installed\n", version)
		return
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	bin := fmt.Sprintf("%s/.sol/bin", homeDir)
	symlink, err := os.Readlink(bin)
	if err != nil {
		log.Fatal(err)
	}

	dest := fmt.Sprintf("%s/.sol/versions/v%s", homeDir, version)
	nodeBin := fmt.Sprintf("%s/bin", dest)
	if nodeBin == symlink {
		if err := os.Remove(bin); err != nil {
			log.Fatal(err)
		}
	}

	if err := os.RemoveAll(dest); err != nil {
		log.Fatal(err)
	}
}

func use(version string) {
	if !isInstalled(version) {
		fmt.Printf("Version %s is not installed\n", version)
		return
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	bin := fmt.Sprintf("%s/.sol/bin", homeDir)
	nodeBin := fmt.Sprintf("%s/.sol/versions/v%s/bin", homeDir, version)
	if err := os.Remove(bin); err != nil {
		log.Fatal(err)
	}

	if err := os.Symlink(nodeBin, bin); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Node.js version %s is now in use\n", version)
}

func list() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	versionsDir := fmt.Sprintf("%s/.sol/versions", homeDir)
	entries, err := os.ReadDir(versionsDir)
	if os.IsNotExist(err) {
		fmt.Println("No versions installed")
		return
	}

	if err != nil {
		log.Fatal(err)
	}

	if len(entries) == 0 {
		fmt.Println("No versions installed")
		return
	}

	currentVersion := ""
	bin := fmt.Sprintf("%s/.sol/bin", homeDir)
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
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	dest := fmt.Sprintf("%s/.sol/versions/v%s", homeDir, version)
	if _, err := os.Stat(dest); err == nil {
		return true
	}

	return false
}
