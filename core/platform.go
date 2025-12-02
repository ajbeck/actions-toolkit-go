package core

import (
	"os/exec"
	"runtime"
	"strings"
)

// OsVersionInfo holds operating system name and version information.
type OsVersionInfo struct {
	Name    string
	Version string
}

// GetWindowsOsInfo retrieves the Windows operating system name and version.
// The function executes PowerShell commands to get the OS details.
func GetWindowsOsInfo() OsVersionInfo {
	windowsVersionCmd := exec.Command("powershell", "-command", "\"(Get-CimInstance -ClassName Win32_OperatingSystem).Version\"")
	versionStr, _ := windowsVersionCmd.Output()

	windowsNameCmd := exec.Command("powershell", "-command", "\"(Get-CimInstance -ClassName Win32_OperatingSystem).Caption\"")
	name, _ := windowsNameCmd.Output()

	return OsVersionInfo{
		Name:    strings.TrimSpace(string(name)),
		Version: strings.TrimSpace(string(versionStr)),
	}
}

// GetMacOsInfo retrieves the macOS operating system name and version.
// The function executes sw_vers commands to get the OS details.
func GetMacOsInfo() OsVersionInfo {
	nameInfoCmd := exec.Command("sw_vers", "-productName")
	nameStr, _ := nameInfoCmd.Output()

	versionInfoCmd := exec.Command("sw_vers", "-productVersion")
	versionStr, _ := versionInfoCmd.Output()

	return OsVersionInfo{
		Name:    strings.TrimSpace(string(nameStr)),
		Version: strings.TrimSpace(string(versionStr)),
	}
}

// GetLinuxOsInfo retrieves the Linux operating system name and version.
// The function executes lsb_release commands to get the OS details.
func GetLinuxOsInfo() OsVersionInfo {
	infoCmd := exec.Command("lsb_release", "-i", "-r", "-s")
	infoStr, _ := infoCmd.Output()

	parts := strings.Split(strings.TrimSpace(string(infoStr)), "\n")

	return OsVersionInfo{
		Name:    strings.TrimSpace(parts[0]),
		Version: strings.TrimSpace(parts[1]),
	}
}

// PlatformInfo holds detailed information about the current platform.
type PlatformInfo struct {
	Name      string
	Platform  string
	Arch      string
	Version   string
	isWindows bool
	isMacOS   bool
	isLinux   bool
}

// GetPlatformInfo retrieves detailed information about the current platform,
// including name, platform type, architecture, and version.
func GetPlatformInfo() PlatformInfo {
	if runtime.GOOS == "windows" {
		winInfo := GetWindowsOsInfo()
		return PlatformInfo{
			Name:      winInfo.Name,
			Platform:  "windows",
			Arch:      runtime.GOARCH,
			Version:   winInfo.Version,
			isWindows: true,
		}
	} else if runtime.GOOS == "darwin" {
		macInfo := GetMacOsInfo()
		return PlatformInfo{
			Name:     macInfo.Name,
			Platform: "macos",
			Arch:     runtime.GOARCH,
			Version:  macInfo.Version,
			isMacOS:  true,
		}
	} else {
		linuxInfo := GetLinuxOsInfo()
		return PlatformInfo{
			Name:     linuxInfo.Name,
			Platform: "linux",
			Arch:     runtime.GOARCH,
			Version:  linuxInfo.Version,
			isLinux:  true,
		}
	}
}
