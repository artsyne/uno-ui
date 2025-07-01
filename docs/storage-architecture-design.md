# 多云盘存储架构设计方案

## 1. 项目背景

### 1.1 项目概述

- **项目名称**: idiary.ai - 智能日记管理系统
- **核心竞争力**: 支持多云盘个人存储，用户可以将日记内容存储到自己的云盘中
- **技术栈**: Go + Gin + MySQL + 多云盘API集成

### 1.2 业务特点

- **多用户系统**: 支持大量用户注册和使用
- **个人云盘**: 用户使用自己的云盘授权，非平台统一存储
- **多内容类型**: 支持文字、图片、视频、链接等多种内容形式
- **灵活存储**: 每个日记本可独立配置存储策略
- **数据加密**: 支持内容加密，保护用户隐私

## 2. 模块架构设计

### 2.1 整体架构

```
┌─────────────────────────────────────────────────────────┐
│                    Web层 (Gin Router)                   │
├─────────────────────────────────────────────────────────┤
│  User模块    │  Diary模块   │  Post模块    │  Task模块   │
│  用户管理    │  日记本管理   │  内容管理    │  任务处理   │
│  云盘授权    │  存储配置    │  多存储支持   │  后台作业   │
├─────────────────────────────────────────────────────────┤
│              存储抽象层 (Storage Layer)                  │
│    Database    │   Dropbox   │  GoogleDrive │  AliDrive │
├─────────────────────────────────────────────────────────┤
│                数据持久层 (MySQL)                        │
└─────────────────────────────────────────────────────────┘
```

### 2.2 模块职责划分

#### User模块

- **基础职责**: 用户注册、登录、认证、个人资料管理
- **扩展职责**:
  - 云盘授权管理（OAuth流程、Token刷新）
  - 用户存储偏好设置
  - 存储配额监控
  - 认证中间件提供

#### Diary模块

- **基础职责**: 日记本CRUD操作、成员管理
- **扩展职责**:
  - 存储策略配置管理
  - 调用Post模块获取内容列表（解决循环依赖）
  - 按内容类型配置不同存储策略

#### Post模块

- **基础职责**: 内容CRUD操作（文字、图片、视频、链接等）
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
User模块 (独立模块，提供服务)
   ↑
Post模块 (依赖User获取云盘授权)
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
ALTER TABLE br_user ADD COLUMN storage_quota_limit bigint unsigned DEFAULT 1073741824 COMMENT '存储配额限制(字节，默认1GB)';
```

#### br_diary (日记本表)

```sql
-- 日记本信息表，扩展存储配置字段
ALTER TABLE br_diary ADD COLUMN storage_config text DEFAULT NULL COMMENT '存储配置JSON，包含存储类型、云盘授权ID、加密设置、按内容类型的存储策略等';

-- storage_config JSON结构示例:
-- {
--   "default_type": "cloud",                    // 默认存储类型: database/cloud
--   "cloud_auth_id": 123,                       // 默认云盘授权ID
--   "encryption_enabled": true,                 // 是否启用加密
--   "content_strategies": {                     // 按内容类型的存储策略
--     "text": {                                 // 文字内容
--       "storage": "database",                  // 存储位置
--       "encryption": true                      // 是否加密
--     },
--     "image": {                                // 图片内容
--       "storage": "cloud",                     // 存储位置
--       "encryption": false,                    // 是否加密
--       "thumbnail": true,                      // 是否生成缩略图
--       "compression": 80                       // 压缩质量
--     },
--     "video": {                                // 视频内容
--       "storage": "cloud",                     // 存储位置
--       "encryption": false,                    // 是否加密
--       "transcoding": true                     // 是否转码
--     }
--   }
-- }
```

#### br_post (内容表)

```sql
-- 多媒体内容表，扩展存储相关字段
ALTER TABLE br_post ADD COLUMN storage_type varchar(20) DEFAULT 'database' COMMENT '存储类型: database/cloud';
ALTER TABLE br_post ADD COLUMN storage_path varchar(500) DEFAULT '' COMMENT '云盘存储路径或文件ID';
ALTER TABLE br_post ADD COLUMN is_encrypted tinyint unsigned DEFAULT 0 COMMENT '是否加密存储: 0否 1是';
ALTER TABLE br_post ADD COLUMN file_size bigint unsigned DEFAULT 0 COMMENT '文件大小(字节)';
ALTER TABLE br_post ADD COLUMN mime_type varchar(100) DEFAULT '' COMMENT 'MIME类型，如image/jpeg、video/mp4等';
ALTER TABLE br_post ADD COLUMN thumbnail_path varchar(500) DEFAULT '' COMMENT '缩略图路径(针对图片/视频)';

