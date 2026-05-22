package cmd

import (
	"reflect"
	"testing"

	"github.com/marcosnils/bin/pkg/config"
	"github.com/marcosnils/bin/pkg/providers"
)

type mockProvider struct {
	providers.Provider
	latestVersion    string
	latestVersionURL string
	err              error
	gotCooldown      *int
}

func (m *mockProvider) GetLatestVersion(opts *providers.LatestVersionOpts) (string, string, error) {
	if opts != nil {
		days := opts.CooldownPeriodDays
		m.gotCooldown = &days
	}
	return m.latestVersion, m.latestVersionURL, m.err
}

func (m *mockProvider) GetID() string {
	return "github"
}

func TestGetLatestVersion(t *testing.T) {
	type mockValues struct {
		latestVersion    string
		latestVersionURL string
		err              error
	}
	cases := []struct {
		in  *config.Binary
		m   mockValues
		out *updateInfo
	}{
		{
			&config.Binary{
				Path:       "/home/user/bin/launchpad",
				Version:    "1.1.0",
				URL:        "https://github.com/Mirantis/launchpad/releases/download/1.1.0/launchpad-linux-x64",
				RemoteName: "launchpad-linux-x64",
				Provider:   "github",
			},
			mockValues{"1.1.1", "https://github.com/Mirantis/launchpad/releases/download/1.1.1/launchpad-linux-x64", nil},
			&updateInfo{
				version: "1.1.1",
				url:     "https://github.com/Mirantis/launchpad/releases/download/1.1.1/launchpad-linux-x64",
			},
		},
		{
			&config.Binary{
				Path:       "/home/user/bin/launchpad",
				Version:    "1.2.0-rc.1",
				URL:        "https://github.com/Mirantis/launchpad/releases/download/1.2.0-rc.1/launchpad-linux-x64",
				RemoteName: "launchpad-linux-x64",
				Provider:   "github",
			},
			mockValues{"1.1.1", "https://github.com/Mirantis/launchpad/releases/download/1.1.1/launchpad-linux-x64", nil},
			nil,
		},
	}

	for _, c := range cases {
		p := &mockProvider{latestVersion: c.m.latestVersion, latestVersionURL: c.m.latestVersionURL, err: c.m.err}
		if v, err := getLatestVersion(c.in, p, 0); err != nil {
			t.Fatalf("Error during getLatestVersion(%#v, %#v): %v", c.in, p, err)
		} else if !reflect.DeepEqual(v, c.out) {
			t.Fatalf("For case %#v: %#v does not match %#v", c.in, v, c.out)
		}
	}

}

func TestGetLatestVersionPassesCooldown(t *testing.T) {
	p := &mockProvider{latestVersion: "1.1.1", latestVersionURL: "https://example.test/v1.1.1"}
	b := &config.Binary{Path: "/home/user/bin/launchpad", Version: "1.1.0", URL: p.latestVersionURL, Provider: "github"}
	if _, err := getLatestVersion(b, p, 7); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.gotCooldown == nil {
		t.Fatalf("expected provider to receive LatestVersionOpts, got nil")
	}
	if *p.gotCooldown != 7 {
		t.Fatalf("expected cooldown 7, got %d", *p.gotCooldown)
	}
}
