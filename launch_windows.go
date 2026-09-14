//go:build windows

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
	"unsafe"
)

// Prefer a Chromium-based browser in --app mode: gives a clean borderless
// window with no address bar, and its own profile so it never touches the
// user's normal browsing session.
func browserPaths() []string {
	pf := os.Getenv("ProgramFiles")
	pf86 := os.Getenv("ProgramFiles(x86)")
	local := os.Getenv("LOCALAPPDATA")
	return []string{
		filepath.Join(pf86, `Microsoft\Edge\Application\msedge.exe`),
		filepath.Join(pf, `Microsoft\Edge\Application\msedge.exe`),
		filepath.Join(local, `Microsoft\Edge\Application\msedge.exe`),
		filepath.Join(pf86, `Google\Chrome\Application\chrome.exe`),
		filepath.Join(pf, `Google\Chrome\Application\chrome.exe`),
		filepath.Join(local, `Google\Chrome\Application\chrome.exe`),
	}
}

// browserCmd is the game window, so an update can close it before the
// installer tries to overwrite files.
var browserCmd *exec.Cmd

// runInstallerAndExit starts the downloaded installer in silent mode and quits.
// The installer waits a moment for this process to die, replaces the files and
// relaunches the game itself.
func runInstallerAndExit(path string) {
	c := exec.Command(path, "/S")
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := c.Start(); err != nil {
		msgBox("The update could not be started:\n\n" + err.Error())
		return
	}
	if browserCmd != nil && browserCmd.Process != nil {
		browserCmd.Process.Kill()
	}
	os.Exit(0)
}

func main() {
	url, err := startServer()
	if err != nil {
		msgBox("Could not start the local game server:\n\n" + err.Error())
		return
	}

	go checkForUpdate() // never blocks startup; offline is fine

	profile := filepath.Join(os.TempDir(), appName+"Profile")
	var cmd *exec.Cmd
	for _, p := range browserPaths() {
		if _, err := os.Stat(p); err != nil {
			continue
		}
		c := exec.Command(p,
			"--app="+url,
			"--user-data-dir="+profile,
			"--window-size=1320,880",
			"--no-first-run",
			"--no-default-browser-check",
			"--disable-features=Translate,TranslateUI",
		)
		c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		if err := c.Start(); err == nil {
			cmd = c
			browserCmd = c
			break
		}
	}

	if cmd == nil {
		// No Chromium-based browser found: fall back to the default browser.
		exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	}

	// Wait on the PAGE, not on the browser process. Edge exits the process we
	// launched whenever it hands off to an instance already using this profile,
	// which used to shut the server down while the game was still open - every
	// later save then failed silently and the player lost their progress.
	WaitForPageGone(14*time.Second, 90*time.Second)
	time.Sleep(1200 * time.Millisecond) // let a final save land
}

func msgBox(msg string) {
	mb := syscall.NewLazyDLL("user32.dll").NewProc("MessageBoxW")
	t, _ := syscall.UTF16PtrFromString(appName)
	m, _ := syscall.UTF16PtrFromString(msg)
	mb.Call(0, uintptr(unsafe.Pointer(m)), uintptr(unsafe.Pointer(t)), 0x10)
}
