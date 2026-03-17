package config

import (
	"flag"
	"os"

	"github.com/lichenglife/issue2md/internal/models"
)

// Load 从命令行参数和环境变量加载配置
// 优先级：命令行参数 > 环境变量 > 默认值
func Load(args []string) (*models.Config, error) {
	fs := flag.NewFlagSet("issue2md", flag.ContinueOnError)

	outputDir := fs.String("o", ".", "输出目录")
	userLinks := fs.Bool("user-links", false, "将用户名渲染为 GitHub 主页链接")
	verbose := fs.Bool("v", false, "详细日志模式")
	token := fs.String("token", "", "GitHub Token")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	// 如果 token 未通过命令行参数提供，尝试从环境变量获取
	if *token == "" {
		*token = os.Getenv("GITHUB_TOKEN")
	}

	return &models.Config{
		OutputDir: *outputDir,
		UserLinks: *userLinks,
		Verbose:   *verbose,
		Token:     *token,
	}, nil
}

// GetTokenFromEnv 从环境变量获取 GitHub Token
func GetTokenFromEnv() string {
	return os.Getenv("GITHUB_TOKEN")
}
