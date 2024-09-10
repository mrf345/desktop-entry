package desktopEntry_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	desktopEntry "github.com/mrf345/desktop-entry"
	"github.com/stretchr/testify/assert"
)

func TestCreateEntryInDev(t *testing.T) {
	assert := assert.New(t)
	testEntry, _, _, does_content_match := getTestEntry()
	defer os.RemoveAll(testEntry.AppsPath)

	os.Args[0] = "/tmp"
	err := testEntry.Create()

	assert.Nil(err)
	assert.False(does_content_match())
}

func TestCreateEntry(t *testing.T) {
	assert := assert.New(t)
	testEntry, iconPath, entryPath, doesContentMatch := getTestEntry()
	defer os.RemoveAll(testEntry.AppsPath)

	os.Args[0] = "/test-dir"
	err := testEntry.Create()
	_, iconExistErr := os.Stat(iconPath)
	_, entryExistErr := os.Stat(entryPath)

	assert.Nil(err)
	assert.Nil(iconExistErr)
	assert.Nil(entryExistErr)
	assert.True(doesContentMatch())
}

func TestNotUpdatingExistingEntry(t *testing.T) {
	assert := assert.New(t)
	testEntry, iconPath, entryPath, does_content_match := getTestEntry()
	testEntry.UpdateIfChanged = false
	entryFile, _ := os.Create(entryPath)
	entryContent := []byte("not changed")
	_, _ = entryFile.Write(entryContent)
	entryFile.Close()
	defer os.RemoveAll(testEntry.AppsPath)
	defer os.Remove(entryFile.Name())

	os.Args[0] = "/test-dir"
	err := testEntry.Create()
	_, iconExistErr := os.Stat(iconPath)
	_, entryExistErr := os.Stat(entryPath)
	content, _ := os.ReadFile(entryFile.Name())

	assert.Nil(err)
	assert.Nil(iconExistErr)
	assert.Nil(entryExistErr)
	assert.Equal(entryContent, content)
	assert.False(does_content_match())
}

func TestUpdateExistingEntry(t *testing.T) {
	assert := assert.New(t)
	testEntry, iconPath, entryPath, does_content_match := getTestEntry()
	entryFile, _ := os.Create(entryPath)
	entryContent := []byte("not changed")
	_, _ = entryFile.Write(entryContent)
	entryFile.Close()
	defer os.RemoveAll(testEntry.AppsPath)
	defer os.Remove(entryFile.Name())

	os.Args[0] = "/test-dir"
	err := testEntry.Create()
	_, iconExistErr := os.Stat(iconPath)
	_, entryExistErr := os.Stat(entryPath)
	content, _ := os.ReadFile(entryFile.Name())

	assert.Nil(err)
	assert.Nil(iconExistErr)
	assert.Nil(entryExistErr)
	assert.NotEqual(entryContent, content)
	assert.True(does_content_match())
}

func getTestEntry() (
	entry *desktopEntry.DesktopEntry,
	iconPath,
	entryPath string,
	doesContentMatch func() bool,
) {
	var tempDir string
	var err error

	if tempDir, err = os.MkdirTemp("", "desktop-entry"); err != nil {
		panic(err)
	}

	name := "testing"
	version := "0.0.1"
	id := name
	entry = desktopEntry.New(
		name,
		version,
		[]byte{},
	)
	entry.RerunIfChanged = false
	entry.AppsPath = tempDir
	entry.IconsPath = tempDir

	entryPath = filepath.Join(tempDir, id+".desktop")
	iconPath = filepath.Join(tempDir, id+".png")

	doesContentMatch = func() (match bool) {
		rawContent, _ := os.ReadFile(entryPath)
		content := string(rawContent)
		lines := strings.Split(content, "\n")
		execPath, _ := os.Executable()
		execLin := fmt.Sprintf("Exec=sh -c '%s %%F'", execPath)

		return (lines[0] == "[Desktop Entry]" &&
			strings.Contains(content, "Type="+entry.Type) &&
			strings.Contains(content, "Name="+entry.Name) &&
			strings.Contains(content, "Icon="+iconPath) &&
			strings.Contains(content, "StartupWMClass="+filepath.Base(os.Args[0])) &&
			strings.Contains(content, execLin))
	}

	return
}
