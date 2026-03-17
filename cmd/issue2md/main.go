package main

import (
	"fmt"
	"os"

	"github.com/lichenglife/issue2md/internal/cli"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "错误：%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// 处理 --help 或 -h
	for _, arg := range os.Args[1:] {
		if arg == "-help" || arg == "--help" || arg == "-h" {
			printUsage()
			return nil
		}
	}

	return cli.RunWithArgs(os.Args[1:])
}

func printUsage() {
	fmt.Println(`issue2md - 将 GitHub Issue/PR/Discussion 转换为 Markdown 文档

用法:
  issue2md [选项] <GitHub URL>

选项:
  -o, --output <目录>   输出目录 (默认：当前目录)
  -user-links          将用户名渲染为 GitHub 主页链接
  -v, --verbose        详细日志模式
  -token <token>       GitHub Token (也可使用 GITHUB_TOKEN 环境变量)
  -h, --help           显示帮助信息

示例:
  issue2md https://github.com/owner/repo/issues/123
  issue2md -o ./output -user-links https://github.com/owner/repo/pull/456
  GITHUB_TOKEN=ghp_xxx issue2md https://github.com/owner/repo/issues/789

输出:
  成功转换后，Markdown 文件将保存到指定的输出目录。
  文件名基于 Issue/PR 标题生成。`)
}
