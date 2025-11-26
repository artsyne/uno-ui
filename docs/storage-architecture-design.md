# iDiary.ai Hebe 存储架构设计

## 项目概述

### 项目背景

- **项目名称**: iDiary.ai Hebe - 智能日记管理系统
- **核心价值**: 支持多云存储服务，用户可将日记内容存储到个人云盘
- **技术框架**: Go 1.25 + DDD架构 + 泛型DI容器 + 多云存储集成

### 业务特性

- **多云存储**: 支持 Dropbox、Google Drive、OneDrive、阿里云盘
- **混合存储模式**: 日记元信息存储于 iDiary.ai 数据库，日记内容（Markdown 文件）存储于用户个人云盘
- **多内容类型**: 文字、图片、视频、链接等富媒体内容
- **用户自主**: 用户控制自己的存储授权和数据
- **智能配置**: 每个日记本可独立配置存储策略
- **安全与验证**: 通过验证串保证云端文件内容不被篡改

## DDD 存储架构设计

### 当前架构现状

基于 iDiary.ai Hebe 的 DDD 分层架构，存储功能分布在以下层次：

```
┌─────────────────────────────────────────────────────────┐
│            Website Layer (HTML Templates)               │
├─────────────────────────────────────────────────────────┤
│              App Layer (Controllers/ViewService)        │
├─────────────────────────────────────────────────────────┤
│    Domain Layer (User/Diary/Post/Task Domains)         │
├─────────────────────────────────────────────────────────┤
│         Infrastructure Layer (Repositories)            │
│    UserRepo    │  DiaryRepo   │  PostRepo   │  TaskRepo  │
├─────────────────────────────────────────────────────────┤
│              Core Layer (Database/DI/Config)           │
│      MySQL       │    Redis     │   SQLx     │           │
└─────────────────────────────────────────────────────────┘

[多云存储扩展规划]
┌─────────────────────────────────────────────────────────┐
│         Cloud Storage Abstraction (Future)             │
│   Dropbox  │ GoogleDrive │   OneDrive  │   阿里云盘      │
└─────────────────────────────────────────────────────────┘
```

### 存储现状分析

#### 已实现的存储功能

1. **MySQL 主存储**

   - 用户信息、日记元数据、文章内容
   - 基于 SQLx 的类型安全数据库操作
   - Repository 模式实现数据访问层

2. **Redis 缓存存储**

   - 用户会话管理
   - 临时数据缓存
   - 提高系统响应速度

3. **本地文件存储**
   - 静态资源（CSS、JS、图片）
   - 上传文件的临时存储
   - 日志文件存储

#### 云存储功能规划

1. **云存储接口设计**

   ```go
   // CloudStorage 定义了与云存储服务交互的统一接口
   type CloudStorage interface {
       // Upload 将日记内容上传到指定的云端路径
       Upload(ctx context.Context, remotePath string, content io.Reader) (*UploadResult, error)
       // Download 从指定的云端路径下载日记内容
       Download(ctx context.Context, remotePath string) (io.ReadCloser, error)
       // Delete 删除指定的云端文件
       Delete(ctx context.Context, remotePath string) error
       // GetFileMetadata 获取云端文件的元信息，用于校验
       GetFileMetadata(ctx context.Context, remotePath string) (*FileMetadata, error)
   }

   // UploadResult 包含上传操作的结果
   type UploadResult struct {
       FileID   string // 云端文件的唯一标识
       Path     string // 文件的云端路径
       Version  string // 文件的版本或 ETag，用于后续的一致性检查
       Size     int64  // 文件大小
   }

   // FileMetadata 包含云端文件的元信息
   type FileMetadata struct {
       FileID   string
       Path     string
       Version  string
       Size     int64
       Modified time.Time // 文件最后修改时间
   }
   ```

2. **存储策略配置**

   - 每个日记本独立配置存储后端（即选择哪个云盘）
   - 用户自主授权云存储访问
   - 支持存储迁移和备份

3. **安全与验证**
   - 通过文件内容哈希生成的验证串，防止云端文件被意外篡改
   - OAuth 2.0 协议安全管理用户授权

