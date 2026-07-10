# SubNexusV2 现服务器适配测试计划

## 安全边界

- 首次启动只允许连接生产 PostgreSQL 的可恢复副本，禁止直接连接生产数据库。
- 使用独立 Redis 或独立数据库编号/键前缀，禁止复用生产缓存和任务队列。
- 使用不同监听端口，不接入生产负载均衡、域名、支付回调、邮件或消息通知入口。
- 不在文档、日志或命令输出中暴露数据库密码、JWT、OAuth Secret、加密密钥或 API Key。
- 保留当前线上镜像、配置和数据库备份，整个验证阶段必须具备快速回切能力。

## 第一阶段：只读盘点

1. 记录现服务器的部署方式、镜像版本、PostgreSQL/Redis 版本、实例数量和反向代理拓扑。
2. 仅记录环境变量名称是否存在，不输出值；重点核对数据库、Redis、`JWT_SECRET`、TOTP/支付/OAuth 加密密钥和时区配置。
3. 导出当前 `schema_migrations` 的完整文件名与 checksum，用于和本项目嵌入迁移逐项比较。
4. 记录核心表行数和关键业务聚合基线：用户、API Key、账号、分组、渠道、订阅、账单与 usage logs。
5. 确认 `ACTIVITY_CONFIG`、`ACTIVITY_CENTER_CONFIG` 和模型广场相关设置键存在情况，但不输出可能包含敏感内容的完整值。

建议只读 SQL：

```sql
SELECT version();
SELECT filename, checksum, applied_at FROM schema_migrations ORDER BY filename;
SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' ORDER BY table_name;
SELECT key FROM settings WHERE key IN ('ACTIVITY_CONFIG', 'ACTIVITY_CENTER_CONFIG') ORDER BY key;
```

## 第二阶段：创建隔离副本

1. 使用 `pg_dump -Fc` 备份生产数据库，并用 `pg_restore --list` 验证备份可读取。
2. 将备份恢复到隔离 PostgreSQL 实例，禁止使用生产数据库名称和账号。
3. 为测试实例创建独立 Redis；不复制生产任务队列，避免重复执行异步任务。
4. 复制线上配置时保持认证/加密密钥一致以验证历史数据可解密，但关闭所有外部副作用，并限制配置文件权限。
5. 在防火墙或测试网络层阻止测试实例接收真实生产流量。

## 第三阶段：迁移与启动

1. 只启动一个 SubNexusV2 实例，观察迁移日志，任何 checksum mismatch 或 SQL 错误都立即停止。
2. 对比启动前后的 `schema_migrations`，确认只新增预期文件，已有迁移 checksum 不变化。
3. 对比核心表行数和聚合基线；`174_disable_removed_subnexus_activity_entries.sql` 只能停用四个旧玩法条目，不得删除业务数据。
4. 检查启动日志中数据库、Redis、密钥解密、后台 Worker 和依赖注入错误。
5. 确认健康检查通过后再开放仅管理员可访问的测试地址。

## 第四阶段：功能回归

1. 登录与会话：现有管理员和普通用户登录、Token 刷新、权限与历史 OAuth 绑定正常。
2. 核心数据：账号、渠道、代理、分组、API Key、订阅、订单、余额和 usage logs 可正确读取。
3. GPT-5.6：普通/流式请求、缓存计费、Usage 完整性、模型映射和错误回退符合线上基线。
4. 新上游功能：用户级 OpenAI Fast/Flex 策略按 API Key 所属用户生效；Grok Responses 的 `reasoning_effort` 正确透传与记录。
5. 二开功能：模型广场、活动中心和跑马灯的开关、用户权限、管理员 CRUD 与历史数据兼容。
6. 首页与 UI：默认首页、自定义 URL/HTML 首页、Logo、文档链接和移动端布局正常。
7. 已移除功能：签到、排行榜、创意工坊和四个旧玩法没有可达路由、菜单或后台 Worker。

## 第五阶段：灰度与回切

1. 先导入脱敏请求样本做压测，比较错误率、P95/P99 延迟、数据库连接和 Redis 命中情况。
2. 生产切换前再次生成并验证数据库备份，冻结结构变更窗口。
3. 仅灰度少量内部流量，持续观察迁移、计费、Usage、任务队列和数据库指标。
4. 达到观察窗口且无异常后逐步扩大流量；任何计费、数据一致性或鉴权异常立即回切旧实例。
5. 回切只切换应用流量，不回滚已成功执行的迁移；数据库修复必须通过新的前向迁移完成。

## 通过标准

- 所有迁移 checksum 一致且没有非预期数据删除或结构回退。
- 后端/前端自动化测试、生产构建和隔离环境冒烟测试全部通过。
- 现有用户、密钥、账号、账单与 Usage 数据前后一致。
- GPT-5.6、新上游功能和三项二开功能均通过权限、计费和错误路径测试。
- 已验证旧版本可快速回切，并完成至少一次回切演练。
