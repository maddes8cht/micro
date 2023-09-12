//go:build linux || darwin || dragonfly || openbsd_amd64 || freebsd || windows
// +build linux darwin dragonfly openbsd_amd64 freebsd windows

package action

import (
	"os"
	"os/exec"
	"runtime"

	shellquote "github.com/kballard/go-shellquote"
	"github.com/zyedidia/micro/v2/internal/shell"
)

// TermEmuSupported is a constant that marks if the terminal emulator is supported
const TermEmuSupported = true

// RunTermEmulator starts a terminal emulator from a bufpane with the given input (command)
// if wait is true it will wait for the user to exit by pressing enter once the executable has terminated
// if getOutput is true it will redirect the stdout of the process to a pipe which will be passed to the
// callback which is a function that takes a string and a list of optional user arguments
func RunTermEmulator(h *BufPane, input string, wait bool, getOutput bool, callback func(out string, userargs []interface{}), userargs []interface{}) error {
	args, err := shellquote.Split(input)
	if err != nil {
		return err
	}
	if len(args) == 0 {
		return nil
	}

	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", input)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err := cmd.Run()
		if err != nil {
			return err
		}
	} else {
		t := new(shell.Terminal)
		err = t.Start(args, getOutput, wait, callback, userargs)
		if err != nil {
			return err
		}

		h.AddTab()
		id := MainTab().Panes[0].ID()

		v := h.GetView()

		tp, err := NewTermPane(v.X, v.Y, v.Width, v.Height, t, id, MainTab())
		if err != nil {
			return err
		}
		MainTab().Panes[0] = tp
		MainTab().SetActive(0)
	}

	return nil
}