## 技术实现路径

### 阶段一：基础架构完善（已完成）

- ✅ DDD 分层架构
- ✅ Go 1.25 泛型 DI 容器
- ✅ Repository 模式
- ✅ MySQL + Redis 存储

### 阶段二：云存储抽象层（规划中）

- 📋 设计 CloudStorage 接口
- 📋 实现 Dropbox 适配器
- 📋 实现 Google Drive 适配器
- 📋 实现阿里云盘适配器

### 阶段三：高级功能（后续）

- 📋 数据加密和解密
- 📋 存储迁移工具
- 📋 同步和备份机制
- 📋 存储用量监控

## 总结

iDiary.ai Hebe 的存储架构采用分层设计，当前已实现稳定的本地存储基础，为未来的多云存储扩展奠定了坚实基础。通过 DDD 架构和泛型 DI 容器，系统具备良好的扩展性和可维护性。

### 2.2 模块职责划分

#### User模块

- **基础职责**: 用户注册、登录、认证、个人资料管理
- **扩展职责**:
  - 云盘授权管理（OAuth流程、Token刷新）
  - 用户存储偏好设置
  - 存储配额监控
  - 提供统一的认证和授权中间件

#### Diary模块

- **基础职责**: 日记本CRUD操作、成员管理
- **扩展职责**:
  - 存储策略配置管理
  - 调用Post模块获取内容列表（解决循环依赖）
  - 按内容类型配置不同存储策略

#### Post模块

- **基础职责**: 内容CRUD操作
  - 支持多种内容类型（文字、图片、视频、链接、文件等）
- **核心职责**:
  - 实现多存储后端支持
  - 内容类型识别和处理
  - 加密/解密功能
  - 文件上传和存储路由

#### Task模块

- **职责**: 后台任务处理、邮件发送、定时作业
- **保持现状**: 继续使用数据库存储，职责不变

### 2.3 依赖关系

```
User模块 (独立模块，提供认证和授权服务)
   ↑
Post模块 (依赖User获取云盘授权，并提供统一的存储服务)
   ↑
Diary模块 (依赖Post获取内容，传递存储配置)
```

**关键设计原则**:

- Post模块不主动获取Diary配置，通过参数接收
- Diary模块作为编排者，传递完整存储配置给Post模块
- User模块提供云盘授权服务，供Post模块调用

## 3. 数据库设计

### 3.1 现有表结构说明

#### br_user (用户表)

```sql
-- 用户基础信息表，保持现有结构并扩展存储相关字段
ALTER TABLE br_user ADD COLUMN storage_quota_used bigint unsigned DEFAULT 0 COMMENT '已使用存储空间(字节)';
ALTER TABLE br_user ADD COLUMN storage_quota_limit bigint unsigned DEFAULT 5368709120 COMMENT '存储配额限制(字节，默认5GB)';

```

#### br_diary (日记本表)

```sql
-- 日记本信息表，扩展存储配置字段
ALTER TABLE br_diary ADD COLUMN storage_config text DEFAULT NULL COMMENT '存储配置JSON，包含存储类型和云盘授权ID';

-- storage_config JSON结构示例:
-- {
--   "storage_type": "cloud",             // 存储类型: database/cloud
--   "cloud_auth_id": 123                 // 云盘授权ID
-- }
```

#### br_post (内容表)

```sql
-- 多媒体内容表，扩展存储相关字段
ALTER TABLE br_post ADD COLUMN storage_type varchar(20) DEFAULT 'database' COMMENT '存储类型: database/cloud';
ALTER TABLE br_post ADD COLUMN storage_path varchar(500) DEFAULT '' COMMENT '云盘存储路径或文件ID';
ALTER TABLE br_post ADD COLUMN content_hash varchar(64) DEFAULT '' COMMENT '文件内容的SHA-256哈希值，用于验证文件完整性';
ALTER TABLE br_post ADD COLUMN file_size bigint unsigned DEFAULT 0 COMMENT '文件大小(字节)';
ALTER TABLE br_post ADD COLUMN mime_type varchar(100) DEFAULT 'text/markdown' COMMENT 'MIME类型';

-- info_type字段说明 (已有字段，重新定义含义):
-- 0: 文字日记 (Markdown)
-- 1-99: 预留扩展
```

