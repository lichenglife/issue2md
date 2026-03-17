package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		setupEnv    func(t *testing.T)
		wantOutputDir string
		wantUserLinks bool
		wantVerbose bool
		wantToken   string
	}{
		{
			name:        "default values",
			args:        []string{},
			wantOutputDir: ".",
			wantUserLinks: false,
			wantVerbose:   false,
			wantToken:     "",
		},
		{
			name:        "custom output dir",
			args:        []string{"-o", "/tmp/output"},
			wantOutputDir: "/tmp/output",
			wantUserLinks: false,
			wantVerbose:   false,
			wantToken:     "",
		},
		{
			name:        "user links enabled",
			args:        []string{"-user-links"},
			wantOutputDir: ".",
			wantUserLinks: true,
			wantVerbose:   false,
			wantToken:     "",
		},
		{
			name:        "verbose mode",
			args:        []string{"-v"},
			wantOutputDir: ".",
			wantUserLinks: false,
			wantVerbose:   true,
			wantToken:     "",
		},
		{
			name:        "token from command line",
			args:        []string{"-token", "ghp_test123"},
			wantOutputDir: ".",
			wantUserLinks: false,
			wantVerbose:   false,
			wantToken:     "ghp_test123",
		},
		{
			name:        "token from environment",
			args:        []string{},
			setupEnv:    func(t *testing.T) { os.Setenv("GITHUB_TOKEN", "ghp_env456") },
			wantOutputDir: ".",
			wantUserLinks: false,
			wantVerbose:   false,
			wantToken:     "ghp_env456",
		},
		{
			name:        "command line token overrides environment",
			args:        []string{"-token", "ghp_cli789"},
			setupEnv:    func(t *testing.T) { os.Setenv("GITHUB_TOKEN", "ghp_env000") },
			wantOutputDir: ".",
			wantUserLinks: false,
			wantVerbose:   false,
			wantToken:     "ghp_cli789",
		},
		{
			name:        "all options combined",
			args:        []string{"-o", "/tmp/md", "-user-links", "-v", "-token", "ghp_all"},
			wantOutputDir: "/tmp/md",
			wantUserLinks: true,
			wantVerbose:   true,
			wantToken:     "ghp_all",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean environment before each test
			os.Unsetenv("GITHUB_TOKEN")

			if tt.setupEnv != nil {
				tt.setupEnv(t)
			}
			defer os.Unsetenv("GITHUB_TOKEN")

			got, err := Load(tt.args)
			if err != nil {
				t.Fatalf("Load() error = %v", err)
			}

			if got.OutputDir != tt.wantOutputDir {
				t.Errorf("OutputDir = %q, want %q", got.OutputDir, tt.wantOutputDir)
			}
			if got.UserLinks != tt.wantUserLinks {
				t.Errorf("UserLinks = %v, want %v", got.UserLinks, tt.wantUserLinks)
			}
			if got.Verbose != tt.wantVerbose {
				t.Errorf("Verbose = %v, want %v", got.Verbose, tt.wantVerbose)
			}
			if got.Token != tt.wantToken {
				t.Errorf("Token = %q, want %q", got.Token, tt.wantToken)
			}
		})
	}
}

func TestGetTokenFromEnv(t *testing.T) {
	tests := []struct {
		name     string
		setupEnv func(t *testing.T)
		want     string
	}{
		{
			name:     "token set",
			setupEnv: func(t *testing.T) { os.Setenv("GITHUB_TOKEN", "ghp_test") },
			want:     "ghp_test",
		},
		{
			name:     "token not set",
			setupEnv: func(t *testing.T) { os.Unsetenv("GITHUB_TOKEN") },
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv(t)
			defer os.Unsetenv("GITHUB_TOKEN")

			got := GetTokenFromEnv()
			if got != tt.want {
				t.Errorf("GetTokenFromEnv() = %q, want %q", got, tt.want)
			}
		})
	}
}
