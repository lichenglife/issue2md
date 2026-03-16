# 核心功能需求规格说明书

**Spec ID**: 001-core-functionality
**版本**: 1.0
**创建日期**: 2026-03-16
**状态**: Approved

---

## 一、用户故事 (User Stories)

### US-001: 命令行工具 - 基础转换

> 作为一名**开发者**，
> 我希望**输入一个 GitHub Issue/PR/Discussion URL，就能得到一个格式化的 Markdown 文件**，
> 以便于**在本地归档或导入到我的知识库系统中**。

**验收场景**:
```bash
# 场景 1: 快速归档一个 Issue
$ issue2md https://github.com/owner/repo/issues/123
✓ 生成文件：repo-123.md

# 场景 2: 归档到指定目录
$ issue2md https://github.com/owner/repo/issues/123 --output ./docs/
✓ 生成文件：./docs/repo-123.md

# 场景 3: 需要用户链接格式
$ issue2md https://github.com/owner/repo/issues/123 --user-links
✓ 用户名渲染为可点击的 GitHub 主页链接
```

---

### US-002: 命令行工具 - 详细日志模式

> 作为一名**高级用户**，
> 我希望**在遇到问题时能看到详细的执行日志**，
> 以便于**诊断问题或了解工具的执行过程**。

**验收场景**:
```bash
$ issue2md https://github.com/owner/repo/issues/123 --verbose
[INFO] 正在解析 URL...
[INFO] URL 解析成功：owner=issue2md, repo=issue2md, type=issue, number=1
[INFO] 正在获取 Issue 数据...
[INFO] 获取成功，评论数：5
[INFO] 正在生成 Markdown...
[INFO] 文件已保存：issue2md-1.md
```

---

### US-003: 命令行工具 - 访问私有仓库

> 作为一名**团队开发者**，
> 我希望**能够提供 GitHub Token 访问私有仓库**，
> 以便于**归档私有项目的讨论记录**。

**验收场景**:
```bash
# 场景 1: 通过命令行参数提供 Token
$ issue2md https://github.com/owner/private-repo/issues/42 --token $GH_TOKEN
✓ 成功访问私有仓库

# 场景 2: Token 无效时
$ issue2md https://github.com/owner/repo/issues/1 --token invalid_token
✗ 错误：Token 无效或已过期
```

---

### US-004: Web 服务（未来扩展）

> 作为一名**非技术用户**，
> 我希望**通过网页界面粘贴 URL 并下载 Markdown 文件**，
> 以便于**无需安装命令行工具即可使用本服务**。

**验收场景**:
```
场景 1: 网页转换
1. 用户访问 Web 界面
2. 在输入框中粘贴 GitHub URL
3. 点击"转换"按钮
4. 预览生成的 Markdown
5. 点击"下载"按钮获取文件
```

---

## 二、功能性需求 (Functional Requirements)

### FR-001: URL 自动识别

| 需求 ID | 描述 | 优先级 |
|---------|------|--------|
| FR-001.1 | 必须支持识别 Issue URL: `https://github.com/owner/repo/issues/{number}` | P0 |
| FR-001.2 | 必须支持识别 PR URL: `https://github.com/owner/repo/pull/{number}` | P0 |
| FR-001.3 | 必须支持识别 Discussion URL: `https://github.com/owner/repo/discussions/{number}` | P0 |
| FR-001.4 | 无法识别的 URL 必须返回友好的错误提示 | P0 |

**解析结果结构**:
```go
type ParsedURL struct {
    Owner  string // 仓库所有者
    Repo   string // 仓库名称
    Type   string // "issue" | "pull" | "discussion"
    Number int    // 编号
}
```

---

### FR-002: GitHub 数据获取

| 需求 ID | 描述 | 优先级 |
|---------|------|--------|
| FR-002.1 | 必须通过 GitHub REST API v3 获取 Issue/PR 数据 | P0 |
| FR-002.2 | 必须通过 GitHub GraphQL API 获取 Discussion 数据 | P0 |
| FR-002.3 | 支持通过 `--token` 参数提供认证 Token | P1 |
| FR-002.4 | 必须获取完整的元数据（标题、作者、时间、标签、里程碑等） | P0 |
| FR-002.5 | 必须获取所有评论并按时间升序排列 | P0 |

**必须获取的字段**:
```go
type Issue struct {
    Title      string
    Body       string
    User       User
    CreatedAt  time.Time
    State      string       // "open" | "closed"
    ClosedAt   *time.Time   // 可选
    Labels     []Label
    Assignees  []User
    Milestone  *Milestone   // 可选
    Comments   []Comment
}
```