### 3.2 新增表结构

#### br_cloud_authorization (云盘授权表)

```sql
CREATE TABLE br_cloud_authorization (
  auth_id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '授权ID，主键',
  user_id bigint unsigned NOT NULL COMMENT '用户ID，关联br_user.user_id',
  provider varchar(50) NOT NULL COMMENT '云盘提供商: dropbox/googledrive/alidrive/baidupan/onedrive/webdav等',
  provider_user_id varchar(255) DEFAULT '' COMMENT '云盘平台的用户ID',
  provider_email varchar(255) DEFAULT '' COMMENT '云盘平台的用户邮箱',
  access_token text COMMENT '访问令牌，AES加密存储',
  refresh_token text COMMENT '刷新令牌，AES加密存储',
  expires_at bigint unsigned DEFAULT 0 COMMENT '令牌过期时间(UTC时间戳)',
  scope varchar(500) DEFAULT '' COMMENT '授权范围，如read、write、delete等',
  is_active tinyint unsigned DEFAULT 1 COMMENT '是否激活: 0停用 1启用',
  nickname varchar(100) DEFAULT '' COMMENT '用户给此授权起的别名，如"工作Dropbox"',
  quota_total bigint unsigned DEFAULT 0 COMMENT '云盘总容量(字节)',
  quota_used bigint unsigned DEFAULT 0 COMMENT '云盘已使用容量(字节)',
  last_sync_time bigint unsigned DEFAULT 0 COMMENT '最后同步时间(UTC时间戳)',
  created_utc_time bigint unsigned NOT NULL COMMENT '创建时间(UTC时间戳)',
  updated_utc_time bigint unsigned DEFAULT 0 COMMENT '更新时间(UTC时间戳)',
  PRIMARY KEY (auth_id),
  KEY idx_user_provider (user_id, provider) COMMENT '用户-提供商索引',
  KEY idx_user_active (user_id, is_active) COMMENT '用户激活状态索引',
  FOREIGN KEY (user_id) REFERENCES br_user(user_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户云盘授权信息表';
```

#### br_user_storage_preference (用户存储偏好表)

```sql
CREATE TABLE br_user_storage_preference (
  user_id bigint unsigned NOT NULL COMMENT '用户ID，主键，关联br_user.user_id',
  default_storage_type varchar(20) DEFAULT 'database' COMMENT '默认存储类型: database/cloud/hybrid',
  default_cloud_auth_id bigint unsigned DEFAULT NULL COMMENT '默认云盘授权ID，关联br_cloud_authorization.auth_id',
  auto_backup_enabled tinyint unsigned DEFAULT 0 COMMENT '是否启用自动备份: 0否 1是',
  backup_cloud_auth_id bigint unsigned DEFAULT NULL COMMENT '备份云盘授权ID，关联br_cloud_authorization.auth_id',
  encryption_enabled tinyint unsigned DEFAULT 0 COMMENT '是否默认启用加密: 0否 1是',
  sync_frequency varchar(20) DEFAULT 'realtime' COMMENT '同步频率: realtime/hourly/daily/manual',
  compression_quality tinyint unsigned DEFAULT 80 COMMENT '图片压缩质量(1-100)',
  video_quality varchar(20) DEFAULT 'medium' COMMENT '视频质量: low/medium/high/original',
  created_utc_time bigint unsigned NOT NULL COMMENT '创建时间(UTC时间戳)',
  updated_utc_time bigint unsigned DEFAULT 0 COMMENT '更新时间(UTC时间戳)',
  PRIMARY KEY (user_id),
  FOREIGN KEY (user_id) REFERENCES br_user(user_id) ON DELETE CASCADE,
  FOREIGN KEY (default_cloud_auth_id) REFERENCES br_cloud_authorization(auth_id) ON DELETE SET NULL,
  FOREIGN KEY (backup_cloud_auth_id) REFERENCES br_cloud_authorization(auth_id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户存储偏好设置表';
```

