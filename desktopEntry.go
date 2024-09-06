// Generate and update .desktop (desktop entry) files for Go binaries automatically.
//
// # Install
//
// Add it to your project
//
//	go get https://github.com/mrf345/desktop-entry@latest
//
// # How it works
//
// With the default settings shown in [desktopEntry.DesktopEntry] the method [desktopEntry.DesktopEntry.Create]
// will check your [desktopEntry.DesktopEntry.AppsPath] for a .desktop file, that matches your
// [desktopEntry.DesktopEntry.Name]-[desktopEntry.DesktopEntry.Version], if it can't find it, it'll create a new one.
// That will later on be updated it only when the binary path changes.
// See test [example].
//
// [example]: https://pkg.go.dev/github.com/mrf345/desktop-entry#example-DesktopEntry.Create
package desktopEntry

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strings"
)

type DesktopEntry struct {
	// Application name (required)
	Name string
	// Application version (required)
	Version string
	// Application .png icon data (required)
	Icon []byte
	// Executable type (default: Application)
	Type string
	// Semicolon separated list of categories (default: '')
	Categories string
	// Description of the app (default: '')
	Comment string
	// Architecture (default: x86_64)
	Arch string
	// Desktop applications path (default: ~/.local/share/applications)
	AppsPath string
	// Desktop icons path (default: ~/.icons)
	IconsPath string
	// Default permission for created files and directories (default: 0776)
	Perm fs.FileMode
	// Supported operating systems (default: []string{"linux"})
	OSs []string
	// Update if executable path has changed (default: true)
	UpdateIfChanged bool
	// Rerun the app if desktop entry has changed (default: true)
	RerunIfChanged bool
}

// Create a new [desktopEntry.DesktopEntry] instance with the default options
func New(name, version string, icon []byte) *DesktopEntry {
	return &DesktopEntry{
		Name:            name,
		Version:         version,
		Icon:            icon,
		Type:            "Application",
		Arch:            "x86_64",
		Perm:            0776,
		AppsPath:        fmt.Sprintf("%s/.local/share/applications", os.Getenv("HOME")),
		IconsPath:       fmt.Sprintf("%s/.icons", os.Getenv("HOME")),
		OSs:             []string{"linux"},
		UpdateIfChanged: true,
		RerunIfChanged:  true,
	}
}

// Creates a new desktop entry or updates an existing one if the executable paths mismatch
func (de *DesktopEntry) Create() (err error) {
	var changed bool

	isDevBuild := strings.HasPrefix(os.Args[0], os.TempDir())
	isSupportedOs := slices.Contains(de.OSs, runtime.GOOS)

	if isDevBuild || !isSupportedOs {
		return
	}

	if err = de.createPaths(); err != nil {
		err = fmt.Errorf("failed to create app or icon paths > %w", err)
		return
	}

	if err = de.createIcon(); err != nil {
		err = fmt.Errorf("failed to create icon file > %w", err)
		return
	}

	if changed, err = de.createEntry(); err != nil {
		err = fmt.Errorf("failed to create or update desktop entry file > %w", err)
		return
	}

	if changed && de.RerunIfChanged {
		err = de.restart()
	}

	return
}

func (de *DesktopEntry) createPaths() (err error) {
	for _, path := range []string{de.AppsPath, de.IconsPath} {
		if _, err = os.Stat(path); os.IsNotExist(err) {
			if err = os.MkdirAll(path, de.Perm); err != nil {
				return
			}
			err = nil
		} else if err != nil {
			return
		}
	}

	return
}

func (de *DesktopEntry) createIcon() (err error) {
	var iconPath = de.getIconPath()

	if _, err = os.Stat(iconPath); !os.IsNotExist(err) {
		return
	}

	return os.WriteFile(iconPath, de.Icon, de.Perm)
}

func (de *DesktopEntry) getIconPath() string {
	return filepath.Join(de.IconsPath, de.getID()+".png")
}

func (de *DesktopEntry) getID() string {
	return fmt.Sprintf("%s-%s", de.Name, de.Version)
}

func (de *DesktopEntry) createEntry() (changed bool, err error) {
	var entryPath = filepath.Join(de.AppsPath, de.getID()+".desktop")
	var entryData string
	var doUpdate = de.UpdateIfChanged

	if _, err = os.Stat(entryPath); err != nil && !os.IsNotExist(err) {
		return
	}

	if _, err = os.Stat(entryPath); err == nil && doUpdate {
		if doUpdate, err = de.shouldUpdate(entryPath); err != nil {
			return
		}
	} else if !os.IsNotExist(err) {
		return
	}

	if doUpdate {
		if entryData, err = de.getEntryContent(); err != nil {
			return
		}

		return true, os.WriteFile(entryPath, []byte(entryData), de.Perm)
	}

	return
}

func (de *DesktopEntry) shouldUpdate(entryPath string) (yes bool, err error) {
	var entryFile *os.File
	var execRegex, classRegex *regexp.Regexp
	var existingData []byte
	var execLine string

	if execRegex, err = regexp.Compile("Exec=sh -c '.*'"); err != nil {
		return
	}

	if classRegex, err = regexp.Compile("StartupWMClass=.*"); err != nil {
		return
	}

	if entryFile, err = os.Open(entryPath); err != nil {
		return
	}
	defer entryFile.Close()

	if existingData, err = io.ReadAll(entryFile); err != nil {
		return
	}

	if execLine, err = de.getExecLine(); err != nil {
		return
	}

	if match := execRegex.Find(existingData); match == nil || string(match) != execLine {
		yes = true
	}

	if match := classRegex.Find(existingData); match == nil || string(match) != de.getStartupClassLine() {
		yes = true
	}

	return
}

func (de *DesktopEntry) getExecLine() (execPath string, err error) {
	if execPath, err = os.Executable(); err != nil {
		return
	}

	return fmt.Sprintf("Exec=sh -c '%s'", execPath), nil
}

func (de *DesktopEntry) getStartupClassLine() string {
	return "StartupWMClass=" + filepath.Base(os.Args[0])
}

func (de *DesktopEntry) getEntryContent() (content string, err error) {
	var execLine string

	if execLine, err = de.getExecLine(); err != nil {
		return
	}

	lines := []string{
		"[Desktop Entry]",
		"Type=" + de.Type,
		"Name=" + de.Name,
		execLine,
		"Icon=" + de.getIconPath(),
		de.getStartupClassLine(),
	}

	if de.Categories != "" {
		lines = append(lines, "Categories="+de.Categories)
	}

	if de.Comment != "" {
		lines = append(lines, "Comment="+de.Comment)
	}

	return strings.Join(lines, "\n"), nil
}

func (de *DesktopEntry) restart() (err error) {
	var cmd *exec.Cmd

	if len(os.Args) > 1 {
		cmd = exec.Command(os.Args[0], os.Args[1:]...)
	} else {
		cmd = exec.Command(os.Args[0])
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err = cmd.Run(); err != nil {
		return
	}

	os.Exit(0)
	return
}
