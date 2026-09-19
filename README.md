# Broker Universal Shielded Execution Platform — v1

第一版目标：可启动的本地功能闭环骨架。

## 已实现

- Go Broker API：健康检查、目标注册、执行租约、文本脱敏。
- TypeScript Browser Controller：Playwright/Chromium 启动、导航、click、fill、网络事件监听。
- Python mitmproxy addon：请求/响应元数据适配器。
- 统一的目标与租约模型，为 Browser lane/API lane 后续接入预留边界。

## 运行

```bash
go run ./cmd/broker-api
# 另一个终端
curl http://127.0.0.1:8080/healthz
curl http://127.0.0.1:8080/v1/targets
curl -X POST http://127.0.0.1:8080/v1/redact -H 'content-type: application/json' \\
  -d '{"text":"Authorization: Bearer secret123","rules":{"secret123":"TOKEN_1"}}'
```

浏览器控制器：

```bash
cd browser-controller
npm install
npx playwright install chromium
npm run dev
```

mitmproxy：

```bash
mitmdump -s mitm-addon/addon.py
```

## 当前状态

这是 v1 开发基线，不是最终生产门禁版本。当前存储为内存实现，目标校验、Browser lease 与 egress relay 仍需下一迭代接入；任何真实目标使用前必须在授权环境中配置代理、目标策略和密钥隔离。
