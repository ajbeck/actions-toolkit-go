package core

import (
  "fmt"
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
func GetWindowsOsInfo() (OsVersionInfo, error) {
  windowsVersionCmd := exec.Command("powershell", "-command", "\"(Get-CimInstance -ClassName Win32_OperatingSystem).Version\"")
  versionStr, err := windowsVersionCmd.Output()
  if err != nil {
    return OsVersionInfo{}, fmt.Errorf("error getting Windows OS Version Info: %w", err)
  }

  windowsNameCmd := exec.Command("powershell", "-command", "\"(Get-CimInstance -ClassName Win32_OperatingSystem).Caption\"")
  name, err := windowsNameCmd.Output()
  if err != nil {
    return OsVersionInfo{}, fmt.Errorf("error getting Windows OS Name Info: %w", err)
  }

  return OsVersionInfo{
    Name:    strings.TrimSpace(string(name)),
    Version: strings.TrimSpace(string(versionStr)),
  }, nil
}

// GetMacOsInfo retrieves the macOS operating system name and version.
// The function executes sw_vers commands to get the OS details.
func GetMacOsInfo() (OsVersionInfo, error) {
  nameInfoCmd := exec.Command("sw_vers", "-productName")
  nameStr, err := nameInfoCmd.Output()
  if err != nil {
    return OsVersionInfo{}, fmt.Errorf("error getting macOS Name Info: %w", err)
  }

  versionInfoCmd := exec.Command("sw_vers", "-productVersion")
  versionStr, err := versionInfoCmd.Output()
  if err != nil {
    return OsVersionInfo{}, fmt.Errorf("error getting macOS Version Info: %w", err)
  }

  return OsVersionInfo{
    Name:    strings.TrimSpace(string(nameStr)),
    Version: strings.TrimSpace(string(versionStr)),
  }, nil
}

// GetLinuxOsInfo retrieves the Linux operating system name and version.
// The function executes lsb_release commands to get the OS details.
func GetLinuxOsInfo() (OsVersionInfo, error) {
  infoCmd := exec.Command("lsb_release", "-i", "-r", "-s")
  infoStr, err := infoCmd.Output()
  if err != nil {
    return OsVersionInfo{}, fmt.Errorf("error getting Linux OS Info: %w", err)
  }

  parts := strings.Split(strings.TrimSpace(string(infoStr)), "\n")

  return OsVersionInfo{
    Name:    strings.TrimSpace(parts[0]),
    Version: strings.TrimSpace(parts[1]),
  }, nil
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
func GetPlatformInfo() (PlatformInfo, error) {
  if runtime.GOOS == "windows" {
    winInfo, err := GetWindowsOsInfo()

    if err != nil {
      return PlatformInfo{}, err
    }

    return PlatformInfo{
      Name:      winInfo.Name,
      Platform:  "windows",
      Arch:      runtime.GOARCH,
      Version:   winInfo.Version,
      isWindows: true,
    }, nil
  } else if runtime.GOOS == "darwin" {
    macInfo, err := GetMacOsInfo()

    if err != nil {
      return PlatformInfo{}, err
    }

    return PlatformInfo{
      Name:     macInfo.Name,
      Platform: "macos",
      Arch:     runtime.GOARCH,
      Version:  macInfo.Version,
      isMacOS:  true,
    }, nil
  } else {
    linuxInfo, err := GetLinuxOsInfo()

    if err != nil {
      return PlatformInfo{}, err
    }

    return PlatformInfo{
      Name:     linuxInfo.Name,
      Platform: "linux",
      Arch:     runtime.GOARCH,
      Version:  linuxInfo.Version,
      isLinux:  true,
    }, nil
  }
}
