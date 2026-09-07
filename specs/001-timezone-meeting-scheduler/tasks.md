---

description: "Implementation tasks for 时区时间查询与跨时区会议安排 (timezone-meeting-scheduler)"
---

# Tasks: 时区时间查询与跨时区会议安排（Timezone Time Lookup & Meeting Scheduler）

**Input**: Design documents from `/specs/001-timezone-meeting-scheduler/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: 依据项目宪法"测试标准"原则（TDD 强制），本任务清单包含测试任务；每条用户故事先写测试、确认失败，再实现。

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- Go 单仓库：`cmd/timezone-utility/`（入口）+ `internal/`（各包）+ `data/`（内嵌数据集）+ `tests/integration/`（集成测试）
- 单元测试与包同目录（`*_test.go`）；集成测试在 `tests/integration/`

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [X] T001 Initialize Go module: create `go.mod` (module `timezone-utility`, Go ≥1.22) and `go.sum`
- [X] T002 [P] Create directory skeleton per plan.md: `cmd/timezone-utility/`, `internal/{cli,lookup,convert,meeting,data,output,clock}/`, `data/`, `tests/integration/`
- [X] T003 [P] Add `Makefile` with targets `build`, `test`, `vet`, `fmt`, `cross` (multi-OS/arch) and add `.gitignore`
- [X] T004 [P] Add `.goreleaser.yml` for linux/darwin/windows × amd64/arm64 release builds

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T005 Define `Clock` interface plus `RealClock` and `FixedClock` implementations in `internal/clock/clock.go`
- [X] T006 [P] Define data-model types `Location`, `Timezone`, `QueryResult`, `Meeting`, `Window` in `internal/data/types.go`
- [X] T007 [P] Implement output formatter (human-readable text + JSON) and unified error format in `internal/output/output.go`
- [X] T008 Implement embedded dataset loading in `internal/data/load.go` plus seed files `data/cities.json` and `data/postcodes.json`
- [X] T009 Implement CLI entry + subcommand dispatch skeleton (`now`/`convert`/`meeting`/`overlap`, `--json`, `--version`, `--help`) in `cmd/timezone-utility/main.go` and `internal/cli/cli.go`; wire `_ "time/tzdata"` import for embedded IANA tzdata

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - 按时区查询当前日期时间 (Priority: P1) 🎯 MVP

**Goal**: 用户通过标准时区标识符（如 `UTC`、`Asia/Shanghai`、`America/New_York`）查询该时区当前日期与时间。

**Independent Test**: 运行 `timezone-utility now Asia/Shanghai`，返回上海当前日期、当地时间、时区标识符与 `+08:00` 偏移。

### Tests for User Story 1 ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [X] T010 [P] [US1] Unit test: timezone → current time with fixed clock (verify `Asia/Shanghai` date/time/offset) in `internal/lookup/lookup_test.go`
- [X] T011 [P] [US1] Unit test: invalid timezone returns error in `internal/lookup/lookup_test.go`
- [X] T012 [US1] Integration test: `now Asia/Shanghai` end-to-end in `tests/integration/now_test.go`

### Implementation for User Story 1

- [X] T013 [US1] Implement timezone resolution (validate + `time.LoadLocation`) in `internal/lookup/resolve.go`
- [X] T014 [US1] Implement current-time lookup for a timezone using `Clock` in `internal/lookup/current.go`
- [X] T015 [US1] Implement `now` subcommand (single timezone) wiring resolver + formatter in `internal/cli/now.go`
- [X] T016 [US1] Wire unified error handling and exit code `2` for invalid timezone in `internal/cli/now.go`

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - 按地点名称查询当前日期时间 (Priority: P1)

**Goal**: 用户通过地点名称（中/英文，如 `北京`/`Beijing`、`东京`/`Tokyo`）查询该地点当前日期与时间。

**Independent Test**: 运行 `timezone-utility now 北京`，返回北京当前日期与时间（时区 `Asia/Shanghai`）。

### Tests for User Story 2 ⚠️

- [X] T017 [P] [US2] Unit test: place alias resolution (`北京`/`Beijing` → `Asia/Shanghai`) in `internal/lookup/resolve_test.go`
- [X] T018 [P] [US2] Unit test: ambiguous / unrecognized place name returns actionable error in `internal/lookup/resolve_test.go`
- [X] T019 [US2] Integration test: `now 北京` end-to-end in `tests/integration/now_test.go`

### Implementation for User Story 2

- [X] T020 [US2] Extend resolver: place name (en + zh alias) → timezone in `internal/lookup/resolve.go`
- [X] T021 [US2] Add place alias seed data (北京/上海/东京/纽约 etc.) to `data/cities.json`
- [X] T022 [US2] Extend `now` to accept place names in `internal/cli/now.go`

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - 按邮政编码查询当前日期时间 (Priority: P2)

**Goal**: 用户通过邮政编码（中国六位邮编与国际邮编，如 `100000`、`10001`）查询对应地点当前日期与时间。

**Independent Test**: 运行 `timezone-utility now 100000`，返回北京当前日期与时间。

### Tests for User Story 3 ⚠️

- [X] T023 [P] [US3] Unit test: postal resolution (`100000` → 北京, `10001` → 纽约) in `internal/lookup/resolve_test.go`
- [X] T024 [P] [US3] Unit test: invalid postal code returns error in `internal/lookup/resolve_test.go`
- [X] T025 [US3] Integration test: `now 100000` end-to-end in `tests/integration/now_test.go`

### Implementation for User Story 3

- [X] T026 [US3] Extend resolver: postal code → location → timezone in `internal/lookup/resolve.go`
- [X] T027 [US3] Add postal seed data (中国六位邮编 + 国际邮编) to `data/postcodes.json`
- [X] T028 [US3] Extend `now` to accept postal codes in `internal/cli/now.go`

**Checkpoint**: At this point, User Stories 1–3 should all work independently

---

## Phase 6: User Story 4 - 跨时区时间换算与多地点对比 (Priority: P2)

**Goal**: 用户可将指定时间点从某地换算到其他地点，也可一次查询多个地点并排对比当地时间。

**Independent Test**: 运行 `convert "2026-09-08 15:00" --from 北京 --to 纽约 伦敦` 与 `now 北京 上海 东京`，返回正确的换算结果与多地点列表。

### Tests for User Story 4 ⚠️

- [X] T029 [P] [US4] Unit test: time conversion (`北京 15:00` → `纽约`/`伦敦`) in `internal/convert/convert_test.go`
- [X] T030 [P] [US4] Unit test: DST-aware conversion at boundary in `internal/convert/convert_test.go`
- [X] T031 [US4] Integration test: `convert` command in `tests/integration/convert_test.go`
- [X] T032 [US4] Integration test: `now 北京 上海 东京` multi-location in `tests/integration/now_test.go`

### Implementation for User Story 4

- [X] T033 [US4] Implement time conversion logic in `internal/convert/convert.go`
- [X] T034 [US4] Implement `convert` subcommand (`--from`/`--to`) in `internal/cli/convert.go`
- [X] T035 [US4] Extend `now` to support multiple queries (multi-location listing) in `internal/cli/now.go`

**Checkpoint**: At this point, User Stories 1–4 should all work independently

---

## Phase 7: User Story 5 - 跨时区会议安排 (Priority: P3)

**Goal**: 用户提供参与者地点与提议时间，工具展示各参与者当地时间；也可提供可用时间窗口，找出所有参与者时间重叠的候选时段。

**Independent Test**: 运行 `meeting "2026-09-10 09:00" --from 北京 伦敦 纽约` 与 `overlap --window "北京=09:00-12:00" --window "伦敦=09:00-12:00"`，返回正确的本地时间与重叠时段。

### Tests for User Story 5 ⚠️

- [X] T036 [P] [US5] Unit test: meeting local times for participants in `internal/meeting/meeting_test.go`
- [X] T037 [P] [US5] Unit test: overlap window finding in `internal/meeting/overlap_test.go`
- [X] T038 [US5] Integration test: `meeting` command in `tests/integration/meeting_test.go`
- [X] T039 [US5] Integration test: `overlap` command in `tests/integration/overlap_test.go`

### Implementation for User Story 5

- [X] T040 [US5] Implement meeting local-time display logic in `internal/meeting/meeting.go`
- [X] T041 [US5] Implement overlap window finding in `internal/meeting/overlap.go`
- [X] T042 [US5] Implement `meeting` subcommand in `internal/cli/meeting.go`
- [X] T043 [US5] Implement `overlap` subcommand in `internal/cli/overlap.go`

**Checkpoint**: All user stories should now be independently functional

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [X] T044 [P] Add JSON output coverage and validate against `contracts/cli-contract.md` schema in `internal/output/output.go`
- [X] T045 [P] Add benchmark tests for lookup/conversion in `internal/lookup/bench_test.go` and `internal/convert/bench_test.go`
- [X] T046 Run `gofmt`, `go vet`, and cleanup across all packages
- [X] T047 Run quickstart.md end-to-end validation (all 8 scenarios) and confirm coverage threshold
- [X] T048 [P] Write `README.md` usage documentation (commands, examples, cross-platform build instructions)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed) or sequentially in priority order (P1 → P2 → P3)
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P1)**: Can start after Foundational (Phase 2) - shares resolver with US1 but independently testable
- **User Story 3 (P2)**: Can start after Foundational (Phase 2) - shares resolver with US1/US2 but independently testable
- **User Story 4 (P2)**: Depends on resolver (US1) for location resolution; adds `convert` + multi-location `now`
- **User Story 5 (P3)**: Depends on US4 (conversion) for meeting local-time computation

### Within Each User Story

- Tests MUST be written and FAIL before implementation (TDD)
- Resolution logic before subcommand wiring
- Core implementation before integration
- Story complete before moving to next priority

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational tasks marked [P] can run in parallel (within Phase 2)
- Once Foundational phase completes, US1/US2/US3 can start in parallel (if staffed)
- All tests for a user story marked [P] can run in parallel
- Different user stories can be worked on in parallel by different team members

---

## Parallel Example: User Story 1

```bash
# Launch all tests for User Story 1 together:
Task: "Unit test: timezone → current time with fixed clock in internal/lookup/lookup_test.go"
Task: "Unit test: invalid timezone returns error in internal/lookup/lookup_test.go"
Task: "Integration test: now Asia/Shanghai end-to-end in tests/integration/now_test.go"

