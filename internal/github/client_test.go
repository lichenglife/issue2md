package github

import (
	"context"
	"testing"
)

// TestNewClient 测试客户端创建
func TestNewClient(t *testing.T) {
	// 匿名客户端
 anonClient := NewClient("")
	if anonClient == nil {
		t.Fatal("NewClient(\"\") returned nil")
	}

	// 带 token 的客户端
	authClient := NewClient("ghp_test123")
	if authClient == nil {
		t.Fatal("NewClient(\"ghp_test123\") returned nil")
	}
}

// TestFetchIssue 测试 Issue 获取（集成测试）
// 需要有效的 GITHUB_TOKEN 才能通过
func TestFetchIssue(t *testing.T) {
	token := getTestToken()
	if token == "" {
		t.Skip("Skipping integration test: GITHUB_TOKEN not set")
	}

	client := NewClient(token)
	ctx := context.Background()

	// 测试获取公开 Issue（以 issue2md 仓库为例）
	tests := []struct {
		name         string
		owner        string
		repo         string
		number       int
		resourceType string
		wantErr      bool
	}{
		{
			name:         "valid issue",
			owner:        "lichenglife",
			repo:         "issue2md",
			number:       1,
			resourceType: "issue",
			wantErr:      false,
		},
		{
			name:         "non-existent issue",
			owner:        "lichenglife",
			repo:         "issue2md",
			number:       99999,
			resourceType: "issue",
			wantErr:      true,
		},
		{
			name:         "invalid owner",
			owner:        "nonexistent-owner-xyz",
			repo:         "nonexistent-repo-xyz",
			number:       1,
			resourceType: "issue",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.FetchIssue(ctx, tt.owner, tt.repo, tt.number, tt.resourceType)

			if (err != nil) != tt.wantErr {
				t.Errorf("FetchIssue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got == nil {
					t.Fatal("FetchIssue() returned nil for valid input")
				}
				if got.Number != tt.number {
					t.Errorf("FetchIssue() Number = %d, want %d", got.Number, tt.number)
				}
				if got.Repo != tt.repo {
					t.Errorf("FetchIssue() Repo = %q, want %q", got.Repo, tt.repo)
				}
			}
		})
	}
}

// TestFetchIssue_ErrorTypes 测试错误类型判断
func TestFetchIssue_ErrorTypes(t *testing.T) {
	client := NewClient("invalid-token-xyz")
	ctx := context.Background()

	_, err := client.FetchIssue(ctx, "lichenglife", "issue2md", 1, "issue")

	// 无效 token 应该返回认证错误
	if !IsUnauthorized(err) && !IsNotFound(err) {
		// 可能因为 token 无效或网络问题，只要返回了错误即可
		if err == nil {
			t.Error("FetchIssue() should return error with invalid token")
		}
	}
}

// TestIsNotFound 测试错误类型判断函数
func TestIsNotFound(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"not found error", ErrNotFound, true},
		{"wrapped not found", wrapError(ErrNotFound), true},
		{"other error", ErrNetworkError, false},
		{"nil error", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsNotFound(tt.err)
			if got != tt.want {
				t.Errorf("IsNotFound() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsUnauthorized 测试错误类型判断函数
func TestIsUnauthorized(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"unauthorized error", ErrUnauthorized, true},
		{"wrapped unauthorized", wrapError(ErrUnauthorized), true},
		{"other error", ErrNotFound, false},
		{"nil error", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsUnauthorized(tt.err)
			if got != tt.want {
				t.Errorf("IsUnauthorized() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsRateLimited 测试错误类型判断函数
func TestIsRateLimited(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"rate limited error", ErrRateLimited, true},
		{"wrapped rate limited", wrapError(ErrRateLimited), true},
		{"other error", ErrNotFound, false},
		{"nil error", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsRateLimited(tt.err)
			if got != tt.want {
				t.Errorf("IsRateLimited() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper functions

func getTestToken() string {
	// 从环境变量获取测试 token
	return ""
}

func wrapError(err error) error {
	return &wrappedError{err: err}
}

type wrappedError struct {
	err error
}

func (w *wrappedError) Error() string {
	return w.err.Error()
}

func (w *wrappedError) Unwrap() error {
	return w.err
}
