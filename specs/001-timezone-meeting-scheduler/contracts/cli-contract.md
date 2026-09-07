# CLI 接口契约（CLI Contract）

本文档定义 `timezone-utility` 的命令行接口，是"用户体验一致性"原则的落地契约。

## 全局约定

- 命令格式：`timezone-utility <subcommand> [flags] [args...]`（--json 可出现在任意位置；其余标志须在位置参数之前）
- 输出格式：默认人可读文本；`--json` 输出结构化 JSON。
- 退出码：`0` 成功；`1` 用法错误；`2` 查询/解析失败（无此地点/时区/邮编）。
- 错误信息统一格式 `error: <原因>（含可操作建议）`，写入 stderr。
- 未提供子命令或参数时，输出本机当前日期与时间及其时区。
- 时间输入格式：`YYYY-MM-DD HH:mm` 或 RFC 3339（`2026-09-08T15:00:00+08:00`）。

## 子命令

### now —— 查询当前日期时间

`timezone-utility now [query...]`

- query 可为：IANA 时区（`UTC`、`Asia/Shanghai`）、地点名（`Beijing`、`北京`）、邮编（`100000`、`10001`）。
- 多个 query 同时列出各地点当前时间；无 query 时返回本机时间。

### convert —— 跨时区时间换算

`timezone-utility convert --from <loc> --to <loc> [--to <loc> ...] <time>`

- 将 `<time>`（在 `--from` 地点的当地时间）换算为各 `--to` 地点的当地时间。

### meeting —— 会议时间展示

`timezone-utility meeting --from <loc> <time> <loc> [<loc> ...]`

- 展示提议时间 `<time>`（在 `--from` 地点的当地时间）对应的每个参与者地点当地时间与日期。

### overlap —— 重叠可用时段查找

`timezone-utility overlap --window <loc>=<HH:mm>-<HH:mm> [--window ...]`

- 找出所有参与者可用时间窗口的重叠时段，并按各自本地时间展示。

## 统一 JSON 输出结构

```json
{
  "results": [
    {
      "query": "北京",
      "resolved": {
        "name": "Beijing",
        "country": "CN",
        "timezone": "Asia/Shanghai",
        "postal_codes": ["100000"]
      },
      "local_date": "2026-09-08",
      "local_time": "10:00:00",
      "utc_offset": "+08:00",
      "dst": false
    }
  ]
}
```

`meeting` 与 `overlap` 在此基础上扩展 `participants` / `overlap_candidates` 字段，保持字段命名与结构一致。
