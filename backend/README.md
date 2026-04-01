# duptwo v1.0.0 数据登记平台

## 简介

duptwo 是一个轻量级数据登记平台，支持多渠道数据采集、后台管理、统计和报告导出。

## 快速开始

### 1. 解压

```bash
tar -xzf duptwo-v1.0.0-linux-amd64.tar.gz
cd duptwo-v1.0.0-linux-amd64
```

### 2. 配置

编辑 `conf/app.yaml`：

```yaml
base:
  port: 18421          # 服务端口
  serve_web: true       # 提供前端

database:
  type: sqlite         # sqlite / mysql / postgres
  path: ./data/registry.db

jwt:
  secret: CHANGE_THIS_TO_A_LONG_RANDOM_SECRET  # ⚠️ 必须修改
```

### 3. 创建目录并启动

```bash
mkdir -p data logs
./duptwo --config conf/app.yaml
```

### 4. 访问

- **Web UI**: http://localhost:18421
- **Swagger API**: http://localhost:18421/swagger
- **健康检查**: http://localhost:18421/health
- **默认账号**: `admin` / `admin123`

## 管理命令

```bash
# 重置 admin 密码
./duptwo reset-admin <新密码>

# 查看帮助
./duptwo --help
```

## 数据存储

| 数据库 | 配置 | 适用场景 |
|--------|------|---------|
| SQLite | `type: sqlite` | 开发、小规模 |
| MySQL | `type: mysql` | 生产环境 |
| PostgreSQL | `type: postgres` | 生产环境 |

## 生产部署建议

1. **修改密钥**
   ```yaml
   jwt:
     secret: 使用 `openssl rand -base64 32` 生成
   session:
     secret: 使用 `openssl rand -base64 32` 生成
   ```

2. **设置 HTTPS**
   - 使用 Nginx 反向代理
   - 或使用 Docker 部署（已包含 Nginx 配置）

3. **数据备份**
   ```bash
   # SQLite 备份
   cp data/registry.db backup/registry.db.$(date +%Y%m%d)

   # MySQL 备份
   mysqldump -u root -p duptwo > backup/duptwo.$(date +%Y%m%d).sql
   ```

## Docker 部署

详见 [DEPLOY.md](./DEPLOY.md)

## 技术支持

- GitHub: https://github.com/your-org/duptwo
