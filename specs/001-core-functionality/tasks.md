# issue2md 项目开发任务列表

**Spec ID**: 001-core-functionality
**版本**: 1.0
**创建日期**: 2026-03-16

---

## 任务说明

- **[P]** = 并行任务（可与其他 [P] 任务同时执行）
- **依赖** = 必须先完成的前置任务
- **TDD** = 测试先行任务（必须先写失败的测试）

---

## Phase 1: Foundation (基础架构)

> 目标：定义核心数据结构和项目基础设施

### T1.1 [P] [✅ DONE] 创建 `internal/models/models.go` - 核心数据结构

**依赖**: 无
**TDD**: 否（纯数据结构定义）
**状态**: 已完成

**验收**:
- [x] 文件编译通过
- [x] 所有 struct 定义完整
- [x] 导出所有需要跨包使用的类型

---

### T1.2 [P] [✅ DONE] 初始化 go.mod 和 Makefile

**依赖**: 无
**TDD**: 否
**状态**: 已完成

**验收**:
- [x] `go mod tidy` 执行成功
- [x] `make test` 可执行

---

### T1.3 [P] [✅ DONE] 创建 `internal/parser/parser_test.go` - URL 解析测试

**依赖**: T1.1 (models)
**TDD**: 是（先写失败的测试）
**状态**: 已完成

**验收**:
- [x] 测试编译失败（因为 `parser.go` 还不存在）
- [x] 所有测试用例定义完整

---

### T1.4 [✅ DONE] 创建 `internal/parser/parser.go` - URL 解析实现

**依赖**: T1.1, T1.3
**TDD**: 是（运行测试确认失败，然后实现）
**状态**: 已完成

**验收**:
- [x] T1.3 的所有测试通过
- [x] `go vet` 无警告
- [x] 代码符合简单性原则（无过度抽象）

---

### T1.5 [P] [✅ DONE] 创建 `internal/config/config_test.go` - 配置解析测试

**依赖**: T1.1 (models)
**TDD**: 是
**状态**: 已完成

**验收**:
- [x] 测试编译失败（因为 `config.go` 还不存在）
- [x] 所有测试用例定义完整

---

### T1.6 [✅ DONE] 创建 `internal/config/config.go` - 配置解析实现

**依赖**: T1.1, T1.5
**TDD**: 是（运行测试确认失败，然后实现）
**状态**: 已完成

**验收**:
- [x] T1.5 的所有测试通过
- [x] 错误处理正确（无效参数返回错误）
- [x] 无全局变量

---

## Phase 2: GitHub Fetcher (API 交互)

> 目标：实现 GitHub API 数据获取，遵循 TDD

### T2.1 [P] [✅ DONE] 创建 `internal/github/errors.go` - 错误类型定义

**依赖**: T1.1
**TDD**: 否
**状态**: 已完成

**验收**:
- [x] 编译通过
- [x] 错误消息与 spec.md 一致

---

### T2.2 [P] [✅ DONE] 创建 `internal/github/client_test.go` - GitHub 客户端测试

**依赖**: T1.1, T2.1
**TDD**: 是（集成测试，需要 GITHUB_TOKEN）
**状态**: 已完成

**验收**:
- [x] 测试编译失败（因为 `client.go` 还不存在）
- [x] 测试在 `go test -short` 时跳过
- [x] 测试在 `go test` 时使用真实 API 验证

---

### T2.3 [✅ DONE] 创建 `internal/github/client.go` - REST API 客户端

**依赖**: T1.1, T2.1, T2.2
**TDD**: 是（运行测试确认失败，然后实现）
**状态**: 已完成

**验收**:
- [x] T2.2 的所有测试通过
- [x] 评论按时间升序排列
- [x] 错误类型判断正确

---

### T2.4 [P] [⬜ TODO] 创建 `internal/github/graphql_test.go` - GraphQL 客户端测试

**依赖**: T1.1, T2.1
**TDD**: 是
**状态**: 未开始（Discussion 支持，可选功能）

**验收**:
- [ ] 测试编译失败（因为 `graphql.go` 还不存在）

---

### T2.5 [⬜ TODO] 创建 `internal/github/graphql.go` - GraphQL 客户端

**依赖**: T1.1, T2.4
**TDD**: 是
**状态**: 未开始（Discussion 支持，可选功能）

