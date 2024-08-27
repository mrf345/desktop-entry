<h2></h2>
<h1>
desktop-entry
<a href='https://github.com/mrf345/desktop-entry/actions/workflows/ci.yml'>
  <img src='https://github.com/mrf345/desktop-entry/actions/workflows/ci.yml/badge.svg'>
</a>
<a href="https://pkg.go.dev/github.com/mrf345/desktop-entry/safelock">
  <img src="https://pkg.go.dev/badge/github.com/mrf345/desktop-entry/.svg" alt="Go Reference">
</a>
</h1>

Generate and update .desktop (desktop entry) files for Go binaries automatically.

### Install

To add it to your project

```shell
go get https://github.com/mrf345/desktop-entry@latest
```

### How it works

With the default settings `desktopEntry.Create()` will check your `~/.local/share/applications` for a .desktop file, that matches your apps name, if it can't find it, it'll create a new one. That will later on be updated it only when the binary path changes.

### Example

```go
package main

import (
	_ "embed"
	"fmt"
	"os"

	desktopEntry "github.com/mrf345/desktop-entry"
)

// assuming you have an existing icon.png file to embed
//go:embed icon.png
var icon []byte

func main() {
	// Create an instance and pass the required settings
	appName := "Desktop Entry"
	appVersion := "0.0.1"
	entry := desktopEntry.New(appName, appVersion, icon)

	// Some optional settings (check https://pkg.go.dev/github.com/mrf345/desktop-entry#DesktopEntry)
	entry.Comment = "package to help creating desktop entry file for Go"
	entry.Categories = "accessories;Development;"
	entry.Arch = "arm64"

	// Make sure to always run it at the beginning of your main function
	if err := entry.Create(); err != nil {
		panic(err)
	}
}
```
