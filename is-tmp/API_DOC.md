# 数据上传记录管理平台 - API 文档

## 服务地址
- 后端服务: http://localhost:18421
- 前端页面: http://localhost:18421

## 公开接口（无需登录）

### 1. 上传记录
- **URL**: `POST /public/upload-records`
- **说明**: 创建一条上传记录

**请求示例 (curl):**
```bash
curl -X POST http://localhost:18421/public/upload-records \
  -H "Content-Type: application/json" \
  -d '{
    "dataType": "disk_backup",
    "filePath": "/mnt/storage/backup_001",
    "fileSize": 5368709120,
    "uploader": "张三",
    "status": "completed",
    "remark": "服务器日常备份"
  }'
```

**请求参数:**
| 字段 | 必填 | 类型 | 说明 | 示例 |
|------|------|------|------|------|
| dataType | ✅ | string | 数据类型 | disk_backup, db_sync, file_upload |
| filePath | ✅ | string | 上传路径 | /mnt/storage/data |
| uploader | ✅ | string | 上传人 | 张三 |
| fileSize | ❌ | int64 | 文件大小(bytes) | 1073741824 |
| recordCount | ❌ | int | 记录条数 | 1000 |
| status | ❌ | string | 状态 | pending/processing/completed/failed |
| remark | ❌ | string | 备注 | 备注信息 |

**响应示例:**
```json
{
  "code": 200,
  "message": "创建成功",
  "data": {
    "id": 1,
    "serialNo": "DISK_BACKUP-20260322-0001",
    "dataType": "disk_backup",
    "filePath": "/mnt/storage/backup_001",
    "fileSize": 5368709120,
    "recordCount": 0,
    "uploader": "张三",
    "status": "completed",
    "remark": "服务器日常备份",
    "metadata": "",
    "createdAt": "2026-03-22T15:46:39+08:00",
    "updatedAt": "2026-03-22T15:46:39+08:00"
  }
}
```

### 2. 按流水号查询
- **URL**: `GET /public/upload-records/:serialNo`
- **说明**: 根据流水号查询上传记录

**请求示例 (curl):**
```bash
curl http://localhost:18421/public/upload-records/DISK_BACKUP-20260322-0001
```

## 管理接口

### 3. 获取上传记录列表
- **URL**: `GET /api/upload-records`
- **说明**: 分页查询上传记录

**请求参数:**
| 参数 | 类型 | 说明 | 示例 |
|------|------|------|------|
| page | int | 页码 | 1 |
| pageSize | int | 每页数量 | 20 |
| dataType | string | 数据类型筛选 | disk_backup |
| status | string | 状态筛选 | completed |
| uploader | string | 上传人筛选 | 张三 |
| startDate | string | 开始日期 | 2026-03-01 |
| endDate | string | 结束日期 | 2026-03-31 |
| keyword | string | 关键词搜索 | backup |

**请求示例 (curl):**
```bash
# 查询所有记录
curl http://localhost:18421/api/upload-records

# 带筛选条件
curl "http://localhost:18421/api/upload-records?page=1&pageSize=10&status=completed"

# 按日期范围筛选
curl "http://localhost:18421/api/upload-records?startDate=2026-03-01&endDate=2026-03-31"
```

### 4. 获取统计数据
- **URL**: `GET /api/upload-records/statistics`
- **说明**: 获取上传记录的统计信息

**请求示例 (curl):**
```bash
curl http://localhost:18421/api/upload-records/statistics
```

**响应示例:**
```json
{
  "code": 200,
  "data": {
    "todayCount": 5,
    "todaySize": 5368709120,
    "todaySizeStr": "5.00 GB",
    "weekCount": 20,
    "monthCount": 50,
    "totalCount": 100,
    "totalSize": 107374182400,
    "totalSizeStr": "100.00 GB",
    "trend": [
      {"date": "2026-03-22", "count": 5, "totalSize": 5368709120}
    ],
    "byStatus": [
      {"status": "completed", "count": 80},
      {"status": "pending", "count": 20}
    ],
    "byDataType": [
      {"dataType": "disk_backup", "count": 60},
      {"dataType": "db_sync", "count": 40}
    ]
  }
}
```

### 5. 获取最近记录
- **URL**: `GET /api/upload-records/recent`
- **说明**: 获取最近的上传记录

**请求示例 (curl):**
```bash
curl "http://localhost:18421/api/upload-records/recent?limit=10"
```

## 浏览器控制台测试

打开 http://localhost:18421 后，按 F12 打开开发者工具，在 Console 中执行：

```javascript
// 上传记录
fetch('http://localhost:18421/public/upload-records', {
  method: 'POST',
  headers: {'Content-Type': 'application/json'},
  body: JSON.stringify({
    dataType: 'disk_backup',
    filePath: '/mnt/storage/test',
    fileSize: 1024000,
    uploader: '测试用户',
    status: 'completed'
  })
}).then(r => r.json()).then(console.log)

// 查询列表
fetch('http://localhost:18421/api/upload-records')
  .then(r => r.json()).then(console.log)

// 查看统计
fetch('http://localhost:18421/api/upload-records/statistics')
  .then(r => r.json()).then(console.log)
```

## Postman / Apifox 测试

- **URL**: http://localhost:18421/public/upload-records
- **Method**: POST
- **Headers**: Content-Type: application/json
- **Body** (raw JSON):
```json
{
  "dataType": "disk_backup",
  "filePath": "/mnt/storage/backup",
  "fileSize": 5368709120,
  "uploader": "张三",
  "status": "completed",
  "remark": "备注"
}
```