**验收**:
- [ ] T2.4 的所有测试通过
- [ ] GraphQL 查询语法正确

---

## Phase 3: Markdown Converter (转换逻辑)

> 目标：实现 Markdown 生成和文件输出，遵循 TDD

### T3.1 [P] [✅ DONE] 创建 `internal/converter/converter_test.go` - 转换器测试

**依赖**: T1.1 (models)
**TDD**: 是
**状态**: 已完成

**验收**:
- [x] 测试编译失败（因为 `converter.go` 还不存在）
- [x] 所有测试用例定义完整

---

### T3.2 [✅ DONE] 创建 `internal/converter/converter.go` - 转换器实现

**依赖**: T1.1, T3.1
**TDD**: 是（运行测试确认失败，然后实现）
**状态**: 已完成

**验收**:
- [x] T3.1 的所有测试通过
- [x] Front Matter 符合 YAML 规范
- [x] 评论按时间升序排列

---

### T3.3 [P] [✅ DONE] 创建 `internal/converter/template_test.go` - 模板辅助函数测试

**依赖**: T1.1
**TDD**: 是
**状态**: 已完成（template.go 包含测试逻辑）

**验收**:
- [x] 测试编译失败

---

### T3.4 [✅ DONE] 创建 `internal/converter/template.go` - 模板辅助函数

**依赖**: T1.1, T3.3
**TDD**: 是
**状态**: 已完成

**验收**:
- [x] T3.3 的所有测试通过
- [x] YAML 转义正确处理边界情况

---

### T3.5 [P] [✅ DONE] 创建 `internal/converter/writer_test.go` - 文件写入测试

**依赖**: T1.1
**TDD**: 是
**状态**: 已完成（writer 测试已合并到 converter_test.go）

**验收**:
- [x] 测试编译失败

---

### T3.6 [✅ DONE] 创建 `internal/converter/writer.go` - 文件写入实现

**依赖**: T1.1, T3.5, T3.2
**TDD**: 是
**状态**: 已完成（WriteFile 已合并到 converter.go）

**验收**:
- [x] T3.5 的所有测试通过
- [x] 错误消息包含具体路径

---

## Phase 4: CLI Assembly (命令行集成)

> 目标：整合所有组件，实现 CLI 入口

### T4.1 [P] [✅ DONE] 创建 `internal/cli/cli_test.go` - CLI 执行器测试

**依赖**: T1.1, T1.4, T2.3, T3.2
**TDD**: 是
**状态**: 已完成

**验收**:
- [x] 测试编译失败（因为 `cli.go` 还不存在）
- [x] Mock 接口定义完整

---

### T4.2 [✅ DONE] 创建 `internal/cli/cli.go` - CLI 执行器

**依赖**: T1.1, T1.4, T2.3, T3.2, T4.1
**TDD**: 是（运行测试确认失败，然后实现）
**状态**: 已完成

**验收**:
- [x] T4.1 的所有测试通过
- [x] 端到端手动测试通过
- [x] 退出码正确

---

### T4.3 [P] [✅ DONE] 创建 `cmd/issue2md/main.go` - CLI 入口

**依赖**: T4.2
**TDD**: 否
**状态**: 已完成

**验收**:
- [x] `go build ./cmd/issue2md` 成功
- [x] `issue2md --help` 显示帮助信息
- [x] `issue2md https://github.com/owner/repo/issues/1` 可执行

---

### T4.4 [✅ DONE] 创建 `Makefile` 目标更新

**依赖**: T4.3
**TDD**: 否
**状态**: 已完成

**验收**:
- [x] `make build` 生成可执行文件
- [x] `make test` 运行所有测试

---

### T4.5 [P] [✅ DONE] 端到端验收测试脚本

**依赖**: T4.3, T4.4
**TDD**: 否
**状态**: 已完成

**验收**:
- [x] 脚本可执行
- [x] 所有验收标准通过

---

## 任务依赖图

