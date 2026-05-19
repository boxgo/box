# Alert Rules Generator

自动生成针对特定 namespace 和 job 的 Prometheus 告警规则。

## 功能特性

- 📦 基于模板自动生成定制化告警规则
- 🎯 支持按 namespace 和 job 过滤
- 🔧 Shell 脚本实现，无需额外依赖
- ✅ 自动添加标签过滤器到所有 PromQL 表达式
- 🔍 内置验证功能，自动检查生成的规则

## 使用方法

```bash
# 基本用法（自动验证）
./scripts/generate_alerts.sh <namespace> <job>

# 指定输出文件
./scripts/generate_alerts.sh <namespace> <job> <output_file>

# 跳过验证
./scripts/generate_alerts.sh <namespace> <job> --no-verify
```

**示例：**

```bash
# 为 prod 命名空间的 api-service 生成告警规则（自动验证）
./scripts/generate_alerts.sh prod api-service

# 输出: docs/prod_api-service_alerts.yaml

# 自定义输出路径
./scripts/generate_alerts.sh prod api-service alerts/production/api.yaml

# 快速生成，跳过验证
./scripts/generate_alerts.sh prod api-service --no-verify
```

## 参数说明

| 参数 | 说明 | 必需 |
|------|------|------|
| `namespace` | Kubernetes 命名空间 | 是 |
| `job` | 服务名称（Job） | 是 |
| `output_file` | 输出文件路径 | 否，默认: `docs/${namespace}_${job}_alerts.yaml` |
| `--no-verify` | 跳过自动验证 | 否 |
| 模板文件 | 告警规则模板 | - (固定: `docs/prometheus_alerts_template.yaml`) |


## 生成的文件内容

生成的告警规则文件会：

1. **添加标签过滤器**：所有 PromQL 表达式都会添加 `namespace` 和 `job` 过滤条件
2. **保留原有结构**：保持原模板的告警组织结构
3. **添加生成信息**：文件头部包含生成参数和时间

**示例对比：**

**原模板：**
```yaml
- alert: HttpServerQpsHigh
  expr: sum by (namespace, job) (rate(http_server_requests_total[1m])) > 1000
```

**生成后（namespace=prod, job=api-service）：**
```yaml
- alert: HttpServerQpsHigh
  expr: sum by (namespace, job) (rate(http_server_requests_total{namespace="prod",job="api-service"}[1m])) > 1000
```

## 部署到 Prometheus

生成告警规则后，有以下几种部署方式：

### 1. Kubernetes ConfigMap 方式

```bash
# 创建 ConfigMap
kubectl create configmap prod-api-alerts \
  --from-file=docs/prod_api-service_alerts.yaml \
  -n monitoring

# 在 Prometheus 配置中引用
# prometheus.yml:
# rule_files:
#   - '/etc/prometheus/rules/prod_api-service_alerts.yaml'
```

### 2. 直接文件挂载

```yaml
# prometheus-deployment.yaml
volumeMounts:
  - name: alert-rules
    mountPath: /etc/prometheus/rules
volumes:
  - name: alert-rules
    configMap:
      name: prod-api-alerts
```

### 3. Prometheus Operator 方式

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: prod-api-alerts
  namespace: monitoring
spec:
  groups:
    - name: box-http-server-alerts
      interval: 30s
      rules:
        # 粘贴生成的告警规则
```

## 目录结构

```
scripts/
├── README.md                    # 本文档
└── generate_alerts.sh           # 告警生成和验证脚本

