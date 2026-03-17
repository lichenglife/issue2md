package parser

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/lichenglife/issue2md/internal/models"
)

// Parse 解析 GitHub URL
// 支持的格式:
//   - https://github.com/owner/repo/issues/{number}
//   - https://github.com/owner/repo/pull/{number}
//   - https://github.com/owner/repo/discussions/{number}
func Parse(rawURL string) (*models.ParsedURL, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("无效的 GitHub URL: %w", err)
	}

	// 验证域名
	if u.Host != "github.com" {
		return nil, fmt.Errorf("无效的 GitHub URL: 域名必须是 github.com")
	}

	// 解析路径
	// 路径格式：/owner/repo/type/number
	path := strings.TrimSuffix(u.Path, "/")
	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")

	if len(parts) < 4 {
		return nil, fmt.Errorf("无效的 GitHub URL: 路径不完整")
	}

	owner := parts[0]
	repo := parts[1]
	resourceType := parts[2]
	numberStr := parts[3]

	// 验证资源类型
	var parsedType string
	switch resourceType {
	case "issues":
		parsedType = "issue"
	case "pull":
		parsedType = "pull"
	case "discussions":
		parsedType = "discussion"
	default:
		return nil, fmt.Errorf("不支持的资源类型：%s", resourceType)
	}

	// 解析数字
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		return nil, fmt.Errorf("无效的数字：%s", numberStr)
	}

	return &models.ParsedURL{
		Owner:  owner,
		Repo:   repo,
		Type:   parsedType,
		Number: number,
	}, nil
}

// IsValid 检查 URL 是否为有效的 GitHub Issue/PR/Discussion URL
func IsValid(rawURL string) bool {
	_, err := Parse(rawURL)
	return err == nil
}