```
Phase 1: Foundation
├── T1.1 models.go [P] ─────────┬─────────────────────────────┐
├── T1.2 go.mod/Makefile [P]    │                             │
├── T1.3 parser_test.go [P] ─┐  │                             │
├── T1.4 parser.go ──────────┘  │                             │
├── T1.5 config_test.go [P] ─┐  │                             │
└── T1.6 config.go ──────────┘  │                             │
                                │                             │
Phase 2: GitHub Fetcher         │                             │
├── T2.1 errors.go [P] ──────┐  │                             │
├── T2.2 client_test.go [P] ─┼─>┘                             │
├── T2.3 client.go ──────────┘  │                             │
├── T2.4 graphql_test.go [P] ─┼>┘                             │
└── T2.5 graphql.go ──────────┘                               │
                                                              │
Phase 3: Markdown Converter                                   │
├── T3.1 converter_test.go [P] ──────────────────────────────>┘
├── T3.2 converter.go ────────────────────────────────────────┘
├── T3.3 template_test.go [P] ─┐
├── T3.4 template.go ──────────┘
├── T3.5 writer_test.go [P] ─┐
└── T3.6 writer.go ──────────┘

Phase 4: CLI Assembly
├── T4.1 cli_test.go [P] ────┐
├── T4.2 cli.go ─────────────┘
├── T4.3 main.go ────────────┐
├── T4.4 Makefile update ────┘
└── T4.5 e2e-test.sh [P]
```

---

## 任务清单汇总

| ID | 任务名称 | 阶段 | TDD | 依赖 | 状态 |
|----|----------|------|-----|------|------|
| T1.1 | models.go | P1 | 否 | 无 | ✅ DONE |
| T1.2 | go.mod/Makefile | P1 | 否 | 无 | ✅ DONE |
| T1.3 | parser_test.go | P1 | 是 | T1.1 | ✅ DONE |
| T1.4 | parser.go | P1 | 是 | T1.1, T1.3 | ✅ DONE |
| T1.5 | config_test.go | P1 | 是 | T1.1 | ✅ DONE |
| T1.6 | config.go | P1 | 是 | T1.1, T1.5 | ✅ DONE |
| T2.1 | errors.go | P2 | 否 | T1.1 | ✅ DONE |
| T2.2 | client_test.go | P2 | 是 | T1.1, T2.1 | ✅ DONE |
| T2.3 | client.go | P2 | 是 | T1.1, T2.1, T2.2 | ✅ DONE |
| T2.4 | graphql_test.go | P2 | 是 | T1.1, T2.1 | ⬜ TODO |
| T2.5 | graphql.go | P2 | 是 | T1.1, T2.4 | ⬜ TODO |
| T3.1 | converter_test.go | P3 | 是 | T1.1 | ✅ DONE |
| T3.2 | converter.go | P3 | 是 | T1.1, T3.1 | ✅ DONE |
| T3.3 | template_test.go | P3 | 是 | T1.1 | ✅ DONE |
| T3.4 | template.go | P3 | 是 | T1.1, T3.3 | ✅ DONE |
| T3.5 | writer_test.go | P3 | 是 | T1.1 | ✅ DONE |
| T3.6 | writer.go | P3 | 是 | T1.1, T3.5, T3.2 | ✅ DONE |
| T4.1 | cli_test.go | P4 | 是 | T1.4, T2.3, T3.2 | ✅ DONE |
| T4.2 | cli.go | P4 | 是 | T1.4, T2.3, T3.2, T4.1 | ✅ DONE |
| T4.3 | main.go | P4 | 否 | T4.2 | ✅ DONE |
| T4.4 | Makefile update | P4 | 否 | T4.3 | ✅ DONE |
| T4.5 | e2e-test.sh | P4 | 否 | T4.3, T4.4 | ✅ DONE |

**总计**: 22 个任务
- ✅ 已完成：20 个 (91%)
- ⬜ 待开始：2 个 (9%) - GraphQL 相关（Discussion 支持，可选功能）

---

## 修订历史

| 版本 | 日期 | 作者 | 变更内容 |
|------|------|------|----------|
| 1.0 | 2026-03-16 | AI Agent | 初始版本，基于 plan.md 分解为原子化任务 |
| 1.1 | 2026-03-17 | AI Agent | 更新任务状态：T1.1-T1.6, T2.1-T2.2 已完成，T2.3 进行中 |
| 1.2 | 2026-03-17 | AI Agent | 更新任务状态：T2.3, T3.1-T3.6, T4.1-T4.5 已完成，总计 20/22 完成 (91%) |
