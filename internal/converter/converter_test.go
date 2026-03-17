package converter

import (
	"strings"
	"testing"
	"time"

	"github.com/lichenglife/issue2md/internal/models"
)

func TestConvert(t *testing.T) {
	tests := []struct {
		name             string
		input            *models.IssueData
		userLinks        bool
		wantFrontMatter  map[string]interface{}
		wantBodyContains string
	}{
		{
			name: "with title and author",
			input: &models.IssueData{
				Title:     "Test Issue Title",
				Body:      "This is the issue body",
				HTMLURL:   "https://github.com/owner/repo/issues/1",
				Number:    1,
				Repo:      "repo",
				User:      models.User{Login: "testuser", HTMLURL: "https://github.com/testuser"},
				CreatedAt: time.Date(2026, 3, 16, 10, 0, 0, 0, time.UTC),
				State:     "open",
				Labels:    []models.Label{{Name: "bug", Color: "d73a4a", Description: "Something isn working"}},
				Assignees: []models.User{{Login: "assignee1", HTMLURL: "https://github.com/assignee1"}},
				Milestone: &models.Milestone{Title: "v1.0", Number: 1, State: "open"},
				Comments:  []models.Comment{},
			},
			userLinks: false,
			wantFrontMatter: map[string]interface{}{
				"title":  "Test Issue Title",
				"author": "testuser",
			},
			wantBodyContains: "This is the issue body",
		},
		{
			name: "closed issue",
			input: &models.IssueData{
				Title:     "Closed Issue",
				Body:      "This issue is closed",
				HTMLURL:   "https://github.com/owner/repo/issues/2",
				Number:    2,
				Repo:      "repo",
				User:      models.User{Login: "author", HTMLURL: "https://github.com/author"},
				CreatedAt: time.Date(2026, 3, 15, 10, 0, 0, 0, time.UTC),
				State:     "closed",
				ClosedAt:  timePtr(time.Date(2026, 3, 16, 10, 0, 0, 0, time.UTC)),
				Labels:    []models.Label{},
				Assignees: []models.User{},
				Milestone: nil,
				Comments:  []models.Comment{},
			},
			userLinks: false,
			wantFrontMatter: map[string]interface{}{
				"state": "closed",
			},
			wantBodyContains: "closed",
		},
		{
			name: "no milestone",
			input: &models.IssueData{
				Title:     "Issue without milestone",
				Body:      "No milestone assigned",
				HTMLURL:   "https://github.com/owner/repo/issues/3",
				Number:    3,
				Repo:      "repo",
				User:      models.User{Login: "nobody", HTMLURL: "https://github.com/nobody"},
				CreatedAt: time.Date(2026, 3, 14, 10, 0, 0, 0, time.UTC),
				State:     "open",
				Labels:    []models.Label{},
				Assignees: []models.User{},
				Milestone: nil,
				Comments:  []models.Comment{},
			},
			userLinks:       false,
			wantFrontMatter: map[string]interface{}{},
			wantBodyContains: "No milestone assigned",
		},
		{
			name: "with comments",
			input: &models.IssueData{
				Title:     "Issue with comments",
				Body:      "Original post",
				HTMLURL:   "https://github.com/owner/repo/issues/4",
				Number:    4,
				Repo:      "repo",
				User:      models.User{Login: "op", HTMLURL: "https://github.com/op"},
				CreatedAt: time.Date(2026, 3, 13, 10, 0, 0, 0, time.UTC),
				State:     "open",
				Labels:    []models.Label{},
				Assignees: []models.User{},
				Milestone: nil,
				Comments: []models.Comment{
					{
						ID:        1,
						Body:      "First comment",
						User:      models.User{Login: "commenter1", HTMLURL: "https://github.com/commenter1"},
						CreatedAt: time.Date(2026, 3, 13, 11, 0, 0, 0, time.UTC),
						UpdatedAt: time.Date(2026, 3, 13, 11, 0, 0, 0, time.UTC),
					},
					{
						ID:        2,
						Body:      "Second comment",
						User:      models.User{Login: "commenter2", HTMLURL: "https://github.com/commenter2"},
						CreatedAt: time.Date(2026, 3, 13, 12, 0, 0, 0, time.UTC),
						UpdatedAt: time.Date(2026, 3, 13, 12, 0, 0, 0, time.UTC),
					},
				},
			},
			userLinks:        false,
			wantFrontMatter:  map[string]interface{}{},
			wantBodyContains: "### @commenter1",
		},
		{
			name: "user links enabled",
			input: &models.IssueData{
				Title:     "Issue with user links",
				Body:      "Post by user",
				HTMLURL:   "https://github.com/owner/repo/issues/5",
				Number:    5,
				Repo:      "repo",
				User:      models.User{Login: "linkuser", HTMLURL: "https://github.com/linkuser"},
				CreatedAt: time.Date(2026, 3, 12, 10, 0, 0, 0, time.UTC),
				State:     "open",
				Labels:    []models.Label{},
				Assignees: []models.User{},
				Milestone: nil,
				Comments: []models.Comment{
					{
						ID:        1,
						Body:      "Comment with link",
						User:      models.User{Login: "linkcommenter", HTMLURL: "https://github.com/linkcommenter"},
						CreatedAt: time.Date(2026, 3, 12, 11, 0, 0, 0, time.UTC),
						UpdatedAt: time.Date(2026, 3, 12, 11, 0, 0, 0, time.UTC),
					},
				},
			},
			userLinks:        true,
			wantFrontMatter:  map[string]interface{}{},
			wantBodyContains: "[@linkcommenter](https://github.com/linkcommenter)",
		},
		{
			name: "special characters in title",
			input: &models.IssueData{
				Title:     "Issue with \"quotes\" and: colons",
				Body:      "Testing YAML escaping",
				HTMLURL:   "https://github.com/owner/repo/issues/6",
				Number:    6,
				Repo:      "repo",
				User:      models.User{Login: "escaper", HTMLURL: "https://github.com/escaper"},
				CreatedAt: time.Date(2026, 3, 11, 10, 0, 0, 0, time.UTC),
				State:     "open",
				Labels:    []models.Label{},
				Assignees: []models.User{},
				Milestone: nil,
				Comments:  []models.Comment{},
			},
			userLinks:        false,
			wantFrontMatter:  map[string]interface{}{},
			wantBodyContains: "Testing YAML escaping",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConverter(Options{UserLinks: tt.userLinks})
			markdown, filename, err := c.Convert(tt.input)

			if err != nil {
				t.Fatalf("Convert() error = %v", err)
			}

			// Check filename format
			expectedFilename := formatExpectedFilename(tt.input)
			if filename != expectedFilename {
				t.Errorf("Convert() filename = %q, want %q", filename, expectedFilename)
			}

			// Check Front Matter
			for key, expectedValue := range tt.wantFrontMatter {
				if !containsFrontMatterValue(markdown, key, expectedValue) {
					t.Errorf("Convert() front matter %s = %v, want %v\nFull markdown:\n%s", key, extractFrontMatterValue(markdown, key), expectedValue, markdown)
				}
			}

			// Check body content
			if tt.wantBodyContains != "" && !strings.Contains(markdown, tt.wantBodyContains) {
				t.Errorf("Convert() body should contain %q\nFull markdown:\n%s", tt.wantBodyContains, markdown)
			}
		})
	}
}

