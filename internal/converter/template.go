package converter

import (
	"strings"
)

// escapeYAMLString 转义 YAML 字符串中的特殊字符
// 用于 Front Matter 中的 title 等字段
func escapeYAMLString(s string) string {
	// 如果字符串包含特殊字符，使用双引号包裹并转义
	if needsQuoting(s) {
		// 转义双引号和反斜杠
		s = strings.ReplaceAll(s, `\`, `\\`)
		s = strings.ReplaceAll(s, `"`, `\"`)
		return `"` + s + `"`
	}
	return s
}

// needsQuoting 判断 YAML 字符串是否需要引号包裹
func needsQuoting(s string) bool {
	if s == "" {
		return true
	}

	// 检查是否包含需要转义的特殊字符
	specialChars := []string{
		`"`, `'`, `:`, `#`, `[`, `]`, `{`, `}`, `,`, `&`, `*`, `?`, `|`, `-`, `<`, `>`, `=`, `!`, `%`, `@`, `\`,
	}

	for _, char := range specialChars {
		if strings.Contains(s, char) {
			return true
		}
	}

	// 检查是否以特殊字符开头或结尾
	if strings.HasPrefix(s, " ") || strings.HasSuffix(s, " ") {
		return true
	}

	// 检查是否是 YAML 关键字
	yamlKeywords := []string{
		"true", "false", "yes", "no", "null", "~",
	}
	lower := strings.ToLower(s)
	for _, keyword := range yamlKeywords {
		if lower == keyword {
			return true
		}
	}

	return false
}

// formatUserLink 格式化用户链接
func formatUserLink(login string, htmlURL string, userLinks bool) string {
	if userLinks {
		return "[@" + login + "](" + htmlURL + ")"
	}
	return "@" + login
}