docs/
├── prometheus_alerts_template.yaml  # 告警规则模板
├── prod_api-service_alerts.yaml     # 生成的告警文件示例
└── ...
```

## 支持的指标

脚本会自动为以下所有指标添加 `namespace` 和 `job` 过滤器：

| 组件 | 指标 |
|------|------|
| **HTTP Server** | `http_server_requests_total` |
| | `http_server_request_duration_seconds_bucket` |
| | `http_server_request_duration_seconds_count` |
| **HTTP Client** | `http_client_requests_total` |
| | `http_client_requests_inflight` |
| | `http_client_request_duration_seconds_bucket` |
| **gRPC Server** | `grpc_server_requests_total` |
| | `grpc_server_panics_total` |
| **Database** | `db_client_request_duration_seconds_count` |
| | `db_client_request_duration_seconds_bucket` |
| | `db_client_connections_in_use` |
| | `db_client_connections_max_open` |
| **Redis** | `redis_client_requests_total` |
| | `redis_client_request_duration_seconds_bucket` |
| **MongoDB** | `mongo_client_requests_total` |
| | `mongo_client_request_duration_seconds_bucket` |
| **Schedule** | `schedule_jobs_total` |
| **Go Runtime** | `go_goroutines`, `go_threads` |
| | `go_memstats_sys_bytes`, `go_memstats_heap_alloc_bytes` |
| | `go_gc_duration_seconds`, `go_gc_duration_seconds_count` |
| | `go_memstats_gc_cpu_fraction` |

## 使用场景

### 场景 1: 多环境部署

为不同环境生成独立的告警规则：

```bash
./scripts/generate_alerts.sh prod api-service
./scripts/generate_alerts.sh staging api-service
./scripts/generate_alerts.sh dev api-service
```

### 场景 2: 微服务架构

为每个微服务生成专属告警：

```bash
for service in user-service order-service payment-service; do
  ./scripts/generate_alerts.sh prod $service
done
```

### 场景 3: 批量生成脚本

创建批量生成脚本：

```bash
#!/bin/bash
# batch_generate.sh

NAMESPACE="prod"
SERVICES=(
  "api-service"
  "user-service"
  "order-service"
  "payment-service"
)

for service in "${SERVICES[@]}"; do
  echo "Generating alerts for $service..."
  if ./scripts/generate_alerts.sh "$NAMESPACE" "$service"; then
    echo "✓ $service - OK"
  else
    echo "✗ $service - FAILED"
    exit 1
  fi
done

