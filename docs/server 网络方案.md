## 最终版：网络系统设计方案 (The Ultimate Version)

### 一、 宿主机环境预调优 (Host Tuning)

```bash
# 开启转发、BBR，增大 UDP 缓冲区
cat <<EOF | sudo tee /etc/sysctl.d/99-singbox.conf
net.core.rmem_max=16777216
net.core.wmem_max=16777216
net.ipv4.tcp_congestion_control=bbr
net.core.default_qdisc=fq
EOF
sudo sysctl --system
sudo modprobe tun
```

---

### 二、 最终版 Sing-box 配置文件 (`config.json`)

```json
{
  "log": { "level": "info", "timestamp": true },
  "dns": {
    "servers": [
      {
        "tag": "dns-cn-main",
        "type": "https",
        "server": "223.5.5.5",
        "detour": "direct"
      },
      {
        "tag": "dns-cn-backup",
        "type": "https",
        "server": "120.53.53.53",
        "detour": "direct"
      },
      {
        "tag": "dns-proxy-main",
        "type": "https",
        "server": "1.1.1.1",
        "detour": "proxy"
      },
      {
        "tag": "dns-proxy-backup",
        "type": "https",
        "server": "8.8.8.8",
        "detour": "proxy"
      }
    ],
    "rules": [
      {
        "domain_suffix": [".ts.net", ".home.arpa", ".local"],
        "action": "route",
        "server": "dns-cn-main"
      },
      {
        "rule_set": ["geosite-cn"],
        "action": "route",
        "server": "dns-cn-main"
      },
      {
        "rule_set": ["geosite-gfw"],
        "action": "route",
        "server": "dns-proxy-main"
      },
      { "fallback_status": [2], "action": "route", "server": "dns-proxy-main" }
    ],
    "final": "dns-proxy-main",
    "strategy": "prefer_ipv4",
    "independent_cache": true
  },
  "inbounds": [
    {
      "type": "tun",
      "inet4_address": "172.19.0.1/30",
      "auto_route": true,
      "strict_route": true,
      "sniff": true,
      "sniff_override_destination": true
    }
  ],
  "outbounds": [
    { "type": "direct", "tag": "direct" },
    {
      "type": "hysteria2",
      "tag": "proxy",
      "server": "YOUR_VPS_IP",
      "server_port": 443,
      "password": "YOUR_PASSWORD",
      "tls": { "enabled": true, "server_name": "YOUR_DOMAIN" }
    },
    {
      "type": "fallback",
      "tag": "smart",
      "outbounds": ["direct", "proxy"],
      "url": "https://www.gstatic.com/generate_204",
      "interval": "5m",
      "tolerance": 100
    },
    { "type": "dns", "tag": "dns-out" }
  ],
  "route": {
    "auto_detect_interface": true,
    "rules": [
      { "protocol": "dns", "outbound": "dns-out" },
      {
        "ip_cidr": [
          "10.0.0.0/8",
          "172.16.0.0/12",
          "192.168.0.0/16",
          "100.64.0.0/10"
        ],
        "outbound": "direct"
      },
      { "rule_set": ["geosite-gfw"], "outbound": "proxy" },
      { "rule_set": ["geosite-cn", "geoip-cn"], "outbound": "direct" }
    ],
    "final": "smart",
    "rule_set": [
      {
        "tag": "geosite-cn",
        "type": "remote",
        "format": "binary",
        "url": "https://raw.githubusercontent.com/SagerNet/sing-geosite/rule-set/geosite-cn.srs"
      },
      {
        "tag": "geosite-gfw",
        "type": "remote",
        "format": "binary",
        "url": "https://raw.githubusercontent.com/SagerNet/sing-geosite/rule-set/geosite-gfw.srs"
      },
      {
        "tag": "geoip-cn",
        "type": "remote",
        "format": "binary",
        "url": "https://raw.githubusercontent.com/SagerNet/sing-geoip/rule-set/geoip-cn.srs"
      }
    ]
  }
}
```

---

### 三、 终版网络逻辑总结 (System Behavior)

| 场景                     | 最终行为         | 逻辑路径                                                         |
| :----------------------- | :--------------- | :--------------------------------------------------------------- |
| **Tailscale / 内网**     | **绝对直连**     | 匹配 `100.64.0.0/10` 或 `.ts.net` 直接出站，不进代理。           |
| **国内已知网站**         | **极速直连**     | 匹配 `geosite-cn` 使用阿里 DoH，走 `direct` 出站。               |
| **被墙已知网站**         | **无污染代理**   | 匹配 `geosite-gfw` 使用 1.1.1.1 解析，走 `proxy` 出站。          |
| **国内解析异常**         | **自动修复**     | 若阿里 DNS 返回解析错误，自动 fallback 到 1.1.1.1 获取正确 IP。  |
| **国外未封站 (如 Bing)** | **智能加速**     | 命中 `smart`：若直连延迟 < 100ms 则直连，否则自动切换 Hy2 代理。 |
| **未知流量兜底**         | **代理优先解析** | 域名解析默认走代理 DNS，路由默认由 `smart` 判断最优路径。        |

---

### 四、 部署 Checklist

- **Docker**: 确保 `network_mode: host` 和 `privileged: true`。
- **Hysteria 2**: 确保 VPS 的 443 端口（UDP）在防火墙中已开放。
- **Tailscale**: 如果这台机器是 Tailscale 节点，建议在 Tailscale 管理面板开启 `Accept Routes`，以配合 Sing-box 的 TUN 模式。

### 结语

这套方案目前是 **Linux 透明网关** 领域最稳健的架构之一。它兼顾了**解析速度**（国内 DoH）、**安全性**（代理 DNS 兜底）、**容错性**（Smart Fallback）以及**私有组网兼容性**（Tailscale 直连）。

**方案已完全就绪，可以执行。** 你想让我为你生成一个 **docker-compose.yml** 示例，将这些特权模式和网路挂载逻辑直接转化为可运行的命令吗？
