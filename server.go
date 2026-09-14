package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"
)

//go:embed game.html
var assets embed.FS

const appName = "CaseClicker"

func saveFile() string {
	base, err := os.UserConfigDir()
	if err != nil || base == "" {
		base, _ = os.Getwd()
	}
	dir := filepath.Join(base, appName)
	os.MkdirAll(dir, 0o755)
	return filepath.Join(dir, "save.json")
}

// lastPing is the last time the game page told us it was still open, as a Unix
// milli timestamp. The page is the authority on its own lifetime: we cannot
// rely on the browser process we spawned, because Chromium/Edge hand off to an
// already-running instance that shares the same --user-data-dir and then exit
// immediately, which would look identical to the window being closed.
var lastPing atomic.Int64
var everSeen atomic.Bool

func touch() {
	lastPing.Store(time.Now().UnixMilli())
	everSeen.Store(true)
}

// WaitForPageGone blocks until the page has stopped checking in. graceStart is
// how long to wait for the window to appear at all before giving up.
func WaitForPageGone(silence, graceStart time.Duration) {
	start := time.Now()
	for {
		time.Sleep(400 * time.Millisecond)
		if !everSeen.Load() {
			if time.Since(start) > graceStart {
				return // the window never opened; nothing to serve
			}
			continue
		}
		if time.Since(time.UnixMilli(lastPing.Load())) > silence {
			return // window closed (or crashed) - stop serving
		}
	}
}

// startServer binds a free loopback port, serves the game and the save API,
// and returns the URL to open.
func startServer() (string, error) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", err
	}
	port := ln.Addr().(*net.TCPAddr).Port

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		b, _ := assets.ReadFile("game.html")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		w.Write(b)
	})
	// the app window asks for a favicon; answer it so the console stays clean
	// ---- update endpoints ----
	mux.HandleFunc("/alive", func(w http.ResponseWriter, r *http.Request) {
		touch()
		w.Header().Set("Cache-Control", "no-store")
		w.Write([]byte("ok"))
	})
	mux.HandleFunc("/update/status", func(w http.ResponseWriter, r *http.Request) {
		upd.mu.Lock()
		out := map[string]any{"checked": upd.Checked, "available": upd.Available,
			"current": upd.Current, "busy": upd.Busy, "release": upd.Release}
		upd.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(out)
	})
	mux.HandleFunc("/update/skip", func(w http.ResponseWriter, r *http.Request) {
		upd.mu.Lock()
		if upd.Release != nil {
			os.WriteFile(skipFile(), []byte(upd.Release.Version), 0o644)
		}
		upd.Available = false
		upd.mu.Unlock()
		w.Write([]byte("ok"))
	})
	mux.HandleFunc("/update/apply", func(w http.ResponseWriter, r *http.Request) {
		upd.mu.Lock()
		rel, busy := upd.Release, upd.Busy
		if !busy {
			upd.Busy = true
		}
		upd.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		if rel == nil || busy {
			json.NewEncoder(w).Encode(map[string]string{"error": "nothing to install"})
			return
		}
		path, err := downloadInstaller(rel)
		if err != nil {
			upd.mu.Lock()
			upd.Busy = false
			upd.Error = err.Error()
			upd.mu.Unlock()
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"ok": "installing"})
		// Hand off after the response is flushed: the installer needs this
		// process (and its window) gone before it can replace the .exe.
		go func() { time.Sleep(500 * time.Millisecond); runInstallerAndExit(path) }()
	})
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Write([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32"><rect width="32" height="32" rx="5" fill="#1b1f28"/><path d="M5 11l11-5 11 5v13l-11 5-11-5z" fill="#4a5568"/><path d="M5 11l11 5 11-5-11-5z" fill="#5d6b80"/><rect x="13" y="12" width="6" height="7" rx="1" fill="#f0b90b"/></svg>`))
	})
	mux.HandleFunc("/load", func(w http.ResponseWriter, r *http.Request) {
		b, err := os.ReadFile(saveFile())
		if err != nil {
			http.Error(w, "no save", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	})
	// lets the UI show the real save location on this machine
	mux.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"savePath": saveFile()})
	})
	mux.HandleFunc("/save", func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(io.LimitReader(r.Body, 32<<20))
		if err != nil || len(b) == 0 {
			http.Error(w, "bad", http.StatusBadRequest)
			return
		}
		// write to a temp file first so a crash mid-write can't corrupt the save
		tmp := saveFile() + ".tmp"
		if err := os.WriteFile(tmp, b, 0o644); err == nil {
			os.Rename(tmp, saveFile())
		}
		w.Write([]byte("ok"))
	})

	// Any request at all proves the window is still there.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		touch()
		mux.ServeHTTP(w, r)
	})
	go (&http.Server{Handler: handler}).Serve(ln)
	return fmt.Sprintf("http://127.0.0.1:%d/", port), nil
}
