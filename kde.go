//+build linux

package wallpaper

import (
	"bufio"
	"errors"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

func parseKDEConfig() (*[]string, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	filename := filepath.Join(usr.HomeDir, ".config", "plasma-org.kde.plasma.desktop-appletsrc")
	if err != nil {
		return nil, err
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var imagePaths []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) >= 6 && line[:6] == "Image=" {
			imagePaths = append(imagePaths, strings.TrimSpace(removeProtocol(line[6:])))
		}
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	err = file.Close()
	if err != nil {
		return nil, err
	}

	if len(imagePaths) == 0 {
		return nil, errors.New("kde image not found")
	}

	return &imagePaths, nil
}

func setKDEBackground(uri string) error {
	return exec.Command("qdbus", "org.kde.plasmashell", "/PlasmaShell", "org.kde.PlasmaShell.evaluateScript", `
		const monitors = desktops()
		for (var i = 0; i < monitors.length; i++) {
			monitors[i].currentConfigGroup = ["Wallpaper"], ["org.kde.image"], ["General"]
			monitors[i].writeConfig("Image", `+strconv.Quote(uri)+`)
		}
	`).Run()
}
