package sysinfo

import (
	"os"
	"os/exec"
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
	// lookup table for shells
	terminals := map[string]string{
		"cmd.exe":        "cmd",
		"powershell.exe": "powershell",
		"pwsh.exe":       "powershell",
		"bash":           "bash",
		"zsh":            "zsh",
		"fish":           "fish",
		"sh":             "sh",
	}

	// windows-based shell detection
	if runtime.GOOS == "windows" {
		// first check if we are inside PowerShell
		if os.Getenv("PSModulePath") != "" {
			return "powershell"
		}

		out, err := exec.Command("wmic", "process", "get", "name").Output()
		if err == nil {
			processes := strings.ToLower(string(out))
			if strings.Contains(processes, "powershell.exe") {
				return "powershell"
			}
		}

		// check ComSpec as a fallback
		comspec := os.Getenv("ComSpec")
		for key, value := range terminals {
			if strings.Contains(strings.ToLower(comspec), strings.ToLower(key)) {
				return value
			}
		}

		return "cmd" // default to cmd if nothing else matches
	}

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
