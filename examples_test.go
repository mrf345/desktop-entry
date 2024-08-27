package desktopEntry_test

import (
	_ "embed"
	"fmt"
	"os"

	desktopEntry "github.com/mrf345/desktop-entry"
)

// assuming you have an existing icon.png file to embed
//
//go:embed icon.png
var icon []byte

func ExampleDesktopEntry_Create() {
	// Create an instance and pass the required settings
	appName := "Desktop Entry"
	appVersion := "0.0.1"
	entry := desktopEntry.New(appName, appVersion, icon)

	// Some optional settings
	entry.Comment = "package to help creating desktop entry file for Go"
	entry.Categories = "accessories;Development;"
	entry.Arch = "arm64"

	// Changing the apps and icons path to `/tmp` for the test (should ignore this)
	os.Args[0] = "/tmp"
	tempDir, _ := os.MkdirTemp("", "de_example")
	entry.AppsPath = tempDir
	entry.IconsPath = tempDir
	defer os.RemoveAll(tempDir)

	// Make sure to always run it at the beginning of your main function
	if err := entry.Create(); err != nil {
		fmt.Println(err)
	}

	// Output:
}
