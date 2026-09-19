# Broker Universal Shielded Execution Platform v1.1

Broker 是通用屏蔽执行平台：不削弱 Playwright、mitmproxy 或自定义适配器的功能；它把真实凭据、浏览器状态和原始响应留在执行平面，并向 Agent/AI 暴露可控、可脱敏的结果。

> 仅在已获授权的目标和数据上使用。Broker 不替代法律审查、授权确认或目标方规则。

## 本版已实现的闭环

- **Go API**：健康检查、目标注册与 Origin 校验、持久化状态、短期租约、租约撤销、预算计数、事件接收、事件查询、文本脱敏。
- **Browser lane**：Playwright/Chromium 生命周期、导航/click/fill/wait/screenshot、请求/响应事件回传到租约事件接口。
- **API lane 基础**：统一目标和租约模型，可供 API 适配器提交事件；目标 Origin 不匹配的事件会被拒绝。
- **MCP gateway**：标准 JSON-RPC stdio 的 initialize、tools/list、tools/call；提供 health、targets、redact 三个工具。
- **mitmproxy addon**：请求/响应元数据事件适配器。
- **本地部署**：Docker Compose、Nginx 测试目标、Go 单元测试。
- **协议与阶段基线（Phase 1）**：引入 Protobuf 统一协议源（`proto/broker/v1/broker.proto`）、Buf 规则配置（`buf.yaml`）与缺口关闭基线说明（`docs/phase1-gap-closure.md`）。

## 快速运行

```bash
export BROKER_STATE=./broker-state.json
go test ./...
go run ./cmd/broker-api
```

```bash
curl http://127.0.0.1:8080/healthz
curl http://127.0.0.1:8080/v1/targets
curl -X POST http://127.0.0.1:8080/v1/leases \
  -H 'content-type: application/json' \
  -d '{"target_id":"demo-web","ttl_seconds":300,"budget":100}'
```

事件必须属于租约，且 URL 必须匹配该目标的 Origin：

```bash
curl -X POST http://127.0.0.1:8080/v1/leases/LEASE_ID/events \
  -H 'content-type: application/json' \
  -d '{"type":"request","method":"GET","url":"http://127.0.0.1:4173/"}'
```

启动浏览器控制器：

```bash
cd browser-controller
npm install
npx playwright install chromium
BROKER_LEASE=LEASE_ID npm run dev
```

启动 MCP stdio 网关：

```bash
cd mcp-gateway
npm install
npm run dev
```

启动代理适配器：

```bash
mitmdump -s mitm-addon/addon.py
```

Protobuf 语法与规范校验：

```bash
buf lint
```

## API 约定

- `POST /v1/targets` 注册 `{id,name,origins[]}`；只接受 `http/https` Origin。
- `POST /v1/leases` 创建 `{target_id,ttl_seconds,budget}`；TTL 最大 24 小时，预算最大 100000。
- `POST /v1/leases/{id}/revoke` 立即撤销租约。
- `POST /v1/leases/{id}/events` 接收请求、响应和适配器事件；每个事件消耗一个预算单位。
- `GET /v1/leases/{id}/events` 查询事件。
- `POST /v1/redact` 执行显式规则和 Authorization/Cookie 基础脱敏。

## 状态与边界

已完成的是可运行的 v1 功能基线，不是生产安全门禁。尚未声称完成：Browser Relay 的完整请求级代理、真正的密钥服务/KMS、不可篡改审计、多租户身份系统、WebCrypto/WASM 沙箱、HTTP/3 专用适配器、生产级策略引擎和大体量 artifact 存储。这些应作为后续版本独立交付，而不能由当前元数据事件适配器冒充完成。

生产接入前至少应替换：持久化层的访问控制、认证授权、密钥存储、脱敏规则管理、审计存储、出网策略和浏览器隔离配置。