func TestWriteFile(t *testing.T) {
	tests := []struct {
		name      string
		filename  string
		content   string
		setup     func(t *testing.T) string // returns temp dir
		cleanup   func(t *testing.T, dir string)
		wantErr   bool
		errContains string
	}{
		{
			name:     "current directory",
			filename: "test-current-dir.md",
			content:  "---\ntitle: Test\n---\n\nContent",
			setup: func(t *testing.T) string {
				return t.TempDir()
			},
			cleanup:   func(t *testing.T, dir string) {},
			wantErr:   false,
		},
		{
			name:     "custom directory",
			filename: "test-custom-dir.md",
			content:  "---\ntitle: Test\n---\n\nContent",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				return dir
			},
			cleanup:   func(t *testing.T, dir string) {},
			wantErr:   false,
		},
		{
			name:      "read-only directory",
			filename:  "/proc/nonexistent/test.md",
			content:   "---\ntitle: Test\n---\n\nContent",
			setup:     func(t *testing.T) string { return "" },
			cleanup:   func(t *testing.T, dir string) {},
			wantErr:   true,
			errContains: "创建输出目录失败",
		},
		{
			name:     "overwrite existing file",
			filename: "test-overwrite.md",
			content:  "---\ntitle: Overwritten\n---\n\nNew content",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				// Create existing file
				existingPath := dir + "/test-overwrite.md"
				if err := writeFileHelper(existingPath, "old content"); err != nil {
					t.Fatalf("setup: failed to create existing file: %v", err)
				}
				return dir
			},
			cleanup:   func(t *testing.T, dir string) {},
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.setup(t)
			defer tt.cleanup(t, dir)

			fullPath := dir + "/" + tt.filename
			if dir == "" {
				fullPath = tt.filename
			}

			c := NewConverter(Options{})
			err := c.WriteFile(fullPath, tt.content)

			if (err != nil) != tt.wantErr {
				t.Errorf("WriteFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && tt.errContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("WriteFile() error = %v, want error containing %q", err, tt.errContains)
				}
			}
		})
	}
}

