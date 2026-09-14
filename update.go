package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// appVersion must match the VERSION in installer.nsi.
const appVersion = "1.3.1"

// How the update feed is located, in order:
//  1. update.txt sitting next to the .exe
//  2. the updateFeedDefault constant below
//
// The file may contain either "owner/repo" (GitHub Releases) or a full
// https:// URL to an update.json. Empty or missing means updates are off.
const updateFeedDefault = ""

// Release describes the newest build, as published in update.json.
type Release struct {
	Version string `json:"version"` // e.g. "1.3.0"
	Notes   string `json:"notes"`   // shown to the player
	URL     string `json:"url"`     // https:// link to the installer .exe
	SHA256  string `json:"sha256"`  // hex digest of that installer
	Date    string `json:"date"`
}

type updateState struct {
	mu        sync.Mutex
	Checked   bool     `json:"checked"`
	Available bool     `json:"available"`
	Current   string   `json:"current"`
	Busy      bool     `json:"busy"`
	Error     string   `json:"error,omitempty"`
	Release   *Release `json:"release,omitempty"`
}

var upd = &updateState{Current: appVersion}

// ---------------------------------------------------------------- feed config

func exeDir() string {
	p, err := os.Executable()
	if err != nil {
		return "."
	}
	return filepath.Dir(p)
}

func feedURL() string {
	cfg := updateFeedDefault
	if b, err := os.ReadFile(filepath.Join(exeDir(), "update.txt")); err == nil {
		for _, line := range strings.Split(string(b), "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue // allow comments in the file
			}
			cfg = line
			break
		}
	}
	cfg = strings.TrimSpace(cfg)
	if cfg == "" {
		return ""
	}
	if strings.HasPrefix(cfg, "https://") {
		return cfg
	}
	if strings.HasPrefix(cfg, "http://") {
		return "" // refuse plaintext: an update feed must be authenticated transport
	}
	// "owner/repo" -> GitHub's permanent "latest release asset" URL
	if parts := strings.Split(cfg, "/"); len(parts) == 2 && parts[0] != "" && parts[1] != "" {
		return "https://github.com/" + cfg + "/releases/latest/download/update.json"
	}
	return ""
}

// ---------------------------------------------------------------- versions

// newerThan reports whether version a is strictly newer than b.
// Handles plain dotted numeric versions; anything unparsable sorts as 0.
func newerThan(a, b string) bool {
	pa, pb := splitVer(a), splitVer(b)
	for i := 0; i < 4; i++ {
		if pa[i] != pb[i] {
			return pa[i] > pb[i]
		}
	}
	return false
}

func splitVer(v string) [4]int {
	var out [4]int
	v = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(v), "v"))
	for i, p := range strings.Split(v, ".") {
		if i > 3 {
			break
		}
		n := 0
		for _, r := range p {
			if r < '0' || r > '9' {
				break
			}
			n = n*10 + int(r-'0')
		}
		out[i] = n
	}
	return out
}

// ---------------------------------------------------------------- checking

var httpClient = &http.Client{Timeout: 15 * time.Second}

func checkForUpdate() {
	feed := feedURL()
	upd.mu.Lock()
	upd.Checked = true
	upd.mu.Unlock()
	if feed == "" {
		return // updates not configured; stay silent
	}

	req, err := http.NewRequest("GET", feed, nil)
	if err != nil {
		setUpdErr(err)
		return
	}
	req.Header.Set("User-Agent", "CaseClicker/"+appVersion)
	req.Header.Set("Accept", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		setUpdErr(err) // offline is normal, never surfaced to the player
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		setUpdErr(fmt.Errorf("feed returned %d", resp.StatusCode))
		return
	}

	var rel Release
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&rel); err != nil {
		setUpdErr(err)
		return
	}
	if rel.Version == "" || rel.URL == "" {
		setUpdErr(errors.New("update.json is missing version or url"))
		return
	}
	if !isSafeDownload(rel.URL) {
		setUpdErr(errors.New("installer url is not https"))
		return
	}
	if !newerThan(rel.Version, appVersion) {
		return // already current
	}
	if v, _ := os.ReadFile(skipFile()); strings.TrimSpace(string(v)) == rel.Version {
		return // player chose to skip this one
	}

	upd.mu.Lock()
	upd.Available, upd.Release, upd.Error = true, &rel, ""
	upd.mu.Unlock()
}

func setUpdErr(err error) {
	upd.mu.Lock()
	upd.Error = err.Error()
	upd.mu.Unlock()
}

func isSafeDownload(raw string) bool {
	u, err := url.Parse(raw)
	return err == nil && u.Scheme == "https" && u.Host != ""
}

func skipFile() string { return filepath.Join(filepath.Dir(saveFile()), "skipped-version.txt") }

// ---------------------------------------------------------------- applying

// downloadInstaller fetches the new installer and refuses to return a path
// unless its SHA-256 matches the digest published in update.json. Without
// that check, anything able to answer the feed could run code as the player.
func downloadInstaller(rel *Release) (string, error) {
	if !isSafeDownload(rel.URL) {
		return "", errors.New("refusing non-https download")
	}
	if len(rel.SHA256) != 64 {
		return "", errors.New("update.json has no valid sha256 - refusing to run it")
	}
	req, _ := http.NewRequest("GET", rel.URL, nil)
	req.Header.Set("User-Agent", "CaseClicker/"+appVersion)
	client := &http.Client{Timeout: 10 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("download returned %d", resp.StatusCode)
	}

	dir := filepath.Join(os.TempDir(), appName+"Update")
	os.MkdirAll(dir, 0o755)
	dst := filepath.Join(dir, "CaseClickerSetup-"+strconv.FormatInt(time.Now().Unix(), 10)+".exe")
	f, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	h := sha256.New()
	if _, err := io.Copy(io.MultiWriter(f, h), io.LimitReader(resp.Body, 256<<20)); err != nil {
		f.Close()
		os.Remove(dst)
		return "", err
	}
	f.Close()

	got := hex.EncodeToString(h.Sum(nil))
	if !strings.EqualFold(got, rel.SHA256) {
		os.Remove(dst)
		return "", errors.New("checksum mismatch - the download was corrupted or tampered with")
	}
	return dst, nil
}
