package system

import (
	"os/exec"
	"runtime"
	"strings"
)

// SystemInfo holds detailed information about the host system.
type SystemInfo struct {
	OS              string // e.g. linux, darwin, windows
	Arch            string // e.g. amd64, arm64
	KernelVersion   string
	PlatformVersion string // e.g. macOS version, Ubuntu version, etc.
	Platform        string // e.g. macOS, Ubuntu, CentOS, etc.

	Other map[string]string // for any extra info
}

// GetSystemInfo gathers as much system information as possible and returns it as a SystemInfo struct.
func GetSystemInfo() SystemInfo {
	info := SystemInfo{
		OS:    runtime.GOOS,
		Arch:  runtime.GOARCH,
		Other: make(map[string]string),
	}

	// Try to get kernel version
	if out, err := exec.Command("uname", "-r").Output(); err == nil {
		info.KernelVersion = strings.TrimSpace(string(out))
	}

	// Try to get platform/version (macOS, Linux, etc.)
	if info.OS == "darwin" {
		if out, err := exec.Command("sw_vers", "-productName").Output(); err == nil {
			info.Platform = strings.TrimSpace(string(out))
		}
		if out, err := exec.Command("sw_vers", "-productVersion").Output(); err == nil {
			info.PlatformVersion = strings.TrimSpace(string(out))
		}
	} else if info.OS == "linux" {
		if out, err := exec.Command("lsb_release", "-ds").Output(); err == nil {
			info.Platform = strings.Trim(strings.ReplaceAll(string(out), "\"", ""), "\n ")
		} else if out, err := exec.Command("cat", "/etc/os-release").Output(); err == nil {
			lines := strings.Split(string(out), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "PRETTY_NAME=") {
					info.Platform = strings.Trim(line[len("PRETTY_NAME="):], "\" ")
				}
			}
		}
	}

	return info
}
