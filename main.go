package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func getWallpapers(wallpapers []string, wallpaperPath string) []string {

	err := filepath.Walk(wallpaperPath, func(path string, info fs.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))

		if ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
			wallpapers = append(wallpapers, path)
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	return wallpapers

}

func rofi(wallpapers []string) string {

	input := strings.Join(wallpapers, "\n")

	cmd := exec.Command("rofi", "-dmenu", "-i", "-p", "Wallpaper")
	cmd.Stdin = strings.NewReader(input)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		panic(err)
	}

	selected := strings.TrimSpace(out.String())

	return selected

}

func main() {
	wallpaperPath := os.Getenv("WALLPAPER_PATH")

	fmt.Println(wallpaperPath)

	var wallpapers []string

	wallpapers = getWallpapers(wallpapers, wallpaperPath)

	selected := rofi(wallpapers)

	print(selected)
	setKdeWallpaper(selected)

}

func setKdeWallpaper(path string) {
script := fmt.Sprintf(`var allDesktops = desktops(); for (var i=0; i<allDesktops.length; i++) { var d = allDesktops[i]; d.wallpaperPlugin = "org.kde.image"; d.currentConfigGroup = ["Wallpaper", "org.kde.image", "General"]; d.writeConfig("Image", "file://%s"); }`, path)
	exec.Command(
		"qdbus6",
		"org.kde.plasmashell",
		"/PlasmaShell",
		"org.kde.PlasmaShell.evaluateScript",
		script,
	).Run()
}
