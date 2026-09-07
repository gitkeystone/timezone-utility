# timezone-utility

一个单二进制命令行工具，用于查询任意地点的当前日期与时间，并支持跨时区会议安排。

## 功能

- 按时区、地点名称（中/英文）或邮政编码查询当前日期与时间
- 跨时区时间换算与多地点对比
- 跨时区会议时间展示与重叠可用时段查找
- 内嵌 IANA 时区数据与地点数据，完全离线可用
- 单二进制、零运行时依赖，支持交叉编译

## 安装与构建

需要 Go ≥ 1.22。

```bash
go build -o timezone-utility ./cmd/timezone-utility
```

## 用法

```bash
# 按时区查询
timezone-utility now Asia/Shanghai

# 按地点查询（中英文）
timezone-utility now 北京
timezone-utility now Tokyo

# 按邮编查询
timezone-utility now 100000     # 北京
timezone-utility now 10001      # 纽约

# 一次查询多个地点
timezone-utility now 北京 上海 东京

# 跨时区时间换算
timezone-utility convert --from 北京 --to 纽约 --to 伦敦 "2026-09-08 15:00"

# 会议时间展示
timezone-utility meeting --from 北京 "2026-09-10 09:00" 伦敦 纽约

# 查找重叠可用时段
timezone-utility overlap --window "北京=16:00-20:00" --window "伦敦=09:00-12:00"

# JSON 输出（--json 可出现在任意位置）
timezone-utility now 北京 --json
```

退出码：`0` 成功；`1` 用法错误；`2` 查询/解析失败。

## 跨平台构建

```bash
# 使用 Makefile 一次性构建全部平台
make cross

# 或手动交叉编译
GOOS=windows GOARCH=amd64 go build -o timezone-utility.exe ./cmd/timezone-utility
GOOS=darwin  GOARCH=arm64  go build -o timezone-utility     ./cmd/timezone-utility

# 或使用 goreleaser 生成发布产物
goreleaser release --snapshot
```

支持 linux / darwin / windows × amd64 / arm64。

## 测试与质量

```bash
go test ./...        # 运行全部单元与集成测试
go vet ./...         # 静态检查
make coverage        # 查看覆盖率报告
make check-coverage  # 强制整体覆盖率 ≥ 70% 阈值（不达标即失败）
```

覆盖率阈值（项目宪法"测试标准"原则）：整体语句覆盖率不得低于 70%。

## 设计文档

- 功能规范: `specs/001-timezone-meeting-scheduler/spec.md`
- 实现计划: `specs/001-timezone-meeting-scheduler/plan.md`
- 任务清单: `specs/001-timezone-meeting-scheduler/tasks.md`
