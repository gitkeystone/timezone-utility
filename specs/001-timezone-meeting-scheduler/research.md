# Phase 0 Research: 时区时间查询与跨时区会议安排

## 决策摘要

| # | 主题 | 决策 | 理由 |
|---|------|------|------|
| R1 | 时区数据处理 | 使用 Go 标准库 `time` + `time/tzdata` 内嵌 IANA 数据库 | 离线、DST 正确、零第三方依赖 |
| R2 | 地点/邮编→时区映射 | 内嵌精选静态数据集（城市 + 中国邮编 + 国际邮编） | 单二进制、离线、数据可控 |
| R3 | CLI 解析 | 优先标准库 `flag`；复杂度增长时评估 cobra/urfave | YAGNI，最小依赖 |
| R4 | 跨平台发布 | GOOS/GOARCH 交叉编译 + goreleaser（或 Makefile） | 生成多 OS/架构版本 |
| R5 | 时钟注入 | `internal/clock` 接口，默认系统时钟，测试注入固定时钟 | 确定性测试 |

## R1 时区数据处理

- Decision: 使用 `time` 包配合 `import _ "time/tzdata"`，将 IANA 时区数据库静态链接进二进制。
- Rationale: 无需系统时区数据库或网络，二进制可在任意环境离线正确计算本地时间与夏令时。
- Alternatives considered:
  - 依赖宿主系统 tzdata —— 跨平台不一致、可能缺失。
  - 调用在线时区 API —— 需网络且引入延迟与外部依赖。

## R2 地点/邮编→时区映射

- Decision: 内嵌精选静态数据集（JSON，`go:embed`），覆盖常见国际城市、中国主要城市及其邮编；英文名与中文别名映射到 IANA 时区。
- Rationale: 满足规范"离线可用"假设与"覆盖 90% 常见城市"目标，无 API 密钥与网络依赖。
- Alternatives considered:
  - 在线地理编码 API（如 Nominatim）—— 需网络、限流、密钥。
  - 完整地名库（GeoNames）—— 数据量大、二进制膨胀。
  - 用户自定义映射文件 —— 作为后续扩展项，v1 内置即可。

## R3 CLI 解析

- Decision: 优先使用标准库 `flag` 与子命令分发；输出格式（文本/JSON）通过全局标志控制。
- Rationale: 遵循宪法"简单优先/YAGNI"，v1 子命令数量有限，标准库即可满足。
- Alternatives considered:
  - cobra / urfave/cli —— 功能丰富但引入第三方依赖；仅当子命令与补全需求显著增长时再引入。

## R4 跨平台发布

- Decision: 使用 Go 内建交叉编译（GOOS/GOARCH）产出 linux/darwin/windows × amd64/arm64 组合；以 goreleaser 或 Makefile 脚本统一构建、校验与发布。
- Rationale: Go 原生支持交叉编译，零运行时依赖；goreleaser 是社区标准的发布工具。
- Alternatives considered:
  - 手写各平台构建脚本 —— 可行但重复、易漏。
  - CI 矩阵 —— 与 goreleaser 互补而非替代。

## R5 时钟注入

- Decision: 定义 `internal/clock` 接口（`Now() time.Time`），生产用系统时钟，测试注入固定时钟。
- Rationale: 满足宪法"测试必须确定性、时间相关逻辑必须注入可控制时钟"。
- Alternatives considered:
  - 直接调用 `time.Now()` —— 测试无法固定，DST/日期边界用例不可重复。