echo "All services processed!"
```

### 场景 4: CI/CD 集成

#### GitLab CI

```yaml
generate-alerts:
  stage: build
  script:
    - ./scripts/generate_alerts.sh ${CI_ENVIRONMENT_NAME} ${SERVICE_NAME}
  artifacts:
    paths:
      - docs/*_alerts.yaml
```

#### GitHub Actions

```yaml
- name: Generate alerts
  run: |
    ./scripts/generate_alerts.sh ${{ vars.NAMESPACE }} ${{ vars.SERVICE }}
```

## 脚本输出说明

### 标准输出（带验证）

```bash
$ ./scripts/generate_alerts.sh prod api-service

╔════════════════════════════════════════════════════════════╗
║  Prometheus Alert Rules Generator & Verifier              ║
╚════════════════════════════════════════════════════════════╝

📋 Configuration:
  Namespace: prod
  Job: api-service
  Template: docs/prometheus_alerts_template.yaml
  Output: docs/prod_api-service_alerts.yaml
  Verify: Enabled

🔨 Step 1: Generating alert rules...
✓ Alert rules generated successfully!

🔍 Step 2: Verifying alert rules...
✓ All metrics are correctly filtered!

╔════════════════════════════════════════════════════════════╗
║  Summary                                                   ║
╚════════════════════════════════════════════════════════════╝

✓ Generation completed
✓ Verification passed

📄 Output file: docs/prod_api-service_alerts.yaml

📝 Next steps:
  1. Review the generated file
  2. Validate with promtool
  3. Deploy to Kubernetes
```

### 快速模式输出（跳过验证）

```bash
$ ./scripts/generate_alerts.sh prod api-service --no-verify

╔════════════════════════════════════════════════════════════╗
║  Prometheus Alert Rules Generator & Verifier              ║
╚════════════════════════════════════════════════════════════╝

📋 Configuration:
  Namespace: prod
  Job: api-service
  Verify: Disabled

🔨 Step 1: Generating alert rules...
✓ Alert rules generated successfully!

⚠ Verification skipped (--no-verify flag)

✓ Generation completed
```

## 验证生成的规则

### 自动验证（默认启用）

脚本会自动验证生成的告警规则，确保所有指标都正确添加了过滤器。

如需跳过验证（快速生成）：

```bash
./scripts/generate_alerts.sh prod api-service --no-verify
```

### 使用 promtool 验证语法

```bash
# 使用 promtool 验证语法
promtool check rules docs/prod_api-service_alerts.yaml

# 或使用 Docker
docker run --rm -v $(pwd):/workspace prom/prometheus:latest \
  promtool check rules /workspace/docs/prod_api-service_alerts.yaml
```

### 完整的生成和部署流程

```bash
# 1. 生成和验证告警规则（自动验证）
./scripts/generate_alerts.sh prod api-service

# 2. （可选）使用 promtool 验证语法
promtool check rules docs/prod_api-service_alerts.yaml

# 3. 部署
kubectl create configmap prod-api-alerts \
  --from-file=docs/prod_api-service_alerts.yaml \
  -n monitoring
```

## 常见问题

### Q1: 如何确保所有指标都添加了过滤器？

脚本默认会自动验证生成的规则：

```bash
./scripts/generate_alerts.sh prod api-service
# 会自动验证所有指标
```

如果有指标遗漏过滤器，脚本会明确指出并返回错误。

### Q2: 生成的规则不生效？

检查以下几点：
1. Prometheus 配置中是否正确引用了规则文件
2. 规则文件的 YAML 格式是否正确
3. Prometheus 是否成功重载了配置（查看日志）

```bash
# 重载 Prometheus 配置
curl -X POST http://prometheus:9090/-/reload
```

### Q3: 如何修改告警阈值？

两种方式：
1. 修改模板文件 `docs/prometheus_alerts_template.yaml`，然后重新生成
2. 直接编辑生成的文件（不推荐，因为会在下次生成时被覆盖）

### Q4: 支持批量生成吗？

是的，可以通过循环实现：

```bash
# 为多个服务批量生成
for service in api-service user-service order-service; do
  ./scripts/generate_alerts.sh prod $service
done
```

## 快速参考

```bash
# 基本用法
./scripts/generate_alerts.sh <namespace> <job>

# 自定义输出
./scripts/generate_alerts.sh <namespace> <job> <output_file>

# 快速模式（跳过验证）
./scripts/generate_alerts.sh <namespace> <job> --no-verify

# 批量生成
for svc in api user order; do
  ./scripts/generate_alerts.sh prod ${svc}-service
done

# 验证语法
promtool check rules docs/prod_api-service_alerts.yaml

# 部署
kubectl create configmap <name> \
  --from-file=<alert_file> \
  -n monitoring
```

## 故障排查

### 问题：验证失败

```bash
✗ Found unfiltered metric: http_server_requests_total
```

**解决方案：**
1. 检查模板文件是否被修改
2. 重新生成文件
3. 如果问题持续，联系维护者

### 问题：告警未触发

**检查步骤：**

```bash
# 1. 检查规则是否加载
curl http://prometheus:9090/api/v1/rules | jq

# 2. 检查指标是否存在
curl 'http://prometheus:9090/api/v1/query?query=http_server_requests_total{namespace="prod",job="api-service"}' | jq

# 3. 重载 Prometheus
curl -X POST http://prometheus:9090/-/reload
```

## 最佳实践

1. **默认使用验证**：生产环境务必验证
2. **版本控制**：将生成的文件提交到 Git
3. **定期更新**：模板更新后重新生成所有文件
4. **命名规范**：使用 `<namespace>_<job>_alerts.yaml` 格式
5. **测试先行**：在测试环境验证后再部署生产

## 贡献

如需改进脚本或添加新功能，请：

1. Fork 项目
2. 创建功能分支
3. 提交 Pull Request

## 相关文档

- [prometheus_alerts_template.yaml](../docs/prometheus_alerts_template.yaml) - 告警规则模板
- [metric.md](../docs/metric.md) - 指标和看板文档

## 许可

与主项目保持一致。