-- info_type字段说明 (已有字段，重新定义含义):
-- 0: 文字日记
-- 1: 图片内容
-- 2: 视频内容
-- 3: 链接分享
-- 4: 音频内容
-- 5: 文件附件
-- 6-99: 预留扩展
```

### 3.2 新增表结构

#### br_cloud_authorization (云盘授权表)

```sql
CREATE TABLE br_cloud_authorization (
  auth_id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '授权ID，主键',
  user_id bigint unsigned NOT NULL COMMENT '用户ID，关联br_user.user_id',
  provider varchar(50) NOT NULL COMMENT '云盘提供商: dropbox/googledrive/alidrive/baidupan/onedrive等',
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
  default_storage_type varchar(20) DEFAULT 'database' COMMENT '默认存储类型: database/cloud',
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
Service Layer (业务逻辑层)
    ↓
Storage Manager (存储管理器)
    ↓
Content Processor (内容处理器) → Encryption Layer (加密层)
    ↓
Storage Adapter (存储适配器)
    ↓
Storage Backend (存储后端: DB/Cloud)
```

### 4.2 存储策略矩阵

| 内容类型 | 推荐存储 | 加密支持 | 处理功能         | 备注             |
| -------- | -------- | -------- | ---------------- | ---------------- |
| 文字日记 | Database | ✅       | Markdown渲染     | 快速检索         |
| 图片     | Cloud    | ✅       | 缩略图生成、压缩 | 节省服务器空间   |
| 视频     | Cloud    | ❌       | 转码、截图       | 大文件必须云存储 |
| 音频     | Cloud    | ✅       | 格式转换         | 支持多种音频格式 |
| 文件     | Cloud    | ✅       | 病毒扫描         | 通用文件支持     |
| 链接     | Database | ❌       | 元数据抓取       | 链接预览功能     |

### 4.3 存储适配器设计

#### Database Adapter

- **职责**: 处理数据库存储的内容
- **支持内容**: 文字、链接、小文件
- **特点**: 快速访问、事务支持、全文检索

#### Cloud Adapters

每个云盘一个独立适配器:

- **Dropbox Adapter**

  - API版本: v2
  - 支持功能: 文件上传/下载、分享链接、版本管理
  - 限制: 单文件150MB，API限流

- **Google Drive Adapter**

  - API版本: v3
  - 支持功能: 文件操作、OCR识别、协作编辑
  - 限制: 日请求限制、文件大小限制

- **阿里云盘 Adapter**
  - API版本: 最新版
  - 支持功能: 大文件上传、秒传、分享
  - 特点: 国内访问快、容量大

### 4.4 内容处理器设计

#### Text Processor

- **功能**: Markdown解析、敏感词过滤、全文索引
- **加密**: 支持AES-256加密
- **存储**: 优先数据库，超大文本可存云盘

#### Image Processor

- **功能**:
  - 格式转换 (HEIC→JPEG)
  - 尺寸调整和压缩
  - EXIF信息清理
  - 缩略图生成 (多种尺寸)
- **加密**: 支持图片加密 (性能考虑可选)
- **存储**: 原图云盘，缩略图本地缓存

#### Video Processor

- **功能**:
  - 格式标准化 (转换为MP4 H.264)
  - 视频压缩 (多码率)
  - 截图生成 (封面和预览)
  - 时长和分辨率检测
- **加密**: 暂不支持 (处理复杂度高)
- **存储**: 必须云盘存储

#### File Processor

- **功能**:
  - 文件类型检测
  - 病毒扫描 (集成ClamAV)
  - 文档预览生成 (PDF/Office)
  - 压缩包内容扫描
- **加密**: 支持通用文件加密
- **存储**: 根据大小选择存储位置

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

```go
func (m *StorageManager) RouteStorage(post *Post, config *StorageConfig) StorageAdapter {
    contentType := post.InfoType
    strategy := config.ContentStrategies[contentType]

    switch strategy.Storage {
    case "database":
        return m.databaseAdapter
    case "cloud":
        return m.getCloudAdapter(config.CloudAuthID)
    default:
        return m.databaseAdapter // 默认回退
    }
}
```

### 5.4 加密实现

```go
// 使用AES-256-GCM加密
func (e *Encryptor) Encrypt(data []byte, userKey string) ([]byte, error) {
    // 1. 派生密钥 (PBKDF2)
    key := pbkdf2.Key([]byte(userKey), salt, 10000, 32, sha256.New)

    // 2. AES-GCM加密
    block, _ := aes.NewCipher(key)
    gcm, _ := cipher.NewGCM(block)

    nonce := make([]byte, gcm.NonceSize())
    rand.Read(nonce)

    ciphertext := gcm.Seal(nonce, nonce, data, nil)
    return ciphertext, nil
}
```

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
    free_user: 1073741824 # 1GB
    premium_user: 10737418240 # 10GB

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
- 支持Dropbox和Google Drive
- 基础的图片和视频处理

### 9.2 中期扩展 (6个月)

- 新增国内云盘支持 (阿里云盘、百度网盘)
- 高级内容处理功能 (OCR、语音识别)
- 协作功能增强

### 9.3 长期扩展 (12个月)

- AI内容分析和推荐
- 企业级功能 (团队协作、权限管理)
- 开放API和第三方集成

---

**文档版本**: v1.0  
**创建时间**: 2024年  
**最后更新**: 需要时更新  
**维护者**: 开发团队