---

### FR-003: Markdown 生成

| 需求 ID | 描述 | 优先级 |
|---------|------|--------|
| FR-003.1 | 必须生成符合 YAML Front Matter 规范的文件头 | P0 |
| FR-003.2 | 必须保留原始 Markdown 内容的格式 | P0 |
| FR-003.3 | 图片/附件链接必须保留原始 URL，不下载 | P0 |
| FR-003.4 | 评论必须按创建时间升序排列 | P0 |
| FR-003.5 | `--user-links` 开启时，用户名必须渲染为 GitHub 主页链接 | P1 |
| FR-003.6 | 已删除的评论必须保留原文 | P1 |

---

### FR-004: 命令行参数

| Flag | 类型 | 默认值 | 描述 | 优先级 |
|------|------|--------|------|--------|
| `--output`, `-o` | string | 当前目录 | 输出文件的目标目录 | P0 |
| `--user-links` | bool | false | 将用户名渲染为 GitHub 主页链接 | P1 |
| `--verbose`, `-v` | bool | false | 输出详细日志 | P1 |
| `--token` | string | 空 | GitHub Personal Access Token | P1 |

---

### FR-005: 文件输出

| 需求 ID | 描述 | 优先级 |
|---------|------|--------|
| FR-005.1 | 文件名格式必须为 `{repo}-{number}.md` | P0 |
| FR-005.2 | 默认输出到当前工作目录 | P0 |
| FR-005.3 | 支持通过 `--output` 指定输出目录 | P0 |
| FR-005.4 | 目标文件已存在时必须直接覆盖 | P0 |

---

### FR-006: 错误处理

| 需求 ID | 描述 | 优先级 |
|---------|------|--------|
| FR-006.1 | URL 格式无效时必须返回错误码 1 并提示友好信息 | P0 |
| FR-006.2 | HTTP 404 时必须提示"Issue 不存在或仓库为私有" | P0 |
| FR-006.3 | 网络请求失败时必须返回具体错误信息 | P0 |
| FR-006.4 | 输出目录不存在时必须提示用户创建目录 | P0 |
| FR-006.5 | GitHub API 限流时必须提示用户稍后重试 | P1 |

---

## 三、非功能性需求 (Non-Functional Requirements)

### NFR-001: 架构解耦

| 需求 ID | 描述 |
|---------|------|
| NFR-001.1 | CLI 层必须只负责参数解析和结果输出 |
| NFR-001.2 | 数据获取逻辑必须封装在独立的 `fetcher` 包中 |
| NFR-001.3 | Markdown 转换逻辑必须封装在独立的 `converter` 包中 |
| NFR-001.4 | 数据结构必须定义在独立的 `models` 包中 |
| NFR-001.5 | 禁止使用全局变量，所有依赖必须通过函数参数或结构体注入 |

---

### NFR-002: 错误处理

| 需求 ID | 描述 |
|---------|------|
| NFR-002.1 | 所有错误必须显式处理，禁止忽略 |
| NFR-002.2 | 错误传递时必须使用 `fmt.Errorf("...: %w", err)` 包装 |
| NFR-002.3 | 用户可见的错误信息必须友好、具体 |
| NFR-002.4 | 详细错误信息仅在 `--verbose` 模式下输出 |

---

### NFR-003: 代码质量

| 需求 ID | 描述 |
|---------|------|
| NFR-003.1 | 必须使用 Go 1.24+ 版本 |
| NFR-003.2 | 必须优先使用标准库，禁止引入不必要的依赖 |
| NFR-003.3 | 所有核心功能必须有单元测试覆盖 |
| NFR-003.4 | 单元测试必须采用表格驱动测试风格 |

---

## 四、验收标准 (Acceptance Criteria)

### AC-001: URL 解析测试

| 测试用例 | 输入 | 期望输出 |
|----------|------|----------|
| AC-001.1 | `https://github.com/owner/repo/issues/123` | `{Owner: owner, Repo: repo, Type: issue, Number: 123}` |
| AC-001.2 | `https://github.com/owner/repo/pull/456` | `{Owner: owner, Repo: repo, Type: pull, Number: 456}` |
| AC-001.3 | `https://github.com/owner/repo/discussions/789` | `{Owner: owner, Repo: repo, Type: discussion, Number: 789}` |
| AC-001.4 | `https://invalid-url.com` | 错误：无法解析 URL |
| AC-001.5 | `https://github.com/owner/repo` | 错误：无法解析 URL |

