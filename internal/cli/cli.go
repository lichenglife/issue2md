package cli

import (
	"context"
	"fmt"

	"github.com/lichenglife/issue2md/internal/config"
	"github.com/lichenglife/issue2md/internal/converter"
	"github.com/lichenglife/issue2md/internal/github"
	"github.com/lichenglife/issue2md/internal/models"
	"github.com/lichenglife/issue2md/internal/parser"
)

// Executor CLI 执行器
type Executor struct {
	cfg *models.Config
}

// NewExecutor 创建 CLI 执行器
func NewExecutor(cfg *models.Config) *Executor {
	return &Executor{cfg: cfg}
}

// Run 执行 CLI 命令
// url: GitHub Issue/PR/Discussion URL
func (e *Executor) Run(rawURL string) error {
	// 1. 解析 URL
	parsed, err := parser.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("解析 URL 失败：%w", err)
	}

	if e.cfg.Verbose {
		fmt.Printf("解析结果：所有者=%s, 仓库=%s, 类型=%s, 编号=%d\n",
			parsed.Owner, parsed.Repo, parsed.Type, parsed.Number)
	}

	// 2. 获取数据
	client := github.NewClient(e.cfg.Token)
	ctx := context.Background()

	var data *models.IssueData
	switch parsed.Type {
	case "issue", "pull":
		data, err = client.FetchIssue(ctx, parsed.Owner, parsed.Repo, parsed.Number, parsed.Type)
	case "discussion":
		data, err = client.FetchDiscussion(ctx, parsed.Owner, parsed.Repo, parsed.Number)
	default:
		return fmt.Errorf("不支持的资源类型：%s", parsed.Type)
	}

	if err != nil {
		return fmt.Errorf("获取数据失败：%w", err)
	}

	if e.cfg.Verbose {
		fmt.Printf("获取到数据：标题=%s, 作者=%s, 评论数=%d\n",
			data.Title, data.User.Login, len(data.Comments))
	}

	// 3. 生成 Markdown
	converterOpts := converter.Options{
		UserLinks: e.cfg.UserLinks,
		OutputDir: e.cfg.OutputDir,
	}
	conv := converter.NewConverter(converterOpts)

	markdown, filename, err := conv.Convert(data)
	if err != nil {
		return fmt.Errorf("转换失败：%w", err)
	}

	if e.cfg.Verbose {
		fmt.Printf("生成 Markdown 完成：文件名=%s\n", filename)
	}

	// 4. 写入文件
	outputPath := e.cfg.OutputDir + "/" + filename
	if e.cfg.OutputDir == "." {
		outputPath = filename
	}

	if err := conv.WriteFile(outputPath, markdown); err != nil {
		return fmt.Errorf("写入文件失败：%w", err)
	}

	fmt.Printf("成功生成：%s\n", outputPath)
	return nil
}

// RunWithArgs 从命令行参数运行
func RunWithArgs(args []string) error {
	// 分离 flag 参数和非 flag 参数
	var urlArgs []string
	var flagArgs []string

	for _, arg := range args {
		if len(arg) > 0 && arg[0] == '-' {
			flagArgs = append(flagArgs, arg)
		} else {
			urlArgs = append(urlArgs, arg)
		}
	}

	// 加载配置
	cfg, err := config.Load(flagArgs)
	if err != nil {
		return fmt.Errorf("加载配置失败：%w", err)
	}

	if len(urlArgs) == 0 {
		return fmt.Errorf("请提供 GitHub Issue/PR/Discussion URL")
	}

	if len(urlArgs) > 1 {
		return fmt.Errorf("目前只支持处理单个 URL")
	}

	// 创建执行器并运行
	executor := NewExecutor(cfg)
	return executor.Run(urlArgs[0])
}
