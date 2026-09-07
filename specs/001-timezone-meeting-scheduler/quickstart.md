# Quickstart 验证指南

本指南用于端到端验证功能，非完整实现说明。实现细节见 `tasks.md` 与实现阶段。

## 前置条件

- 已安装 Go（≥1.22）。
- 构建二进制：`go build ./cmd/timezone-utility`。

## 验证场景

### 1. 按时区查询（US1）

命令：`./timezone-utility now Asia/Shanghai`
预期：返回上海当前日期与时间，标注 `Asia/Shanghai`、`+08:00`。

### 2. 按地点查询（US2）

命令：`./timezone-utility now 北京`
预期：返回北京当前日期与时间，时区 `Asia/Shanghai`。

### 3. 按邮编查询（US3）

命令：`./timezone-utility now 10001`
预期：返回纽约当前日期与时间。

### 4. 跨时区时间换算（US4）

命令：`./timezone-utility convert --from 北京 --to 纽约 --to 伦敦 "2026-09-08 15:00"`
预期：返回纽约与伦敦在该时刻的当地时间。

### 5. 会议时间展示（US5）

命令：`./timezone-utility meeting --from 北京 "2026-09-10 09:00" 伦敦 纽约`
预期：返回各参与者当地时间与日期。

### 6. 重叠时段查找（US5）

命令：`./timezone-utility overlap --window "北京=09:00-12:00" --window "伦敦=09:00-12:00"`
预期：返回双方空闲的重叠时段。

### 7. 错误处理

命令：`./timezone-utility now Asia/Atlantis`
预期：stderr 输出清晰错误，退出码为 `2`。

### 8. 结构化输出

命令：`./timezone-utility now 北京 --json`
预期：输出符合 `contracts/cli-contract.md` 的 JSON 结构。

## 运行测试与检查

- `go test ./...` —— 运行全部单元与集成测试。
- `go vet ./...` —— 静态检查。
- `go test -cover ./...` —— 查看覆盖率（满足覆盖率阈值）。

## 跨平台构建

- 单平台：`go build ./cmd/timezone-utility`
- 多平台：`GOOS=linux GOARCH=amd64 go build ...`（同理 darwin/windows × amd64/arm64），或使用 goreleaser 统一发布。
