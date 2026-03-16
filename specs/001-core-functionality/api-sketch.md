# API 接口设计草图

**Spec ID**: 001-core-functionality
**版本**: 1.0
**创建日期**: 2026-03-16
**状态**: Draft

---

本文档定义 `internal/converter` 和 `internal/github` 两个核心包的对外接口，
作为后续开发和测试的参考。

---

## 一、`internal/github` 包

负责与 GitHub API 交互，封装 REST API 和 GraphQL API 的调用细节。

### 1.1 数据结构

```go
package github

// User GitHub 用户信息
type User struct {
    Login  string `json:"login"`
    HTMLURL string `json:"html_url"`
}

// Label Issue/PR 标签
type Label struct {
    Name        string `json:"name"`
    Color       string `json:"color"`
    Description string `json:"description"`
}

// Milestone 里程碑信息
type Milestone struct {
    Title       string     `json:"title"`
    Number      int        `json:"number"`
    State       string     `json:"state"`
    Description string     `json:"description"`
    DueOn       *time.Time `json:"due_on"`
}

// Comment 评论信息
type Comment struct {
    ID        int64     `json:"id"`
    Body      string    `json:"body"`
    User      User      `json:"user"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// Issue Issue/PR 数据
type Issue struct {
    Title      string      `json:"title"`
    Body       string      `json:"body"`
    User       User        `json:"user"`
    CreatedAt  time.Time   `json:"created_at"`
    State      string      `json:"state"` // "open" | "closed"
    ClosedAt   *time.Time  `json:"closed_at"`
    Labels     []Label     `json:"labels"`
    Assignees  []User      `json:"assignees"`
    Milestone  *Milestone  `json:"milestone"`
    Comments   []Comment   `json:"comments"`
    HTMLURL    string      `json:"html_url"`
    Number     int         `json:"number"`
}

// Discussion Discussion 数据（通过 GraphQL API 获取）
type Discussion struct {
    Title      string      `json:"title"`
    Body       string      `json:"body"`
    User       User        `json:"user"`
    CreatedAt  time.Time   `json:"created_at"`
    State      string      `json:"state"` // "open" | "closed"
    ClosedAt   *time.Time  `json:"closed_at"`
    Labels     []Label     `json:"labels"`
    Assignees  []User      `json:"assignees"`
    Milestone  *Milestone  `json:"milestone"`
    Comments   []Comment   `json:"comments"`
    HTMLURL    string      `json:"html_url"`
    Number     int         `json:"number"`
}

// Resource 统一的资源类型（Issue 或 Discussion）
type Resource interface {
    GetTitle() string
    GetBody() string
    GetUser() User
    GetCreatedAt() time.Time
    GetState() string
    GetClosedAt() *time.Time
    GetLabels() []Label
    GetAssignees() []User
    GetMilestone() *Milestone
    GetComments() []Comment
    GetHTMLURL() string
    GetNumber() int
}
```

### 1.2 Client 接口

```go
// Client GitHub API 客户端
type Client struct {
    token string
    // 内部字段：HTTP 客户端、GraphQL 端点等
}

// NewClient 创建新的 GitHub API 客户端
// token 为空字符串时使用匿名访问（仅限公开仓库）
func NewClient(token string) *Client

// FetchIssue 获取 Issue/PR 数据
// resourceType: "issue" 或 "pull"
// 返回的 Issue.Comments 已按 CreatedAt 升序排列
func (c *Client) FetchIssue(ctx context.Context, owner, repo string, number int, resourceType string) (*Issue, error)

// FetchDiscussion 获取 Discussion 数据
// 返回的 Discussion.Comments 已按 CreatedAt 升序排列
func (c *Client) FetchDiscussion(ctx context.Context, owner, repo string, number int) (*Discussion, error)
```

### 1.3 错误类型

```go
// 预定义的错误，便于调用方进行错误处理
var (
    ErrUnauthorized      = errors.New("github: token 无效或已过期")
    ErrNotFound          = errors.New("github: 资源不存在或仓库为私有")
    ErrRateLimited       = errors.New("github: API 请求已限流，请稍后重试")
    ErrNetworkError      = errors.New("github: 网络请求失败")
    ErrParseResponse     = errors.New("github: 解析 API 响应失败")
)
```

---

## 二、`internal/converter` 包

负责将 GitHub 数据转换为 Markdown 格式。

### 2.1 配置选项

```go
package converter

// Options Markdown 生成选项
type Options struct {
    // UserLinks 是否将用户名渲染为 GitHub 主页链接
    // 默认：false（纯文本 @username）
    UserLinks bool

    // OutputDir 输出目录
    // 默认：当前工作目录
    OutputDir string
}
```

### 2.2 Converter 接口

```go
// Converter Markdown 转换器
type Converter struct {
    opts Options
}

// NewConverter 创建新的 Markdown 转换器
func NewConverter(opts Options) *Converter

// Convert 将 Issue 数据转换为 Markdown 内容
// 返回 Markdown 字符串和生成的文件名（不含目录）
func (c *Converter) Convert(issue *github.Issue) (markdown string, filename string, err error)

// ConvertDiscussion 将 Discussion 数据转换为 Markdown 内容
func (c *Converter) ConvertDiscussion(discussion *github.Discussion) (markdown string, filename string, err error)

// WriteFile 将 Markdown 内容写入文件
// 目标文件已存在时会直接覆盖
// 输出目录不存在时返回错误
func (c *Converter) WriteFile(filename string, content string) error
```

### 2.3 生成的 Markdown 格式

```markdown
---
title: "Issue 标题"
url: "https://github.com/owner/repo/issues/1"
author: "username"
created_at: "2026-03-10T08:30:00Z"
state: "open"
labels: ["label1", "label2"]
assignees: ["user1", "user2"]
milestone: "v1.0.0"
---

Issue 正文内容...

---

### @commenter - 2026-03-10T09:00:00Z

评论内容...
```

**UserLinks 开启时**:
```markdown
### [@commenter](https://github.com/commenter) - 2026-03-10T09:00:00Z
```

---

## 三、包依赖关系

```
cmd/issue2md ─┬─> internal/cli ─────┬─> internal/parser
              │                     ├─> internal/github
              │                     └─> internal/converter
              │
              └─> internal/config

internal/parser   : 无外部依赖
internal/github   : 无外部依赖（仅使用 net/http）
internal/converter: 依赖 internal/github（数据结构）
internal/cli      : 依赖 parser, github, converter, config
internal/config   : 无外部依赖
```

---

## 四、使用示例

### 4.1 基础用法

```go
package main

import (
    "context"
    "github.com/issue2md/internal/github"
    "github.com/issue2md/internal/converter"
)

func main() {
    // 1. 创建 GitHub 客户端
    client := github.NewClient("your-token-here")

    // 2. 获取 Issue 数据
    issue, err := client.FetchIssue(context.Background(), "owner", "repo", 123, "issue")
    if err != nil {
        // 处理错误
    }

    // 3. 创建转换器
    conv := converter.NewConverter(converter.Options{
        UserLinks: true,
        OutputDir: "./output",
    })

    // 4. 转换为 Markdown
    markdown, filename, err := conv.Convert(issue)
    if err != nil {
        // 处理错误
    }

    // 5. 写入文件
    err = conv.WriteFile(filename, markdown)
    if err != nil {
        // 处理错误
    }
}
```

---

## 五、待办事项

- [ ] 确认 Discussion 的 GraphQL 查询结构
- [ ] 确认 PR 是否需要单独的 Fetch 方法（目前复用 FetchIssue）
- [ ] 确认是否需要支持批量转换（多个 URL 一次性处理）
