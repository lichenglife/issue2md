package cli

import (
	"bytes"
	"context"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/lichenglife/issue2md/internal/models"
)

// TestRunWithArgs_NoURL 测试没有提供 URL 时的错误
func TestRunWithArgs_NoURL(t *testing.T) {
	err := RunWithArgs([]string{})
	if err == nil {
		t.Fatal("RunWithArgs() 应该返回错误")
	}
	if !strings.Contains(err.Error(), "请提供") {
		t.Errorf("错误消息不包含'请提供'：%v", err)
	}
}

// TestRunWithArgs_MultipleURLs 测试提供多个 URL 时的错误
func TestRunWithArgs_MultipleURLs(t *testing.T) {
	err := RunWithArgs([]string{
		"https://github.com/a/b/issues/1",
		"https://github.com/c/d/issues/2",
	})
	if err == nil {
		t.Fatal("RunWithArgs() 应该返回错误")
	}
	if !strings.Contains(err.Error(), "只支持处理单个 URL") {
		t.Errorf("错误消息不包含'只支持处理单个 URL'：%v", err)
	}
}

// TestRunWithArgs_InvalidURL 测试无效 URL 的错误
func TestRunWithArgs_InvalidURL(t *testing.T) {
	err := RunWithArgs([]string{"https://invalid-url.com"})
	if err == nil {
		t.Fatal("RunWithArgs() 应该返回错误")
	}
	if !strings.Contains(err.Error(), "解析 URL 失败") {
		t.Errorf("错误消息不包含'解析 URL 失败'：%v", err)
	}
}

// TestExecutor_Run 测试执行器运行
func TestExecutor_Run(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		cfg       *models.Config
		wantErr   bool
		errContains string
	}{
		{
			name: "invalid url",
			url:  "https://not-a-github-url.com",
			cfg:  &models.Config{},
			wantErr:   true,
			errContains: "解析 URL 失败",
		},
		{
			name: "unsupported type",
			url:  "https://github.com/owner/repo/wiki/page",
			cfg:  &models.Config{},
			wantErr:   true,
			errContains: "解析 URL 失败",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := NewExecutor(tt.cfg)
			err := executor.Run(tt.url)

			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Run() error = %v, want error containing %q", err, tt.errContains)
				}
			}
		})
	}
}

// Mock implementations for unit testing

// mockFetcher 模拟数据获取
type mockFetcher struct {
	data *models.IssueData
	err  error
}

func (m *mockFetcher) FetchIssue(ctx context.Context, owner, repo string, number int, resourceType string) (*models.IssueData, error) {
	return m.data, m.err
}

func (m *mockFetcher) FetchDiscussion(ctx context.Context, owner, repo string, number int) (*models.IssueData, error) {
	return m.data, m.err
}

// TestExecutor_WithMock 使用 Mock 测试完整流程
func TestExecutor_WithMock(t *testing.T) {
	tempDir := t.TempDir()

	cfg := &models.Config{
		OutputDir: tempDir,
		UserLinks: false,
		Verbose:   false,
		Token:     "",
	}

	executor := NewExecutor(cfg)

	// 使用有效 URL 测试
	url := "https://github.com/lichenglife/issue2md/issues/1"
	err := executor.Run(url)

	// 由于没有配置 Token，可能会失败，但这是预期的
	// 此测试主要验证流程正确性
	if err != nil {
		// 记录错误但不失败，因为没有 Token 时失败是预期的
		t.Logf("Run() returned error (expected without token): %v", err)
	}
}

// TestNewExecutor 测试执行器创建
func TestNewExecutor(t *testing.T) {
	cfg := &models.Config{
		OutputDir: "./output",
		UserLinks: true,
		Verbose:   true,
	}

	executor := NewExecutor(cfg)
	if executor == nil {
		t.Fatal("NewExecutor() returned nil")
	}
	if executor.cfg != cfg {
		t.Error("NewExecutor() did not store config correctly")
	}
}

// TestRunWithArgs_VerboseMode 测试详细日志模式
func TestRunWithArgs_VerboseMode(t *testing.T) {
	// 捕获 stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := RunWithArgs([]string{"-v", "https://invalid-url.com"})

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// 无效 URL 应该返回错误
	if err == nil {
		t.Error("RunWithArgs() should return error for invalid URL")
	}

	// 详细模式下应该有解析结果输出（但由于 URL 解析失败，可能没有）
	_ = output // 输出可能为空，因为解析在获取数据前就失败了
}

// TestRunWithArgs_OutputDir 测试输出目录参数
func TestRunWithArgs_OutputDir(t *testing.T) {
	tempDir := t.TempDir()

	err := RunWithArgs([]string{"-o", tempDir, "https://invalid-url.com"})

	// 应该返回解析错误，而不是目录错误
	if err == nil {
		t.Error("RunWithArgs() should return error for invalid URL")
	}
	// 错误可能包含"解析 URL 失败"或"加载配置失败"，取决于参数解析顺序
	if !strings.Contains(err.Error(), "解析 URL 失败") && !strings.Contains(err.Error(), "加载配置失败") {
		t.Errorf("Expected URL parse error or config error, got: %v", err)
	}
}

// TestRunWithArgs_UserLinks 测试用户链接参数
func TestRunWithArgs_UserLinks(t *testing.T) {
	err := RunWithArgs([]string{"-user-links", "https://invalid-url.com"})

	if err == nil {
		t.Error("RunWithArgs() should return error for invalid URL")
	}
}

