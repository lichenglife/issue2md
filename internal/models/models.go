package models

import "time"

// ParsedURL URL 解析结果
type ParsedURL struct {
	Owner  string // 仓库所有者
	Repo   string // 仓库名称
	Type   string // "issue" | "pull" | "discussion"
	Number int    // 编号
}

// User GitHub 用户
type User struct {
	Login   string `json:"login"`    // 用户名
	HTMLURL string `json:"html_url"` // 用户主页 URL
}

// Label 标签
type Label struct {
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description"`
}

// Milestone 里程碑
type Milestone struct {
	Title       string     `json:"title"`
	Number      int        `json:"number"`
	State       string     `json:"state"`
	Description string     `json:"description"`
	DueOn       *time.Time `json:"due_on"`
}

// Comment 评论
type Comment struct {
	ID        int64     `json:"id"`
	Body      string    `json:"body"`
	User      User      `json:"user"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// IssueData Issue/PR/Discussion 统一数据结构
// 包含 spec.md FR-002 要求的所有字段
type IssueData struct {
	// 基础信息
	Title     string `json:"title"`
	Body      string `json:"body"`
	HTMLURL   string `json:"html_url"`
	Number    int    `json:"number"`
	Repo      string `json:"repo"` // 用于生成文件名

	// 作者信息
	User User `json:"user"`

	// 时间信息
	CreatedAt time.Time  `json:"created_at"`
	State     string     `json:"state"` // "open" | "closed"
	ClosedAt  *time.Time `json:"closed_at"` // 可选，仅 closed 状态有值

	// 分类信息
	Labels    []Label     `json:"labels"`
	Assignees []User      `json:"assignees"`
	Milestone *Milestone  `json:"milestone"` // 可选

	// 评论列表（按 CreatedAt 升序排列）
	Comments []Comment `json:"comments"`
}

// Config 应用配置
type Config struct {
	OutputDir string // 输出目录
	UserLinks bool   // 用户名渲染为链接
	Verbose   bool   // 详细日志
	Token     string // GitHub Token
}
