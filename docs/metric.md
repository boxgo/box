- [Metrics Documentation](#metrics-documentation)
  - [1. 核心指标定义 (Definitions)](#1-核心指标定义-definitions)
    - [1.1 HTTP Server (Gin)](#11-http-server-gin)
    - [1.2 HTTP Client (Wukong)](#12-http-client-wukong)
    - [1.3 gRPC Server](#13-grpc-server)
    - [1.4 Redis Client](#14-redis-client)
    - [1.5 Database Client (GORM)](#15-database-client-gorm)
    - [1.6 MongoDB Client](#16-mongodb-client)
    - [1.7 Schedule (定时任务)](#17-schedule-定时任务)
    - [1.8 Go Runtime (运行时)](#18-go-runtime-运行时)
      - [基础指标](#基础指标)
      - [内存分配统计](#内存分配统计)
      - [堆内存统计](#堆内存统计)
      - [栈内存统计](#栈内存统计)
      - [MSpan / MCache 统计](#mspan--mcache-统计)
      - [GC 统计](#gc-统计)
      - [其他系统内存](#其他系统内存)
  - [2. 推荐看板指标 (Grafana PromQL)](#2-推荐看板指标-grafana-promql)
    - [2.1 📊 概览 (Overview)](#21--概览-overview)
      - [Apdex Score](#apdex-score)
      - [关键指标卡片](#关键指标卡片)
    - [2.2 🌐 HTTP Server](#22--http-server)
    - [2.3 📤 HTTP Client](#23--http-client)
    - [2.4 🔌 gRPC Server](#24--grpc-server)
    - [2.5 🐹 Go Runtime](#25--go-runtime)
    - [2.6 🔴 Redis](#26--redis)
    - [2.7 🗄️ Database (DB)](#27-️-database-db)
    - [2.8 🍃 MongoDB](#28--mongodb)
  - [3. 告警规则 (Alerting Rules)](#3-告警规则-alerting-rules)
  - [4. Go Runtime 指标解读](#4-go-runtime-指标解读)
    - [4.1 内存指标关系](#41-内存指标关系)
    - [4.2 关键指标说明](#42-关键指标说明)
      - [Goroutine 监控](#goroutine-监控)
      - [内存监控](#内存监控)
      - [GC 监控](#gc-监控)
  - [5. 错误分类说明 (Error Classification)](#5-错误分类说明-error-classification)
    - [5.1 错误分类原则](#51-错误分类原则)
    - [5.2 HTTP Client 错误分类](#52-http-client-错误分类)
    - [5.3 Redis Client 错误分类](#53-redis-client-错误分类)
    - [5.4 Database Client 错误分类](#54-database-client-错误分类)
    - [5.5 错误分类性能影响](#55-错误分类性能影响)
  - [6. 常见问题诊断 (Troubleshooting)](#6-常见问题诊断-troubleshooting)
    - [6.1 Go Runtime 问题](#61-go-runtime-问题)
      - [问题 1: Goroutine 泄漏](#问题-1-goroutine-泄漏)
      - [问题 2: 内存泄漏](#问题-2-内存泄漏)
      - [问题 3: GC 压力过大](#问题-3-gc-压力过大)
      - [问题 4: 线程数异常增长](#问题-4-线程数异常增长)
    - [6.2 中间件与服务问题](#62-中间件与服务问题)
      - [问题 5: 数据库连接池耗尽](#问题-5-数据库连接池耗尽)
      - [问题 6: Redis 延迟抖动](#问题-6-redis-延迟抖动)
      - [问题 7: Context Cancelled / Timeout](#问题-7-context-cancelled--timeout)
      - [问题 8: 定时任务堆积](#问题-8-定时任务堆积)

# Metrics Documentation

本文档记录了 `box` 框架中各组件暴露的 Prometheus 监控指标、推荐的 Grafana 看板配置以及告警规则。

## 1. 核心指标定义 (Definitions)

### 1.1 HTTP Server (Gin)

| 指标名称                               | 类型      | Labels                               | 说明                                        |
| :------------------------------------- | :-------- | :----------------------------------- | :------------------------------------------ |
| `http_server_requests_inflight`        | Gauge     | `method`, `url`                      | 当前正在处理的 HTTP 请求数 (饱和度)         |
| `http_server_requests_total`           | Counter   | `method`, `url`, `status`, `errcode` | 处理的 HTTP 请求总数 (流量 & 错误)          |
| `http_server_request_duration_seconds` | Histogram | `method`, `url`, `status`, `errcode` | HTTP 请求耗时分布 (延迟)，桶：.005s - 10s   |
| `http_server_request_size_bytes`       | Histogram | `method`, `url`                      | HTTP 请求体大小分布 (流量)，桶：1KB - 100MB |
| `http_server_response_size_bytes`      | Histogram | `method`, `url`, `status`, `errcode` | HTTP 响应体大小分布 (流量)，桶：1KB - 100MB |

### 1.2 HTTP Client (Wukong)

| 指标名称                               | 类型      | Labels                                        | 说明                           |
| :------------------------------------- | :-------- | :-------------------------------------------- | :----------------------------- |
| `http_client_requests_inflight`        | Gauge     | `method`, `baseUrl`, `url`                    | 当前正在进行的下游 HTTP 请求数 |
| `http_client_requests_total`           | Counter   | `method`, `baseUrl`, `url`, `status`, `error` | 发起的 HTTP 请求总数           |
| `http_client_request_duration_seconds` | Histogram | `method`, `baseUrl`, `url`, `status`, `error` | HTTP 请求耗时分布              |

**错误分类 (`error` 标签值)**:

- `` - 成功（无错误）
- `timeout_error` - 超时错误（context 超时、I/O 超时等）
- `connection_error` - 连接错误（连接被拒绝、连接丢失等）
- `dns_error` - DNS 解析错误
- `tls_error` - TLS/SSL 错误（证书错误、握手失败等）
- `other_error` - 其他未分类错误

**注意**: HTTP 状态码通过 `status` 标签单独上报，`error` 标签仅用于底层网络/协议错误。

### 1.3 gRPC Server

| 指标名称                               | 类型      | Labels                   | 说明                       |
| :------------------------------------- | :-------- | :----------------------- | :------------------------- |
| `grpc_server_requests_inflight`        | Gauge     | `method`, `type`         | 当前正在处理的 gRPC 请求数 |
| `grpc_server_requests_total`           | Counter   | `method`, `type`, `code` | 处理的 gRPC 请求总数       |
| `grpc_server_request_duration_seconds` | Histogram | `method`, `type`, `code` | gRPC 请求耗时分布          |
| `grpc_server_panics_total`             | Counter   | `method`                 | gRPC 服务 Panic 总次数     |

### 1.4 Redis Client

| 指标名称                                | 类型      | Labels                                                 | 说明                   |
| :-------------------------------------- | :-------- | :----------------------------------------------------- | :--------------------- |
| `redis_client_requests_total`           | Counter   | `address`, `db`, `masterName`, `pipe`, `cmd`, `result` | Redis 命令执行总数     |
| `redis_client_request_duration_seconds` | Histogram | `address`, `db`, `masterName`, `pipe`, `cmd`, `result` | Redis 命令执行耗时分布 |

**错误分类 (`result` 标签值)**:

- `success` - 成功（包括 `redis.Nil`，键不存在是正常情况）
- `timeout_error` - 超时错误（context 超时、I/O 超时等）
- `connection_error` - 连接错误（连接被拒绝、连接丢失等）
- `command_error` - Redis 命令错误（WRONGTYPE、未知命令、参数错误等）
- `transaction_error` - 事务错误（事务失败、WATCH 失败等）
- `auth_error` - 权限/认证错误（NOAUTH、认证失败等）
- `oom_error` - 内存不足错误（OOM、内存限制等）
- `cluster_error` - 集群相关错误（MOVED、ASK、CLUSTERDOWN 等）
- `other_error` - 其他未分类错误

### 1.5 Database Client (GORM)

| 指标名称                             | 类型      | Labels                                 | 说明                       |
| :----------------------------------- | :-------- | :------------------------------------- | :------------------------- |
| `db_client_connections_idle`         | Gauge     | `driver`, `database`                   | 连接池空闲连接数           |
| `db_client_connections_in_use`       | Gauge     | `driver`, `database`                   | 连接池正在使用的连接数     |
| `db_client_connections_open`         | Gauge     | `driver`, `database`                   | 连接池当前打开的总连接数   |
| `db_client_connections_max_open`     | Gauge     | `driver`, `database`                   | 连接池最大允许打开的连接数 |
| `db_client_connections_wait_total`   | Gauge     | `driver`, `database`                   | 等待连接的总次数           |
| `db_client_connections_wait_seconds` | Gauge     | `driver`, `database`                   | 等待连接的总耗时           |
| `db_client_requests_total`           | Counter   | `driver`, `database`, `type`, `result` | 数据库请求执行总数         |
| `db_client_request_duration_seconds` | Histogram | `driver`, `database`, `type`, `result` | SQL 执行耗时分布           |

**错误分类 (`result` 标签值)**:

- `success` - 成功（包括 `gorm.ErrRecordNotFound`，记录不存在是正常情况）
- `timeout_error` - 超时错误（context 超时、查询超时等）
- `connection_error` - 连接错误（连接被拒绝、连接丢失、连接池耗尽等）
- `constraint_error` - 约束错误（唯一键冲突、外键约束、非空约束等）
- `syntax_error` - SQL 语法错误（语法错误、未知列/表等）
- `transaction_error` - 事务相关错误（死锁、锁等待超时等）
- `other_error` - 其他未分类错误

### 1.6 MongoDB Client

| 指标名称                                | 类型      | Labels              | 说明                          |
| :-------------------------------------- | :-------- | :------------------ | :---------------------------- |
| `mongo_client_requests_total`           | Counter   | `command`, `result` | MongoDB 命令执行总数          |
| `mongo_client_request_duration_seconds` | Histogram | `command`, `result` | MongoDB 命令耗时分布          |
| `mongo_client_sessions_inflight`        | Gauge     | -                   | 当前正在进行的 MongoDB 会话数 |

### 1.7 Schedule (定时任务)

| 指标名称                        | 类型      | Labels           | 说明                 |
| :------------------------------ | :-------- | :--------------- | :------------------- |
| `schedule_jobs_total`           | Counter   | `task`, `result` | 定时任务执行总数     |
| `schedule_job_duration_seconds` | Histogram | `task`, `result` | 定时任务执行耗时分布 |

### 1.8 Go Runtime (运行时)

#### 基础指标

| 指标名称        | 类型  | Labels    | 说明                |
| :-------------- | :---- | :-------- | :------------------ |
| `go_info`       | Gauge | `version` | Go 版本信息         |
| `go_goroutines` | Gauge | -         | 当前 Goroutine 数量 |
| `go_threads`    | Gauge | -         | 当前 OS 线程数量    |

#### 内存分配统计

| 指标名称                        | 类型    | Labels | 说明                                     |
| :------------------------------ | :------ | :----- | :--------------------------------------- |
| `go_memstats_alloc_bytes`       | Gauge   | -      | 已分配且仍在使用的堆内存字节数           |
| `go_memstats_alloc_bytes_total` | Counter | -      | 累计分配的堆内存总字节数（包括已释放的） |
| `go_memstats_sys_bytes`         | Gauge   | -      | 从操作系统获取的内存总字节数             |
| `go_memstats_lookups_total`     | Counter | -      | 指针查找总次数（通常为 0）               |
| `go_memstats_mallocs_total`     | Counter | -      | 累计内存分配次数                         |
| `go_memstats_frees_total`       | Counter | -      | 累计内存释放次数                         |

#### 堆内存统计

| 指标名称                          | 类型  | Labels | 说明                                 |
| :-------------------------------- | :---- | :----- | :----------------------------------- |
| `go_memstats_heap_alloc_bytes`    | Gauge | -      | 堆内存已分配字节数（已分配且在使用） |
| `go_memstats_heap_sys_bytes`      | Gauge | -      | 从系统获取的堆内存字节数             |
| `go_memstats_heap_idle_bytes`     | Gauge | -      | 堆内存空闲字节数（等待被使用）       |
| `go_memstats_heap_inuse_bytes`    | Gauge | -      | 堆内存正在使用的字节数               |
| `go_memstats_heap_released_bytes` | Gauge | -      | 已释放回操作系统的堆内存字节数       |
| `go_memstats_heap_objects`        | Gauge | -      | 堆中已分配的对象数量                 |
| `go_memstats_next_gc_bytes`       | Gauge | -      | 下次 GC 触发时的堆内存目标字节数     |

#### 栈内存统计

| 指标名称                        | 类型  | Labels | 说明                     |
| :------------------------------ | :---- | :----- | :----------------------- |
| `go_memstats_stack_inuse_bytes` | Gauge | -      | 栈分配器正在使用的字节数 |
| `go_memstats_stack_sys_bytes`   | Gauge | -      | 从系统获取的栈内存字节数 |

#### MSpan / MCache 统计

| 指标名称                         | 类型  | Labels | 说明                           |
| :------------------------------- | :---- | :----- | :----------------------------- |
| `go_memstats_mspan_inuse_bytes`  | Gauge | -      | MSpan 结构体正在使用的字节数   |
| `go_memstats_mspan_sys_bytes`    | Gauge | -      | 从系统获取的 MSpan 内存字节数  |
| `go_memstats_mcache_inuse_bytes` | Gauge | -      | MCache 结构体正在使用的字节数  |
| `go_memstats_mcache_sys_bytes`   | Gauge | -      | 从系统获取的 MCache 内存字节数 |

#### GC 统计

| 指标名称                           | 类型    | Labels     | 说明                                               |
| :--------------------------------- | :------ | :--------- | :------------------------------------------------- |
| `go_gc_duration_seconds`           | Summary | `quantile` | GC 暂停耗时分布（quantile: 0, 0.25, 0.5, 0.75, 1） |
| `go_memstats_gc_sys_bytes`         | Gauge   | -          | GC 元数据使用的内存字节数                          |
| `go_memstats_gc_cpu_fraction`      | Gauge   | -          | 程序启动以来 GC 使用的 CPU 时间占比                |
| `go_memstats_last_gc_time_seconds` | Gauge   | -          | 上次 GC 的 Unix 时间戳（秒）                       |

#### 其他系统内存

| 指标名称                          | 类型  | Labels | 说明                           |
| :-------------------------------- | :---- | :----- | :----------------------------- |
| `go_memstats_buck_hash_sys_bytes` | Gauge | -      | 性能分析哈希表使用的内存字节数 |
| `go_memstats_other_sys_bytes`     | Gauge | -      | 其他系统分配的内存字节数       |

---

## 2. 推荐看板指标 (Grafana PromQL)

以下 PromQL 假设你有一个 Dashboard 变量 `$namespace`、`$service` 和 `$instance`。

看板结构分为以下板块：

- **概览** - Apdex、关键指标概览
- **HTTP Server** - HTTP 服务器详细指标
- **HTTP Client** - HTTP 客户端详细指标
- **gRPC Server** - gRPC 服务器详细指标
- **Go Runtime** - Go 运行时详细指标
- **Database (DB)** - 数据库详细指标
- **Redis** - Redis 详细指标
- **MongoDB** - MongoDB 详细指标

### 2.1 📊 概览 (Overview)

#### Apdex Score

| 面板名称        | 说明                 | PromQL                                                                                                                                                                                                                                                                                                                                                                                                              |
| :-------------- | :------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| **Apdex Score** | 用户满意度 (T=250ms) | `(sum(rate(http_server_request_duration_seconds_bucket{namespace=~"$namespace",job=~"$service",instance=~"$instance", le="0.25"}[5m])) * 0.5 + sum(rate(http_server_request_duration_seconds_bucket{namespace=~"$namespace",job=~"$service",instance=~"$instance", le="1"}[5m])) * 0.5) / sum(rate(http_server_request_duration_seconds_count{namespace=~"$namespace",job=~"$service",instance=~"$instance"}[5m]))` |

**Apdex 计算说明**：

- **Satisfied（满意）**: 响应时间 ≤ T (250ms)
- **Tolerating（可容忍）**: T < 响应时间 ≤ 4T (250ms < t ≤ 1s)
- **Frustrated（失望）**: 响应时间 > 4T (> 1s)
- **公式**: `Apdex = (Satisfied + Tolerating/2) / Total`
- **取值范围**: 0 到 1，越接近 1 表示用户体验越好
- **无请求时**: 当 Total = 0 时（无流量），Apdex 分数显示为 **N/A**（不可用），因为没有用户访问就无法评估用户体验

**评价标准与阈值区域**（Grafana 看板会自动显示评级与颜色）：

| Apdex 分数 | 评级                           | 颜色   | 阈值 | 用户体验       | 建议措施                       |
| :--------- | :----------------------------- | :----- | :--- | :------------- | :----------------------------- |
| 0.94-1.00  | **Excellent** (优秀) 🟢        | 绿色   | 0.94 | 极佳，用户满意 | 保持现状，持续监控             |
| 0.85-0.94  | **Good** (良好) 🟡             | 黄色   | 0.85 | 良好，可接受   | 关注趋势，优化慢请求           |
| 0.70-0.85  | **Fair** (一般) 🟠             | 橙色   | 0.70 | 一般，需改进   | 排查性能瓶颈，优化关键路径     |
| 0.50-0.70  | **Poor** (较差) 🔴             | 红色   | 0.50 | 较差，影响体验 | 立即介入，分析慢查询和依赖服务 |
| 0.00-0.50  | **Unacceptable** (不可接受) ⚫ | 深红色 | 0.00 | 不可接受，严重 | 紧急处理，可能需要扩容或限流   |

**Grafana 阈值配置**：

```json
{
  "thresholds": {
    "mode": "absolute",
    "steps": [
      { "color": "dark-red", "value": null },
      { "color": "red", "value": 0.5 },
      { "color": "orange", "value": 0.7 },
      { "color": "yellow", "value": 0.85 },
      { "color": "green", "value": 0.94 }
    ]
  }
}
```

**PromQL 实现**：

基于 Prometheus Histogram 的累积桶特性，Apdex 公式实现为：

```promql
(
  sum(rate(http_server_request_duration_seconds_bucket{le="0.25"}[5m])) * 0.5 +
  sum(rate(http_server_request_duration_seconds_bucket{le="1"}[5m])) * 0.5
) / sum(rate(http_server_request_duration_seconds_count[5m]))
```

**注意事项**：

- `le="0.25"` 桶包含 ≤250ms 的所有请求（Satisfied）
- `le="1"` 桶包含 ≤1s 的所有请求（Satisfied + Tolerating）
- 由于桶的累积特性，需要用 `le="1"` 的值减去 `le="0.25"` 来计算 Tolerating 部分
- 公式简化为：`(Satisfied * 0.5 + (Satisfied + Tolerating) * 0.5) / Total`
- 结果等价于标准 Apdex 公式：`(Satisfied + Tolerating/2) / Total`

#### 关键指标卡片

| 面板名称              | 说明                   | PromQL                                                                                                                                                                                                                                     |
| :-------------------- | :--------------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **HTTP QPS**          | 当前请求速率           | `sum(rate(http_server_requests_total{namespace=~"$namespace",job=~"$service",instance=~"$instance"}[1m]))`                                                                                                                                 |
| **HTTP P99 Latency**  | P99 延迟               | `histogram_quantile(0.99, sum(rate(http_server_request_duration_seconds_bucket{namespace=~"$namespace",job=~"$service",instance=~"$instance"}[1m])) by (le))`                                                                              |
| **HTTP Success Rate** | 成功率                 | `sum(rate(http_server_requests_total{namespace=~"$namespace",job=~"$service",instance=~"$instance", status!~"5.."}[1m])) / sum(rate(http_server_requests_total{namespace=~"$namespace",job=~"$service",instance=~"$instance"}[1m])) * 100` |
| **Today's Peak QPS**  | 今日 QPS 峰值          | `max_over_time(sum(rate(http_server_requests_total{namespace=~"$namespace",job=~"$service",instance=~"$instance"}[1m]))[1d:])`                                                                                                             |
| **Goroutines**        | 协程总数（多实例聚合） | `sum(go_goroutines{namespace=~"$namespace",job=~"$service",instance=~"$instance"})`                                                                                                                                                        |
| **Memory InUse**      | 内存总量（多实例聚合） | `sum(go_memstats_heap_inuse_bytes{namespace=~"$namespace",job=~"$service",instance=~"$instance"})`                                                                                                                                         |
| **HTTP Inflight**     | 并发请求数             | `sum(http_server_requests_inflight{namespace=~"$namespace",job=~"$service",instance=~"$instance"})`                                                                                                                                        |

**注意事项**：

- 概览区域的 **Goroutines** 和 **Memory InUse** 面板使用 `sum()` 聚合显示所有实例的总和，适合快速了解整体资源使用情况
- 如需查看单个实例的详细情况，请访问 **Go Runtime** 板块，其中的时序图按 `instance` 分组显示每个实例的详细趋势

### 2.2 🌐 HTTP Server

| 面板名称                               | 说明         | PromQL                                                                                                                                                                                                                                                                                                                                                             |
| :------------------------------------- | :----------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **HTTP QPS by Endpoint**               | 按端点的 QPS | `sum(rate(http_server_requests_total{namespace=~"$namespace",job=~"$service",instance=~"$instance"}[1m])) by (method, url)`                                                                                                                                                                                                                                        |
| **HTTP Latency (P99/P95) by Endpoint** | 按端点的延迟 | P99: `histogram_quantile(0.99, sum(rate(http_server_request_duration_seconds_bucket{namespace=~"$namespace",job=~"$service",instance=~"$instance"}[1m])) by (le, method, url))`<br>P95: `histogram_quantile(0.95, sum(rate(http_server_request_duration_seconds_bucket{namespace=~"$namespace",job=~"$service",instance=~"$instance"}[1m])) by (le, method, url))` |
| **HTTP Status Codes**                  | 状态码分布   | `sum(rate(http_server_requests_total{namespace=~"$namespace",job=~"$service",instance=~"$instance"}[1m])) by (status)`                                                                                                                                                                                                                                             |
| **HTTP Network Traffic**               | 网络流量     | Request: `sum(rate(http_server_request_size_bytes_sum{namespace=~"$namespace",job=~"$service",instance=~"$instance"}[1m]))`<br>Response: `sum(rate(http_server_response_size_bytes_sum{namespace=~"$namespace",job=~"$service",instance=~"$instance"}[1m]))`                                                                                                       |

### 2.3 📤 HTTP Client

| 面板名称                       | 说明           | PromQL                                                                                                                                                                      |
| :----------------------------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **HTTP Client QPS**            | 客户端请求 QPS | `sum(rate(http_client_requests_total{namespace=~"$namespace",job=~"$service",instance=~"$instance"}[1m])) by (baseUrl, url)`                                                |
| **HTTP Client Latency (P99)**  | 客户端延迟     | `histogram_quantile(0.99, sum(rate(http_client_request_duration_seconds_bucket{namespace=~"$namespace",job=~"$service",instance=~"$instance"}[1m])) by (le, baseUrl, url))` |
| **HTTP Client Errors**         | 客户端错误     | `sum(rate(http_client_requests_total{namespace=~"$namespace",job=~"$service",instance=~"$instance",error!=""}[1m])) by (baseUrl, url, error)`                               |
| **HTTP Client Errors by Type** | 按错误类型分类 | `sum(rate(http_client_requests_total{namespace=~"$namespace",job=~"$service",instance=~"$instance",error!=""}[1m])) by (error)`                                             |

### 2.4 🔌 gRPC Server

| 面板名称                      | 说明             | PromQL                                                                                                                                                                |
| :---------------------------- | :--------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **gRPC Server QPS**           | gRPC 调用量      | `sum(rate(grpc_server_requests_total{namespace=~"$namespace",job=~"$service",instance=~"$instance"}[1m])) by (method, type)`                                          |
| **gRPC Server Latency (P99)** | gRPC 接口延迟    | `histogram_quantile(0.99, sum(rate(grpc_server_request_duration_seconds_bucket{namespace=~"$namespace",job=~"$service",instance=~"$instance"}[1m])) by (le, method))` |
| **gRPC Server Errors**        | gRPC 错误数      | `sum(rate(grpc_server_requests_total{namespace=~"$namespace",job=~"$service",instance=~"$instance", code!="OK"}[1m])) by (method, code)`                              |
| **gRPC Server Panics**        | Panic 发生的次数 | `increase(grpc_server_panics_total{namespace=~"$namespace",job=~"$service",instance=~"$instance"}[1m])`                                                               |

### 2.5 🐹 Go Runtime

| 面板名称                   | 说明            | PromQL                                                                                                                                                                                                                                                                                                                          |
| :------------------------- | :-------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| **Goroutines**             | Goroutine 数量  | `go_goroutines{namespace=~"$namespace",job=~"$service",instance=~"$instance"}`                                                                                                                                                                                                                                                  |
| **Heap Memory**            | 堆内存使用情况  | Alloc: `go_memstats_heap_alloc_bytes{namespace=~"$namespace",job=~"$service",instance=~"$instance"}`<br>InUse: `go_memstats_heap_inuse_bytes{namespace=~"$namespace",job=~"$service",instance=~"$instance"}`<br>Sys: `go_memstats_heap_sys_bytes{namespace=~"$namespace",job=~"$service",instance=~"$instance"}`                |
| **GC Duration**            | GC 耗时         | Avg: `rate(go_gc_duration_seconds_sum{namespace=~"$namespace",job=~"$service",instance=~"$instance"}[1m]) / rate(go_gc_duration_seconds_count{namespace=~"$namespace",job=~"$service",instance=~"$instance"}[1m])`<br>Max: `go_gc_duration_seconds{namespace=~"$namespace",job=~"$service",instance=~"$instance",quantile="1"}` |
| **GC Rate**                | GC 执行频率     | `rate(go_gc_duration_seconds_count{namespace=~"$namespace",job=~"$service",instance=~"$instance"}[1m])`                                                                                                                                                                                                                         |
| **GC CPU Fraction**        | GC CPU 占用比例 | `go_memstats_gc_cpu_fraction{namespace=~"$namespace",job=~"$service",instance=~"$instance"}`                                                                                                                                                                                                                                    |
| **Memory Allocation Rate** | 内存分配速率    | `rate(go_memstats_alloc_bytes_total{namespace=~"$namespace",job=~"$service",instance=~"$instance"}[1m])`                                                                                                                                                                                                                        |

### 2.6 🔴 Redis

| 面板名称                        | 说明     | PromQL                                                                                                                                                              |
| :------------------------------ | :------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| **Redis Command QPS**           | 命令 QPS | `sum(rate(redis_client_requests_total{namespace=~"$namespace",job=~"$service",instance=~"$instance"}[1m])) by (cmd)`                                                |
| **Redis Command Latency (P99)** | 命令延迟 | `histogram_quantile(0.99, sum(rate(redis_client_request_duration_seconds_bucket{namespace=~"$namespace",job=~"$service",instance=~"$instance"}[1m])) by (le, cmd))` |
| **Redis Command Errors**        | 命令错误 | `sum(rate(redis_client_requests_total{namespace=~"$namespace",job=~"$service",instance=~"$instance",result!="success"}[1m])) by (cmd)`                              |

### 2.7 🗄️ Database (DB)

| 面板名称                   | 说明           | PromQL                                                                                                                                                                                                                                                                                                         |
| :------------------------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **DB Connection Pool**     | 连接池状态     | Open: `db_client_connections_open{namespace=~"$namespace",job=~"$service",instance=~"$instance"}`<br>InUse: `db_client_connections_in_use{namespace=~"$namespace",job=~"$service",instance=~"$instance"}`<br>Idle: `db_client_connections_idle{namespace=~"$namespace",job=~"$service",instance=~"$instance"}` |
| **DB Query Latency (P99)** | 查询延迟       | `histogram_quantile(0.99, sum(rate(db_client_request_duration_seconds_bucket{namespace=~"$namespace",job=~"$service",instance=~"$instance"}[1m])) by (le, type, database))`                                                                                                                                    |
| **DB Query QPS**           | 查询 QPS       | `sum(rate(db_client_requests_total{namespace=~"$namespace",job=~"$service",instance=~"$instance"}[1m])) by (type, database)`                                                                                                                                                                   |
| **DB Query Errors**        | 查询错误       | `sum(rate(db_client_requests_total{namespace=~"$namespace",job=~"$service",instance=~"$instance",result!="success"}[1m])) by (type, database, result)`                                                                                                                                         |
| **DB Errors by Type**      | 按错误类型分类 | `sum(rate(db_client_requests_total{namespace=~"$namespace",job=~"$service",instance=~"$instance",result!="success"}[1m])) by (result)`                                                                                                                                                         |

### 2.8 🍃 MongoDB

| 面板名称                          | 说明     | PromQL                                                                                                                                                                  |
| :-------------------------------- | :------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **MongoDB Command QPS**           | 命令 QPS | `sum(rate(mongo_client_requests_total{namespace=~"$namespace",job=~"$service",instance=~"$instance"}[1m])) by (command)`                                                |
| **MongoDB Command Latency (P99)** | 命令延迟 | `histogram_quantile(0.99, sum(rate(mongo_client_request_duration_seconds_bucket{namespace=~"$namespace",job=~"$service",instance=~"$instance"}[1m])) by (le, command))` |
| **MongoDB Command Errors**        | 命令错误 | `sum(rate(mongo_client_requests_total{namespace=~"$namespace",job=~"$service",instance=~"$instance",result="error"}[1m])) by (command)`                                 |

---

## 3. 告警规则 (Alerting Rules)

以下是基于 Prometheus 的推荐告警规则配置，涵盖了可用性、延迟、错误率、资源饱和度及运行时异常。

[prometheus_alerts_template](./prometheus_alerts_template.yaml)

---

## 4. Go Runtime 指标解读

### 4.1 内存指标关系

```
go_memstats_sys_bytes (从系统获取的总内存)
├── go_memstats_heap_sys_bytes (堆内存)
│   ├── go_memstats_heap_inuse_bytes (使用中的堆内存)
│   │   └── go_memstats_heap_alloc_bytes (已分配的堆内存)
│   └── go_memstats_heap_idle_bytes (空闲堆内存)
│       └── go_memstats_heap_released_bytes (已释放给 OS 的内存)
├── go_memstats_stack_sys_bytes (栈内存)
├── go_memstats_mspan_sys_bytes (MSpan 元数据)
├── go_memstats_mcache_sys_bytes (MCache 元数据)
├── go_memstats_buck_hash_sys_bytes (性能分析哈希表)
├── go_memstats_gc_sys_bytes (GC 元数据)
└── go_memstats_other_sys_bytes (其他系统内存)
```

### 4.2 关键指标说明

#### Goroutine 监控

- **正常范围**: 取决于业务负载，通常在 100-1000 之间
- **泄漏迹象**: 持续增长且不下降，或增长速率过快 (>100/s)
- **优化建议**: 确保所有 goroutine 都有退出机制，避免永久阻塞

#### 内存监控

- **heap_alloc**: 实际使用的堆内存，频繁上下波动是正常的（GC 会回收）
- **heap_inuse**: 包含已分配和待回收的内存，通常比 heap_alloc 大
- **heap_sys**: 从系统申请的堆内存，增长后不会轻易释放
- **泄漏迹象**: `heap_alloc` 持续增长、`heap_sys` 不断扩大且 GC 无法回收

#### GC 监控

- **正常 GC 耗时**: P99 应在 10-100ms 之内（依赖堆大小）
- **正常 GC 频率**: 每分钟几次到几十次（依赖分配速率）
- **GC CPU 占比**: 通常在 5%-25% 之间
- **异常情况**:
  - GC 耗时过长 (>1s): 可能堆太大或存在大对象
  - GC 频率过高 (>5 次/s): 分配速率过快，考虑对象池复用
  - GC CPU 占比过高 (>30%): 严重影响业务性能

## 5. 错误分类说明 (Error Classification)

为了在保留有用错误信息的同时避免指标爆炸（cardinality explosion），框架对错误进行了分类汇总。

### 5.1 错误分类原则

1. **避免指标爆炸**: 将错误归类为有限的几个类别（通常 5-10 个），而不是每个错误一个指标
2. **保留有用信息**: 通过类别区分常见错误类型，便于监控和告警
3. **性能优化**: 错误分类开销极小（< 200ns），相对于网络/IO 操作可忽略不计

### 5.2 HTTP Client 错误分类

HTTP 客户端错误通过 `error` 标签分类：

| 错误类型           | 说明           | 常见场景                              |
| ------------------ | -------------- | ------------------------------------- |
| `success`          | 成功（无错误） | 请求成功完成                          |
| `timeout_error`    | 超时错误       | context 超时、I/O 超时、网络超时      |
| `connection_error` | 连接错误       | 连接被拒绝、连接丢失、EOF、网络不可达 |
| `dns_error`        | DNS 解析错误   | 主机未找到、DNS 查询失败              |
| `tls_error`        | TLS/SSL 错误   | 证书错误、握手失败、X509 验证失败     |
| `other_error`      | 其他错误       | 未分类的错误                          |

**注意**: HTTP 状态码（如 404、500）通过 `status` 标签单独上报，`error` 标签仅用于底层网络/协议错误。

**示例 PromQL**:

```promql
# 查看超时错误
sum(rate(http_client_requests_total{error="timeout_error"}[5m])) by (baseUrl, url)

# 查看连接错误
sum(rate(http_client_requests_total{error="connection_error"}[5m])) by (baseUrl, url)

# 查看 DNS 错误
sum(rate(http_client_requests_total{error="dns_error"}[5m])) by (baseUrl, url)
```

### 5.3 Redis Client 错误分类

Redis 客户端错误通过 `result` 标签分类：

| 错误类型            | 说明     | 常见场景                                |
| ------------------- | -------- | --------------------------------------- |
| `success`           | 成功     | 命令执行成功（包括 `redis.Nil`）        |
| `timeout_error`     | 超时错误 | context 超时、I/O 超时                  |
| `connection_error`  | 连接错误 | 连接被拒绝、连接丢失、连接关闭          |
| `command_error`     | 命令错误 | WRONGTYPE、未知命令、参数错误、NOSCRIPT |
| `transaction_error` | 事务错误 | 事务失败、WATCH 失败、EXECABORT         |
| `auth_error`        | 权限错误 | NOAUTH、认证失败、ACL 权限错误          |
| `oom_error`         | 内存不足 | OOM、内存限制                           |
| `cluster_error`     | 集群错误 | MOVED、ASK、CLUSTERDOWN、跨槽错误       |
| `other_error`       | 其他错误 | 未分类的错误                            |

**示例 PromQL**:

```promql
# 查看连接错误
sum(rate(redis_client_requests_total{result="connection_error"}[5m])) by (cmd)

# 查看命令错误（可能是代码问题）
sum(rate(redis_client_requests_total{result="command_error"}[5m])) by (cmd)

# 查看内存不足错误（紧急）
sum(rate(redis_client_requests_total{result="oom_error"}[5m]))
```

### 5.4 Database Client 错误分类

数据库客户端错误通过 `result` 标签分类：

| 错误类型            | 说明         | 常见场景                                  |
| ------------------- | ------------ | ----------------------------------------- |
| `success`           | 成功         | 查询成功（包括 `gorm.ErrRecordNotFound`） |
| `timeout_error`     | 超时错误     | context 超时、查询超时、I/O 超时          |
| `connection_error`  | 连接错误     | 连接被拒绝、连接丢失、连接池耗尽          |
| `constraint_error`  | 约束错误     | 唯一键冲突、外键约束、非空约束            |
| `syntax_error`      | SQL 语法错误 | 语法错误、未知列/表、表不存在             |
| `transaction_error` | 事务错误     | 死锁、锁等待超时、事务回滚                |
| `other_error`       | 其他错误     | 未分类的错误                              |

**示例 PromQL**:

```promql
# 查看连接错误
sum(rate(db_client_requests_total{result="connection_error"}[5m])) by (database)

# 查看超时错误
sum(rate(db_client_requests_total{result="timeout_error"}[5m])) by (database)

# 查看约束错误（可能是业务逻辑问题）
sum(rate(db_client_requests_total{result="constraint_error"}[5m])) by (database)

# 查看死锁错误（紧急）
sum(rate(db_client_requests_total{result="transaction_error"}[5m])) by (database)
```

### 5.5 错误分类性能影响

错误分类的性能开销极小：

- **绝对开销**: 50-200 纳秒（ns）
- **相对开销**: < 0.2%（相对于网络/IO 操作）
- **内存开销**: 每次 1-2 个字符串分配（< 100 bytes）

详细性能分析请参考：[错误分类性能分析文档](./error_classification_performance_analysis.md)

---

## 6. 常见问题诊断 (Troubleshooting)

### 6.1 Go Runtime 问题

#### 问题 1: Goroutine 泄漏

**症状**: `go_goroutines` 持续增长
**排查**:

```promql
# 查看 Goroutine 增长速率
rate(go_goroutines[5m])

# 对比不同实例
go_goroutines by (instance)
```

**解决方案**:

- 使用 `pprof` 工具分析 goroutine 堆栈
- 检查 channel 是否正确关闭
- 确保 context 取消信号正确传递

#### 问题 2: 内存泄漏

**症状**: `go_memstats_heap_alloc_bytes` 持续增长，GC 无法回收
**排查**:

```promql
# 查看内存分配速率
rate(go_memstats_alloc_bytes_total[1m])

# 查看存活对象数
go_memstats_mallocs_total - go_memstats_frees_total
```

**解决方案**:

- 使用 `pprof` 工具分析内存分配
- 检查是否有全局变量持续引用对象
- 排查 map、slice 等容器是否及时清理

#### 问题 3: GC 压力过大

**症状**: `go_gc_duration_seconds` 过高或 `go_memstats_gc_cpu_fraction` 过高
**排查**:

```promql
# GC 执行频率
rate(go_gc_duration_seconds_count[1m])

# GC 平均耗时
rate(go_gc_duration_seconds_sum[1m]) / rate(go_gc_duration_seconds_count[1m])
```

**解决方案**:

- 优化内存分配，使用对象池（sync.Pool）
- 减少小对象分配，批量处理
- 调大 `GOGC` 环境变量（默认 100）
- 考虑使用 Go 1.19+ 的 Soft Memory Limit 特性

#### 问题 4: 线程数异常增长

**症状**: `go_threads` 持续增长，甚至导致程序 Crash (达到系统限制)。
**排查**:

```promql
# 查看线程数趋势
go_threads
```

**原因**:

- Go runtime 在进行系统调用（System Call）或 CGO 调用时，如果被阻塞，会创建新的 OS 线程来调度其他 Goroutine。
- 典型的阻塞场景：DNS 查询慢、文件 IO 阻塞、锁竞争。

**解决方案**:

- 优化阻塞的系统调用，使用非阻塞 IO
- 限制并发度
- 检查 CGO 代码逻辑

### 6.2 中间件与服务问题

#### 问题 5: 数据库连接池耗尽

**症状**: 数据库操作延迟增加，出现 `driver: bad connection` 或连接等待超时错误。
**排查**:

```promql
# 查看连接池使用率
sum(db_client_connections_in_use) by (database) / sum(db_client_connections_max_open) by (database)

# 查看连接等待次数
rate(db_client_connections_wait_total[1m])
```

**解决方案**:

- 调大 `SetMaxOpenConns`（需考虑 DB 服务端承载能力）
- 检查是否存在慢 SQL 长期占用连接
- 确保事务在所有路径（包括错误处理）中都能正确 `Commit` 或 `Rollback`

#### 问题 6: Redis 延迟抖动

**症状**: `redis_client_request_duration_seconds` P99 偶尔飙升，影响接口响应。
**排查**:

```promql
# 按命令查看延迟
histogram_quantile(0.99, sum(rate(redis_client_request_duration_seconds_bucket[5m])) by (le, cmd))
```

**解决方案**:

- 检查是否使用了 `KEYS`、`HGETALL` 等 O(N) 复杂度的命令
- 检查是否存在 Big Key（Value 过大），导致网络传输和序列化耗时增加
- 检查 Redis 服务端是否有慢查询日志

#### 问题 7: Context Cancelled / Timeout

**症状**: 客户端收到大量 `context canceled` 或 `deadline exceeded` 错误。
**排查**:

```promql
# 查看 gRPC/HTTP 错误码分布
sum(rate(grpc_server_requests_total{code="Canceled"}[1m]))
sum(rate(grpc_server_requests_total{code="DeadlineExceeded"}[1m]))
```

**原因**:

- 上游服务设置的超时时间过短
- 当前服务处理过慢（检查延迟指标）
- 客户端在请求完成前主动断开了连接

**解决方案**:

- 检查链路超时配置，确保下游超时时间 < 上游超时时间
- 优化接口性能
- 增加重试机制（需配合指数退避）

#### 问题 8: 定时任务堆积

**症状**: `schedule_jobs_total` 正常，但任务执行时间超过了调度间隔，导致上一轮未结束下一轮又开始。
**排查**:

```promql
# 查看任务执行耗时
histogram_quantile(0.99, sum(rate(schedule_job_duration_seconds_bucket[5m])) by (le, task))
```

**解决方案**:

- 增加分布式锁，确保同一时刻只有一个实例执行任务
- 优化任务逻辑，减少单次执行时间
- 调整调度间隔或使用消息队列异步处理
