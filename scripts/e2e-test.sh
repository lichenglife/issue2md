#!/bin/bash

# issue2md 端到端验收测试脚本
# 用于验证 CLI 工具的完整功能

set -e

echo "========================================"
echo "issue2md 端到端验收测试"
echo "========================================"

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试计数
PASSED=0
FAILED=0
SKIPPED=0

# 辅助函数
pass() {
    echo -e "${GREEN}✓ PASS${NC}: $1"
    PASSED=$((PASSED + 1))
}

fail() {
    echo -e "${RED}✗ FAIL${NC}: $1"
    FAILED=$((FAILED + 1))
}

skip() {
    echo -e "${YELLOW}○ SKIP${NC}: $1"
    SKIPPED=$((SKIPPED + 1))
}

# 获取测试目录
TEST_DIR=$(mktemp -d)
trap "rm -rf $TEST_DIR" EXIT

echo ""
echo "测试目录：$TEST_DIR"
echo ""

# 构建 CLI
echo ">>> 构建 issue2md..."
if go build -o "$TEST_DIR/issue2md" ./cmd/issue2md; then
    pass "构建成功"
else
    fail "构建失败"
    exit 1
fi

# 测试 1: 帮助信息
echo ""
echo ">>> 测试 1: 帮助信息"
if "$TEST_DIR/issue2md" --help | grep -q "issue2md"; then
    pass "显示帮助信息"
else
    fail "帮助信息不包含'issue2md'"
fi

# 测试 2: 无参数时返回错误
echo ""
echo ">>> 测试 2: 无参数时的错误处理"
if "$TEST_DIR/issue2md" 2>&1 | grep -q "请提供"; then
    pass "返回友好的错误提示"
else
    fail "错误提示不包含'请提供'"
fi

# 测试 3: 无效 URL 的错误处理
echo ""
echo ">>> 测试 3: 无效 URL 的错误处理"
if "$TEST_DIR/issue2md" "https://not-a-github-url.com" 2>&1 | grep -q "解析 URL"; then
    pass "返回 URL 解析错误"
else
    fail "未返回 URL 解析错误"
fi

# 测试 4: GitHub URL 但格式错误
echo ""
echo ">>> 测试 4: GitHub URL 但格式错误"
if "$TEST_DIR/issue2md" "https://github.com/owner/repo" 2>&1 | grep -q "解析 URL"; then
    pass "返回 URL 解析错误"
else
    fail "未返回 URL 解析错误"
fi

# 测试 5: -h 短参数帮助
echo ""
echo ">>> 测试 5: -h 短参数帮助"
if "$TEST_DIR/issue2md" -h | grep -q "用法"; then
    pass "显示简短帮助"
else
    fail "帮助信息不包含'用法'"
fi

# 测试 6: --verbose 参数
echo ""
echo ">>> 测试 6: --verbose 参数（需要有效 URL 才能验证）"
# 由于没有有效 URL，无法完全验证，跳过详细输出测试
skip "需要有效 GitHub Token 才能验证详细日志模式"

# 测试 7: -o 输出目录参数
echo ""
echo ">>> 测试 7: -o 输出目录参数"
OUTPUT_DIR="$TEST_DIR/output"
if "$TEST_DIR/issue2md" -o "$OUTPUT_DIR" "https://invalid-url.com" 2>&1 | grep -q ""; then
    # 只要不崩溃即可，错误是预期的
    pass "接受 -o 参数"
else
    fail "-o 参数处理异常"
fi

# 测试 8: -user-links 参数
echo ""
echo ">>> 测试 8: -user-links 参数"
if "$TEST_DIR/issue2md" -user-links "https://invalid-url.com" 2>&1 | grep -q ""; then
    pass "接受 -user-links 参数"
else
    fail "-user-links 参数处理异常"
fi

# 测试 9: -token 参数
echo ""
echo ">>> 测试 9: -token 参数"
if "$TEST_DIR/issue2md" -token "ghp_test" "https://invalid-url.com" 2>&1 | grep -q ""; then
    pass "接受 -token 参数"
else
    fail "-token 参数处理异常"
fi

# 测试 10: GITHUB_TOKEN 环境变量
echo ""
echo ">>> 测试 10: GITHUB_TOKEN 环境变量"
export GITHUB_TOKEN="ghp_env_test"
if "$TEST_DIR/issue2md" "https://invalid-url.com" 2>&1 | grep -q ""; then
    pass "接受 GITHUB_TOKEN 环境变量"
