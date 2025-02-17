package sysinfo

import (
	"os"
	"runtime"
	"strings"
)

type User struct {
	shell string
}

func NewUser() *User {
	return &User{
		shell: GetShell(),
	}
}

func (u *User) GetUserShell() string {
	return u.shell
}

func GetShell() string {
	// define a lookup table for popular shells and terminal emulators
	terminals := map[string]string{
		"cmd.exe":          "cmd",
		"powershell.exe":   "powershell",
		"Windows Terminal": "windows_terminal",
		"Git Bash":         "git_bash",
		"bash":             "bash",
		"zsh":              "zsh",
		"fish":             "fish",
		"sh":               "sh",
		"iTerm":            "iterm",
		"Terminal":         "mac_terminal",
	}

	// windows-based terminals
	if runtime.GOOS == "windows" {
		comspec := os.Getenv("ComSpec")
		for key, value := range terminals {
			if strings.Contains(strings.ToLower(comspec), strings.ToLower(key)) {
				return value
			}
		}

		termProgram := os.Getenv("TERM_PROGRAM")
		for key, value := range terminals {
			if strings.Contains(strings.ToLower(termProgram), strings.ToLower(key)) {
				return value
			}
		}

		return "windows_terminal"
	}

	// Unix-based terminals (Linux/macOS)
	shell := os.Getenv("SHELL")
	for key, value := range terminals {
		if strings.Contains(strings.ToLower(shell), strings.ToLower(key)) {
			return value
		}
	}

	termProgram := os.Getenv("TERM_PROGRAM")
	for key, value := range terminals {
		if strings.Contains(strings.ToLower(termProgram), strings.ToLower(key)) {
			return value
		}
	}

	return "terminal"
}