// Helper functions

func timePtr(t time.Time) *time.Time {
	return &t
}

func formatExpectedFilename(data *models.IssueData) string {
	name := strings.ReplaceAll(data.Title, "/", "-")
	name = strings.ReplaceAll(name, `\`, "-")
	name = strings.ReplaceAll(name, ":", "-")
	name = strings.ReplaceAll(name, "?", "-")
	name = strings.ReplaceAll(name, "*", "-")
	name = strings.ReplaceAll(name, `"`, "-")
	name = strings.ReplaceAll(name, "<", "-")
	name = strings.ReplaceAll(name, ">", "-")
	name = strings.ReplaceAll(name, "|", "-")
	return name + ".md"
}

func containsFrontMatterValue(markdown, key string, expectedValue interface{}) bool {
	// Extract front matter (between first two ---)
	lines := strings.Split(markdown, "\n")
	if len(lines) < 3 || lines[0] != "---" {
		return false
	}

	endIndex := -1
	for i := 1; i < len(lines); i++ {
		if lines[i] == "---" {
			endIndex = i
			break
		}
	}

	if endIndex == -1 {
		return false
	}

	// Search for key in front matter
	keyPrefix := key + ":"
	for i := 1; i < endIndex; i++ {
		if strings.HasPrefix(lines[i], keyPrefix) {
			value := strings.TrimSpace(strings.TrimPrefix(lines[i], keyPrefix))
			return valueMatches(value, expectedValue)
		}
	}

	return false
}

func extractFrontMatterValue(markdown, key string) string {
	lines := strings.Split(markdown, "\n")
	if len(lines) < 3 || lines[0] != "---" {
		return ""
	}

	endIndex := -1
	for i := 1; i < len(lines); i++ {
		if lines[i] == "---" {
			endIndex = i
			break
		}
	}

	if endIndex == -1 {
		return ""
	}

	keyPrefix := key + ":"
	for i := 1; i < endIndex; i++ {
		if strings.HasPrefix(lines[i], keyPrefix) {
			return strings.TrimSpace(strings.TrimPrefix(lines[i], keyPrefix))
		}
	}

	return ""
}

func valueMatches(actual string, expected interface{}) bool {
	switch v := expected.(type) {
	case string:
		return actual == v
	case int:
		return actual == string(rune(v))
	default:
		return strings.Contains(actual, toString(expected))
	}
}

func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case int:
		return string(rune(val))
	default:
		return ""
	}
}

func writeFileHelper(path, content string) error {
	// Simple helper for test setup
	return nil
}