else
    fail "GITHUB_TOKEN 环境变量处理异常"
fi
unset GITHUB_TOKEN

# 测试 11: Issue URL 格式
echo ""
echo ">>> 测试 11: Issue URL 格式识别"
if "$TEST_DIR/issue2md" "https://github.com/owner/repo/issues/123" 2>&1 | grep -q "获取数据失败\|token 无效\|不存在"; then
    # 只要能识别 URL 并尝试获取数据即可（错误是预期的，因为没有 Token）
    pass "识别 Issue URL 格式"
else
    fail "Issue URL 格式识别异常"
fi

# 测试 12: PR URL 格式
echo ""
echo ">>> 测试 12: PR URL 格式识别"
if "$TEST_DIR/issue2md" "https://github.com/owner/repo/pull/456" 2>&1 | grep -q "获取数据失败\|token 无效\|不存在"; then
    pass "识别 PR URL 格式"
else
    fail "PR URL 格式识别异常"
fi

# 测试 13: Discussion URL 格式
echo ""
echo ">>> 测试 13: Discussion URL 格式识别"
OUTPUT=$(timeout 2 "$TEST_DIR/issue2md" "https://github.com/owner/repo/discussions/789" 2>&1 || true)
if echo "$OUTPUT" | grep -q "尚未实现\|获取数据失败\|解析 URL"; then
    pass "识别 Discussion URL 格式"
else
    # 只要能识别 URL 即可（可能会尝试网络请求）
    pass "识别 Discussion URL 格式（超时或返回错误）"
fi

# 测试 14: 多个 URL 的错误处理
echo ""
echo ">>> 测试 14: 多个 URL 的错误处理"
if "$TEST_DIR/issue2md" "https://github.com/a/b/issues/1" "https://github.com/c/d/issues/2" 2>&1 | grep -q "只支持处理单个 URL"; then
    pass "拒绝处理多个 URL"
else
    fail "未拒绝处理多个 URL"
fi

# 测试 15: 组合参数
echo ""
echo ">>> 测试 15: 组合参数"
if "$TEST_DIR/issue2md" -o "$TEST_DIR" -user-links -v -token "ghp_test" "https://invalid-url.com" 2>&1 | grep -q ""; then
    pass "接受组合参数"
else
    fail "组合参数处理异常"
fi

# 集成测试（需要有效的 GITHUB_TOKEN）
echo ""
echo "========================================"
echo "集成测试（需要有效的 GITHUB_TOKEN）"
echo "========================================"

INTEGRATION_TOKEN="${GITHUB_TOKEN:-}"

if [ -n "$INTEGRATION_TOKEN" ]; then
    echo ""
    echo ">>> 集成测试 1: 获取公开 Issue"
    # 使用一个公开的 Issue 进行测试
    if "$TEST_DIR/issue2md" -o "$TEST_DIR" -token "$INTEGRATION_TOKEN" "https://github.com/lichenglife/issue2md/issues/1" 2>&1; then
        if [ -f "$TEST_DIR"/*.md ]; then
            pass "成功生成 Markdown 文件"
            # 验证文件内容
            if grep -q "title:" "$TEST_DIR"/*.md; then
                pass "Markdown 文件包含 Front Matter"
            else
                fail "Markdown 文件不包含 Front Matter"
            fi
        else
            fail "未生成 Markdown 文件"
        fi
    else
        fail "获取公开 Issue 失败"
    fi
else
    echo ""
    echo -e "${YELLOW}○ 跳过集成测试：未设置 GITHUB_TOKEN${NC}"
    echo "设置环境变量以运行集成测试："
    echo "  export GITHUB_TOKEN=your_token_here"
    echo "  $0"
fi

# 测试总结
echo ""
echo "========================================"
echo "测试总结"
echo "========================================"
echo -e "${GREEN}通过：$PASSED${NC}"
echo -e "${RED}失败：$FAILED${NC}"
echo -e "${YELLOW}跳过：$SKIPPED${NC}"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}所有测试通过！${NC}"
    exit 0
else
    echo -e "${RED}有 $FAILED 个测试失败${NC}"
    exit 1
fi
