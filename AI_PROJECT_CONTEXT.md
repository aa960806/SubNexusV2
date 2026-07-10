# SubNexusV2 项目上下文

## 项目基线

- 本项目是一个独立的新目录，不包含 `.git`，不会修改 `sub2api` 或 `SubNexus` 原项目。
- 代码基线已同步到上游 `Wei-Shaw/sub2api` 提交 `5260a42a`（2026-07-10 17:31:55 +08:00）。
- 已完整保留 `5260a42a` 及之前的上游功能，包括 GPT-5.6 计费、Usage 完整性、用户级 OpenAI Fast/Flex 策略和 Grok reasoning effort 修复。

## 保留的二开功能

仅回植以下三项功能：

1. 模型广场：用户展示页、管理员配置页及对应 API。
2. 活动中心：用户展示页、管理员配置与活动条目 CRUD。
3. 跑马灯：登录用户全局滚动公告、管理员开关与公告 CRUD。

签到、排行榜、创意工坊及旧活动玩法均不接入路由、Handler、Service 或后台 Worker。

## 数据库兼容策略

- `SubNexus` 的 `151` 到 `164` 共 14 个迁移原样保留，文件名和 SHA256 均未改变，避免已部署数据库发生迁移校验和冲突。
- 旧表、旧列和旧迁移记录不删除，未保留功能只停止代码接线。
- 上游 `173_allow_cyber_blocked_usage_request_type.sql` 原样保留。
- `174_disable_removed_subnexus_activity_entries.sql` 仅将四个旧玩法入口设为停用，不删除任何业务数据。
- 活动中心服务端还会强制隐藏 `daily_spin`、`invite_lottery`、`recharge_wheel`、`invite_milestone` 类型，防止旧数据误暴露无效页面。
- 跑马灯继续复用 `ACTIVITY_CONFIG`，更新开关时以 JSON 合并方式保留所有未知字段。

## 功能配置

- 模型广场使用原二开设置键和 `/settings/model-plaza`、`/admin/settings/model-plaza` 接口。
- 活动中心使用 `ACTIVITY_CENTER_CONFIG`；默认关闭，管理员启用后用户菜单才显示。
- 跑马灯使用 `ACTIVITY_CONFIG.broadcast_enabled`；字段不存在时保持旧版默认启用行为。

## 首页与视觉移植

- 仅将 `SubNexus` 的默认首页视觉、动画和 `frontend/public/cxk.gif` 移植到当前上游基线。
- 未使用旧版覆盖 `AppHeader.vue`、`AppSidebar.vue` 或 `App.vue`；上游 Header 业务逻辑保持原样，Sidebar 和 App 只增加三项保留功能所需接线。
- 站点 Logo、文档地址和自定义首页 URL 均通过 `sanitizeUrl` 清洗；自定义 HTML 首页保持原有配置行为。
- 默认首页动画支持 `prefers-reduced-motion`，切换到自定义首页或组件卸载时会清理 RAF、定时器和鼠标监听。
- `cxk.gif` SHA256：`51E827CB161122A3D3A608418A334AAB157478B5901D1BE38F03DD4CE258B1D8`。

## 发布前安全步骤

1. 备份生产 PostgreSQL 数据库并验证备份可恢复。
2. 先在生产数据库副本上启动本项目，确认全部迁移成功且无 checksum mismatch。
3. 验证现有账号、渠道、分组、账单和 GPT-5.6 请求链路。
4. 验证三项二开功能的开关、用户权限和管理员 CRUD。
5. 观察日志与数据库指标后再灰度切换生产流量，并保留旧服务的快速回切能力。

## 已完成验证

- 上游 `9a2f11b4..6dd3274a` 的 41 个变更文件均与目标提交 Git blob 完全一致。
- 上游 `6dd3274a..5260a42a` 的 17 个变更文件均与目标提交 Git blob 完全一致；该范围没有数据库迁移变更。
- `AppHeader.vue` 与 `5260a42a` 对应 Git blob 完全一致；`App.vue` 和 `AppSidebar.vue` 仅包含三项保留功能的接线差异。
- `SubNexus` 的 14 个历史迁移已按完整文件名逐一核对，文件名和 SHA256 全部一致；迁移器以完整文件名而不是数字前缀作为主键，同号上游迁移可并存。
- 后端 Wire 依赖注入已重新生成。
- 后端 `go test -vet=off -p 1 ./...` 全量通过。
- 跑马灯配置 JSON 合并、旧活动类型隐藏和迁移非删除行为的回归测试通过。
- 首页安全与定制行为的 3 个定向测试文件共 14 个测试通过。
- 前端 Vitest 全量测试、`vue-tsc` 类型检查和 Vite 生产构建全部通过。
- `frontend/pnpm-lock.yaml` 保持上游 SHA256：`876891E718D780818038047228CA276378C9E68965519B82919B433B8A0E969C`。
- 服务器适配必须按 `SERVER_COMPATIBILITY_TEST_PLAN.md` 先在生产数据库副本与隔离 Redis 上执行，不允许首次验证直接连接生产数据库。