# Then implement sequentially (same files, ordered):
Task: "Implement timezone resolution in internal/lookup/resolve.go"
Task: "Implement current-time lookup in internal/lookup/current.go"
Task: "Implement now subcommand in internal/cli/now.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1 (按时区查询)
4. **STOP and VALIDATE**: Test User Story 1 independently
5. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 (时区查询) → Test → Demo (MVP!)
3. Add User Story 2 (地点查询) → Test → Demo
4. Add User Story 3 (邮编查询) → Test → Demo
5. Add User Story 4 (换算/多地点) → Test → Demo
6. Add User Story 5 (会议安排) → Test → Demo
7. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 (时区查询)
   - Developer B: User Story 2 (地点查询)
   - Developer C: User Story 3 (邮编查询)
3. Then User Story 4 (换算) and User Story 5 (会议) build on the resolver

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing (宪法 TDD 要求)
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence

---

## Phase 9: Convergence

**Purpose**: Close the gap between the feature's spec/plan/tasks and the implemented code.

- [X] T049 Remove dead code: delete unused `data.Window` type in `internal/data/types.go` and unused `Store.PlaceNames()` in `internal/data/load.go` per Constitution 一 (contradicts)
- [X] T050 Add benchmark tests for timezone lookup and conversion in `internal/lookup/bench_test.go` and `internal/convert/bench_test.go` per T045 / Constitution 四 (missing)
- [X] T051 Define a coverage threshold and add direct unit tests for `internal/cli`, `internal/data`, and `internal/output` (currently 0% direct coverage); enforce the threshold in CI/Makefile per Constitution 二 (partial)
- [X] T052 Add a test asserting `--json` output matches the `contracts/cli-contract.md` schema in `tests/integration/json_test.go` per T044 (partial)
- [X] T053 Detect and clearly report DST "skipped time" (non-existent local time, e.g. spring-forward 02:30) in `parseTimeIn` in `internal/cli/flags.go` per spec edge case (partial)