#### br_cloud_sync_log (云盘同步日志表)

```sql
CREATE TABLE br_cloud_sync_log (
  log_id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '日志ID，主键',
  user_id bigint unsigned NOT NULL COMMENT '用户ID',
  auth_id bigint unsigned NOT NULL COMMENT '云盘授权ID',
  post_id bigint unsigned DEFAULT NULL COMMENT '内容ID，关联br_post.post_id',
  operation varchar(20) NOT NULL COMMENT '操作类型: upload/download/delete/sync',
  status varchar(20) NOT NULL COMMENT '状态: success/failed/pending',
  file_path varchar(500) DEFAULT '' COMMENT '文件路径',
  file_size bigint unsigned DEFAULT 0 COMMENT '文件大小(字节)',
  error_message text COMMENT '错误信息(如果操作失败)',
  duration_ms int unsigned DEFAULT 0 COMMENT '操作耗时(毫秒)',
  created_utc_time bigint unsigned NOT NULL COMMENT '创建时间(UTC时间戳)',
  PRIMARY KEY (log_id),
  KEY idx_user_time (user_id, created_utc_time) COMMENT '用户时间索引',
  KEY idx_auth_status (auth_id, status) COMMENT '授权状态索引',
  FOREIGN KEY (user_id) REFERENCES br_user(user_id) ON DELETE CASCADE,
  FOREIGN KEY (auth_id) REFERENCES br_cloud_authorization(auth_id) ON DELETE CASCADE,
  FOREIGN KEY (post_id) REFERENCES br_post(post_id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='云盘同步操作日志表';
```

## 4. 存储架构设计

### 4.1 存储层级架构

```
App Layer (应用层)
    ↓
Domain Layer (领域层)
    ↓
Content Processor (内容处理器) → Encryption Layer (加密层)
    ↓
Storage Adapter (存储适配器)
    ↓
Storage Backend (存储后端: DB/Cloud)
```

### 4.2 存储策略矩阵

| 内容类型 | 推荐存储 | 加密支持 | 处理功能         | 备注                       |
| -------- | -------- | -------- | ---------------- | -------------------------- |
| 文字日记 | Database | ✅       | Markdown渲染     | 快速检索                   |
| 图片     | Cloud    | ✅       | 缩略图生成、压缩 | 节省服务器空间             |
| 视频     | Cloud    | ✅       | 转码、截图       | 大文件优先云存储，支持加密 |
| 音频     | Cloud    | ✅       | 格式转换         | 支持多种音频格式           |
| 文件     | Cloud    | ✅       | 病毒扫描         | 通用文件支持               |
| 链接     | Database | ❌       | 元数据抓取       | 链接预览功能               |

### 4.3 存储适配器设计

存储适配器层将遵循统一的 `CloudStorage` 接口，为上层应用提供一致的云存储操作体验。每个适配器负责将通用接口调用转换为特定云盘服务的 API 请求。

#### Database Adapter

- **职责**: 处理存储在 iDiary.ai 本地数据库中的日记内容。
- **支持内容**: 适用于未配置或禁用云存储的日记本。
- **特点**: 访问速度快，依赖内置数据库，无需外部网络请求。

#### Cloud Adapters

每个云盘服务都将有一个独立的适配器实现。

- **Dropbox Adapter**
  - **API**: Dropbox API v2
  - **职责**: 实现文件的上传、下载和删除操作。

- **Google Drive Adapter**
  - **API**: Google Drive API v3
  - **职责**: 实现文件的上传、下载和删除操作。

- **OneDrive Adapter**
  - **API**: Microsoft Graph API
  - **职责**: 实现文件的上传、下载和删除操作。

- **阿里云盘 Adapter**
  - **API**: 官方 API
  - **职责**: 实现文件的上传、下载和删除操作。

### 4.5 云端存储结构设计

#### 目录结构

