# Implementation Plan: 时区时间查询与跨时区会议安排（Timezone Time Lookup & Meeting Scheduler）

**Branch**: `001-timezone-meeting-scheduler` | **Date**: 2026-09-08 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/001-timezone-meeting-scheduler/spec.md`

## Summary

构建一个跨平台单二进制命令行工具：用户可通过时区标识符、地点名称（中/英文）或邮政编码查询某地当前日期与时间，支持跨时区时间换算与多地点对比，并提供跨时区会议排程（提议时间→各参与者当地时间、重叠可用时段查找）。技术方案采用 Go 语言，依赖标准库时间与时区能力、内嵌静态地点数据，通过交叉编译生成适用于 Linux/macOS/Windows 及 amd64/arm64 等平台与架构的版本。

## Technical Context

**Language/Version**: Go（Golang），最新稳定版本（≥1.22）

**Primary Dependencies**:
- Go 标准库 `time` + `time/tzdata`（内嵌 IANA 时区数据，保证离线与 DST 正确性）
- 内嵌静态地点↔时区↔邮编数据集（通过 `embed` 打包，离线可用）
- CLI 解析优先使用标准库 `flag`（若子命令复杂度增加，再评估引入 CLI 框架）

**Storage**: N/A（无状态单二进制；数据集以 `embed` 内嵌）

**Testing**: Go 标准库 `testing`（表驱动单元测试 + 端到端集成测试）

**Target Platform**: Linux / macOS / Windows；amd64 / arm64（交叉编译）

**Project Type**: CLI（命令行工具，单二进制）

**Performance Goals**: 单次查询在本地运行 < 1 秒返回（规范要求 3 秒内）；二进制启动快速

**Constraints**: 单静态二进制、离线可用（内嵌时区与地点数据）、可交叉编译至多 OS/架构

**Scale/Scope**: 单用户本地 CLI，无服务端；内置常见国际城市 + 中国主要城市数据集

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| 宪法原则 | 计划符合性 |
|----------|-----------|
| 一、代码质量优先 | 采用 Go 惯用分层（cmd/internal）；`go fmt`、`go vet`、静态检查纳入 CI；单一职责、公开 API 文档 |
| 二、测试标准 | TDD：先写表驱动测试；核心时间换算逻辑注入可控制时钟（`internal/clock`）保证确定性；单元 + 集成测试，覆盖率阈值 |
| 三、用户体验一致性 | 统一 CLI 输出格式（文本/JSON）、统一错误信息与非零退出码、全局一致参数约定（见 contracts/cli-contract.md） |
| 四、性能要求 | 明确 <1s 查询预算；基准测试纳入 CI；内嵌数据 + 线性查找避免低效；无 N+1 类问题 |

**Gate 结论**: 通过。无违反项，无需复杂性豁免。

**Post-design 复核**: Phase 1 的 data-model 与 contracts 保持与宪法一致性、可测试性要求相符，无新增违反项。

## Project Structure

### Documentation (this feature)

```text
specs/001-timezone-meeting-scheduler/
├── plan.md              # 本文件（/speckit-plan 输出）
├── research.md          # Phase 0 输出（/speckit-plan 输出）
├── data-model.md        # Phase 1 输出（/speckit-plan 输出）
├── quickstart.md        # Phase 1 输出（/speckit-plan 输出）
├── contracts/           # Phase 1 输出（/speckit-plan 输出）
│   └── cli-contract.md  # CLI 接口契约
└── tasks.md             # Phase 2 输出（/speckit-tasks 生成，非本命令创建）
```

### Source Code (repository root)

```text
cmd/
└── timezone-utility/
    └── main.go          # 入口，装配并启动 CLI

internal/
├── cli/                 # 子命令解析与分发（now / convert / meeting / overlap）
├── lookup/              # 时区/地点/邮编查询核心逻辑
├── convert/             # 跨时区时间换算
├── meeting/             # 会议本地时间展示与重叠时段查找
├── data/                # 内嵌数据集加载与解析（地点↔时区↔邮编）
├── output/              # 统一输出格式（文本/JSON）
└── clock/               # 可注入时钟（测试确定性）

data/                    # 内嵌静态数据集（cities、postcodes）
tests/
└── integration/         # 端到端集成测试（构建并调用真实二进制）
```

**Structure Decision**: 采用 Go 社区惯用的 `cmd/` + `internal/` 单仓库布局；查询、换算、会议、输出各自独立包，遵循单一职责；`internal/clock` 抽象时钟以满足宪法"测试确定性"要求。

## Complexity Tracking

> 无宪法违规，无需记录。
