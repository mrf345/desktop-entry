package desktopEntry

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func getExecLine() (execPath string, err error) {
	if execPath, err = os.Executable(); err != nil {
		return
	}

	return fmt.Sprintf("Exec=sh -c '%s %%F'", execPath), nil
}

func getStartupClassLine() string {
	return "StartupWMClass=" + filepath.Base(os.Args[0])
}

func restart() (err error) {
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
