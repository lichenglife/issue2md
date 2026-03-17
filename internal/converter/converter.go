package converter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lichenglife/issue2md/internal/models"
)

// Options Markdown 生成选项
type Options struct {
	UserLinks bool   // 用户名渲染为 GitHub 主页链接
	OutputDir string // 输出目录
}

// Converter Markdown 转换器
type Converter struct {
	opts Options
}

// NewConverter 创建转换器
func NewConverter(opts Options) *Converter {
	return &Converter{opts: opts}
}

// Convert 将 IssueData 转换为 Markdown 内容
// 返回：markdown 内容，文件名（不含路径），错误
func (c *Converter) Convert(data *models.IssueData) (markdown string, filename string, err error) {
	// 生成文件名
	filename = generateFilename(data)

	// 构建 Markdown 内容
	var sb strings.Builder

	// 1. Front Matter
	sb.WriteString(c.generateFrontMatter(data))

	// 2. 正文内容
	sb.WriteString(c.generateBody(data))

	return sb.String(), filename, nil
}

// generateFilename 生成 Markdown 文件名
func generateFilename(data *models.IssueData) string {
	// 使用标题作为文件名基础，替换非法字符
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

// generateFrontMatter 生成 YAML Front Matter
func (c *Converter) generateFrontMatter(data *models.IssueData) string {
	var sb strings.Builder

	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("title: %s\n", escapeYAMLString(data.Title)))
	sb.WriteString(fmt.Sprintf("author: %s\n", escapeYAMLString(data.User.Login)))
	sb.WriteString(fmt.Sprintf("url: %s\n", escapeYAMLString(data.HTMLURL)))
	sb.WriteString(fmt.Sprintf("number: %d\n", data.Number))
	sb.WriteString(fmt.Sprintf("created_at: %s\n", data.CreatedAt.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("state: %s\n", data.State))

	// ClosedAt (可选)
	if data.ClosedAt != nil {
		sb.WriteString(fmt.Sprintf("closed_at: %s\n", data.ClosedAt.Format(time.RFC3339)))
	}

	// Milestone (可选)
	if data.Milestone != nil {
		sb.WriteString(fmt.Sprintf("milestone: %s\n", escapeYAMLString(data.Milestone.Title)))
	}

	// Labels
	if len(data.Labels) > 0 {
		sb.WriteString("labels:\n")
		for _, label := range data.Labels {
			sb.WriteString(fmt.Sprintf("  - %s\n", escapeYAMLString(label.Name)))
		}
	}

	// Assignees
	if len(data.Assignees) > 0 {
		sb.WriteString("assignees:\n")
		for _, assignee := range data.Assignees {
			sb.WriteString(fmt.Sprintf("  - %s\n", escapeYAMLString(assignee.Login)))
		}
	}

	sb.WriteString("---\n\n")

	return sb.String()
}

// generateBody 生成 Markdown 正文
func (c *Converter) generateBody(data *models.IssueData) string {
	var sb strings.Builder

	// 标题
	sb.WriteString(fmt.Sprintf("# %s\n\n", data.Title))

	// 作者信息
	sb.WriteString(fmt.Sprintf("**作者**: %s\n\n", formatUserLink(data.User.Login, data.User.HTMLURL, c.opts.UserLinks)))

	// 状态信息
	if data.State == "closed" && data.ClosedAt != nil {
		sb.WriteString(fmt.Sprintf("**状态**: 已关闭 (关闭时间：%s)\n\n", data.ClosedAt.Format("2006-01-02 15:04:05")))
	} else {
		sb.WriteString(fmt.Sprintf("**状态**: %s\n\n", data.State))
	}

	// 标签
	if len(data.Labels) > 0 {
		sb.WriteString("**标签**: ")
		var labelNames []string
		for _, label := range data.Labels {
			labelNames = append(labelNames, label.Name)
		}
		sb.WriteString(strings.Join(labelNames, ", "))
		sb.WriteString("\n\n")
	}

	// 里程碑
	if data.Milestone != nil {
		sb.WriteString(fmt.Sprintf("**里程碑**: %s\n\n", data.Milestone.Title))
	}

	// 分配给
	if len(data.Assignees) > 0 {
		sb.WriteString("**分配给**: ")
		var assigneeLinks []string
		for _, assignee := range data.Assignees {
			assigneeLinks = append(assigneeLinks, formatUserLink(assignee.Login, assignee.HTMLURL, c.opts.UserLinks))
		}
		sb.WriteString(strings.Join(assigneeLinks, ", "))
		sb.WriteString("\n\n")
	}

	// 正文内容
	sb.WriteString("---\n\n")
	sb.WriteString("## 正文\n\n")
	sb.WriteString(data.Body)
	sb.WriteString("\n\n")

	// 评论
	if len(data.Comments) > 0 {
		sb.WriteString("---\n\n")
		sb.WriteString("## 评论\n\n")

		for _, comment := range data.Comments {
			sb.WriteString(fmt.Sprintf("### %s\n\n", formatUserLink(comment.User.Login, comment.User.HTMLURL, c.opts.UserLinks)))
			sb.WriteString(fmt.Sprintf("*%s*\n\n", comment.CreatedAt.Format("2006-01-02 15:04:05")))
			sb.WriteString(comment.Body)
			sb.WriteString("\n\n")
		}
	}

	// 原始链接
	sb.WriteString("---\n\n")
	sb.WriteString(fmt.Sprintf("*本文档由 [issue2md](https://github.com/lichenglife/issue2md) 自动生成 | [原始 Issue](%s)*\n", data.HTMLURL))

	return sb.String()
}

// WriteFile 将 Markdown 内容写入文件
// 目录不存在时返回错误
// 文件已存在时直接覆盖
func (c *Converter) WriteFile(filename, content string) error {
	// 确保目录存在
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败：%w", err)
	}

	// 写入文件
	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		return fmt.Errorf("写入文件失败：%w", err)
	}

	return nil
}
