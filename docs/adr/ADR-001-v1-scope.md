# ADR-001: v1 scope

第一版先实现本地可运行闭环，不把生产门禁伪装成已完成：

- Go API 是唯一状态入口；
- Playwright 只负责浏览器生命周期和观测；
- mitmproxy 只做轻量流量事件适配；
- 脱敏在 Broker worker/API 侧集中实现；
- 浏览器与 API lane 的完整 relay、持久化、KMS、WASM 实验和不可篡改审计列入下一阶段；
- 真实目标必须由使用者确认授权并在独立环境配置。
