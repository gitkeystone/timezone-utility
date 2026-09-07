# Data Model: 时区时间查询与跨时区会议安排

## 实体（Entities）

### Location（地点）

代表一个可查询的地理位置。

| 字段 | 类型 | 说明 |
|------|------|------|
| id | string | 稳定唯一标识 |
| name | string | 规范英文名（唯一） |
| aliases | string[] | 别名（含中文名，如 `北京`） |
| country | string | 国家/地区代码（如 `CN`、`US`） |
| timezone | string | IANA 时区标识符（如 `Asia/Shanghai`） |
| postal_codes | string[] | 关联邮编（一个地点可对应多个邮编） |

### Timezone（时区）

| 字段 | 类型 | 说明 |
|------|------|------|
| iana_id | string | IANA 时区标识符 |
| utc_offset | string | 当前 UTC 偏移（如 `+08:00`，随 DST 变化） |
| dst | bool | 查询时刻是否处于夏令时 |

### QueryResult（查询结果）

| 字段 | 类型 | 说明 |
|------|------|------|
| query | string | 用户原始输入 |
| resolved_location | Location? | 解析后的地点（仅按时区查询时可空） |
| local_date | string | 当地日期（`YYYY-MM-DD`） |
| local_time | string | 当地时间（`HH:mm:ss`） |
| timezone | string | 时区标识符 |
| utc_offset | string | UTC 偏移 |
| dst | bool | 是否处于夏令时 |

### Meeting（会议）

| 字段 | 类型 | 说明 |
|------|------|------|
| proposed_time | time | 提议时间（含来源时区） |
| participants | Location[] | 参与者地点列表 |
| local_times | QueryResult[] | 各参与者当地时间 |
| availability_windows | Window[] | 各参与者可用时间窗口 |
| overlap_candidates | Window[] | 重叠候选时段 |

其中 `Window = { location, start, end }`。

## 关系

- Location 1—1 Timezone（每个地点属于一个时区）。
- Location 1—N PostalCode（一个地点可对应多个邮编）。
- Meeting 1—N Location（一次会议含多个参与者地点）。

## 校验规则（源自功能需求）

- 时区标识符必须是有效的 IANA 时区（FR-001、FR-010）。
- 邮编格式：中国为 6 位数字；国际按对应国家规则校验（FR-003、FR-010）。
- 地点名称必须能唯一解析；歧义时返回可操作提示而非静默猜测（FR-010）。
- 查询结果必须包含当地日期、当地时间、时区标识符、UTC 偏移（FR-004）。
- 偏移量必须反映查询时刻当前生效的 DST 状态（FR-005）。
