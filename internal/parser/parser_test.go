package parser

import (
	"testing"

	"github.com/lichenglife/issue2md/internal/models"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      *models.ParsedURL
		wantErr   bool
		errContains string
	}{
		{
			name:  "valid issue url",
			input: "https://github.com/owner/repo/issues/123",
			want: &models.ParsedURL{
				Owner:  "owner",
				Repo:   "repo",
				Type:   "issue",
				Number: 123,
			},
			wantErr: false,
		},
		{
			name:  "valid pull url",
			input: "https://github.com/owner/repo/pull/456",
			want: &models.ParsedURL{
				Owner:  "owner",
				Repo:   "repo",
				Type:   "pull",
				Number: 456,
			},
			wantErr: false,
		},
		{
			name:  "valid discussion url",
			input: "https://github.com/owner/repo/discussions/789",
			want: &models.ParsedURL{
				Owner:  "owner",
				Repo:   "repo",
				Type:   "discussion",
				Number: 789,
			},
			wantErr: false,
		},
		{
			name:      "invalid domain",
			input:     "https://gitlab.com/owner/repo/issues/123",
			want:      nil,
			wantErr:   true,
			errContains: "域名必须是 github.com",
		},
		{
			name:      "invalid url format",
			input:     "not a url",
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "incomplete path",
			input:     "https://github.com/owner/repo",
			want:      nil,
			wantErr:   true,
			errContains: "路径不完整",
		},
		{
			name:      "unsupported resource type",
			input:     "https://github.com/owner/repo/wiki/Page",
			want:      nil,
			wantErr:   true,
			errContains: "不支持的资源类型",
		},
		{
			name:      "invalid number",
			input:     "https://github.com/owner/repo/issues/abc",
			want:      nil,
			wantErr:   true,
			errContains: "无效的数字",
		},
		{
			name:  "url with trailing slash",
			input: "https://github.com/owner/repo/issues/123/",
			want: &models.ParsedURL{
				Owner:  "owner",
				Repo:   "repo",
				Type:   "issue",
				Number: 123,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)

			if (err != nil) != tt.wantErr {
				t.Fatalf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && tt.errContains != "" {
				if err == nil || !contains(err.Error(), tt.errContains) {
					t.Errorf("Parse() error = %v, want error containing %q", err, tt.errContains)
				}
			}

			if !tt.wantErr && got == nil {
				t.Fatal("Parse() returned nil for valid input")
			}

			if !tt.wantErr {
				if got.Owner != tt.want.Owner || got.Repo != tt.want.Repo ||
					got.Type != tt.want.Type || got.Number != tt.want.Number {
					t.Errorf("Parse() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"valid issue", "https://github.com/owner/repo/issues/123", true},
		{"valid pull", "https://github.com/owner/repo/pull/456", true},
		{"valid discussion", "https://github.com/owner/repo/discussions/789", true},
		{"invalid domain", "https://gitlab.com/owner/repo/issues/123", false},
		{"invalid url", "not a url", false},
		{"incomplete", "https://github.com/owner/repo", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValid(tt.input)
			if got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
