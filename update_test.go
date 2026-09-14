package main

import "testing"

func TestNewerThan(t *testing.T) {
	cases := []struct {
		a, b string
		want bool
	}{
		{"1.2.0", "1.1.0", true},
		{"1.1.0", "1.2.0", false},
		{"1.2.0", "1.2.0", false}, // same version must not prompt
		{"1.10.0", "1.9.0", true}, // numeric, not lexical
		{"2.0.0", "1.99.99", true},
		{"v1.3.0", "1.2.0", true}, // tolerate a leading v
		{"1.2.1", "1.2.0", true},
		{"1.2.0.1", "1.2.0", true},
		{"", "1.0.0", false}, // junk feed must never trigger an update
		{"garbage", "1.0.0", false},
	}
	for _, c := range cases {
		if got := newerThan(c.a, c.b); got != c.want {
			t.Errorf("newerThan(%q,%q) = %v, want %v", c.a, c.b, got, c.want)
		}
	}
}

func TestIsSafeDownload(t *testing.T) {
	bad := []string{
		"http://evil.test/x.exe", // plaintext
		"ftp://evil.test/x.exe",
		"file:///C:/Windows/evil.exe", // local file
		"javascript:alert(1)",
		"",
	}
	for _, u := range bad {
		if isSafeDownload(u) {
			t.Errorf("isSafeDownload(%q) = true, must be false", u)
		}
	}
	if !isSafeDownload("https://github.com/a/b/releases/latest/download/Setup.exe") {
		t.Error("https github url should be allowed")
	}
}