---

### AC-002: Markdown 生成测试

| 测试用例 | 输入数据 | 期望输出 |
|----------|----------|----------|
| AC-002.1 | Issue 包含标题、作者、标签 | Front Matter 包含所有字段 |
| AC-002.2 | Issue 已关闭 | Front Matter 包含 `closed_at` |
| AC-002.3 | Issue 无里程碑 | Front Matter 中 `milestone: null` |
| AC-002.4 | Issue 有 5 条评论 | 正文包含 5 条评论，按时间升序 |
| AC-002.5 | `--user-links` 开启 | 用户名为 `[@login](https://github.com/login)` 格式 |
| AC-002.6 | `--user-links` 关闭 | 用户名为纯文本 `@login` |

---

### AC-003: 错误处理测试

| 测试用例 | 模拟场景 | 期望行为 |
|----------|----------|----------|
| AC-003.1 | URL 格式无效 | 返回错误码 1，输出友好提示 |
| AC-003.2 | HTTP 404 | 返回错误码 1，提示"Issue 不存在或仓库为私有" |
| AC-003.3 | 网络超时 | 返回错误码 1，提示网络错误信息 |
| AC-003.4 | 输出目录不存在 | 返回错误码 1，提示创建目录 |
| AC-003.5 | GitHub API 限流 | 返回错误码 1，提示稍后重试 |

---

### AC-004: 文件输出测试

| 测试用例 | 输入 | 期望输出 |
|----------|------|----------|
| AC-004.1 | URL: `github.com/my/repo/issues/1` | 文件：`repo-1.md` |
| AC-004.2 | `--output ./docs/` | 文件：`./docs/repo-1.md` |
| AC-004.3 | 目标文件已存在 | 直接覆盖，`--verbose` 模式下输出警告 |

---

## 五、输出格式示例

### 5.1 Front Matter 示例

```yaml
---
title: "Add support for dark mode"
url: "https://github.com/issue2md/issue2md/issues/1"
author: "zhangsan"
created_at: "2026-03-10T08:30:00Z"
state: "open"
labels: ["enhancement", "ui"]
assignees: ["lisi", "wangwu"]
milestone: "v1.0.0"
---
```

**已关闭 Issue 的 Front Matter**:
```yaml
---
title: "Fix login bug"
url: "https://github.com/issue2md/issue2md/pull/42"
author: "lisi"
created_at: "2026-03-08T14:00:00Z"
state: "closed"
closed_at: "2026-03-09T10:00:00Z"
labels: ["bug"]
assignees: []
milestone: null
---
```

---

### 5.2 完整 Markdown 输出示例

```markdown
---
title: "Add support for dark mode"
url: "https://github.com/issue2md/issue2md/issues/1"
author: "zhangsan"
created_at: "2026-03-10T08:30:00Z"
state: "open"
labels: ["enhancement", "ui"]
assignees: ["lisi", "wangwu"]
milestone: "v1.0.0"
---

Add support for dark mode in the application.

## Requirements

- [ ] Add dark mode toggle
- [ ] Persist user preference
- [ ] Support system-level dark mode detection

### Design Mockup

![Dark Mode Mockup](https://user-images.githubusercontent.com/123/456.png)

---

### @lisi - 2026-03-10T09:00:00Z

I can help with this! Let me know if you need any assistance.

---

### @wangwu - 2026-03-10T10:30:00Z

Here's a reference implementation from another project:

```go
func isDarkMode() bool {
    // implementation
}
```

---

### @zhangsan - 2026-03-10T11:00:00Z

@lisi That would be great! I'll send you the design specs.
```

---

### 5.3 `--user-links` 开启时的评论格式

```markdown
### [@lisi](https://github.com/lisi) - 2026-03-10T09:00:00Z

I can help with this! Let me know if you need any assistance.

---

### [@wangwu](https://github.com/wangwu) - 2026-03-10T10:30:00Z

Here's a reference implementation from another project:
```

---

## 六、依赖关系

| 依赖项 | 说明 |
|--------|------|
| Go 1.24+ | 运行环境 |
| GitHub REST API v3 | Issue/PR 数据获取 |
| GitHub GraphQL API | Discussion 数据获取 |

---

## 七、修订历史

| 版本 | 日期 | 作者 | 变更内容 |
|------|------|------|----------|
| 1.0 | 2026-03-16 | AI Agent | 初始版本，包含核心功能需求 |