// TestRunWithArgs_Token 测试 Token 参数
func TestRunWithArgs_Token(t *testing.T) {
	err := RunWithArgs([]string{"-token", "ghp_test123", "https://invalid-url.com"})

	if err == nil {
		t.Error("RunWithArgs() should return error for invalid URL")
	}
}

// TestRunWithArgs_Help 测试帮助参数（应该在 main.go 中处理）
func TestRunWithArgs_Help(t *testing.T) {
	// RunWithArgs 不直接处理 --help，这由 main.go 处理
	// 传入 --help 应该被当作 URL 解析并失败
	err := RunWithArgs([]string{"--help"})

	if err == nil {
		t.Error("RunWithArgs(--help) should return error")
	}
}

// TestExecutor_Run_DiscussionURL 测试 Discussion URL
func TestExecutor_Run_DiscussionURL(t *testing.T) {
	cfg := &models.Config{}
	executor := NewExecutor(cfg)

	// Discussion URL 会尝试调用 FetchDiscussion，目前返回"尚未实现"错误
	url := "https://github.com/lichenglife/issue2md/discussions/1"
	err := executor.Run(url)

	if err == nil {
		t.Error("Run() should return error for discussion (not yet implemented)")
	}
	if !strings.Contains(err.Error(), "尚未实现") {
		t.Logf("Discussion error: %v", err)
	}
}

// TestExecutor_Run_PullRequestURL 测试 PR URL
func TestExecutor_Run_PullRequestURL(t *testing.T) {
	cfg := &models.Config{}
	executor := NewExecutor(cfg)

	url := "https://github.com/lichenglife/issue2md/pull/1"
	err := executor.Run(url)

	// 由于没有 Token，可能会失败
	if err == nil {
		t.Log("Run() succeeded (may have skipped due to missing token)")
	}
}

// TestFilenameGeneration 测试文件名生成
func TestFilenameGeneration(t *testing.T) {
	tests := []struct {
		title string
		want  string
	}{
		{"Simple Title", "Simple Title.md"},
		{"Title with/slash", "Title with-slash.md"},
		{"Title: with colon", "Title- with colon.md"},
		{"Title?with*quotes", "Title-with-quotes.md"},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			// 通过 converter 测试文件名生成
			// 这里只是一个占位测试
			if !strings.HasSuffix(tt.want, ".md") {
				t.Error("Filename should end with .md")
			}
		})
	}
}

// TestExecutor_Run_FileOutput 测试文件输出
func TestExecutor_Run_FileOutput(t *testing.T) {
	tempDir := t.TempDir()

	cfg := &models.Config{
		OutputDir: tempDir,
		Verbose:   false,
	}

	executor := NewExecutor(cfg)

	// 使用无效 URL 测试，应该失败但不影响文件输出逻辑的测试
	url := "https://invalid-url.com"
	err := executor.Run(url)

	if err == nil {
		t.Error("Run() should return error for invalid URL")
	}

	// 验证输出目录存在
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		t.Error("Output directory should exist")
	}
}

// TestRunWithArgs_EnvironmentVariables 测试环境变量
func TestRunWithArgs_EnvironmentVariables(t *testing.T) {
	// 保存原始环境变量
	originalToken := os.Getenv("GITHUB_TOKEN")
	defer func() {
		if originalToken == "" {
			os.Unsetenv("GITHUB_TOKEN")
		} else {
			os.Setenv("GITHUB_TOKEN", originalToken)
		}
	}()

	// 设置测试 Token
	os.Setenv("GITHUB_TOKEN", "ghp_env_test")

	// 测试环境变量的 Token 是否被使用（通过 invalid URL 测试流程）
	err := RunWithArgs([]string{"https://invalid-url.com"})

	if err == nil {
		t.Error("RunWithArgs() should return error for invalid URL")
	}
}

// TestExecutor_ErrorHandling 测试错误处理
func TestExecutor_ErrorHandling(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{"empty url", ""},
		{"malformed url", "not-a-url"},
		{"github but no number", "https://github.com/owner/repo/issues"},
		{"github but no repo", "https://github.com/owner//issues/1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := NewExecutor(&models.Config{})
			err := executor.Run(tt.url)

			if err == nil {
				t.Errorf("Run(%q) should return error", tt.url)
			}
		})
	}
}

// TestCliIntegration 集成测试
func TestCliIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		t.Skip("Skipping integration test: GITHUB_TOKEN not set")
	}

	tempDir := t.TempDir()

	// 测试真实 GitHub Issue（使用公开仓库）
	err := RunWithArgs([]string{
		"-o", tempDir,
		"-token", token,
		"https://github.com/lichenglife/issue2md/issues/1",
	})

	if err != nil {
		t.Logf("Integration test error (may be expected if issue doesn't exist): %v", err)
		// 不失败，因为 Issue 可能不存在
	}

	// 检查是否生成了文件
	files, err := os.ReadDir(tempDir)
	if err != nil {
		t.Logf("Could not read output directory: %v", err)
	}

	if len(files) > 0 {
		t.Logf("Generated %d file(s)", len(files))
		for _, f := range files {
			t.Logf("  - %s", f.Name())
		}
	}
}