为了保持云端文件的整洁和高效管理，所有日记文件将存储在指定的目录结构中。在用户授权的云盘根目录下，系统将创建一个名为 `idiary` 的主文件夹。

```
/idiary/{diary_id}/{year}/{month}/{post_id}_{validation_hash}.md
```

- **`/idiary`**: 应用的根目录，用于存放所有与 iDiary.ai 相关的数据。
- **`/{diary_id}`**: 日记本的唯一标识符（ID），每个日记本一个文件夹。
- **`/{year}`**: 四位数的年份，例如 `2024`。
- **`/{month}`**: 两位数的月份，例如 `07`。
- **`/{post_id}_{validation_hash}.md`**: 日记文件名，由三部分组成：
    - `post_id`: 日记内容的唯一标识符（ID）。
    - `validation_hash`: 用于验证文件完整性的哈希串（取哈希值的前8位）。
    - `.md`: 文件扩展名，表示 Markdown 格式。

#### 文件名验证哈希算法

为了防止云端文件被意外修改或篡改，文件名中包含一个验证哈希。此哈希的计算方式如下：

1.  **选择算法**: 使用 `SHA-256` 哈希算法，它在安全性和性能之间取得了良好的平衡。
2.  **计算内容**: 哈希值是根据日记的**完整原始内容**（UTF-8 编码的 Markdown 文本）计算得出的。
3.  **生成验证串**: 将计算出的 `SHA-256` 哈希值（一个64位的十六进制字符串）**截取其前8位**作为文件名中的 `validation_hash`。

**伪代码示例**:

```go
import (
    "crypto/sha256"
    "fmt"
    "io"
)

// calculateValidationHash 计算并返回用于文件名的验证哈希
func calculateValidationHash(content string) string {
    hasher := sha256.New()
    io.WriteString(hasher, content)
    fullHash := fmt.Sprintf("%x", hasher.Sum(nil))
    return fullHash[:8] // 取前8位作为验证串
}
```

**验证流程**:

当从云端下载文件时，系统会重新计算文件内容的 `SHA-256` 哈希值，并将其前8位与文件名中的 `validation_hash`进行比较。如果不匹配，则表明文件可能已损坏或被修改，系统将向用户发出警告。

## 5. 关键技术实现

### 5.1 OAuth授权流程

```
1. 用户点击"授权云盘" → 跳转到云盘授权页面
2. 用户授权后 → 云盘回调到系统callback接口
3. 系统获取authorization_code → 换取access_token和refresh_token
4. 加密存储token → 保存到br_cloud_authorization表
5. 获取用户云盘信息 → 更新quota_total等字段
6. 测试连接有效性 → 设置is_active状态
```

### 5.2 Token刷新机制

```go
// 伪代码示例
func (s *CloudService) RefreshTokenIfNeeded(authID int) error {
    auth := s.GetAuthorization(authID)

    // 检查是否即将过期 (提前30分钟刷新)
    if time.Now().Unix() + 1800 > auth.ExpiresAt {
        newToken := s.callCloudAPI(auth.RefreshToken)
        s.updateTokens(authID, newToken)
    }

    return nil
}
```

### 5.3 存储路由算法

存储路由的核心是根据日记本的配置，决定将日记内容存储到数据库还是云端。这个决策由 `StorageManager` 服务来执行。

```go
// StorageManager 负责根据配置选择合适的存储适配器
type StorageManager struct {
    databaseAdapter StorageAdapter
    cloudAdapters   map[string]StorageAdapter // key: provider_name
}

// RouteStorage 根据日记本的存储配置返回相应的存储适配器
func (m *StorageManager) RouteStorage(config *StorageConfig) (StorageAdapter, error) {
    if config.StorageType == "cloud" {
        // 从 cloudAdapters 中查找对应的云存储适配器
        adapter, ok := m.cloudAdapters[config.CloudProvider]
        if !ok {
            return nil, fmt.Errorf("unsupported cloud provider: %s", config.CloudProvider)
        }
        return adapter, nil
    }
    
    // 默认或明确配置为 database 时，返回数据库适配器
    return m.databaseAdapter, nil
}
```

