# AI Server 磁盘与数据架构

针对 **1TB NVMe + LVM** 的家庭 AI 服务器方案。

核心原则：

```
系统与数据分离
数据集中管理
```

---

# 一、磁盘结构

| 挂载点     | 容量     | 用途                         |
| ------- | ------ | ---------------------------- |
| `/`     | 100GB  | Ubuntu 系统、Docker 引擎、CLI 工具 |
| `/data` | ~820GB | **服务器核心数据区**               |

保留 **约 30GB LVM 空闲空间**，用于未来扩展或快照。

```
NVMe SSD (1TB)
│
├─ /boot/efi     1GB
├─ /boot         2GB
│
└─ LVM
     ├─ /        100GB
     └─ /data    ~820GB
```

---

# 二、`/data` 目录规划

```
/data
│
├─ docker     # Docker 镜像与容器数据
│
├─ svc        # 各服务持久化数据
│   ├─ mysql
│   ├─ redis
│   └─ openclaw
│
├─ dev        # 个人项目代码
│
├─ log        # 服务日志
│
└─ bak        # 本地备份
```

---

## 各目录解释

### `/data/docker`

Docker 默认数据目录在 `/var/lib/docker`，长期运行会撑爆系统盘。

配置 Docker 使用 `/data/docker`，避免污染系统。

---

### `/data/svc`

存放所有服务的数据目录，Docker 容器通过 bind mount 使用：

```
/data/svc/mysql:/var/lib/mysql
/data/svc/redis:/data
/data/svc/openclaw:/app/data
```

优点：数据透明、易备份。

---

### `/data/dev`

个人开发项目：量化系统、AI 工具、脚本、实验代码。

---

### `/data/log`

服务日志。长期运行日志量大，单独目录便于清理和轮转。

---

### `/data/bak`

本地备份：数据库 dump、配置文件、重要数据。

---

# 三、最终结构

```
Linux Server
│
├─ /boot
├─ /boot/efi
│
├─ /
│   ├─ system
│   ├─ docker engine
│   └─ CLI tools
│
└─ /data
    ├─ docker
    ├─ svc
    │   ├─ mysql
    │   ├─ redis
    │   └─ openclaw
    ├─ dev
    ├─ log
    └─ bak
```

---

# 六、方案优势

1. **系统可快速恢复** — 重装系统后 mount `/data` 即可恢复服务
2. **数据集中管理** — 所有核心数据在 `/data`，方便备份和迁移
3. **Docker 不污染系统** — 避免 `/var/lib/docker` 爆满

---

# 七、下一步

完成磁盘规划后，按顺序部署：

1. 安装 Docker
2. 配置 Docker data-root（创建 `/etc/docker/daemon.json`，指向 `/data/docker`）
3. 部署 MySQL / Redis / OpenClaw