**处理流程**:

1.  当用户创建或更新一篇日记时，应用层会获取该日记所属日记本的 `storage_config`。
2.  `StorageManager` 的 `RouteStorage` 方法被调用，并传入配置信息。
3.  如果 `storage_type` 为 `cloud`，则根据 `cloud_provider`（例如 `dropbox`, `googledrive`）从适配器映射中选择对应的云存储适配器。
4.  如果 `storage_type` 为 `database` 或未配置，则默认使用 `databaseAdapter`。
5.  获取到适配器后，执行相应的 `Upload`、`Download` 或 `Delete` 操作。

## 6. 部署和运维

### 6.1 配置管理

```yaml
# config/storage.yaml
storage:
  default_type: 'database'

  encryption:
    algorithm: 'AES-256-GCM'
    key_derivation: 'PBKDF2'
    iterations: 10000

  cloud_providers:
    dropbox:
      app_key: '${DROPBOX_APP_KEY}'
      app_secret: '${DROPBOX_APP_SECRET}'
      redirect_uri: 'https://yourdomain.com/auth/dropbox/callback'

    googledrive:
      client_id: '${GOOGLE_CLIENT_ID}'
      client_secret: '${GOOGLE_CLIENT_SECRET}'
      redirect_uri: 'https://yourdomain.com/auth/google/callback'

  quotas:
    free_user: 5368709120 # 5GB
    premium_user: 53687091200 # 50GB

  processing:
    image:
      max_size: 10485760 # 10MB
      quality: 80
      thumbnail_sizes: [150, 300, 600]

    video:
      max_size: 104857600 # 100MB
      max_duration: 600 # 10分钟
```

### 6.2 监控指标

- **存储使用率**: 各用户的存储配额使用情况
- **API调用**: 各云盘API的调用频率和成功率
- **错误率**: 上传/下载失败率和重试次数
- **性能指标**: 文件处理时间、传输速度
- **成本分析**: 各存储后端的成本对比

### 6.3 容灾备份

- **数据库备份**: 定期备份元数据和配置信息
- **多云备份**: 重要内容可配置多云盘备份
- **本地缓存**: 热点内容本地缓存加速访问
- **降级策略**: 云盘不可用时降级到数据库存储

## 7. 安全考虑

### 7.1 数据安全

- **Token加密**: 所有云盘token使用AES加密存储
- **内容加密**: 用户可选择对敏感内容加密
- **传输安全**: 所有API调用使用HTTPS
- **访问控制**: 严格的权限验证和授权检查

### 7.2 隐私保护

- **数据最小化**: 只获取必要的用户信息
- **用户控制**: 用户完全控制自己的存储选择
- **数据删除**: 支持彻底删除用户数据
- **审计日志**: 记录所有敏感操作

## 8. 性能优化

### 8.1 缓存策略

- **元数据缓存**: 用户配置和授权信息缓存
- **内容缓存**: 经常访问的内容本地缓存
- **CDN加速**: 静态资源和缩略图CDN分发

### 8.2 异步处理

- **文件上传**: 大文件异步上传到云盘
- **内容处理**: 图片/视频处理放入队列异步执行
- **同步任务**: 定期同步云盘配额和状态信息

## 9. 扩展规划

### 9.1 短期扩展 (3个月)

- 完成核心存储架构实现
- 支持主流云盘（Dropbox, Google Drive, 阿里云盘, WebDAV）
- 基础的图片和视频处理

### 9.2 中期扩展 (6个月)

- 支持更多云盘（如百度网盘）
- 高级内容处理功能 (OCR、语音识别)
- 协作功能增强

### 9.3 长期扩展 (12个月)

- AI内容分析和推荐
- 企业级功能 (团队协作、权限管理)
- 开放API和第三方集成
- 支持去中心化存储（IPFS、Filecoin等）

---

**文档版本**: v1.0  
**创建时间**: 2024年  
**最后更新**: 需要时更新  
**维护者**: 开发团队
